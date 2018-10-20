package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "runtime"
    "time"
)

func main() {
    connStr := "user=root dbname=youtube host=localhost port=5432 sslmode=disable"
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

    sqlStr := "SELECT id, serial FROM youtube.entities.channels"

    rows, err := db.Query(sqlStr)
    if err != nil {
        panic(err)
    }

    for rows.Next() {
        var (
            id   uint64
            name string
        )
        if err := rows.Scan(&id, &name); err != nil {
            panic(err)
        } else {
            fmt.Printf("id = %d, name = %s\n", id, name)
        }
    }

    runtime.GC()
    time.Sleep(dur)
}
