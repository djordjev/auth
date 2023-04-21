package api

import (
	"encoding/json"
	"fmt"
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
