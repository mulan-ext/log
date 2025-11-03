package log

import (
	"context"
)

type Service struct {
	cfg    *Config
	logger *Logger
}

func (s *Service) Serve() (err error) { s.logger, err = NewWithConfig(s.cfg); return }

func (s *Service) Shutdown(ctx context.Context) error { return s.logger.Close() }

func (s *Service) Close() error { return s.logger.Close() }

func NewService(cfg *Config) *Service { return &Service{cfg: cfg} }
