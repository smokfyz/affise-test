package urls

import (
	"encoding/json"
	"net/http"

	"github.com/smokfyz/affise-test/pkg/log"
)

const maxNumberOfUrls = 20
const maxSimultaniousRequests = 100
const RequestIDKey contextKey = "requestID"

var requestsSem = make(chan struct{}, maxSimultaniousRequests)

type contextKey string

type indexRequestBody struct {
	Urls []string `json:"urls"`
}

type indexResponseBody struct {
	Results []string `json:"results"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	select {
	case requestsSem <- struct{}{}:
		defer func() {
			<-requestsSem
		}()
	default:
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}

	ctx := r.Context()

	requestID := ctx.Value(RequestIDKey).(string)

	switch r.Method {
	case http.MethodPost:
		var reqBody indexRequestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Debug.Printf("requestID: %s, failed to parse request body: %v", requestID, err)
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}

		if len(reqBody.Urls) > maxNumberOfUrls {
			log.Debug.Printf("requestID: %s, too many urls", requestID)
			http.Error(w, "too many urls", http.StatusBadRequest)
			return
		}

		results, err := requestUrls(ctx, reqBody.Urls)
		if err != nil {
			log.Debug.Printf("requestID: %s, failed to get urls: %v", requestID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(indexResponseBody{results})
		if err != nil {
			log.Debug.Printf("requestID: %s, failed to marshal response body: %v", requestID, err)
			http.Error(w, "failed to marshal response body", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(respBody)
		if err != nil {
			log.Debug.Printf("requestID: %s, failed to write response body: %v", requestID, err)
			http.Error(w, "failed to write response body", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
