package v1

import (
	"context"
	"osoc-dialog/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=../../../mocks/services.go -package=mocks
type (
	// DialogProvider -.
	DialogProvider interface {
		SaveMessage(ctx context.Context, userID int, authorID int, text string) error
		Messages(ctx context.Context, authorID int, userID int) ([]entity.Message, error)
	}
)
