package instrument

import (
	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	instrumentRepo InstrumentRepository
}

type ServiceDeps struct {
	Logger         *zap.Logger
	InstrumentRepo InstrumentRepository
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		logger:         deps.Logger,
		instrumentRepo: deps.InstrumentRepo,
	}
}
