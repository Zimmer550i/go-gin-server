package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080"

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

type TestResult struct {
	TotalRequests   int64
	SuccessRequests int64
	FailedRequests  int64
	TotalDuration   time.Duration
	RequestsPerSec  float64
}

func sendRequest(method string, url string, body any) ([]byte, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return responseBody, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return responseBody, nil
}

func sendSimpleRequest(method string, url string, body any) error {
	_, err := sendRequest(method, url, body)
	return err
}

func sendRequestAllowStatus(method string, url string, body any, allowedStatusCodes ...int) ([]byte, int, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}

		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, 0, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)

	for _, allowedStatusCode := range allowedStatusCodes {
		if res.StatusCode == allowedStatusCode {
			return responseBody, res.StatusCode, nil
		}
	}

	return responseBody, res.StatusCode, fmt.Errorf("unexpected status code: %d", res.StatusCode)
}

func createUser(name string, age int) (int64, error) {
	body := map[string]any{
		"name": name,
		"age":  age,
	}

	responseBody, err := sendRequest(http.MethodPost, baseURL+"/users", body)
	if err != nil {
		return 0, err
	}

	id, ok := extractIDFromJSON(responseBody)
	if !ok {
		return 0, fmt.Errorf("could not extract created user ID from response: %s", string(responseBody))
	}

	return id, nil
}

func extractIDFromJSON(responseBody []byte) (int64, bool) {
	var value any

	if err := json.Unmarshal(responseBody, &value); err != nil {
		return 0, false
	}

	return findID(value)
}

func findID(value any) (int64, bool) {
	switch typedValue := value.(type) {
	case map[string]any:
		if rawID, exists := typedValue["id"]; exists {
			return numberToInt64(rawID)
		}

		for _, nestedValue := range typedValue {
			if id, ok := findID(nestedValue); ok {
				return id, true
			}
		}

	case []any:
		for _, nestedValue := range typedValue {
			if id, ok := findID(nestedValue); ok {
				return id, true
			}
		}
	}

	return 0, false
}

func numberToInt64(value any) (int64, bool) {
	switch typedValue := value.(type) {
	case float64:
		return int64(typedValue), true
	case int64:
		return typedValue, true
	case int:
		return int64(typedValue), true
	default:
		return 0, false
	}
}

func cleanupUsers(t *testing.T, userIDs []int64) {
	t.Helper()

	for _, userID := range userIDs {
		if userID == 0 {
			continue
		}

		deleteURL := fmt.Sprintf("%s/users/%d", baseURL, userID)

		_, statusCode, err := sendRequestAllowStatus(
			http.MethodDelete,
			deleteURL,
			nil,
			http.StatusOK,
			http.StatusNoContent,
			http.StatusNotFound,
		)
		if err != nil {
			t.Logf("failed to cleanup user %d: %v", userID, err)
			continue
		}

		if statusCode == http.StatusNotFound {
			t.Logf("cleanup skipped user %d because it was already missing", userID)
		}
	}
}

func runLoadTest(name string, totalRequests int, concurrency int, requestFn func() error) TestResult {
	var successCount int64
	var failedCount int64
	var totalCount int64

	start := time.Now()

	jobs := make(chan int, totalRequests)

	var wg sync.WaitGroup

	for worker := 0; worker < concurrency; worker++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range jobs {
				err := requestFn()

				atomic.AddInt64(&totalCount, 1)

				if err != nil {
					atomic.AddInt64(&failedCount, 1)
				} else {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}()
	}

	for i := 0; i < totalRequests; i++ {
		jobs <- i
	}

	close(jobs)
	wg.Wait()

	duration := time.Since(start)

	rps := float64(totalCount) / duration.Seconds()

	result := TestResult{
		TotalRequests:   totalCount,
		SuccessRequests: successCount,
		FailedRequests:  failedCount,
		TotalDuration:   duration,
		RequestsPerSec:  rps,
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("API Load Test:", name)
	fmt.Println("========================================")
	fmt.Println("Total Requests:  ", result.TotalRequests)
	fmt.Println("Success Requests:", result.SuccessRequests)
	fmt.Println("Failed Requests: ", result.FailedRequests)
	fmt.Println("Total Duration:  ", result.TotalDuration)
	fmt.Printf("Requests/sec:    %.2f\n", result.RequestsPerSec)
	fmt.Println("========================================")
	fmt.Println()

	return result
}

func TestGetUsersAPIPerformance(t *testing.T) {
	totalRequests := 1000
	concurrency := 50

	result := runLoadTest("GET /users", totalRequests, concurrency, func() error {
		return sendSimpleRequest(http.MethodGet, baseURL+"/users", nil)
	})

	if result.FailedRequests > 0 {
		t.Fatalf("expected 0 failed requests, got %d", result.FailedRequests)
	}
}

func TestCreateUserAPIPerformance(t *testing.T) {
	totalRequests := 500
	concurrency := 25

	var userNumber int64
	createdUserIDs := make([]int64, 0, totalRequests)
	var createdUserIDsMutex sync.Mutex

	defer func() {
		cleanupUsers(t, createdUserIDs)
	}()

	result := runLoadTest("POST /users", totalRequests, concurrency, func() error {
		currentUserNumber := atomic.AddInt64(&userNumber, 1)

		createdUserID, err := createUser(
			fmt.Sprintf("Test User %d", currentUserNumber),
			20+int(currentUserNumber%40),
		)
		if err != nil {
			return err
		}

		createdUserIDsMutex.Lock()
		createdUserIDs = append(createdUserIDs, createdUserID)
		createdUserIDsMutex.Unlock()

		return nil
	})

	if result.FailedRequests > 0 {
		t.Fatalf("expected 0 failed requests, got %d", result.FailedRequests)
	}
}

func TestMixedAPIPerformance(t *testing.T) {
	totalRequests := 1000
	concurrency := 50

	var counter int64
	createdUserIDs := make([]int64, 0, totalRequests/2)
	var createdUserIDsMutex sync.Mutex

	defer func() {
		cleanupUsers(t, createdUserIDs)
	}()

	result := runLoadTest("Mixed API Test", totalRequests, concurrency, func() error {
		current := atomic.AddInt64(&counter, 1)

		switch current % 2 {
		case 0:
			return sendSimpleRequest(http.MethodGet, baseURL+"/users", nil)

		default:
			createdUserID, err := createUser(
				fmt.Sprintf("Mixed User %d", current),
				18+int(current%50),
			)
			if err != nil {
				return err
			}

			createdUserIDsMutex.Lock()
			createdUserIDs = append(createdUserIDs, createdUserID)
			createdUserIDsMutex.Unlock()

			return nil
		}
	})

	if result.FailedRequests > 0 {
		t.Fatalf("expected 0 failed requests, got %d", result.FailedRequests)
	}
}