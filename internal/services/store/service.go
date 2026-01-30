package store

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/Secure-Website-Builder/Backend/internal/database"
	"github.com/Secure-Website-Builder/Backend/internal/models"
	"github.com/Secure-Website-Builder/Backend/internal/storage"
	"github.com/Secure-Website-Builder/Backend/internal/utils"
)

type Service struct {
	db   *database.DB
	site *storage.MinIOStorage
}

func New(db *database.DB, site *storage.MinIOStorage) *Service {
	return &Service{
		db:   db,
		site: site,
	}
}

// IsOwner checks if user is owner of the store
func (s *Service) IsOwner(ctx context.Context, userID, storeID int64) (bool, error) {
	return s.db.Queries.IsStoreOwner(ctx, models.IsStoreOwnerParams{
		StoreOwnerID: userID,
		StoreID:      storeID,
	})
}


// CreateStore creates a store for a given owner or retries store initialization
// if a previous attempt failed.
//
// The function uses a database transaction to ensure that store creation
// and ownership constraints are enforced atomically.
// A store is created (or reused if previously failed) with download_status = 'pending'.
//
// Site configuration upload is executed AFTER the transaction commits.
// This follows a Saga-style approach:
//   - If upload succeeds, the store download_status is updated to 'completed'.
//   - If upload fails, the store download_status is updated to 'failed'.
//   - The same endpoint can be safely retried to complete initialization.
//
// External side effects (file upload) are intentionally excluded from the transaction
// to avoid long-running database locks.
func (s *Service) CreateStore(
	ctx context.Context,
	storeOwnerID int64,
	name string,
	domain string,
	currency string,
	timezone string,
	siteConfig json.RawMessage,
) (int64, error) {

	var (
		store models.Store
		err   error
	)

	// Normalize defaults
	if currency == "" {
		currency = "EGP"
	}
	if timezone == "" {
		timezone = "UTC"
	}

	// Transaction: create or reuse store
	err = s.db.RunInTx(ctx, func(qtx *models.Queries) error {

		// Check if a store already exists for this owner
		existingStore, err := qtx.GetStoreByOwnerID(ctx, storeOwnerID)
		if err == nil {
			switch existingStore.DownloadStatus {
			case "completed":
				return errors.New("store already exists for this owner")
			case "failed", "pending":
				store = existingStore
				return nil
			default:
				return errors.New("invalid store state")
			}
		}

		// Create new store with download_status = 'pending'
		store, err = qtx.CreateStore(ctx, models.CreateStoreParams{
			StoreOwnerID: storeOwnerID,
			Name:         name,
			Domain:       sql.NullString{String: domain, Valid: domain != ""},
			Currency:     sql.NullString{String: currency, Valid: true},
			Timezone:     sql.NullString{String: timezone, Valid: true},
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	// Upload site configuration AFTER transaction commit
	key := generateStoreUploadKey(store.StoreID)

	_, err = s.site.Upload(
		ctx,
		key,
		bytes.NewReader(siteConfig),
		int64(len(siteConfig)),
		"application/json",
	)

	if err != nil {
		// Best-effort status update
		_ = s.db.Queries.UpdateStoreDownloadStatus(ctx, models.UpdateStoreDownloadStatusParams{
			StoreID:        store.StoreID,
			DownloadStatus: "failed",
		})
		return 0, errors.New("failed to upload site config")
	}

	// Mark store initialization as completed
	_ = s.db.Queries.UpdateStoreDownloadStatus(ctx, models.UpdateStoreDownloadStatusParams{
		StoreID:        store.StoreID,
		DownloadStatus: "completed",
	})

	return store.StoreID, nil
}

func (s *Service) GetStore(
	ctx context.Context,
	storeID int64,
) (*models.StoreDTO, error) {

	store, err := s.db.Queries.GetStore(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return &models.StoreDTO{
		StoreID:      store.StoreID,
		StoreOwnerID: store.StoreOwnerID,
		Name:         store.Name,
		Domain:       utils.NullStringToPtr(store.Domain),
		Currency:     utils.NullStringToPtr(store.Currency),
		Timezone:     utils.NullStringToPtr(store.Timezone),
		CreatedAt:    store.CreatedAt,
		UpdatedAt:    store.UpdatedAt,
	}, nil
}
