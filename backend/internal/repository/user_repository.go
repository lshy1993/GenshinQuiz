package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"golang.org/x/crypto/bcrypt"

	"genshin-quiz/internal/models"
	"genshin-quiz/internal/table"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetAll(limit, offset int) ([]models.User, int, error) {
	start := time.Now()
	
	// 查询总数
	var total int64
	countStmt := postgres.SELECT(postgres.COUNT(postgres.STAR)).FROM(table.Users)
	
	err := countStmt.Query(r.db, &total)
	if err != nil {
		r.logger.Printf("Failed to get user count (duration: %v): %v", time.Since(start), err)
		return nil, 0, fmt.Errorf("failed to get user count: %w", err)
	}

	// 分页查询用户
	var users []models.User
	stmt := postgres.SELECT(
		table.Users.ID,
		table.Users.Username,
		table.Users.Email,
		table.Users.DisplayName,
		table.Users.AvatarURL,
		table.Users.TotalScore,
		table.Users.QuizzesCompleted,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	).FROM(
		table.Users,
	).ORDER_BY(
		table.Users.CreatedAt.DESC(),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	err = stmt.Query(r.db, &users)
	if err != nil {
		r.logger.Printf("Failed to query users (limit: %d, offset: %d, duration: %v): %v", limit, offset, time.Since(start), err)
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}

	r.logger.Printf("users retrieved successfully (count: %d, total: %d, duration: %v)", len(users), total, time.Since(start))	return users, int(total), nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	start := time.Now()
	
	var user models.User
	stmt := postgres.SELECT(
		table.Users.ID,
		table.Users.Username,
		table.Users.Email,
		table.Users.DisplayName,
		table.Users.AvatarURL,
		table.Users.TotalScore,
		table.Users.QuizzesCompleted,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	).FROM(
		table.Users,
	).WHERE(
		table.Users.ID.EQ(postgres.Int64(id)),
	)

	err := stmt.Query(r.db, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		r.logger.Printf("Failed to get user by ID (user_id: %d, duration: %v): %v", id, time.Since(start), err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	r.logger.Printf("Retrieved user by ID (user_id: %d, duration: %v)", id, time.Since(start))

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	start := time.Now()
	
	var user models.User
	stmt := postgres.SELECT(
		table.Users.ID,
		table.Users.Username,
		table.Users.Email,
		table.Users.PasswordHash,
		table.Users.DisplayName,
		table.Users.AvatarURL,
		table.Users.TotalScore,
		table.Users.QuizzesCompleted,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	).FROM(
		table.Users,
	).WHERE(
		table.Users.Email.EQ(postgres.String(email)),
	)

	err := stmt.Query(r.db, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by email", 
			zap.String("email", email),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	r.logger.Printf("user retrieved by email successfully (email: %s, user_id: %d, duration: %v)", email, user.ID, time.Since(start))	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	start := time.Now()
	
	var user models.User
	stmt := postgres.SELECT(
		table.Users.ID,
		table.Users.Username,
		table.Users.Email,
		table.Users.DisplayName,
		table.Users.AvatarURL,
		table.Users.TotalScore,
		table.Users.QuizzesCompleted,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	).FROM(
		table.Users,
	).WHERE(
		table.Users.Username.EQ(postgres.String(username)),
	)

	err := stmt.Query(r.db, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by username", 
			zap.String("username", username),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	r.logger.Debug("Retrieved user by username",
		zap.String("username", username),
		zap.Int64("user_id", user.ID),
		zap.Duration("duration", time.Since(start)),
	)

	return &user, nil
}

func (r *UserRepository) Create(req models.CreateUserRequest) (*models.User, error) {
	start := time.Now()
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user
	var user models.User
	stmt := table.Users.INSERT(
		table.Users.Username,
		table.Users.Email,
		table.Users.PasswordHash,
		table.Users.DisplayName,
	).VALUES(
		req.Username,
		req.Email,
		string(hashedPassword),
		req.Username, // Default display name to username
	).RETURNING(
		table.Users.ID,
		table.Users.Username,
		table.Users.Email,
		table.Users.DisplayName,
		table.Users.AvatarURL,
		table.Users.TotalScore,
		table.Users.QuizzesCompleted,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	)

	err = stmt.Query(r.db, &user)
	if err != nil {
		r.logger.Error("Failed to create user", 
			zap.String("username", req.Username),
			zap.String("email", req.Email),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("Created new user",
		zap.Int64("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("email", user.Email),
		zap.Duration("duration", time.Since(start)),
	)

	return &user, nil
}

func (r *UserRepository) Update(id int64, req models.UpdateUserRequest) (*models.User, error) {
	start := time.Now()
	
	updateStmt := table.Users.UPDATE()

	// Build dynamic update statement
	if req.Username != nil {
		updateStmt = updateStmt.SET(table.Users.Username.SET(postgres.String(*req.Username)))
	}
	if req.Email != nil {
		updateStmt = updateStmt.SET(table.Users.Email.SET(postgres.String(*req.Email)))
	}
	if req.DisplayName != nil {
		updateStmt = updateStmt.SET(table.Users.DisplayName.SET(postgres.String(*req.DisplayName)))
	}
	if req.AvatarURL != nil {
		updateStmt = updateStmt.SET(table.Users.AvatarURL.SET(postgres.String(*req.AvatarURL)))
	}

	// Always update the timestamp
	updateStmt = updateStmt.SET(table.Users.UpdatedAt.SET(postgres.Raw("CURRENT_TIMESTAMP")))

	// Add WHERE clause and RETURNING
	var user models.User
	stmt := updateStmt.WHERE(
		table.Users.ID.EQ(postgres.Int64(id)),
	).RETURNING(
		table.Users.ID,
		table.Users.Username,
		table.Users.Email,
		table.Users.DisplayName,
		table.Users.AvatarURL,
		table.Users.TotalScore,
		table.Users.QuizzesCompleted,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	)

	err := stmt.Query(r.db, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		r.logger.Error("Failed to update user", 
			zap.Int64("user_id", id),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	r.logger.Info("Updated user",
		zap.Int64("user_id", id),
		zap.Duration("duration", time.Since(start)),
	)

	return &user, nil
}

func (r *UserRepository) Delete(id int64) error {
	start := time.Now()
	
	stmt := table.Users.DELETE().WHERE(table.Users.ID.EQ(postgres.Int64(id)))

	result, err := stmt.Exec(r.db)
	if err != nil {
		r.logger.Error("Failed to delete user", 
			zap.Int64("user_id", id),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}

	r.logger.Info("Deleted user",
		zap.Int64("user_id", id),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

func (r *UserRepository) UpdateStats(id int64, totalScore, quizzesCompleted int) error {
	start := time.Now()
	
	stmt := table.Users.UPDATE(
		table.Users.TotalScore,
		table.Users.QuizzesCompleted,
		table.Users.UpdatedAt,
	).SET(
		totalScore,
		quizzesCompleted,
		postgres.Raw("CURRENT_TIMESTAMP"),
	).WHERE(
		table.Users.ID.EQ(postgres.Int64(id)),
	)

	result, err := stmt.Exec(r.db)
	if err != nil {
		r.logger.Error("Failed to update user stats", 
			zap.Int64("user_id", id),
			zap.Int("total_score", totalScore),
			zap.Int("quizzes_completed", quizzesCompleted),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to update user stats: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}

	r.logger.Debug("Updated user stats",
		zap.Int64("user_id", id),
		zap.Int("total_score", totalScore),
		zap.Int("quizzes_completed", quizzesCompleted),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

func (r *UserRepository) Search(query string, limit, offset int) ([]models.User, int, error) {
	start := time.Now()
	
	searchPattern := fmt.Sprintf("%%%s%%", query)
	
	// Count total matching users
	var total int64
	countStmt := postgres.SELECT(postgres.COUNT(postgres.STAR)).FROM(table.Users).WHERE(
		table.Users.Username.LIKE(postgres.String(searchPattern)).
			OR(table.Users.Email.LIKE(postgres.String(searchPattern))).
			OR(table.Users.DisplayName.LIKE(postgres.String(searchPattern))),
	)
	
	err := countStmt.Query(r.db, &total)
	if err != nil {
		r.logger.Error("Failed to count search results", 
			zap.String("query", query),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Search users
	var users []models.User
	stmt := postgres.SELECT(
		table.Users.ID,
		table.Users.Username,
		table.Users.Email,
		table.Users.DisplayName,
		table.Users.AvatarURL,
		table.Users.TotalScore,
		table.Users.QuizzesCompleted,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	).FROM(
		table.Users,
	).WHERE(
		table.Users.Username.LIKE(postgres.String(searchPattern)).
			OR(table.Users.Email.LIKE(postgres.String(searchPattern))).
			OR(table.Users.DisplayName.LIKE(postgres.String(searchPattern))),
	).ORDER_BY(
		table.Users.CreatedAt.DESC(),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	err = stmt.Query(r.db, &users)
	if err != nil {
		r.logger.Error("Failed to search users", 
			zap.String("query", query),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	r.logger.Debug("Search completed",
		zap.String("query", query),
		zap.Int("results", len(users)),
		zap.Int64("total", total),
		zap.Duration("duration", time.Since(start)),
	)

	return users, int(total), nil
}
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
