package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "runtime"
    "time"
)

const (
    connStr = "user=root dbname=youtube host=localhost port=5432 sslmode=disable"
)

func insert(db *sql.DB, channel string) {
    sqlInsert := "INSERT INTO youtube.entities.channels (serial) VALUES ($1) ON CONFLICT (serial) DO NOTHING"

    _, err := db.Exec(sqlInsert, channel)
    if err != nil {
        panic(err)
    }
}

func main() {
    dur, err := time.ParseDuration("3s")

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }

    defer func() {
        err := db.Close()
        if err != nil {
            panic(err)
        }
    }()

    sqlStr := "SELECT count(*) FROM youtube.entities.channels"

    var count uint64
    row, err := db.Query(sqlStr)
    if err != nil {
        panic(err)
    }

    if row.Next() {
        err = row.Scan(&count)
        if err != nil {
            panic(err)
        }
    } else {
        panic("No count")
    }

    fmt.Println("Count", count)

    /*for rows.Next() {
        var (
            id   uint64
            name string
        )
        if err := rows.Scan(&id, &name); err != nil {
            panic(err)
        } else {
            fmt.Printf("id = %d, name = %s\n", id, name)
        }
    }*/

    runtime.GC()
    time.Sleep(dur)
}
