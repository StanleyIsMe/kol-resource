package business

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func UseCaesErrorToErrorResp(err error) (int, ErrorResponse) {
	var usecaseErr UseCaseError
	if !errors.As(err, &usecaseErr) {
		return http.StatusInternalServerError, ErrorResponse{
			ErrorCode:    "INTERNAL_SERVER_ERROR",
			ErrorMessage: "Internal Server Error",
		}
	}

	return usecaseErr.HTTPStatusCode(), ErrorResponse{
		ErrorCode:    usecaseErr.ErrorCode(),
		ErrorMessage: usecaseErr.ErrorMsg(),
	}
}
