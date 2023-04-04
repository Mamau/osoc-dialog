package dialog

import (
	"context"
	"osoc-dialog/internal/entity"
	"osoc-dialog/internal/errors"
	"osoc-dialog/pkg/log"
)

type Service struct {
	logger     log.Logger
	repository MessageStorage
}

func NewService(l log.Logger, r MessageStorage) *Service {
	return &Service{
		logger:     l,
		repository: r,
	}
}

func (s *Service) Messages(ctx context.Context, authorID int, userID int) ([]entity.Message, error) {
	data, err := s.repository.GetList(ctx, authorID, userID)
	if err != nil {
		s.logger.Err(err).Msg("error while get messages")
		return nil, errors.SomethingWentWrong
	}
	return data, nil
}

func (s *Service) SaveMessage(ctx context.Context, userID int, authorID int, text string) error {
	message := entity.Message{
		UserID:   userID,
		AuthorID: authorID,
		Text:     text,
	}
	if err := s.repository.Save(ctx, message); err != nil {
		s.logger.Err(err).Msg("error while save dialog message")
		return errors.SomethingWentWrong
	}
	return nil
}
