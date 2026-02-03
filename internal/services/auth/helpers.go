package auth

import (
	"context"
	"database/sql"

	"github.com/Secure-Website-Builder/Backend/internal/models"
	"github.com/google/uuid"
)

func (s *Service) mergeCustomerCartOnLogin(
	ctx context.Context,
	storeID int64,
	customerID int64,
	sessionID uuid.UUID,
) error {

	return s.db.RunInTx(ctx, func(q *models.Queries) error {

		var (
			sessionCart  *models.Cart
			customerCart *models.Cart
		)

		// Get guest/session cart
		cartBySession, err := q.GetCartBySessionForUpdate(ctx, models.GetCartBySessionForUpdateParams{
			StoreID:   storeID,
			SessionID: sessionID,
		})
		if err == nil {
			sessionCart = &cartBySession
		} else if err != sql.ErrNoRows {
			return err
		}

		// Get customer cart
		cartByCustomer, err := q.GetCartByCustomerForUpdate(ctx, models.GetCartByCustomerForUpdateParams{
			StoreID:   storeID,
			CustomerID: sql.NullInt64{Int64: customerID, Valid: true},
		})
		if err == nil {
			customerCart = &cartByCustomer
		} else if err != sql.ErrNoRows {
			return err
		}

		// No carts at all
		if sessionCart == nil && customerCart == nil {
			return nil
		}

		// Only session cart - attach to customer
		if sessionCart != nil && customerCart == nil {
			return q.AttachCartToCustomer(ctx, models.AttachCartToCustomerParams{
				CustomerID: sql.NullInt64{Int64: customerID, Valid: true},
				SessionID:  sessionID,
				CartID:     sessionCart.CartID,
			})
		}

		// Only customer cart - nothing to merge
		if sessionCart == nil && customerCart != nil {
			return nil
		}

		// Both carts exist - merge items
		err = q.MergeCartItems(ctx, models.MergeCartItemsParams{
			FromCartID: sessionCart.CartID,
			ToCartID:   customerCart.CartID,
		})
		if err != nil {
			return err
		}

		// Delete session cart
		if err := q.DeleteCart(ctx, sessionCart.CartID); err != nil {
			return err
		}

		return nil
	})
}
