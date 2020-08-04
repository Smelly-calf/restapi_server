package model

import (
	"fmt"
	"log"
	"restapi_server/pkg/constant"
	"time"
)

//todo
// 用户关系数据集
type RelationshipModel struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	FollowerID int       `json:"follower_id"`
	State      string    `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

var DefaultRelationModel RelationshipModel

func init() {
	DefaultRelationModel = RelationshipModel{}
}

//todo error return
func (r *RelationshipModel) GetRelationshipsByUserID(userID int) ([]*RelationshipModel, error) {
	rows, err := RunoobDB.Query("SELECT * FROM relationships WHERE user_id=$1", userID)
	if err != nil {
		log.Printf(fmt.Sprintf("PG Statements Wrong: %v", err))
		return nil, err
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
		return nil, err
	}
	rows.Close()
	return res, nil
}

func (r *RelationshipModel) GetOneRelation(userID int, followerID int) (*RelationshipModel, error) {
	row := RunoobDB.QueryRow("SELECT * FROM relationships WHERE user_id=$1 and follower_id=$2", userID, followerID)
	var m RelationshipModel
	err := row.Scan(&m.ID, &m.UserID, &m.FollowerID, &m.State, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *RelationshipModel) InsertRelationship(userID int, followerID int, state string) (*RelationshipModel, error) {
	stmt, err := RunoobDB.Prepare("INSERT INTO relationships (user_id, follower_id, state) values ($1,$2,$3)")
	if err != nil {
		log.Printf(fmt.Sprintf("PG Statements Wrong: %v", err))
		return nil, err
	}
	_, err = stmt.Exec(userID, followerID, state)
	if err != nil {
		log.Printf(fmt.Sprintf("PG Statements Exec Wrong: %v", err))
		return nil, err
	}

	if err != nil {
		log.Printf(fmt.Sprintf("PG Statements Exec Wrong: %v", err))
		return nil, err
	}

	return r.GetOneRelation(userID, followerID)
}

func (r *RelationshipModel) UpdateRelation(userID int, followerID int, state string) (*RelationshipModel, error) {
	query := "UPDATE relationships set state=$1 WHERE user_id=$2 AND follower_id=$3"
	// 判断反向关系
	inverse, _ := r.GetOneRelation(followerID, userID)
	// 互相喜欢 == "matched"
	//todo
	if inverse != nil && state == constant.Liked && (inverse.State == constant.Liked || inverse.State == constant.Matched) {
		// 事物: 均更新为 matched
		tx, _ := RunoobDB.Begin()
		state = constant.Matched
		//todo
		tx.Exec(query, state, userID, followerID)
		tx.Exec(query, state, followerID, userID)
		err := tx.Commit()
		if err != nil {
			log.Printf(fmt.Sprintf("PG Statements Exec Wrong: %v", err))
			return nil, err
		}
	} else if inverse != nil && state == constant.DisLiked && inverse.State == constant.Matched {
		// 事物: inverse.State 更新为 liked
		tx, _ := RunoobDB.Begin()
		tx.Exec(query, state, userID, followerID)
		tx.Exec(query, constant.Liked, followerID, userID)
		err := tx.Commit()
		if err != nil {
			log.Printf(fmt.Sprintf("PG Statements Exec Wrong: %v", err))
			return nil, err
		}
	} else {
		// 其余只改一条
		stmt, err := RunoobDB.Prepare(query)
		if err != nil {
			return nil, err
		}

		stmt.Exec(state, userID, followerID)
	}

	return r.GetOneRelation(userID, followerID)
}
