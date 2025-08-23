package usecase

import (
	"chat_service/internal/domain"
	"errors"
	"time"
)

type UseCase struct {
	repo domain.MessageRepository
}

func NewUseCase(repo domain.MessageRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) SendMessage(user domain.User, content string) (domain.Message, error) {
	if content == "" {
		return domain.Message{}, errors.New("message content cannot be empty")
	}

	msg := domain.Message{
		UserID:     user.ID,
		Content:    content,
		Created_at: time.Now(),
	}

	saveMsg, err := uc.repo.SendMessage(msg)
	if err != nil {
		return domain.Message{}, err
	}

	return saveMsg, nil
}

func (uc *UseCase) DeleteMessage(id int64) error {
	return uc.repo.DeleteMessage(id)
}

func (uc *UseCase) GetMessages(limit, offset int) ([]domain.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	return uc.repo.GetMessages(limit, offset)
}
