package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"

	"genshin-quiz-backend/internal/database"
	"genshin-quiz-backend/internal/models"
	"genshin-quiz-backend/internal/table"
)

type QuizRepository struct {
	db *database.DB
}

func NewQuizRepository(db *database.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) GetAll(limit, offset int, category, difficulty string) ([]models.Quiz, int, error) {
	// Build dynamic WHERE conditions using go-jet
	var conditions []postgres.BoolExpression

	if category != "" {
		conditions = append(conditions, table.QuizzesCategory.EQ(postgres.String(category)))
	}
	if difficulty != "" {
		conditions = append(conditions, table.QuizzesDifficulty.EQ(postgres.String(difficulty)))
	}

	// Get total count
	countStmt := postgres.SELECT(postgres.COUNT(postgres.STAR)).FROM(table.Quizzes)
	if len(conditions) > 0 {
		countStmt = countStmt.WHERE(postgres.AND(conditions...))
	}

	var total int
	err := r.db.QueryRowStatement(countStmt).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get quiz count: %w", err)
	}

	// Get quizzes with pagination
	stmt := postgres.SELECT(
		table.QuizzesID,
		table.QuizzesTitle,
		table.QuizzesDescription,
		table.QuizzesCategory,
		table.QuizzesDifficulty,
		table.QuizzesTimeLimit,
		table.QuizzesCreatedBy,
		table.QuizzesCreatedAt,
		table.QuizzesUpdatedAt,
	).FROM(
		table.Quizzes,
	).ORDER_BY(
		table.QuizzesCreatedAt.DESC(),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	if len(conditions) > 0 {
		stmt = stmt.WHERE(postgres.AND(conditions...))
	}

	rows, err := r.db.QueryStatement(stmt)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query quizzes: %w", err)
	}
	defer rows.Close()

	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		err := rows.Scan(
			&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Category,
			&quiz.Difficulty, &quiz.TimeLimit, &quiz.CreatedBy,
			&quiz.CreatedAt, &quiz.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan quiz: %w", err)
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, total, nil
}

func (r *QuizRepository) GetByID(id int64) (*models.Quiz, error) {
	stmt := postgres.SELECT(
		table.QuizzesID,
		table.QuizzesTitle,
		table.QuizzesDescription,
		table.QuizzesCategory,
		table.QuizzesDifficulty,
		table.QuizzesTimeLimit,
		table.QuizzesCreatedBy,
		table.QuizzesCreatedAt,
		table.QuizzesUpdatedAt,
	).FROM(
		table.Quizzes,
	).WHERE(
		table.QuizzesID.EQ(postgres.Int(id)),
	)

	var quiz models.Quiz
	err := r.db.QueryRowStatement(stmt).Scan(
		&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Category,
		&quiz.Difficulty, &quiz.TimeLimit, &quiz.CreatedBy,
		&quiz.CreatedAt, &quiz.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz by ID: %w", err)
	}

	// Get questions for this quiz
	questions, err := r.getQuestionsByQuizID(quiz.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions for quiz: %w", err)
	}
	quiz.Questions = questions

	return &quiz, nil
}

func (r *QuizRepository) getQuestionsByQuizID(quizID int64) ([]models.Question, error) {
	stmt := postgres.SELECT(
		table.QuestionsID,
		table.QuestionsQuizID,
		table.QuestionsQuestionText,
		table.QuestionsQuestionType,
		table.QuestionsOptions,
		table.QuestionsCorrectAnswer,
		table.QuestionsExplanation,
		table.QuestionsPoints,
		table.QuestionsOrderIndex,
		table.QuestionsCreatedAt,
		table.QuestionsUpdatedAt,
	).FROM(
		table.Questions,
	).WHERE(
		table.QuestionsQuizID.EQ(postgres.Int(quizID)),
	).ORDER_BY(
		table.QuestionsOrderIndex.ASC(),
	)

	rows, err := r.db.QueryStatement(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var question models.Question
		var options pq.StringArray

		err := rows.Scan(
			&question.ID, &question.QuizID, &question.QuestionText,
			&question.QuestionType, &options, &question.CorrectAnswer,
			&question.Explanation, &question.Points, &question.OrderIndex,
			&question.CreatedAt, &question.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}

		question.Options = []string(options)
		questions = append(questions, question)
	}

	return questions, nil
}

