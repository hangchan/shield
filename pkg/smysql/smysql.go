package smysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

func Search(mc MysqlConn, user string) string {
	sqlQuery := "SELECT user,host FROM mysql.user WHERE user = ?"
	return query(mc, user, sqlQuery)
}

func DropUser(mc MysqlConn, user string, host string) {
	sqlQuery := fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s'", user, host)
	exec(mc, sqlQuery)
}

func query(mc MysqlConn, user string, query string) string {
	var resUser string
	var resHost string
	var resultsArr []string

	db := dbConn(mc)
	defer db.Close()

	sqlQuery := query
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

func exec(mc MysqlConn, query string) {
	db := dbConn(mc)
	defer db.Close()

	sqlQuery := query
	stmt, err := db.Prepare(sqlQuery)
	util.LogError(err)
	_,err = stmt.Exec()
	util.LogError(err)

	stmt.Close()
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

/*func Delete(mc MysqlConn, user string) string {
	var resUser string
	var resHost string
	var resultsArr []string

	db := dbConn(mc)
	defer db.Close()

	sqlQuery := "DROP USER IF EXISTS ?"
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
}*/

