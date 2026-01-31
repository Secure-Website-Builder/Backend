package errorx

import (
	"database/sql"
	"errors"
	"net/http"
)

type HTTPError struct {
	Status  int
	Message string
}

func Resolve(err error) HTTPError {
	switch {
	case errors.Is(err, ErrInvalidSession):
		return HTTPError{http.StatusUnauthorized, MsgInvalidSession}

	case errors.Is(err, ErrInvalidVariant):
		return HTTPError{http.StatusNotFound, MsgInvalidVariant}

	case errors.Is(err, ErrInsufficientStock),
		errors.Is(err, ErrOutOfStock):
		return HTTPError{http.StatusConflict, MsgInsufficientStock}

	case errors.Is(err, ErrCartNotFound):
		return HTTPError{http.StatusNotFound, MsgCartNotFound}

	case errors.Is(err, ErrCartEmpty):
		return HTTPError{http.StatusBadRequest, MsgCartEmpty}

	case errors.Is(err, ErrInvalidQuantity):
		return HTTPError{http.StatusBadRequest, MsgInvalidQuantity}

	case errors.Is(err, sql.ErrNoRows):
		return HTTPError{http.StatusNotFound, MsgResourceNotFound}

	default:
		return HTTPError{http.StatusInternalServerError, MsgCheckoutFailed}
	}
}
