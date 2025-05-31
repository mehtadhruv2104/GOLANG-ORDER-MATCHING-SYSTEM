package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)


var dB *sql.DB


func ConnectToMySQL()(*sql.DB, error) {
	dbURL := os.Getenv("DB_URL")
	dB, err := sql.Open("mysql", dbURL)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return nil,err
	}
    if err != nil {
        fmt.Printf("Error creating database: %v", err)
    }
	fmt.Println("Connected to the database successfully")
	return dB, nil
}

func GetDB() *sql.DB {

	if dB == nil {
		fmt.Println("Database connection is not established")
		return nil
	}
	return dB
}