func (r *QuizRepository) Create(req models.CreateQuizRequest) (*models.Quiz, error) {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the quiz
	quizQuery := `
		INSERT INTO quizzes (title, description, category, difficulty, time_limit, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, category, difficulty, time_limit, 
		          created_by, created_at, updated_at
	`

	var quiz models.Quiz
	err = tx.QueryRow(quizQuery, req.Title, req.Description, req.Category,
		req.Difficulty, req.TimeLimit, req.CreatedBy).Scan(
		&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Category,
		&quiz.Difficulty, &quiz.TimeLimit, &quiz.CreatedBy,
		&quiz.CreatedAt, &quiz.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create quiz: %w", err)
	}

	// Create questions
	for _, questionReq := range req.Questions {
		questionQuery := `
			INSERT INTO questions (quiz_id, question_text, question_type, options, 
			                     correct_answer, explanation, points, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, err = tx.Exec(questionQuery, quiz.ID, questionReq.QuestionText,
			questionReq.QuestionType, pq.Array(questionReq.Options),
			questionReq.CorrectAnswer, questionReq.Explanation,
			questionReq.Points, questionReq.OrderIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to create question: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Fetch the complete quiz with questions
	return r.GetByID(quiz.ID)
}

func (r *QuizRepository) Update(id int64, req models.UpdateQuizRequest) (*models.Quiz, error) {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Build dynamic update query for quiz
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argCount))
		args = append(args, *req.Title)
		argCount++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *req.Description)
		argCount++
	}
	if req.Category != nil {
		setParts = append(setParts, fmt.Sprintf("category = $%d", argCount))
		args = append(args, *req.Category)
		argCount++
	}
	if req.Difficulty != nil {
		setParts = append(setParts, fmt.Sprintf("difficulty = $%d", argCount))
		args = append(args, *req.Difficulty)
		argCount++
	}
	if req.TimeLimit != nil {
		setParts = append(setParts, fmt.Sprintf("time_limit = $%d", argCount))
		args = append(args, *req.TimeLimit)
		argCount++
	}

	if len(setParts) > 0 {
		setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
		args = append(args, id)

		query := fmt.Sprintf(`
			UPDATE quizzes 
			SET %s
			WHERE id = $%d
		`, strings.Join(setParts, ", "), argCount)

		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to update quiz: %w", err)
		}
	}

	// Update questions if provided
	if req.Questions != nil {
		// Delete existing questions
		_, err = tx.Exec("DELETE FROM questions WHERE quiz_id = $1", id)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing questions: %w", err)
		}

		// Create new questions
		for _, questionReq := range req.Questions {
			questionQuery := `
				INSERT INTO questions (quiz_id, question_text, question_type, options, 
				                     correct_answer, explanation, points, order_index)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			`

			_, err = tx.Exec(questionQuery, id, questionReq.QuestionText,
				questionReq.QuestionType, pq.Array(questionReq.Options),
				questionReq.CorrectAnswer, questionReq.Explanation,
				questionReq.Points, questionReq.OrderIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to create question: %w", err)
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Fetch the updated quiz
	return r.GetByID(id)
}

func (r *QuizRepository) Delete(id int64) error {
	stmt := table.Quizzes.DELETE().
		WHERE(table.QuizzesID.EQ(postgres.Int(id)))

	result, err := r.db.ExecStatement(stmt)
	if err != nil {
		return fmt.Errorf("failed to delete quiz: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quiz not found")
	}

	return nil
}
