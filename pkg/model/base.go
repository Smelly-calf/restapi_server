package model

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"restapi_server/pkg/config"
)

const (
	MaxConns  = 20
	IdleConns = 10
)

var RunoobDB *sql.DB

// db 对象连接池
func init() {
	log.Println("Connecting PostgreSQL....")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.PgHost, config.PgPort, config.PgUser, config.PgPassword, config.PgRunoobdb)
	var err error
	RunoobDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Connect PG Failed: ", err)
	}
	RunoobDB.SetMaxOpenConns(MaxConns)
	RunoobDB.SetMaxIdleConns(IdleConns)

	err = RunoobDB.Ping()
	if err != nil {
		log.Fatal("Ping GP Failed: ", err)
	}
	fmt.Println("PG Successfull Connected!")

}
