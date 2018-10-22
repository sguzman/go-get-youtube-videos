package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    _ "github.com/lib/pq"
    "net/http"
    "runtime"
    "strings"
)

const (
    connStr = "user=postgres dbname=youtube host=192.168.1.63 port=30000 sslmode=disable"
)

func insert(db *sql.DB, serial string) {
    sqlInsert := "INSERT INTO youtube.entities.videos (serial) VALUES ($1) ON CONFLICT (serial) DO NOTHING"

    _, err := db.Exec(sqlInsert, serial)
    if err != nil {
        panic(err)
    }
}

func connection() *sql.DB {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }

    return db
}

func channels() string {
    sqlStr := "SELECT serial FROM youtube.entities.channels ORDER BY RANDOM() LIMIT 1"
    db := connection()
    defer func() {
        err := db.Close()
        if err != nil {
            panic(err)
        }
    }()

    row, err := db.Query(sqlStr)
    if err != nil {
        panic(err)
    }

    var serial string
    if row.Next() {
        err = row.Scan(&serial)
        if err != nil {
            panic(err)
        }
    }

    return serial
}

func doc(channel string) *goquery.Document {
    url := fmt.Sprintf("https://www.youtube.com/channel/%s/videos?flow=grid&view=0", channel)

    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }

    userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
    req.Header.Add("User-Agent", userAgent)

    res, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer func() {
        err := res.Body.Close()
        if err != nil {
            panic(err)
        }
    }()

    if res.StatusCode != 200 {
        panic(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        panic(err)
    }

    return doc
}

func mapJson(d *goquery.Document) interface{} {
    prefix := "\n    window[\"ytInitialData\"] = "
    suffix := ";\n"

    var jsonMap interface{}
    d.Find("script").Each(func(i int, s *goquery.Selection) {
        text := s.Text()
        if strings.HasPrefix(text, prefix) {
            idx := strings.Index(text, suffix)
            subText := text[len(prefix):idx]

            err := json.Unmarshal([]byte(subText), &jsonMap)
            if err != nil {
                panic(err)
            }
        }
    })

    return jsonMap
}

func videoSerials(obj interface{}) []string {
    m1, ok1 := obj.(map[string]interface{})
    if !ok1 {
        panic("Failed type cast")
    }

    m2, ok2 := m1["contents"].(map[string]interface{})
    if !ok2 {
        panic("Failed type cast")
    }

    m3, ok3 := m2["twoColumnBrowseResultsRenderer"].(map[string]interface{})
    if !ok3 {
        panic("Failed type cast")
    }

    m4, ok4 := m3["tabs"].([]interface{})
    if !ok4 {
        panic("Failed type cast")
    }

    m5, ok5 := m4[1].(map[string]interface{})
    if !ok5 {
        panic("Failed type cast")
    }

    m6, ok6 := m5["tabRenderer"].(map[string]interface{})
    if !ok6 {
        panic("Failed type cast")
    }

    m7, ok7 := m6["content"].(map[string]interface{})
    if !ok7 {
        panic("Failed type cast")
    }

    m8, ok8 := m7["sectionListRenderer"].(map[string]interface{})
    if !ok8 {
        panic("Failed type cast")
    }

    m9, ok9 := m8["contents"].([]interface{})
    if !ok9 {
        panic("Failed type cast")
    }

    m10, ok10 := m9[0].(map[string]interface{})
    if !ok10 {
        panic("Failed type cast")
    }

    m11, ok11 := m10["itemSectionRenderer"].(map[string]interface{})
    if !ok11 {
        panic("Failed type cast")
    }

    m12, ok12 := m11["contents"].([]interface{})
    if !ok12 {
        panic("Failed type cast")
    }

    m13, ok13 := m12[0].(map[string]interface{})
    if !ok13 {
        panic("Failed type cast")
    }

    m14, ok14 := m13["gridRenderer"].(map[string]interface{})
    if !ok14 {
        panic("Failed type cast")
    }

    m15, ok15 := m14["items"].([]interface{})
    if !ok15 {
        panic("Failed type cast")
    }

    serials := make([]string, len(m15))
    for i := 0; i < len(m15); i++ {
        item, ok := m15[i].(map[string]interface{})
        if !ok {
            panic("Failed type cast")
        }

        item1, ok16 := item["gridVideoRenderer"].(map[string]interface{})
        if !ok16 {
            panic("Failed type cast")
        }

        item2, ok17 := item1["videoId"].(string)
        if !ok17 {
            panic("Failed type cast")
        }

        serials[i] = item2
        fmt.Println(item2)
    }

    return serials
}

func process() {
    channel := channels()
    d := doc(channel)
    inter := mapJson(d)
    vids := videoSerials(inter)
    conn := connection()

    for i := 0; i < len(vids); i++ {
        v := vids[i]
        insert(conn, v)
    }
}

func main() {
   for {
       process()
       runtime.GC()
   }
}
