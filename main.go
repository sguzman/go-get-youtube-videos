package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "runtime"
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

/*func scrapeChannel(channel string) []string {
    //url := fmt.Sprintf("https://www.youtube.com/channel/%s/videos?flow=grid&view=0", channel.serial)
}*/

func main() {
   for {
       channel := channels()
       fmt.Println(channel)

       runtime.GC()
   }
}
