package storage 

import (
    "database/sql"
    "fmt"
    "time"

   
)

func StoreResult(result int) error {
    db, err := sql.Open("postgres", "postgres://username:password@localhost/database?sslmode=disable")
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    defer db.Close()

    _, err = db.Exec("INSERT INTO results (value, created_at) VALUES ($1, $2)", result, time.Now())
    if err != nil {
        return fmt.Errorf("failed to insert result: %w", err)
    }
    return nil
}
