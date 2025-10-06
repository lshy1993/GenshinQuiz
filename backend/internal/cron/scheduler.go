package cron

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"genshin-quiz/internal/services"
	"genshin-quiz/internal/tasks"
)

// Scheduler handles scheduled tasks
type Scheduler struct {
	userService *services.UserService
	quizService *services.QuizService
	taskClient  *asynq.Client
	logger      *log.Logger
	scheduler   *gocron.Scheduler
}

// NewScheduler creates a new cron scheduler
func NewScheduler(userService *services.UserService, quizService *services.QuizService, taskClient *asynq.Client, logger *log.Logger) *Scheduler {
	return &Scheduler{
		userService: userService,
		quizService: quizService,
		taskClient:  taskClient,
		logger:      logger,
		stopChan:    make(chan struct{}),
	}
}

// Start begins the scheduled task execution
func (s *Scheduler) Start() {
	// Start various scheduled tasks
	go s.runDailyUserStatisticsUpdate()
	go s.runWeeklyQuizAnalytics()
	go s.runDailyDataCleanup()
	go s.runHourlyHealthCheck()
}

// Stop stops all scheduled tasks
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// runDailyUserStatisticsUpdate runs daily at 2:00 AM
func (s *Scheduler) runDailyUserStatisticsUpdate() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Calculate time until next 2:00 AM
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
	time.Sleep(time.Until(next))

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Running daily user statistics update")

			if err := s.performDailyUserStatisticsUpdate(); err != nil {
				s.logger.Error("Failed to perform daily user statistics update", zap.Error(err))
			}

		case <-s.stopChan:
			s.logger.Info("Stopping daily user statistics update scheduler")
			return
		}
	}
}

// runWeeklyQuizAnalytics runs weekly on Sunday at 3:00 AM
func (s *Scheduler) runWeeklyQuizAnalytics() {
	ticker := time.NewTicker(7 * 24 * time.Hour)
	defer ticker.Stop()

	// Calculate time until next Sunday 3:00 AM
	now := time.Now()
	daysUntilSunday := (7 - int(now.Weekday())) % 7
	if daysUntilSunday == 0 && now.Hour() >= 3 {
		daysUntilSunday = 7
	}
	next := time.Date(now.Year(), now.Month(), now.Day()+daysUntilSunday, 3, 0, 0, 0, now.Location())
	time.Sleep(time.Until(next))

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Running weekly quiz analytics")

			if err := s.performWeeklyQuizAnalytics(); err != nil {
				s.logger.Error("Failed to perform weekly quiz analytics", zap.Error(err))
			}

		case <-s.stopChan:
			s.logger.Info("Stopping weekly quiz analytics scheduler")
			return
		}
	}
}

// runDailyDataCleanup runs daily at 1:00 AM
func (s *Scheduler) runDailyDataCleanup() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Calculate time until next 1:00 AM
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 1, 0, 0, 0, now.Location())
	if now.Hour() < 1 {
		next = time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, now.Location())
	}
	time.Sleep(time.Until(next))

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Running daily data cleanup")

			if err := s.performDailyDataCleanup(); err != nil {
				s.logger.Error("Failed to perform daily data cleanup", zap.Error(err))
			}

		case <-s.stopChan:
			s.logger.Info("Stopping daily data cleanup scheduler")
			return
		}
	}
}

// runHourlyHealthCheck runs every hour
func (s *Scheduler) runHourlyHealthCheck() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Running hourly health check")

			if err := s.performHealthCheck(); err != nil {
				s.logger.Error("Health check failed", zap.Error(err))
			}

		case <-s.stopChan:
			s.logger.Info("Stopping hourly health check scheduler")
			return
		}
	}
}

// performDailyUserStatisticsUpdate performs daily user statistics update
func (s *Scheduler) performDailyUserStatisticsUpdate() error {
	// TODO: Implement logic to:
	// 1. Calculate daily active users
	// 2. Update user engagement metrics
	// 3. Generate daily reports
	// 4. Queue individual user statistics update tasks

	taskClient := tasks.NewClient(s.taskClient)

	// Example: Queue user statistics update for all active users
	data := map[string]interface{}{
		"type": "daily_update",
		"date": time.Now().Format("2006-01-02"),
	}

	return taskClient.EnqueueUserStatisticsUpdate(0, "daily_batch_update", data)
}

// performWeeklyQuizAnalytics performs weekly quiz analytics
func (s *Scheduler) performWeeklyQuizAnalytics() error {
	// TODO: Implement logic to:
	// 1. Analyze quiz performance trends
	// 2. Generate weekly reports
	// 3. Update quiz rankings
	// 4. Calculate popular quizzes

	taskClient := tasks.NewClient(s.taskClient)

	data := map[string]interface{}{
		"type":       "weekly_analytics",
		"week_start": time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
		"week_end":   time.Now().Format("2006-01-02"),
	}

	return taskClient.EnqueueQuizAnalytics(0, "weekly_batch_analytics", data)
}

// performDailyDataCleanup performs daily data cleanup
func (s *Scheduler) performDailyDataCleanup() error {
	// TODO: Implement logic to:
	// 1. Clean up expired sessions
	// 2. Remove old temporary files
	// 3. Archive old quiz submissions
	// 4. Clean up unused images

	s.logger.Info("Performing data cleanup",
		zap.String("cleanup_date", time.Now().Format("2006-01-02")),
	)

	// Example cleanup tasks
	cutoffDate := time.Now().AddDate(0, 0, -30) // 30 days ago

	// Clean up expired sessions (implement in user service)
	// Clean up old temporary data
	// Archive old submissions

	s.logger.Info("Data cleanup completed")
	return nil
}

// performHealthCheck performs system health check
func (s *Scheduler) performHealthCheck() error {
	// TODO: Implement logic to:
	// 1. Check database connectivity
	// 2. Check Redis connectivity
	// 3. Check external service availability
	// 4. Monitor system resources
	// 5. Send alerts if issues detected

	s.logger.Info("System health check completed")
	return nil
}
