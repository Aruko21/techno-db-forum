package serviceUsecase

import (
	"github.com/pkg/errors"
	"github.com/soulphazed/techno-db-forum/internal/app/service"
	"github.com/soulphazed/techno-db-forum/internal/model"
)

type ServiceUsecase struct {
	rep service.Repository
}

func NewServiceUsecase(rep service.Repository) service.Usecase {
	return &ServiceUsecase{
		rep: rep,
	}
}


func (s ServiceUsecase) ClearAll() error {
	err := s.rep.ClearAll()

	if err != nil {
		return errors.Wrap(err, "ClearAll()")
	}

	return nil
}

func (s ServiceUsecase) GetStatus() (*model.Status, error) {
	status, err := s.rep.GetStatus()

	if err != nil {
		return nil, errors.Wrap(err, "GetStatus()")
	}

	return status, nil
}
