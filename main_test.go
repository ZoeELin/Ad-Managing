package main_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPublicAPIPerformance(t *testing.T) {
	requestRate := 10000                // 每秒 10,000 次請求
	requestDuration := 10 * time.Second // 測試持續 10 秒
	activeAdsLimit := 1000              // 活躍廣告數量上限
	dailyAdsLimit := 3000               // 每日新增廣告數量上限
	activeAdsCount := 0                 // 活躍廣告計數器
	dailyAdsCount := 0                  // 每日新增廣告計數器
	startTime := time.Now()             // 開始測試時間

	// 模擬 Public API 的請求
	for i := 0; i < requestRate*int(requestDuration.Seconds()); i++ {
		// 模擬廣告新增
		// 假設這裡有一個函數 AddAdvertisement() 可以新增廣告，並返回是否成功
		success := AddAdvertisement()
		if success {
			dailyAdsCount++
		}

		// 檢查活躍廣告數量是否超過上限
		if activeAdsCount >= activeAdsLimit {
			t.Fatalf("活躍廣告數量超過上限：%d", activeAdsCount)
		}

		// 檢查每日新增廣告數量是否超過上限
		if dailyAdsCount > dailyAdsLimit {
			t.Fatalf("每日新增廣告數量超過上限：%d", dailyAdsCount)
		}

		// 假設這裡有一個函數 CheckActiveAds() 可以檢查活躍廣告數量並返回數量
		activeAdsCount = CheckActiveAds()
	}

	// 計算測試結束後的持續時間
	duration := time.Since(startTime)

	// 模擬測試結果檢查
	assert.LessOrEqual(t, activeAdsCount, activeAdsLimit, "活躍廣告數量超過上限")
	assert.LessOrEqual(t, dailyAdsCount, dailyAdsLimit, "每日新增廣告數量超過上限")
	assert.GreaterOrEqual(t, requestRate, int(float64(dailyAdsCount)/duration.Seconds()), "每秒新增廣告速率低於請求速率")
}

// 模擬新增廣告的函數
func AddAdvertisement() bool {
	// 執行新增廣告的操作，返回是否成功
	return true
}

// 模擬檢查活躍廣告數量的函數
func CheckActiveAds() int {
	// 假設這裡可以查詢活躍廣告數量並返回
	return 0
}
