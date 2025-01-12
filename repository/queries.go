package repository

import (
	"database/sql"
	"errors"
	"strings"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) FindUserById(id int) (*User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	return r.find(row)
}

func (r *Repository) FindUserByPseudo(pseudo string) (*User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE pseudo = ?", pseudo)

	return r.find(row)
}

func (r *Repository) find(row *sql.Row) (*User, error) {
	user := new(User)
	err := row.Scan(&user.Id, &user.Pseudo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

type AddOperationParams struct {
	Inputs []float64
	Type   OperationType
	Result float64
	UserId int
}

func (r *Repository) AddOperation(param AddOperationParams) error {
	args := make([]interface{}, 0, len(param.Inputs)+3)
	for _, input := range param.Inputs {
		args = append(args, input)
	}
	args = append(args, param.Type, param.Result, param.UserId)

	_, err := r.db.Exec(
		"INSERT INTO operations (inputs, type, result, user_id) VALUES (JSON_ARRAY("+strings.Repeat("?,", len(param.Inputs))[:len(param.Inputs)*2-1]+"), ?, ?, ?)",
		args...,
	)
	if err != nil {
		return err
	}

	return nil
}
