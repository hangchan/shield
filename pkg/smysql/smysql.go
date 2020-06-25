package smysql

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/hangchan/shield/pkg/util"
	"strings"
)

type MysqlConn struct {
	DbDriver 	string
	DbUser		string
	DbPass		string
	DbName		string
	DbAddress	string
}

func Test(mc MysqlConn, user string) string {
	var resUser string
	var resHost string
	var resultsArr []string

	db := dbConn(mc)
	defer db.Close()

	sqlQuery := "SELECT user,host FROM mysql.user WHERE user = ?"
	stmt, err := db.Prepare(sqlQuery)
	util.LogError(err)
	res, err := stmt.Query(user)
	util.LogError(err)

	for res.Next() {
		err = res.Scan(&resUser, &resHost)
		util.LogError(err)
		resultsArr = append(resultsArr, resHost)
	}

	res.Close()
	stmt.Close()

	results := strings.Join(resultsArr, ",")

	return results
}

func dbConn(mc MysqlConn) (db *sql.DB) {
	dbDriver	:= mc.DbDriver
	dbUser		:= mc.DbUser
	dbPass		:= mc.DbPass
	dbName		:= mc.DbName
	dbAddress	:= mc.DbAddress
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@("+dbAddress+")/"+dbName)
	util.LogError(err)
	return db
}

