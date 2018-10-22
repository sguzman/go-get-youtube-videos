package main

import (
    "database/sql"
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

func conn() *sql.DB {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }

    return db
}

func channels() string {
    sqlStr := "SELECT serial FROM youtube.entities.channels ORDER BY RANDOM() LIMIT 1"
    db := conn()
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

func find_script(d *goquery.Document) {
    d.Find("script").Each(func(i int, s *goquery.Selection) {
        text := s.Text()
        if strings.HasPrefix(text, "\n    window[\"ytInitialData\"]") {
            fmt.Println(text)
        }

    })
}

func process() {
    channel := channels()
    d := doc(channel)
    find_script(d)

    fmt.Println(channel)
}

func main() {
   for {
       process()
       runtime.GC()
   }
}
