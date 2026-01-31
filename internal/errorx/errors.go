package errorx

import "errors"

// Sentinel errors (domain-level)
var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidStoreID		= errors.New("invalid store id")
	ErrInvalidSession   = errors.New("invalid session")
	ErrInvalidSessionID = errors.New("invalid session id")
	ErrMissingSessionID = errors.New("missing session id")
	ErrInvalidVariant   = errors.New("invalid variant")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrCartNotFound     = errors.New("cart not found")
	ErrCartEmpty        = errors.New("cart empty")
	ErrOutOfStock       = errors.New("out of stock")
	ErrInvalidQuantity  = errors.New("invalid quantity")
)
