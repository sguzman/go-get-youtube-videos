package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "runtime"
)

type Channel struct {
    id uint64
    serial string
}

const (
    connStr = "user=root dbname=youtube host=localhost port=5432 sslmode=disable"
)

func insert(db *sql.DB, channel Channel) {
    sqlInsert := "INSERT INTO youtube.entities.videos (id, serial) VALUES ($1, $2) ON CONFLICT (serial) DO NOTHING"

    _, err := db.Exec(sqlInsert, channel.id, channel.serial)
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

func channels() Channel {
    sqlStr := "SELECT id, serial FROM youtube.entities.channels ORDER BY RANDOM() LIMIT 1"
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

    if row.Next() {
        var (
            id uint64
            serial string
        )
        err = row.Scan(&id, &serial)
        if err != nil {
            panic(err)
        }

        var channel Channel
        channel.id = id
        channel.serial = serial

        return channel
    } else {
        panic("No count")
    }
}

func main() {
   for {
       channel := channels()
       fmt.Println(channel.id, channel.serial)

       runtime.GC()
   }
}
