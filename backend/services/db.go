package services

import (
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	Db *sqlx.DB
)

func DbConnect() {
	dbuser := os.Getenv("DB_User")
	dbpassword := os.Getenv("DB_Password")
	dbhost := os.Getenv("DB_Server")
	dbport := os.Getenv("DB_Port")
	dbname := os.Getenv("DB_Name")
	tz := os.Getenv("TZ")

	if dbhost == "" {
		dbhost = "localhost"
	}
	if dbport == "" {
		dbport = "3306"
	}
	if tz == "" {
		tz = "UTC"
	}

	log.Println("Connecting to database with " + dbuser + ":........@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?parseTime=true&loc=" + strings.ReplaceAll(tz, "/", "%2F"))

	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?parseTime=true&loc=" + strings.ReplaceAll(tz, "/", "%2F")
	var err error

	Db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Connected to database!")
	log.Println("")
}
