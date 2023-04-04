package dialog

import (
	"context"
	"osoc-dialog/internal/entity"
)

//go:generate mockgen -source=interface.go -destination=../../mocks/message_store.go -package=mocks
type (
	// MessageStorage -.
	MessageStorage interface {
		Save(ctx context.Context, message entity.Message) error
		GetList(ctx context.Context, authorID int, userID int) ([]entity.Message, error)
	}
)
