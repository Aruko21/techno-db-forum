package service

import "github.com/soulphazed/techno-db-forum/internal/model"

type Usecase interface {
	GetStatus() (*model.Status, error)
	ClearAll() error
}