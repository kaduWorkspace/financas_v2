package users

import "time"

type User struct {
    Id        int       `json:"id" form:"id"`
    Name      string    `json:"name" form:"name"`
    Email     string    `json:"email" form:"email"`
    Password  string    `json:"-" form:"password"` // - means omit in json
    CreatedAt time.Time `json:"created_at" form:"created_at"`
    UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

