package bitbank

import "fmt"

type BitbankError struct {
	Message string
}

func (e *BitbankError) Error() string { return e.Message }

type BitbankAPIError struct {
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error"`
	Message    string `json:"message"`
}

func (e *BitbankAPIError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.ErrorCode, e.Message)
}

type AuthenticationError struct{ BitbankAPIError }
type ForbiddenError struct{ BitbankAPIError }
type NotFoundError struct{ BitbankAPIError }
type InsufficientCreditsError struct{ BitbankAPIError }
type RateLimitError struct{ BitbankAPIError }
type ValidationError struct{ BitbankAPIError }
type ServerError struct{ BitbankAPIError }
type ConnectionError struct{ BitbankError }
type WebSocketError struct{ BitbankError }

func mapError(statusCode int, errorCode, message string) error {
	base := BitbankAPIError{StatusCode: statusCode, ErrorCode: errorCode, Message: message}
	switch statusCode {
	case 400:
		return &ValidationError{base}
	case 401:
		return &AuthenticationError{base}
	case 402:
		return &InsufficientCreditsError{base}
	case 403:
		return &ForbiddenError{base}
	case 404:
		return &NotFoundError{base}
	case 429:
		return &RateLimitError{base}
	default:
		if statusCode >= 500 {
			return &ServerError{base}
		}
		return &base
	}
}
