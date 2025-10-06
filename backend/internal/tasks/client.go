package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// Client wraps asynq.Client for task queuing
type Client struct {
	client *asynq.Client
}

// NewClient creates a new task client
func NewClient(client *asynq.Client) *Client {
	return &Client{client: client}
}

// EnqueueEmailVerification queues an email verification task
func (c *Client) EnqueueEmailVerification(userID int64, email, token string) error {
	payload := EmailVerificationPayload{
		UserID: userID,
		Email:  email,
		Token:  token,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailVerification, data)

	_, err = c.client.Enqueue(task,
		asynq.Queue("default"),
		asynq.MaxRetry(3),
		asynq.Timeout(5*time.Minute),
	)

	return err
}

// EnqueueQuizSubmission queues a quiz submission processing task
func (c *Client) EnqueueQuizSubmission(userID, quizID int64, answers map[string]interface{}, submittedAt time.Time) error {
	payload := QuizSubmissionPayload{
		UserID:      userID,
		QuizID:      quizID,
		Answers:     answers,
		SubmittedAt: submittedAt.Format(time.RFC3339),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeQuizSubmission, data)

	_, err = c.client.Enqueue(task,
		asynq.Queue("default"),
		asynq.MaxRetry(5),
		asynq.Timeout(10*time.Minute),
	)

	return err
}

// EnqueueUserStatisticsUpdate queues a user statistics update task
func (c *Client) EnqueueUserStatisticsUpdate(userID int64, action string, data map[string]interface{}) error {
	payload := UserStatisticsUpdatePayload{
		UserID: userID,
		Action: action,
		Data:   data,
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeUserStatisticsUpdate, payloadData)

	_, err = c.client.Enqueue(task,
		asynq.Queue("low"),
		asynq.MaxRetry(3),
		asynq.Timeout(5*time.Minute),
	)

	return err
}

// EnqueueQuizAnalytics queues a quiz analytics processing task
func (c *Client) EnqueueQuizAnalytics(quizID int64, eventType string, data map[string]interface{}) error {
	payload := QuizAnalyticsPayload{
		QuizID:    quizID,
		EventType: eventType,
		Data:      data,
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeQuizAnalytics, payloadData)

	_, err = c.client.Enqueue(task,
		asynq.Queue("low"),
		asynq.MaxRetry(3),
		asynq.Timeout(5*time.Minute),
	)

	return err
}

// EnqueueImageUpload queues an image upload processing task
func (c *Client) EnqueueImageUpload(userID int64, imageURL, imageType string) error {
	payload := ImageUploadPayload{
		UserID:   userID,
		ImageURL: imageURL,
		Type:     imageType,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeImageUpload, data)

	_, err = c.client.Enqueue(task,
		asynq.Queue("default"),
		asynq.MaxRetry(3),
		asynq.Timeout(15*time.Minute),
	)

	return err
}

// EnqueueCriticalTask queues a task with high priority
func (c *Client) EnqueueCriticalTask(taskType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(taskType, data)

	_, err = c.client.Enqueue(task,
		asynq.Queue("critical"),
		asynq.MaxRetry(5),
		asynq.Timeout(30*time.Minute),
	)

	return err
}

// EnqueueDelayedTask queues a task to be processed after a delay
func (c *Client) EnqueueDelayedTask(taskType string, payload interface{}, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(taskType, data)

	_, err = c.client.Enqueue(task,
		asynq.Queue("default"),
		asynq.ProcessIn(delay),
		asynq.MaxRetry(3),
	)

	return err
}