package repository

import (
	"context"
	"dim_kurs/internal/custom_errors"
	"dim_kurs/internal/domain/model"
	"dim_kurs/internal/domain/request"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUser interface {
	GetUsers(ctx context.Context, req request.GetUsers) ([]model.User, error)
	Get(ctx context.Context, login string) (model.User, error)
	Create(ctx context.Context, user model.User) (uuid.UUID, error)
	Update(ctx context.Context, req model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

const (
	usersTable = "users"
	rolesTable = "roles"
)

type User struct {
	pool *pgxpool.Pool
}

func NewUser(pool *pgxpool.Pool) *User {
	return &User{
		pool: pool,
	}
}

func (r *User) GetUsers(ctx context.Context, req request.GetUsers) ([]model.User, error) {
	query := fmt.Sprintf(`
		SELECT u.id, u.login, u.first_name, u.last_name, u.pass_hash, u.token, r.role_name, u.phone_number, u.email
		FROM %s AS u
		LEFT JOIN %s AS r ON u.role_id = r.id
	`, usersTable, rolesTable)

	whereClauses := []string{}
	args := []interface{}{}
	argID := 1

	if req.Email != nil && *req.Email != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.email = $%d", argID))
		args = append(args, *req.Email)
		argID++
	}
	if req.FirstName != nil && *req.FirstName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.first_name = $%d", argID))
		args = append(args, *req.FirstName)
		argID++
	}
	if req.LastName != nil && *req.LastName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.last_name = $%d", argID))
		args = append(args, *req.LastName)
		argID++
	}
	if req.Role != nil && *req.Role != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("r.role_name = $%d", argID))
		args = append(args, req.Role)
		argID++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		var roleID *string // Для хранения role_id, если он NULL
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Token,
			&roleID,
			&user.Phone,
			&user.Email,
		); err != nil {
			return nil, err
		}

		// Присваиваем роль, если она существует
		if roleID != nil {
			user.Role = *roleID
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (r *User) Get(ctx context.Context, login string) (model.User, error) {
	query := fmt.Sprintf(`
		SELECT u.id, u.login, u.first_name, u.last_name, u.pass_hash, u.token, r.role_name, u.phone_number, u.email, u.status
		FROM %s AS u
		LEFT JOIN %s AS r ON u.role_id = r.id
		WHERE u.login = $1
	`, usersTable, rolesTable)

	row := r.pool.QueryRow(ctx, query, login)

	var user model.User
	var roleID *string // Для хранения role_id, если он NULL
	if err := row.Scan(
		&user.ID,
		&user.Login,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Token,
		&roleID,
		&user.Phone,
		&user.Email,
		&user.Status,
	); err != nil {
		if err == pgx.ErrNoRows {
			return model.User{}, custom_errors.UserNotExist
		}
		return model.User{}, err
	}

	if roleID != nil {
		user.Role = *roleID
	}

	return user, nil
}

func (r *User) Create(ctx context.Context, user model.User) (uuid.UUID, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	var roleID uuid.UUID
	err = tx.QueryRow(ctx, "SELECT id FROM roles WHERE role_name = $1", user.Role).Scan(&roleID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, custom_errors.RoleNotExist
		}
		return uuid.Nil, err
	}

	query := `INSERT INTO users (login, first_name, last_name, pass_hash, token, role_id, phone_number, email, status) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var userID uuid.UUID
	row := tx.QueryRow(ctx, query, user.Login, user.FirstName, user.LastName, user.Password, user.Token, roleID, user.Phone, user.Email, user.Status)
	err = row.Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return uuid.Nil, custom_errors.AlreadyExist
		}
		return uuid.Nil, err
	}

	return userID, nil
}

func (r *User) Update(ctx context.Context, user model.User) error {
	var roleID *uuid.UUID
	if user.Role != "" {
		err := r.pool.QueryRow(ctx, "SELECT id FROM roles WHERE role_name = $1", user.Role).Scan(&roleID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return custom_errors.RoleNotExist
			}
			return err
		}
	}

	query := `
		UPDATE users 
		SET login = COALESCE($1, login), 
		    pass_hash = COALESCE($2, pass_hash), 
		    role_id = COALESCE($3, role_id), 
		    phone_number = COALESCE($4, phone_number),
		    email = COALESCE($5, email),
		    status = COALESCE($6, status)
		WHERE id = $7`

	_, err := r.pool.Exec(ctx, query,
		&user.Login,
		&user.Password,
		roleID,
		&user.Phone,
		&user.Email,
		&user.Status,
		&user.ID,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return custom_errors.AlreadyExist
		}
	}

	return err
}

func (r *User) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users 
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return custom_errors.UserNotExist
		}
	}
	return err
}
