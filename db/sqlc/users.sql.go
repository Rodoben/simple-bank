// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"
)

const createuser = `-- name: Createuser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email, 
  contact
) VALUES (
  $1, $2, $3,$4,$5
) RETURNING username, hashed_password, full_name, email, password_changed_at, created_at, contact
`

type CreateuserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Contact        string `json:"contact"`
}

func (q *Queries) Createuser(ctx context.Context, arg CreateuserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createuser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
		arg.Contact,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.Contact,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, hashed_password, full_name, email, password_changed_at, created_at, contact FROM users WHERE username= $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.Contact,
	)
	return i, err
}
