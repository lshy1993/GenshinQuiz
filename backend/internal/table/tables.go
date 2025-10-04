package table

import "github.com/go-jet/jet/v2/postgres"

// Users table and columns for PostgreSQL with go-jet
var (
	// Users table reference
	Users = postgres.NewTable("public", "users", "")

	// Users table columns
	UsersID               = postgres.StringColumn("id")
	UsersUsername         = postgres.StringColumn("username")
	UsersEmail            = postgres.StringColumn("email")
	UsersDisplayName      = postgres.StringColumn("display_name")
	UsersAvatarURL        = postgres.StringColumn("avatar_url")
	UsersTotalScore       = postgres.IntegerColumn("total_score")
	UsersQuizzesCompleted = postgres.IntegerColumn("quizzes_completed")
	UsersCreatedAt        = postgres.TimestampColumn("created_at")
	UsersUpdatedAt        = postgres.TimestampColumn("updated_at")
)

// Quizzes table and columns for PostgreSQL with go-jet
var (
	// Quizzes table reference
	Quizzes = postgres.NewTable("public", "quizzes", "")

	// Quizzes table columns
	QuizzesID          = postgres.IntegerColumn("id")
	QuizzesTitle       = postgres.StringColumn("title")
	QuizzesDescription = postgres.StringColumn("description")
	QuizzesCategory    = postgres.StringColumn("category")
	QuizzesDifficulty  = postgres.StringColumn("difficulty")
	QuizzesTimeLimit   = postgres.IntegerColumn("time_limit")
	QuizzesCreatedBy   = postgres.IntegerColumn("created_by")
	QuizzesCreatedAt   = postgres.TimestampColumn("created_at")
	QuizzesUpdatedAt   = postgres.TimestampColumn("updated_at")
)

// Questions table and columns for PostgreSQL with go-jet
var (
	// Questions table reference
	Questions = postgres.NewTable("public", "questions", "")

	// Questions table columns
	QuestionsID            = postgres.IntegerColumn("id")
	QuestionsQuizID        = postgres.IntegerColumn("quiz_id")
	QuestionsQuestionText  = postgres.StringColumn("question_text")
	QuestionsQuestionType  = postgres.StringColumn("question_type")
	QuestionsOptions       = postgres.StringColumn("options")
	QuestionsCorrectAnswer = postgres.StringColumn("correct_answer")
	QuestionsExplanation   = postgres.StringColumn("explanation")
	QuestionsPoints        = postgres.IntegerColumn("points")
	QuestionsOrderIndex    = postgres.IntegerColumn("order_index")
	QuestionsCreatedAt     = postgres.TimestampColumn("created_at")
	QuestionsUpdatedAt     = postgres.TimestampColumn("updated_at")
)

// QuizAttempts table and columns for PostgreSQL with go-jet
var (
	// QuizAttempts table reference
	QuizAttempts = postgres.NewTable("public", "quiz_attempts", "")

	// QuizAttempts table columns
	QuizAttemptsID          = postgres.IntegerColumn("id")
	QuizAttemptsUserID      = postgres.IntegerColumn("user_id")
	QuizAttemptsQuizID      = postgres.IntegerColumn("quiz_id")
	QuizAttemptsScore       = postgres.IntegerColumn("score")
	QuizAttemptsMaxScore    = postgres.IntegerColumn("max_score")
	QuizAttemptsTimeTaken   = postgres.IntegerColumn("time_taken")
	QuizAttemptsCompletedAt = postgres.TimestampColumn("completed_at")
	QuizAttemptsCreatedAt   = postgres.TimestampColumn("created_at")
)
