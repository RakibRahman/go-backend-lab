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

func (r *Repository) GetUsers(ctx context.Context, limit int, offset int, term string) (GetUsersResponse, error) {
	query := `SELECT id, name, email from users
				ORDER BY created_at DESC
				LIMIT $1 OFFSET $2
	`

	args := []any{limit, offset}

	if term != "" {
		query = `SELECT id, name, email from users
					WHERE name ILIKE $3 || '%' OR email ILIKE $3 || '%'
				ORDER BY created_at DESC
				LIMIT $1 OFFSET $2
	`
		args = append(args, term)
	}

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return GetUsersResponse{}, err
	}
	defer rows.Close()
	users := make([]User, 0)

	for rows.Next() {
		var user User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
		); err != nil {
			return GetUsersResponse{}, err
		}

		users = append(users, user)

	}

	if err := rows.Err(); err != nil {
		return GetUsersResponse{}, err
	}
	var total int64

	countQuery := `SELECT COUNT(*) FROM users`
	countArgs := []any{}

	if term != "" {
		countQuery = `SELECT COUNT(*) FROM users WHERE name ILIKE $1 || '%' OR email ILIKE $1 || '%'`
		countArgs = append(countArgs, term)
	}

	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return GetUsersResponse{}, err
	}

	hasMore := int64(offset+limit) < total

	return GetUsersResponse{
		Content:       users,
		TotalElements: total,
		HasMore:       hasMore,
	}, nil
}
