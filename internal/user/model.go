package user

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type GetUsersResponse struct {
	Content       []User `json:"content"`
	TotalElements int64  `json:"totalElements"`
	HasMore       bool   `json:"hasMore"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
}
