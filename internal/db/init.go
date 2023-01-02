package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.boomerangapp.ir/back/pg/configs"
)

var DB *sql.DB

func InitDB() error {

	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&tls=false",
		configs.Get().DBConfigs.Mysql.Username,
		configs.Get().DBConfigs.Mysql.Password,
		configs.Get().DBConfigs.Mysql.Host,
		configs.Get().DBConfigs.Mysql.DB)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return err
	}

	DB = db
	return db.Ping()
}
