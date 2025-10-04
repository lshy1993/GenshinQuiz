package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"

	"genshin-quiz-backend/internal/database"
	"genshin-quiz-backend/internal/models"
	"genshin-quiz-backend/internal/table"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(limit, offset int) ([]models.User, int, error) {
	// 使用 go-jet 查询总数
	countStmt := postgres.SELECT(postgres.COUNT(postgres.STAR)).FROM(table.Users)

	var total int
	err := r.db.QueryRowStatement(countStmt).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user count: %w", err)
	}

	// 使用 go-jet 分页查询用户
	stmt := postgres.SELECT(
		table.UsersID,
		table.UsersUsername,
		table.UsersEmail,
		table.UsersDisplayName,
		table.UsersAvatarURL,
		table.UsersTotalScore,
		table.UsersQuizzesCompleted,
		table.UsersCreatedAt,
		table.UsersUpdatedAt,
	).FROM(
		table.Users,
	).ORDER_BY(
		table.UsersCreatedAt.DESC(),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	rows, err := r.db.QueryStatement(stmt)
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
	stmt := postgres.SELECT(
		table.UsersID,
		table.UsersUsername,
		table.UsersEmail,
		table.UsersDisplayName,
		table.UsersAvatarURL,
		table.UsersTotalScore,
		table.UsersQuizzesCompleted,
		table.UsersCreatedAt,
		table.UsersUpdatedAt,
	).FROM(
		table.Users,
	).WHERE(
		table.UsersID.EQ(postgres.String(fmt.Sprintf("%d", id))),
	)

	var user models.User
	err := r.db.QueryRowStatement(stmt).Scan(
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
	stmt := postgres.SELECT(
		table.UsersID,
		table.UsersUsername,
		table.UsersEmail,
		table.UsersDisplayName,
		table.UsersAvatarURL,
		table.UsersTotalScore,
		table.UsersQuizzesCompleted,
		table.UsersCreatedAt,
		table.UsersUpdatedAt,
	).FROM(
		table.Users,
	).WHERE(
		table.UsersUsername.EQ(postgres.String(username)),
	)

	var user models.User
	err := r.db.QueryRowStatement(stmt).Scan(
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
	stmt := postgres.SELECT(
		table.UsersID,
		table.UsersUsername,
		table.UsersEmail,
		table.UsersDisplayName,
		table.UsersAvatarURL,
		table.UsersTotalScore,
		table.UsersQuizzesCompleted,
		table.UsersCreatedAt,
		table.UsersUpdatedAt,
	).FROM(
		table.Users,
	).WHERE(
		table.UsersEmail.EQ(postgres.String(email)),
	)

	var user models.User
	err := r.db.QueryRowStatement(stmt).Scan(
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
	// 处理可选字段
	displayName := ""
	if req.DisplayName != nil {
		displayName = *req.DisplayName
	}
	avatarURL := ""
	if req.AvatarURL != nil {
		avatarURL = *req.AvatarURL
	}

	stmt := table.Users.INSERT(
		table.UsersUsername,
		table.UsersEmail,
		table.UsersDisplayName,
		table.UsersAvatarURL,
	).VALUES(
		req.Username,
		req.Email,
		displayName,
		avatarURL,
	).RETURNING(
		table.UsersID,
		table.UsersUsername,
		table.UsersEmail,
		table.UsersDisplayName,
		table.UsersAvatarURL,
		table.UsersTotalScore,
		table.UsersQuizzesCompleted,
		table.UsersCreatedAt,
		table.UsersUpdatedAt,
	)

	var user models.User
	err := r.db.QueryRowStatement(stmt).Scan(
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
	// Create base update statement
	updateStmt := table.Users.UPDATE()

	// Add fields to update based on what's provided
	if req.Username != nil {
		updateStmt = updateStmt.SET(table.UsersUsername.SET(postgres.String(*req.Username)))
	}
	if req.Email != nil {
		updateStmt = updateStmt.SET(table.UsersEmail.SET(postgres.String(*req.Email)))
	}
	if req.DisplayName != nil {
		updateStmt = updateStmt.SET(table.UsersDisplayName.SET(postgres.String(*req.DisplayName)))
	}
	if req.AvatarURL != nil {
		updateStmt = updateStmt.SET(table.UsersAvatarURL.SET(postgres.String(*req.AvatarURL)))
	}

	// If no fields to update, just return current user
	if req.Username == nil && req.Email == nil && req.DisplayName == nil && req.AvatarURL == nil {
		return r.GetByID(id)
	}

	// Always update the updated_at timestamp using SQL expression
	updateStmt = updateStmt.SET(table.UsersUpdatedAt.SET(postgres.RawTimestamp("CURRENT_TIMESTAMP")))

	// Add WHERE clause and RETURNING
	stmt := updateStmt.WHERE(
		table.UsersID.EQ(postgres.String(fmt.Sprintf("%d", id))),
	).RETURNING(
		table.UsersID,
		table.UsersUsername,
		table.UsersEmail,
		table.UsersDisplayName,
		table.UsersAvatarURL,
		table.UsersTotalScore,
		table.UsersQuizzesCompleted,
		table.UsersCreatedAt,
		table.UsersUpdatedAt,
	)

	var user models.User
	err := r.db.QueryRowStatement(stmt).Scan(
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
	stmt := table.Users.DELETE().
		WHERE(table.UsersID.EQ(postgres.String(fmt.Sprintf("%d", id))))

	result, err := r.db.ExecStatement(stmt)
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
	stmt := table.Users.UPDATE().
		SET(
			table.UsersTotalScore.SET(table.UsersTotalScore.ADD(postgres.Int(int64(scoreIncrement)))),
			table.UsersQuizzesCompleted.SET(table.UsersQuizzesCompleted.ADD(postgres.Int(1))),
		).WHERE(
		table.UsersID.EQ(postgres.String(fmt.Sprintf("%d", userID))),
	)

	_, err := r.db.ExecStatement(stmt)
	if err != nil {
		return fmt.Errorf("failed to update user score: %w", err)
	}

	return nil
}
