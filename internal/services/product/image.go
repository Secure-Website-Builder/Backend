package product

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"

	"github.com/Secure-Website-Builder/Backend/internal/models"
	"github.com/google/uuid"
)

func (s *Service) UploadVariantImage(
    ctx context.Context,
    storeID, productID, variantID int64,
    file multipart.File,
    fileHeader *multipart.FileHeader,
    isPrimary bool,
) (string, error) {

    // Verify ownership
    variant, err := s.q.GetVariant(ctx, variantID)
    if err != nil {
        return "", fmt.Errorf("variant not found")
    }
    if variant.StoreID != storeID || variant.ProductID != productID {
        return "", fmt.Errorf("not your variant")
    }

    // Generate S3 key
    key := fmt.Sprintf(
        "stores/%d/products/%d/variants/%d/%s",
        storeID, productID, variantID, uuid.NewString(),
    )

    // Upload to S3 / MinIO
    url, err := s.storage.Upload(
        ctx,
        key,
        file,
        fileHeader.Size,
        fileHeader.Header.Get("Content-Type"),
    )

    if err != nil {
        return "", err
    }

    // Start DB transaction
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        s.storage.Delete(ctx, key)
        return "", err
    }
    defer tx.Rollback()

    qtx := s.q.WithTx(tx)

    if isPrimary {
        err = qtx.SetPrimaryVariantImage(ctx, models.SetPrimaryVariantImageParams{
            VariantID: variantID,
            PrimaryImageUrl: sql.NullString{
                String: url,
                Valid:  true,
            },
        })
    } else {
        _, err = qtx.InsertVariantImage(ctx, models.InsertVariantImageParams{
            ProductVariantID: variantID,
            ImageUrl:         url,
        })
    }

    if err != nil {
        s.storage.Delete(ctx, key)
        return "", err
    }

    if err := tx.Commit(); err != nil {
        s.storage.Delete(ctx, key)
        return "", err
    }

    return url, nil
}
