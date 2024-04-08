package api_test

import (
	"ad-proj/controllers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAdE2E(t *testing.T) {
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

	// 準備 HTTP POST 請求
	req, _ := http.NewRequest("POST", "/api/v1/ad", bytes.NewBuffer(adInfoJSON))
	req.Header.Set("Content-Type", "application/json")

	// 使用 httptest 庫模擬 HTTP 請求
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateAd)
	handler.ServeHTTP(recorder, req)

	// 驗證 HTTP 狀態碼是否為 201 Created
	assert.Equal(t, http.StatusCreated, recorder.Code)

	// 驗證回傳的 JSON 格式是否包含特定訊息
	expectedResponse := `{"message":"Post the ad successfully."}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}
