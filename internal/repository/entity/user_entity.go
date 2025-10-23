package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Username  string     `gorm:"comment:用户名;type:varchar(18)"`                // A regular string field
	Password  string     `gorm:"comment:用户邮箱;type:varchar(128)"`              // A pointer to a string, allowing for null values
	Email     string     `gorm:"comment:用户邮箱;type:varchar(32)"`               // A pointer to a string, allowing for null values
	Age       uint8      `gorm:"comment:用户年龄;size:8"`                         // An unsigned 8-bit integer
	Birthday  *time.Time `gorm:"comment:用户生日"`                                // A pointer to time.Time, can be null
	PostCount uint8      `gorm:"column:post_count;size:8;comment:用户文章数量统计字段"` // 用户文章数量统计字段
	Posts     []Post     `gorm:"comment:用户文章;foreignKey:UserID"`
	gorm.Model
}
