package seeders

import (
	"fmt"
	"goravel/app/models"
	"time"

	"github.com/goravel/framework/facades"
)

type UserSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *UserSeeder) Signature() string {
	return "UserSeeder"
}

// Run executes the seeder logic.
func (s *UserSeeder) Run() error {
    senha, err := facades.Hash().Make("password")
    if err != nil {
        fmt.Println(err)
        return err
    }
	users := []models.User{
		{
			Username:      "Admin User",
			Email:     "admin@example.com",
			Password: senha,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:      "Regular User",
			Email:     "user@example.com",
			Password:  senha,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		if err := facades.Orm().Query().Create(&user); err != nil {
			return err
		}
	}

	return nil
}
