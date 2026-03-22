package dex

import (
	"dex-aggregator/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// 定义一个 BinanceDEX 结构体，可以存一些配置，比如 API 基础地址
type BinanceDEX struct {
	baseURL string
}

// 提供一个“构造函数”
func NewBinanceDEX() *BinanceDEX {
	return &BinanceDEX{
		baseURL: "https://api.binance.com",
	}
}

// 隐式实现接口的 Name() 方法
func (b *BinanceDEX) Name() string {
	return "🔶 Binance (Real API)"
}

func (b *BinanceDEX) GetQuote(fromToken, toToken string, amount float64) (*model.QuoteResult, error) {
	startTime := time.Now()

	symbol := fromToken + "USDT"

	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", b.baseURL, symbol)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求BinanceAPI失败: %v", err)
	}

	defer resp.Body.Close() //延迟关闭数据流，否则会造成内存泄露

	//准备解析 JSON。币安返回的格式是 {"symbol":"ETHUSDT", "price":"3500.12"}

	var result struct{
		Symbol string `json:"symbol"`
		Price  string `json:"price"` //坑点：API返回的 price 是带引号的字符串格式
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 5. 字符串转数字
	priceFloat, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("价格转换失败: %v", err)
	}
	// 6. 计算最终换出来的钱，并打包发走
	return &model.QuoteResult{
		DEXName:   b.Name(),
		AmountOut: amount * priceFloat,
		Latency:   time.Since(startTime).Milliseconds(),
		Err:       nil, // 没错误就填 nil
	}, nil
}
