package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"

	"genshin-quiz-backend/internal/database"
	"genshin-quiz-backend/internal/models"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(limit, offset int) ([]models.User, int, error) {
	// For now, we'll use raw SQL queries. 
	// After running generate_models.sh, we can replace these with Go-Jet queries

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user count: %w", err)
	}

	// Get users with pagination
	query := `
		SELECT id, username, email, display_name, avatar_url, 
		       total_score, quizzes_completed, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.DisplayName,
			&user.AvatarURL, &user.TotalScore, &user.QuizzesCompleted,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, username, email, display_name, avatar_url, 
		       total_score, quizzes_completed, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.AvatarURL, &user.TotalScore, &user.QuizzesCompleted,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, display_name, avatar_url, 
		       total_score, quizzes_completed, created_at, updated_at
		FROM users 
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.AvatarURL, &user.TotalScore, &user.QuizzesCompleted,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, display_name, avatar_url, 
		       total_score, quizzes_completed, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.AvatarURL, &user.TotalScore, &user.QuizzesCompleted,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(req models.CreateUserRequest) (*models.User, error) {
	query := `
		INSERT INTO users (username, email, display_name, avatar_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, display_name, avatar_url, 
		          total_score, quizzes_completed, created_at, updated_at
	`

	var user models.User
	err := r.db.QueryRow(query, req.Username, req.Email, req.DisplayName, req.AvatarURL).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.AvatarURL, &user.TotalScore, &user.QuizzesCompleted,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if pqErr.Constraint == "users_username_key" {
					return nil, fmt.Errorf("username already exists")
				}
				if pqErr.Constraint == "users_email_key" {
					return nil, fmt.Errorf("email already exists")
				}
			}
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(id int64, req models.UpdateUserRequest) (*models.User, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argCount))
		args = append(args, *req.Username)
		argCount++
	}
	if req.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argCount))
		args = append(args, *req.Email)
		argCount++
	}
	if req.DisplayName != nil {
		setParts = append(setParts, fmt.Sprintf("display_name = $%d", argCount))
		args = append(args, *req.DisplayName)
		argCount++
	}
	if req.AvatarURL != nil {
		setParts = append(setParts, fmt.Sprintf("avatar_url = $%d", argCount))
		args = append(args, *req.AvatarURL)
		argCount++
	}

	if len(setParts) == 0 {
		return r.GetByID(id)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE users 
		SET %s
		WHERE id = $%d
		RETURNING id, username, email, display_name, avatar_url, 
		          total_score, quizzes_completed, created_at, updated_at
	`, postgres.String(setParts).Join(", "), argCount)

	var user models.User
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.AvatarURL, &user.TotalScore, &user.QuizzesCompleted,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "users_username_key" {
				return nil, fmt.Errorf("username already exists")
			}
			if pqErr.Constraint == "users_email_key" {
				return nil, fmt.Errorf("email already exists")
			}
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) UpdateScore(userID int64, scoreIncrement int) error {
	query := `
		UPDATE users 
		SET total_score = total_score + $1,
		    quizzes_completed = quizzes_completed + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.Exec(query, scoreIncrement, userID)
	if err != nil {
		return fmt.Errorf("failed to update user score: %w", err)
	}

	return nil
}