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

func insert(db *sql.DB, id uint64, serial string) {
    sqlInsert := "INSERT INTO youtube.entities.videos (id, serial) VALUES ($1, $2) ON CONFLICT (serial) DO NOTHING"

    _, err := db.Exec(sqlInsert, id, serial)
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

func count() uint64 {
    sqlStr := "SELECT count(*) FROM youtube.entities.channels"
    db := conn()
    defer func() {
        err := db.Close()
        if err != nil {
            panic(err)
        }
    }()

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

        return count
    } else {
        panic("No count")
    }
}

func main() {
    dur, err := time.ParseDuration("3s")
    if err != nil {
        panic(err)
    }

    count := count()

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
