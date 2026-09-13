package user

import (
	"fmt"

	"github.com/google/uuid"
)

var UserList = generateSeedUsers()

func generateSeedUsers() []User {
	users := []User{
		{ID: uuid.NewString(), Name: "Alice", Email: "alice@example.com"},
		{ID: uuid.NewString(), Name: "Bob", Email: "bob@example.com"},
		{ID: uuid.NewString(), Name: "Charlie", Email: "charlie@example.com"},
	}

	for i := len(users) + 1; i <= 50; i++ {
		users = append(users, User{
			ID:    uuid.NewString(),
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
		})
	}

	return users
}
