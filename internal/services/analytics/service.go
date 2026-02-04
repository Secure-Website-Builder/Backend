package analytics

import (

	"github.com/Secure-Website-Builder/Backend/internal/database"
)

type Service struct {
	db *database.DB
}

func New(db *database.DB) *Service {
	return &Service{db: db}
}


