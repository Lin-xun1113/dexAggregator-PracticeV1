package dex

import (
	"dex-aggregator/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type OkxDEX struct {
	baseURL string
}

func NewOkxDEX() *OkxDEX {
	return &OkxDEX{
		baseURL: "https://www.okx.com",
	}
}

func (o *OkxDEX) Name() string {
	return "⚫ OKX (Real API)"
}

func (o *OkxDEX) GetQuote(fromToken, toToken string, amount float64) (*model.QuoteResult, error) {
	startTime := time.Now()

	// OKX 的现货交易对格式是 ETH-USDT
	symbol := fmt.Sprintf("%s-USDT", fromToken)
	url := fmt.Sprintf("%s/api/v5/market/ticker?instId=%s", o.baseURL, symbol)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求OKX API失败: %v", err)
	}
	defer resp.Body.Close()

	// OKX 返回的 JSON 结构
	var result struct {
		Code string `json:"code"`
		Data []struct {
			Last string `json:"last"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	if result.Code != "0" || len(result.Data) == 0 {
		return nil, fmt.Errorf("OKX返回异常或未找到数据")
	}

	priceFloat, err := strconv.ParseFloat(result.Data[0].Last, 64)
	if err != nil {
		return nil, fmt.Errorf("价格转换失败: %v", err)
	}

	return &model.QuoteResult{
		DEXName:   o.Name(),
		AmountOut: amount * priceFloat,
		Latency:   time.Since(startTime).Milliseconds(),
		Err:       nil,
	}, nil
}
