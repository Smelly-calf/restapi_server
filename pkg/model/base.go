package model

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"restapi_server/pkg/config"
	"time"
)

const (
	MaxConns  = 20
	IdleConns = 10
)

var RunoobDB *sql.DB

// db 对象连接池
func init() {
	log.Println("Connecting Postgres....")

	pgInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.PgHost, config.PgPort, config.PgUser, config.PgPassword, config.PgRunoobdb)
	var err error

	start := time.Now()
	RunoobDB, err = sql.Open("postgres", pgInfo)
	if err != nil {
		log.Fatal("Connect PG Failed: ", err)
	}
	RunoobDB.SetMaxOpenConns(MaxConns)
	RunoobDB.SetMaxIdleConns(IdleConns)

	end := time.Now()
	log.Printf("conn time: %d ms", end.Sub(start).Microseconds())

	log.Println("PG Successful Connected!")
}
