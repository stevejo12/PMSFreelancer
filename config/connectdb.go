package config

import (
	"database/sql"

	// _ => make sure we use mysql here
	_ "github.com/go-sql-driver/mysql"
)

// DB => Database connection
var DB *sql.DB

// ConnectToDB => Establish connection to db
func ConnectToDB() {
	var err error

	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "thesis"
	//DB, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	DB, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(159.89.202.223:3306)/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	if DB == nil {
		panic("DB is null")
	}
}
