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
	db, err := sql.Open("mysql", "b50e4aa373e192:b27a80a3@eu-cdbr-west-02.cleardb.net/heroku_ce0c3128d87231f")
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

func (g DbAmbassadors) GetData(sqlRow string, sqlTable string) *sql.Rows {
	sqlQuery := "SELECT " + sqlRow + " FROM " + sqlTable
	rows, err := g.Db.Query(sqlQuery)
	if err != nil {
		panic(err)
	}
	return rows
}
