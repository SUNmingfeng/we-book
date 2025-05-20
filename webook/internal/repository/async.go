package repository

import (
	"context"
	"time"
)

type SmsRepository interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type AsyncSmsRepository struct {
}

func (a AsyncSmsRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	//TODO implement me
	panic("implement me")
}
