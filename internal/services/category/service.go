package category

import (
	"context"

	"github.com/Secure-Website-Builder/Backend/internal/database"
	"github.com/Secure-Website-Builder/Backend/internal/models"
)

type Service struct {
	db *database.DB
}

func New(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ListCategoriesByStore(ctx context.Context, storeID int64) ([]models.ListCategoriesByStoreRow, error) {
	return s.db.Queries.ListCategoriesByStore(ctx, storeID)
}

func (s *Service) ListAttributesByCategory(ctx context.Context, categoryID int64) ([]models.ListCategoryAttributesRow, error) {
	return s.db.Queries.ListCategoryAttributes(ctx, categoryID)
}
