package urls

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/smokfyz/affise-test/pkg/log"
)

const maxSimultaneousRequests = 4
const maxBodySize = 5 << 20 // 5MB
const requestUrlTimeout = 1

type orderedData struct {
	index   int
	content string
}

func requestUrl(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	limitedReader := io.LimitReader(resp.Body, maxBodySize+1)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}

	if len(body) > maxBodySize {
		return "", fmt.Errorf("response body exceeds the limit of %d bytes", maxBodySize)
	}

	return string(body), nil
}

func requestUrls(ctx context.Context, urls []string) ([]string, error) {
	requestID := ctx.Value(RequestIDKey).(string)

	log.Debug.Printf("requestID: %s, requesting urls", requestID)

	ctx, cancel := context.WithCancelCause(ctx)

	client := &http.Client{
		Timeout: time.Second * requestUrlTimeout,
	}

	wg := sync.WaitGroup{}

	results := make([]string, len(urls))

	input := make(chan orderedData, len(urls))
	output := make(chan orderedData, len(urls))
	defer close(input)
	defer close(output)

	for range maxSimultaneousRequests {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case inputData, ok := <-input:
					if !ok {
						return
					}

					url := inputData.content

					body, err := requestUrl(client, url)
					if err != nil {
						log.Debug.Printf("requestID: %s, failed to request url %s: %v", requestID, url, err)
						cancel(fmt.Errorf("failed to request url %s: %w", url, err))
						return
					}

					output <- orderedData{inputData.index, body}
				}
			}
		}()
	}

	for i, url := range urls {
		input <- orderedData{i, url}
	}

	for range len(urls) {
		select {
		case <-ctx.Done():
			log.Debug.Printf("requestID: %s, canceled", requestID)
			wg.Wait()
			return nil, context.Cause(ctx)
		case outputData := <-output:
			results[outputData.index] = outputData.content
		}
	}

	log.Debug.Printf("requestID: %s, all urls requested", requestID)
	return results, nil
}
