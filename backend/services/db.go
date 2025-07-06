package services

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	Db    *sqlx.DB
	Stmts map[string]*sql.Stmt = make(map[string]*sql.Stmt)
)

func DbCloseStmts() {
	fn := "CloseStmts"

	slog.Debug(fmt.Sprintf("%v: Closing statements", fn))

	for key, stmt := range Stmts {
		err := stmt.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("%v: Failed closing statement %s: %v", fn, key, err))
		}
	}

	Stmts = make(map[string]*sql.Stmt)
	slog.Debug(fmt.Sprintf("%v: Closing statements completed", fn))
}

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

func dbPrepareStmt(key string, query string) (*sql.Stmt, error) {
	fn := fmt.Sprintf("dbPrepareStmt(%s)", key)

	stmt, found := Stmts[key]
	if found {
		// log.Printf("%v: Already prepared", fn)
		return stmt, nil
	}

	stmt, err := Db.Prepare(query)
	if err != nil {
		log.Printf("%v: Failed preparing stmt %s: %v", fn, key, err)
		return nil, err
	}
	Stmts[key] = stmt

	// log.Printf("%v: prepared", fn)
	return stmt, nil
}
