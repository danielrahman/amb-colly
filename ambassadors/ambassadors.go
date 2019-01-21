package ambassadors

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "127.0.0.1"
	user     = "root"
	password = "root"
	dbname   = "ambassadors"
)

type DbAmbassadors struct {
	Db *sql.DB
}

func (g *DbAmbassadors) ConnectDatabase() (*sql.DB, error) {
	msqlInfo := fmt.Sprintf("host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)
	db, err := sql.Open("mysql", msqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("not connected!")
		return nil, err
	}
	g.Db = db
	return db, nil
}

func (g DbAmbassadors) UpdateDatabase(sqlStatement string) {
	_, err := g.Db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}
