package service

import "github.com/soulphazed/techno-db-forum/internal/model"

type Repository interface {
	GetStatus() (*model.Status, error)
	ClearAll() error
}
