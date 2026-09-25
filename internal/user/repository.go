package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool // means the repository owns access to the database pool, not the database itself.
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, payload CreateUserRequest) (User, error) {
	user := User{
		ID:    uuid.NewString(),
		Name:  payload.Name,
		Email: payload.Email,
	}

	query := `
	INSERT INTO users(id,name,email)
	values ($1,$2,$3)
	RETURNING id,name,email
	`

	err := r.db.QueryRow(ctx, query, user.ID, user.Name, user.Email).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return User{}, err
	}
	return user, nil
}
