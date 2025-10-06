package azure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.uber.org/zap"

	"genshin-quiz-backend/internal/config"
)

// StorageClient wraps Azure Blob Storage operations
type StorageClient struct {
	client    *azblob.Client
	container string
	logger    *zap.Logger
}

// NewStorageClient creates a new Azure Storage client
func NewStorageClient(cfg *config.Config, logger *zap.Logger) (*StorageClient, error) {
	if cfg.Azure.StorageAccount == "" || cfg.Azure.StorageKey == "" {
		logger.Warn("Azure Storage not configured, file uploads will be disabled")
		return nil, nil
	}

	// Create credentials
	credential, err := azblob.NewSharedKeyCredential(cfg.Azure.StorageAccount, cfg.Azure.StorageKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credentials: %w", err)
	}

	// Create client
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", cfg.Azure.StorageAccount)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure client: %w", err)
	}

	return &StorageClient{
		client:    client,
		container: cfg.Azure.ContainerName,
		logger:    logger,
	}, nil
}

// UploadFile uploads a file to Azure Blob Storage
func (s *StorageClient) UploadFile(ctx context.Context, fileName string, data []byte, contentType string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("Azure Storage client not configured")
	}

	start := time.Now()

	// Generate unique filename
	uniqueFileName := s.generateUniqueFileName(fileName)

	// Upload file
	_, err := s.client.UploadBuffer(ctx, s.container, uniqueFileName, data, &azblob.UploadBufferOptions{
		HTTPHeaders: &azblob.BlobHTTPHeaders{
			BlobContentType: &contentType,
		},
		Metadata: map[string]*string{
			"uploaded_at": &[]string{time.Now().Format(time.RFC3339)}[0],
		},
	})

	if err != nil {
		s.logger.Error("Failed to upload file to Azure Storage",
			zap.String("filename", fileName),
			zap.String("container", s.container),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate URL
	fileURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", 
		s.getStorageAccount(), s.container, uniqueFileName)

	s.logger.Info("File uploaded successfully",
		zap.String("filename", fileName),
		zap.String("unique_filename", uniqueFileName),
		zap.String("url", fileURL),
		zap.Int("size", len(data)),
		zap.Duration("duration", time.Since(start)),
	)

	return fileURL, nil
}

// UploadStream uploads a stream to Azure Blob Storage
func (s *StorageClient) UploadStream(ctx context.Context, fileName string, reader io.Reader, contentType string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("Azure Storage client not configured")
	}

	start := time.Now()

	// Read all data from stream
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read stream: %w", err)
	}

	return s.UploadFile(ctx, fileName, data, contentType)
}

// DeleteFile deletes a file from Azure Blob Storage
func (s *StorageClient) DeleteFile(ctx context.Context, fileName string) error {
	if s.client == nil {
		return fmt.Errorf("Azure Storage client not configured")
	}

	start := time.Now()

	_, err := s.client.DeleteBlob(ctx, s.container, fileName, &azblob.DeleteBlobOptions{})
	if err != nil {
		s.logger.Error("Failed to delete file from Azure Storage",
			zap.String("filename", fileName),
			zap.String("container", s.container),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.Info("File deleted successfully",
		zap.String("filename", fileName),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

// GetFileURL generates a URL for accessing a file
func (s *StorageClient) GetFileURL(fileName string) string {
	if s.client == nil {
		return ""
	}

	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", 
		s.getStorageAccount(), s.container, fileName)
}

// GenerateSASURL generates a Shared Access Signature URL for a file
func (s *StorageClient) GenerateSASURL(ctx context.Context, fileName string, expiry time.Duration) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("Azure Storage client not configured")
	}

	// Generate SAS token
	expiryTime := time.Now().Add(expiry)
	
	// Create SAS URL
	sasURL, err := s.client.ServiceClient().GetSASURL(
		azblob.AccountSASResourceTypes{Container: true, Object: true},
		azblob.AccountSASPermissions{Read: true},
		azblob.AccountSASServices{Blob: true},
		expiryTime,
		&azblob.GetAccountSASURLOptions{})

	if err != nil {
		s.logger.Error("Failed to generate SAS URL",
			zap.String("filename", fileName),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to generate SAS URL: %w", err)
	}

	s.logger.Debug("Generated SAS URL",
		zap.String("filename", fileName),
		zap.Time("expiry", expiryTime),
	)

	return sasURL, nil
}

// generateUniqueFileName generates a unique filename with timestamp
func (s *StorageClient) generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := originalName[:len(originalName)-len(ext)]
	timestamp := time.Now().Format("20060102-150405")
	
	return fmt.Sprintf("%s-%s%s", name, timestamp, ext)
}

// getStorageAccount extracts storage account name from client
func (s *StorageClient) getStorageAccount() string {
	// This is a simplified version - in practice you'd extract from the client URL
	return "genshinquiz" // This should be extracted from config or client
}

// HealthCheck verifies Azure Storage connectivity
func (s *StorageClient) HealthCheck(ctx context.Context) error {
	if s.client == nil {
		return fmt.Errorf("Azure Storage client not configured")
	}

	// Try to list containers as a health check
	pager := s.client.NewListContainersPager(&azblob.ListContainersOptions{
		MaxResults: &[]int32{1}[0],
	})

	if pager.More() {
		_, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("Azure Storage health check failed: %w", err)
		}
	}

	return nil
}