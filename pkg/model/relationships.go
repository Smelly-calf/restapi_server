package model

import (
	"fmt"
	"time"
)

// 用户关系数据集
type RelationshipModel struct {
	ID         int       `json:"-"`
	UserID     int       `json:"-"`
	FollowerID int       `json:"user_id"`
	State      string    `json:"state"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	Type       string    `json:"type",default:"relationship"`
}

func GetRelationshipsByUserID(userID int) []*RelationshipModel {
	rows, err := RunoobDB.Query("SELECT * FROM relationships WHERE user_id=$1", userID)
	if err != nil {
		panic(fmt.Sprintf("PG Statements Wrong: %v", err))
	}
	res := make([]*RelationshipModel, 0)
	for rows.Next() {
		var m RelationshipModel
		if err := rows.Scan(&m.ID, &m.UserID, &m.FollowerID, &m.State, &m.CreatedAt, &m.UpdatedAt); err != nil {
			continue
		}
		res = append(res, &m)
	}
	if err := rows.Err(); err != nil {
		panic(fmt.Sprintf("PG Query Failed: %v", err))
	}
	rows.Close()
	return res
}

func GetOneRelation(userID int, followerID int) *RelationshipModel {
	row := RunoobDB.QueryRow("SELECT * FROM relationships WHERE user_id=$1 and follower_id=$2", userID, followerID)
	var m RelationshipModel
	err := row.Scan(&m.ID, &m.UserID, &m.FollowerID, &m.State, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil
	}
	return &m
}

func InsertRelationship(userID int, followerID int, state string) *RelationshipModel {
	stmt, err := RunoobDB.Prepare("INSERT INTO relationships (user_id, follower_id, state) values ($1,$2,$3)")
	if err != nil {
		panic(fmt.Sprintf("PG Statements Wrong: %v", err))
	}
	_, err = stmt.Exec(userID, followerID, state)
	if err != nil {
		panic(fmt.Sprintf("PG Statements Exec Wrong: %v", err))
	}

	if err != nil {
		panic(fmt.Sprintf("PG Statements Exec Wrong: %v", err))
	}
	return GetOneRelation(userID, followerID)
}

func UpdateRelation(userID int, followerID int, state string) *RelationshipModel {
	stmt, err := RunoobDB.Prepare("UPDATE relationships set state=$1 WHERE user_id=$2 AND follower_id=$3")
	if err != nil {
		panic(fmt.Sprintf("PG Statements Wrong: %v", err))
	}
	// 反向关系
	inverse := GetOneRelation(followerID, userID)
	// 互相喜欢 == "matched"
	if inverse != nil && state == "liked" && (inverse.State == "liked" || inverse.State == "matched") {
		// 事物: 均更新为 matched
		tx, _ := RunoobDB.Begin()
		state := "matched"
		stmt.Exec(state, userID, followerID)
		stmt.Exec(state, followerID, userID)
		err = tx.Commit()
	} else if inverse != nil && state == "disliked" && inverse.State == "matched" {
		// 事物: inverse.State 更新为 liked
		tx, _ := RunoobDB.Begin()
		stmt.Exec(state, userID, followerID)
		stmt.Exec("liked", followerID, userID)
		err = tx.Commit()
	} else {
		_, err = stmt.Exec(state, userID, followerID)
	}

	if err != nil {
		panic(fmt.Sprintf("PG Statements Exec Wrong: %v", err))
	}
	return GetOneRelation(userID, followerID)
}
