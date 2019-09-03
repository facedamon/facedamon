package sql

import (
	"fmt"
	"log"

	"github.com/facedamon/meta-model/pkg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	sec, err := pkg.Cfg.GetSection("database")
	if nil != err {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}
	dbType := sec.Key("type").String()
	dbSchema := sec.Key("schema").String()
	dbUser := sec.Key("user").String()
	dbPwd := sec.Key("pwd").String()
	dbHost := sec.Key("host").String()
	maxId, _ := sec.Key("max_id").Int()
	maxOpen, _ := sec.Key("max_open").Int()

	db, err = sqlx.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPwd, dbHost, dbSchema))
	if nil != err {
		log.Fatalf(" Fail to open the '%s' connection '%v'", dbType, err)
	}
	err = db.Ping()
	if nil != err {
		log.Fatalf("Fail to Ping the '%s' for schema '%s' %v", dbType, dbSchema, err)
	}
	log.Printf("'%s' for schema '%s' Pinged", dbType, dbSchema)

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxId)
}

func CloseDB() {
	defer db.Close()
}
