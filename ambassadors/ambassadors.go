package ambassadors

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type DbAmbassadors struct {
	Db *sql.DB
}

func (g *DbAmbassadors) ConnectDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost)/ambassadors")
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
