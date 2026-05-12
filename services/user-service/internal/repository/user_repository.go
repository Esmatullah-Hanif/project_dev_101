package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/microservices/user-service/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	query := `
		SELECT id, email, first_name, last_name, bio, avatar_url, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	row := r.db.QueryRow(ctx, query, userID)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, limit int, offset int) ([]model.User, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, email, first_name, last_name, bio, avatar_url, is_active, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, rows.Err()
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error) {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, bio = $3, avatar_url = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, email, first_name, last_name, bio, avatar_url, is_active, created_at, updated_at
	`
	row := r.db.QueryRow(ctx, query, req.FirstName, req.LastName, req.Bio, req.AvatarURL, time.Now(), userID)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET deleted_at = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(ctx, query, time.Now(), time.Now(), userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
