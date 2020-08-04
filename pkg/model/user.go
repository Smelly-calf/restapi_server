package model

import (
	"log"
	"time"
)

// 用户数据集
type UserModel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SelectAllUsers() ([]*UserModel, error) {
	rows, err := RunoobDB.Query("SELECT * FROM  \"users\"")
	if err != nil {
		log.Println("PG Statements Wrong: ", err)
		return nil, err
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
		log.Println("PG Query Failed: ", err)
		return nil, err
	}
	rows.Close()
	return res, nil
}

func SelectUserByName(name string) (*UserModel, error) {
	row := RunoobDB.QueryRow("SELECT * FROM users WHERE name=$1", name)
	var m UserModel
	err := row.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		log.Println("PG Query Failed: ", err)
		return nil, err
	}
	return &m, nil
}

func InsertUser(name string) (*UserModel, error) {
	stmt, err := RunoobDB.Prepare("INSERT INTO users (name) values ($1)")
	if err != nil {
		log.Println("PG Statements Wrong: ", err)
		return nil, err
	}
	_, err = stmt.Exec(name)
	if err != nil {
		log.Println("PG Statements Exec Wrong: ", err)
		return nil, err
	}
	return SelectUserByName(name)
}
