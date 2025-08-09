package service

import (
	"encoding/json"
	"fmt"

	"github.com/gajare/redius-go/db"
	"github.com/gajare/redius-go/models"
)

func CreateTeacher(t models.Teacher) error {
	_, err := db.DB.Exec("INSERT INTO teachers(name, email) VALUES($1, $2)", t.Name, t.Email)
	return err
}

func GetTeacher(id int) (models.Teacher, error) {
	var t models.Teacher

	val, err := db.RedisClient.Get(db.Ctx, fmt.Sprintf("teacher:%d", id)).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &t)
		return t, nil
	}

	row := db.DB.QueryRow("SELECT id, name, email FROM teachers WHERE id=$1", id)
	err = row.Scan(&t.ID, &t.Name, &t.Email)
	if err != nil {
		return t, err
	}

	cache, _ := json.Marshal(t)
	db.RedisClient.Set(db.Ctx, fmt.Sprintf("teacher:%d", id), cache, 0)

	return t, nil
}

func UpdateTeacher(t models.Teacher) error {
	_, err := db.DB.Exec("UPDATE teachers SET name=$1, email=$2 WHERE id=$3", t.Name, t.Email, t.ID)
	db.RedisClient.Del(db.Ctx, fmt.Sprintf("teacher:%d", t.ID))
	return err
}

func DeleteTeacher(id int) error {
	_, err := db.DB.Exec("DELETE FROM teachers WHERE id=$1", id)
	db.RedisClient.Del(db.Ctx, fmt.Sprintf("teacher:%d", id))
	return err
}
