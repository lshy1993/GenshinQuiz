package services

import (
	"context"

	"genshin-quiz/config"
	"genshin-quiz/generated/oapi"
)

func GetQuizzes(
	ctx context.Context,
	app *config.App,
	req oapi.GetQuizzesRequestObject,
) error {
	limit := *req.Params.Limit
	offset := *req.Params.Offset
	// category := *req.Params.Category
	// difficulty := *req.Params.Difficulty
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return nil
}
