package database

import (
	"github.com/jackc/pgx/v5"
)

type Storager interface {
	Init() error
	Read() error
	Write() error
}

type Storage struct {
	db pgx.Conn
}

func (s *Storage) Init() error {
	return nil
}

func (s *Storage) Read() error {
	return nil
}

func (s *Storage) Write() error {
	return nil
}
