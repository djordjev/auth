package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/djordjev/auth/internal/domain"
	"net/http"
)

func parseRequest(req *http.Request, target any) error {
	err := json.NewDecoder(req.Body).Decode(target)

	if err != nil {
		path := req.URL.Path
		return NewApiError(fmt.Sprintf("path %s: invalid request format", path), err, http.StatusBadRequest)
	}

	return nil
}

func mustWriteJSONResponse(w http.ResponseWriter, res any) {
	responseBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	_, errWrite := w.Write(responseBytes)
	if errWrite != nil {
		panic(errWrite)
	}
}

func respondWithError(w http.ResponseWriter, err error) {
	var domainError domain.Error

	if errors.As(err, &domainError) {
		isCritical := domainError.IsCritical()
		var statusCode int
		if isCritical {
			statusCode = http.StatusInternalServerError
		} else {
			statusCode = http.StatusBadRequest
		}

		http.Error(w, domainError.Error(), statusCode)
		return
	}

	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
