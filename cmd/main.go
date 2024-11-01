package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"

	"github.com/jayden1905/nexusnova/cmd/api"
	"github.com/jayden1905/nexusnova/config"
	"github.com/jayden1905/nexusnova/db"
)

func main() {
	db, dbErr := db.NewMySQLStorage(mysql.Config{
		User:              config.Envs.DBUser,
		Passwd:            config.Envs.DBPasswd,
		Addr:              config.Envs.DBAddr,
		DBName:            config.Envs.DBName,
		Net:               "tcp",
		AllowOldPasswords: true,
		ParseTime:         true,
	})
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	serverErr := server.Run()
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection established")
}
