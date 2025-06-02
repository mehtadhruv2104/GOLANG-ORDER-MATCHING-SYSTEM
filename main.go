package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/db"
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg"
)

func init(){
	err := godotenv.Load(".env")
    if err != nil {
        fmt.Println("Error loading .env file:", err)
    }	
}

func main(){
	dB,err := db.ConnectToMySQL()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	router := pkg.StartEngine(dB)
	router.Run(":8080")
	fmt.Println("Server is running on port 8080")
}