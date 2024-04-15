package test

import (
	"ad-proj/router"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateAdRouter(t *testing.T) {
	// 模擬廣告資訊
	adInfo := map[string]interface{}{
		"title":   "Test Ad",
		"startAt": "2024-04-10T12:00:00Z",
		"endAt":   "2024-04-20T12:00:00Z",
		"conditions": []map[string]interface{}{
			{
				"ageStart": 20,
				"ageEnd":   40,
				"gender":   "M",
				"country":  []string{"TW", "JP", "US"},
				"platform": []string{"ios"},
			},
		},
	}

	// 將廣告資訊轉換為 JSON 格式
	adInfoJSON, _ := json.Marshal(adInfo)

	// Create a test server to emulate an actual HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Login for processing requests
		if r.URL.Path == "/api/v1/ad" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// 準備 HTTP POST 請求
	r := gin.Default()
	r.POST("/api/v1/ad", router.CreateAd)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/ad", bytes.NewBuffer(adInfoJSON))
	req.Header.Set("Content-Type", "application/json")

	// 使用 httptest 庫模擬 HTTP 請求
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	// 驗證 HTTP 狀態碼是否為 201 Created
	assert.Equal(t, http.StatusCreated, recorder.Code)

	// 驗證回傳的 JSON 格式是否包含特定訊息
	expectedResponse := `{"message": "Post the ad successfully."}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}
