package models

import "time"

type User struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	Account    string    `json:"account"`
	Password   string    `json:"password"`
	Symbol     string    `json:"symbol"`
	Data       string    `json:"data"`
	CreateTime time.Time `json:"create_time"`
}
