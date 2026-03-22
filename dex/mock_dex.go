package dex

import (
	"dex-aggregator/model"
	"math/rand"
	"time"
)

// MockDEX 是一个用于本地测试和并发验证的假 DEX 实现
type MockDEX struct {
	dexName string
	minWait int // 模拟网络延迟的最小毫秒数
	maxWait int // 模拟网络延迟的最大毫秒数
}

// NewMockDEX 创建一个新的 MockDEX 实例
func NewMockDEX(name string, minWait, maxWait int) *MockDEX {
	return &MockDEX{
		dexName: name,
		minWait: minWait,
		maxWait: maxWait,
	}
}

// Name 实现了 DEX 接口的 Name 方法
func (m *MockDEX) Name() string {
	return m.dexName
}

// GetQuote 实现了 DEX 接口的询价方法，它会随机睡眠一段时间并返回一个随机价格
func (m *MockDEX) GetQuote(fromToken, toToken string, amount float64) (*model.QuoteResult, error) {
	// 1. 模拟网络延迟 (Sleep)
	// 计算需要睡眠的随机毫秒数
	sleepMs := rand.Intn(m.maxWait-m.minWait+1) + m.minWait
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)

	// 2. 模拟返回的价格
	// 假设 1 ETH 兑换的 USDC 在 1800 到 1900 之间波动
	basePrice := 1800.0
	randomFluctuation := rand.Float64() * 100.0
	finalPrice := basePrice + randomFluctuation

	amountOut := amount * finalPrice

	// 3. 构建返回结果
	result := &model.QuoteResult{
		DEXName:   m.dexName,
		AmountOut: amountOut,
		Latency:   int64(sleepMs),
		Err:       nil,
	}

	return result, nil
}
