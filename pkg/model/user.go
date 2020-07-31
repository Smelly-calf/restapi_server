package model

import (
	"log"
	"time"
)

// 用户数据集
type UserModel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Type      string    `json:"type",default:"user"`
}

func SelectAllUsers() []*UserModel {
	rows, err := RunoobDB.Query("SELECT * FROM  \"users\"")
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	res := make([]*UserModel, 0)
	for rows.Next() {
		var m UserModel
		if err := rows.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
			continue
		}
		res = append(res, &m)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("PG Query Failed: ", err)
	}
	rows.Close()
	return res
}

func SelectUserByName(name string) *UserModel {
	row := RunoobDB.QueryRow("SELECT * FROM users WHERE name=$1", name)
	var m UserModel
	err := row.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		log.Fatal("PG Query Failed: ", err)
	}
	return &m
}

func InsertUser(name string) *UserModel {
	stmt, err := RunoobDB.Prepare("INSERT INTO users (name) values ($1)")
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatal("PG Statements Exec Wrong: ", err)
	}
	return SelectUserByName(name)
}
