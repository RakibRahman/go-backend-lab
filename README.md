# Go Backend Lab

A hands-on backend learning project using Go and `net/http`.

This repo is used to practice and gradually build:
- HTTP APIs
- JSON handling
- validation
- authentication and authorization
- PostgreSQL
- caching
- background jobs
- SSE
- logging and observability

## Project structure

```
.
├── cmd/
│   └── api/
│       └── main.go        # entry point, route registration
├── internal/
│   └── user/
│       ├── model.go       # User, CreateUserRequest types
│       ├── data.go        # in-memory seed data (UserList)
│       └── handler.go     # GetUserHandler, CreateUserHandler
├── go.mod
└── README.md
```