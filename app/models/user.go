package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
    Id        int       `json:"id" form:"id"`
    Username      string    `json:"username" form:"username"`
    Email     string    `json:"email" form:"email"`
    Password  string    `json:"-" form:"password"` // - means omit in json
    CreatedAt time.Time `json:"created_at" form:"created_at"`
    UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
   	orm.SoftDeletes
}
func (r User) TableName() string {
    return "users"
}
