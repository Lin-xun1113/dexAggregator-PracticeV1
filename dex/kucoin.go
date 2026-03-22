package dex

import (
	"dex-aggregator/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type KucoinDEX struct {
	baseURL string
}

func NewKucoinDEX() *KucoinDEX {
	return &KucoinDEX{
		baseURL: "https://api.kucoin.com",
	}
}

func (k *KucoinDEX) Name() string {
	return "🟢 KuCoin (Real API)"
}

func (k *KucoinDEX) GetQuote(fromToken, toToken string, amount float64) (*model.QuoteResult, error) {
	startTime := time.Now()

	// KuCoin 的现货交易对格式也是 ETH-USDT
	symbol := fmt.Sprintf("%s-USDT", fromToken)
	url := fmt.Sprintf("%s/api/v1/market/orderbook/level1?symbol=%s", k.baseURL, symbol)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求KuCoin API失败: %v", err)
	}
	defer resp.Body.Close()

	// KuCoin 返回的 JSON 结构
	var result struct {
		Code string `json:"code"`
		Data struct {
			Price string `json:"price"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	if result.Code != "200000" {
		return nil, fmt.Errorf("KuCoin返回异常, code: %s", result.Code)
	}

	priceFloat, err := strconv.ParseFloat(result.Data.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("价格转换失败: %v", err)
	}

	return &model.QuoteResult{
		DEXName:   k.Name(),
		AmountOut: amount * priceFloat,
		Latency:   time.Since(startTime).Milliseconds(),
		Err:       nil,
	}, nil
}
