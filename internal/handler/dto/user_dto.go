package dto

import (
	"time"
)

type UserDTO struct {
	ID        uint32    `json:"id"`
	Username  string    `json:"user_name"`  // A regular string field
	Password  string    `json:"password"`   // A pointer to a string, allowing for null values
	Email     string    `json:"email"`      // A pointer to a string, allowing for null values
	Age       uint8     `json:"age"`        // An unsigned 8-bit integer
	Birthday  time.Time `json:"birthday"`   // A pointer to time.Time, can be null
	PostCount uint8     `json:"post_count"` // 用户文章数量统计字段
}
