package main

import (
	"binance_connector"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// https://testnet.binance.vision
const (
	baseURL   = "https://api.binance.com/"
	apiKey    = ""
	secretKey = ""
)

// 账户信息
type AccountInformation struct {
	omitZeroBalances bool
	recvWindow       int
}

// 获取交易特定的交易和符号
type AccountTradeList struct {
	orderId    int
	startTime  int
	endTime    int
	fromId     int
	limit      int
	recvWindow int
}

// 某个token交易所交易规则和symbol信息
type ExchangeInfo struct {
	symbol      string
	symbols     []string
	permissions string
}

// 历史交易
type HistoryOrders struct {
	symbol string
	limit  int
	fromId int
}

// Compressed/Aggregate 订单列表
type AggregateTrades struct {
	fromId    int
	startTime uint64
	endTime   uint64
	limit     int
}

// k线
type Kline struct {
	symbol    string
	interval  string //enum
	startTime int
	endTime   int
	timeZone  string
	limit     int
}

// UIKlines
type UIKlines struct {
	symbol    string
	interval  string //enum
	startTime int
	endTime   int
	timeZone  string
	limit     int
}

// 一个或多个token
type inputTokens struct {
	symbol  string
	symbols []string
}

// 新订单
type NewOrder struct {
	timeInForce             string
	quantity                float64
	quoteOrderQty           float64
	price                   float64
	newClientOrderId        string
	strategyId              int
	strategyType            int
	stopPrice               float64
	trailingDelta           int
	icebergQty              float64
	newOrderRespType        string
	selfTradePreventionMode string
	recvWindow              int
}

// 取消订单
type CancelOrder struct {
	orderId            int
	origClientOrderId  string
	newClientOrderId   string
	cancelRestrictions string
	recvWindow         int
}

// 所有订单
type AllOrders struct {
	orderId    int
	startTime  int
	endTime    int
	limit      int
	recvWindow int
}

// 当前打开的某个token所有未成交订单
type CurrentTokenAllOpenOrders struct {
	symbol     string
	recvWindow int
}

// 获取账户下的订单
type GetMyTrades struct {
	startTime uint64
	endTime   uint64
	fromId    int64
	limit     int
	orderId   int64
}

// 检查一个订单状态
type QueryOrder struct {
	orderId           int64
	origClientOrderId string
	recvWindow        int
}

// 替代
type CancelReplace struct {
	timeInForce             string
	quantity                float64
	quoteOrderQty           float64
	price                   float64
	cancelNewClientOrderId  string
	cancelOrigClientOrderId string
	cancelOrderId           int64
	newClientOrderId        string
	strategyId              int32
	strategyType            int32
	stopPrice               float64
	trailingDelta           int64
	icebergQty              float64
	newOrderRespType        string
	selfTradePreventionMode string
	cancelRestrictions      string
	recvWindow              int
}

func initClient(thisApiKey, thisSecretKey string, proxy map[string]string) *binance_connector.Client {
	client := binance_connector.NewClient(thisApiKey, thisSecretKey, baseURL)
	if client == nil {
		return nil
	}
	return client
}

// 得到账户信息
func getAccountInformation(
	thisApiKey,
	thisSecretKey string,
	timestamp int64,
	ai AccountInformation,
	proxy map[string]string,
) (*binance_connector.AccountResponse, error) {
	reClient := initClient(apiKey, secretKey, proxy)
	accountInformation, err := reClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(accountInformation))
	return accountInformation, err
}

// GetRequestJSON 发送 HTTP GET 请求到指定的 URL，并返回 JSON 格式的响应数据
func GetRequestJSON(url string) (map[string]interface{}, error) {
	// 发送 HTTP GET 请求
	var reqURL = baseURL + url
	response, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 响应数据
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func ping(
	thisApiKey,
	thisSecretKey string,
	proxy map[string]string,
) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// NewPingService
	ping := reClient.NewPingService().Do(context.Background())
	fmt.Println(binance_connector.PrettyPrint(ping))
}

// 得到binance系统时间
func getTime(
	thisApiKey,
	thisSecretKey string,
	proxy map[string]string,
) (*binance_connector.ServerTimeResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// NewServerTimeService
	serverTime, err := reClient.NewServerTimeService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(serverTime))
	return serverTime, err
}

// 得到当前交易所所有token交易规则和symbol信息
func getExchangeInfo(
	thisApiKey,
	thisSecretKey string,
	proxy map[string]string,
) (*binance_connector.ExchangeInfoResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	exchangeInfo, err := reClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(exchangeInfo))
	return exchangeInfo, err
}

// 得到OrderBook深度
func getOrderBookDepth(
	thisApiKey,
	thisSecretKey,
	symbol string,
	limit int,
	proxy map[string]string,
) (*binance_connector.OrderBookResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	orderBook, err := reClient.NewOrderBookService().
		Symbol(symbol).Limit(limit).Do(context.Background())

	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(orderBook))
	return orderBook, err
}

// 近期交易列表
func getRecentTradeList(
	thisApiKey,
	thisSecretKey,
	symbol string,
	limit int,
	proxy map[string]string,
) ([]*binance_connector.RecentTradesListResponse, error) {

	reClient := initClient(thisApiKey, thisSecretKey, proxy)

	// RecentTradesList
	recentTradesList, err := reClient.NewRecentTradesListService().
		Symbol(symbol).Limit(limit).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(recentTradesList))
	return recentTradesList, err
}

// HistoricalTradeLookup
func getHistoryTrades(
	thisApiKey,
	thisSecretKey,
	symbol string,
	fromId int64,
	limit uint,
	proxy map[string]string,
) ([]*binance_connector.RecentTradesListResponse, error) {
	reClient := initClient(apiKey, secretKey, proxy)
	historicalTradeLookup, err := reClient.NewHistoricalTradeLookupService().
		Symbol(symbol).FromId(fromId).Limit(limit).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(historicalTradeLookup))
	return historicalTradeLookup, err
}

// 得到总成交量
func getAggTradesList(
	thisApiKey,
	thisSecretKey,
	symbol string,
	at AggregateTrades,
	proxy map[string]string,
) ([]*binance_connector.AggTradesListResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// AggTradesList
	aggTradesList, err := reClient.NewAggTradesListService().
		Symbol(symbol).FromId(at.fromId).Limit(at.limit).StartTime(at.startTime).
		EndTime(at.endTime).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(aggTradesList))
	return aggTradesList, err
}

// ticker
func getTicker(
	thisApiKey,
	thisSecretKey,
	symbol string,
	proxy map[string]string,
) (*binance_connector.TickerResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// Ticker
	ticker, err := reClient.NewTickerService().
		Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(ticker))
	return ticker, err
}

// 得到k线
// func getKlinesData(k Kline) ([]*binance_connector.KlinesResponse, error) {
// 	reClient := initClient(apiKey, secretKey)
// 	// Klines
// 	klines, err := reClient.NewKlinesService().
// 		Symbol(k.symbol).Interval(k.interval).Limit(k.limit).
// 		StartTime(uint64(k.startTime)).
// 		EndTime(uint64(k.endTime)).Do(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(binance_connector.PrettyPrint(klines))
// 	return klines, err
// }

// UIKlines
// func getUIKlines(uik UIKlines) ([]*binance_connector.UiKlinesResponse, error) {
// 	reClient := initClient(apiKey, secretKey)
// 	// UiKlines
// 	uiKlines, err := reClient.NewUIKlinesService().
// 		Symbol(uik.symbol).Limit(uik.limit).Interval(uik.interval).
// 		StartTime(uint64(uik.startTime)).EndTime(uint64(uik.endTime)).Do(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(binance_connector.PrettyPrint(uiKlines))
// 	return uiKlines, err
// }

// 一个token的当前平均价格。
func getAvgPrice(
	thisApiKey,
	thisSecretKey,
	symbol string,
	proxy map[string]string,
) (*binance_connector.AvgPriceResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// AvgPrice
	avgPrice, err := reClient.NewAvgPriceService().
		Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(avgPrice))
	return avgPrice, err
}

// 24小时滚动窗价格变动统计。
func getTicker24hrPrice(
	thisApiKey,
	thisSecretKey string,
	it inputTokens,
	proxy map[string]string,
) (*binance_connector.Ticker24hrResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// Ticker24hr
	ticker24hr, err := reClient.NewTicker24hrService().
		Symbol(it.symbol).Symbols(it.symbols).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(ticker24hr))
	return ticker24hr, err
}

// 一个或多个股票的最新价格
func getTickersPrice(
	thisApiKey,
	thisSecretKey string,
	it inputTokens,
	proxy map[string]string,
) (*binance_connector.TickerPriceResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)

	TickerPrice, err := reClient.NewTickerPriceService().
		Symbol(it.symbol).Symbols(it.symbols).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(TickerPrice))
	return TickerPrice, err
}

// 一个或多个股票的订单簿上的最佳价格/数量。
func getSymbolOrderBookTicker(
	thisApiKey,
	thisSecretKey string,
	it inputTokens,
	proxy map[string]string,
) ([]*binance_connector.TickerBookTickerResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)

	TickerBookTicker, err := reClient.NewTickerBookTickerService().
		Symbol(it.symbol).Symbols(it.symbols).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(TickerBookTicker))
	return TickerBookTicker, err
}

// 得到所有订单
func getAllOrders(
	thisApiKey,
	thisSecretKey,
	symbol string,
	timestamp int64,
	ao AllOrders,
	proxy map[string]string,
) ([]*binance_connector.NewAllOrdersResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)

	// Binance Get all account orders; active, canceled, or filled - GET /api/v3/allOrders
	getAllOrders, err := reClient.NewGetAllOrdersService().Symbol(symbol).
		OrderId(int64(ao.orderId)).StartTime(uint64(ao.startTime)).
		EndTime(uint64(ao.endTime)).Limit(ao.limit).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(getAllOrders))
	return getAllOrders, err
}

// 得到某个token当前打开的所有未成交订单
func getCurrentOpenOrders(
	thisApiKey,
	thisSecretKey,
	symbol string,
	timestamp int64,
	proxy map[string]string,
) ([]*binance_connector.NewOpenOrdersResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// Binance Get current open orders - GET /api/v3/openOrders
	getCurrentOpenOrders, err := reClient.NewGetOpenOrdersService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(getCurrentOpenOrders))
	return getCurrentOpenOrders, err
}

// 获取特定账户的交易
func getGetMyTrades(
	thisApiKey,
	thisSecretKey,
	symbol string,
	gmt GetMyTrades,
	proxy map[string]string,
) ([]*binance_connector.AccountTradeListResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// Binance Get trades for a specific account and symbol (USER_DATA) - GET /api/v3/myTrades
	getMyTradesService, err := reClient.NewGetMyTradesService().
		Symbol(symbol).StartTime(gmt.startTime).EndTime(gmt.endTime).FromId(gmt.fromId).
		Limit(gmt.limit).OrderId(gmt.orderId).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(getMyTradesService))
	return getMyTradesService, nil
}

// 检查一个订单状态
func getQueryOrder(
	thisApiKey,
	thisSecretKey,
	symbol string,
	timestamp int64,
	qo QueryOrder,
	proxy map[string]string,
) (*binance_connector.GetOrderResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	// Binance Query Order (USER_DATA) - GET /api/v3/order
	queryOrder, err := reClient.NewGetOrderService().Symbol(symbol).OrderId(qo.orderId).
		OrigClientOrderId(qo.origClientOrderId).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(queryOrder))
	return queryOrder, err
}

// 查询当前订单计数使用情况
func QueryCurrentOrderCountUsage(
	thisApiKey,
	thisSecretKey string,
	proxy map[string]string,
) ([]*binance_connector.QueryCurrentOrderCountUsageResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)

	// Query Current Order Count Usage (TRADE)
	getQueryCurrentOrderCountUsageService, err := reClient.NewGetQueryCurrentOrderCountUsageService().
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(getQueryCurrentOrderCountUsageService))
	return getQueryCurrentOrderCountUsageService, err
}

// 创建新订单
func createNewOrder(
	thisApiKey,
	thisSecretKey,
	symbol string,
	side string,
	orderType string,
	timestamp int64,
	no NewOrder,
	proxy map[string]string,
) (interface{}, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)

	newOrder, err := reClient.NewCreateOrderService().Symbol(symbol).
		Side(side).Type(orderType).IcebergQuantity(no.icebergQty).
		NewClientOrderId(no.newClientOrderId).NewOrderRespType(no.newOrderRespType).
		Price(no.price).Quantity(no.quantity).
		QuoteOrderQty(no.quoteOrderQty).SelfTradePreventionMode(no.selfTradePreventionMode).
		StopPrice(no.stopPrice).StrategyId(no.strategyId).StrategyType(no.strategyType).
		TimeInForce(no.timeInForce).TrailingDelta(no.trailingDelta).Do(context.Background())

	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(newOrder))
	return newOrder, err
}

// 取消某个token订单
func cancelOrder(
	thisApiKey,
	thisSecretKey,
	symbol string,
	co CancelOrder,
	proxy map[string]string,
) (*binance_connector.CancelOrderResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	cancelOrder, err := reClient.NewCancelOrderService().Symbol(symbol).
		OrderId(int64(co.orderId)).OrigClientOrderId(co.origClientOrderId).
		NewClientOrderId(co.newClientOrderId).CancelRestrictions(co.cancelRestrictions).Do(context.Background())

	if err != nil {
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(cancelOrder))
	return cancelOrder, err
}

// 取消某个token所有开放的orders
func cancelSymbolAllOpenOrders(
	thisApiKey,
	thisSecretKey,
	symbol string,
	timestamp int64,
	proxy map[string]string,
) ([]*binance_connector.CancelOrderResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	cancelOpenOrders, err := reClient.NewCancelOpenOrdersService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(cancelOpenOrders))
	return cancelOpenOrders, err
}

// 取消某个token下的订单后立即创建一个订单
func cancelReplace(
	thisApiKey,
	thisSecretKey,
	symbol string,
	side string,
	orderType string,
	cancelReplaceMode string,
	timestamp int64,
	cr CancelReplace,
	proxy map[string]string,
) (*binance_connector.CancelReplaceResponse, error) {
	reClient := initClient(thisApiKey, thisSecretKey, proxy)
	cancelReplace, err := reClient.NewCancelReplaceService().
		Symbol(symbol).Side(side).OrderType(orderType).CancelRestrictions(cr.cancelRestrictions).
		CancelReplaceMode(cancelReplaceMode).CancelOrderId(cr.cancelOrderId).
		CancelNewClientOrderId(cr.cancelNewClientOrderId).CancelOrigClientOrderId(cr.cancelOrigClientOrderId).
		TimeInForce(cr.timeInForce).IcebergQty(cr.icebergQty).Quantity(cr.quantity).
		Price(cr.price).NewOrderRespType(cr.newOrderRespType).NewClientOrderId(cr.newClientOrderId).
		SelfTradePreventionMode(cr.selfTradePreventionMode).StrategyId(cr.strategyId).StrategyType(cr.strategyType).
		StopPrice(cr.stopPrice).TimeInForce(cr.timeInForce).TrailingDelta(cr.trailingDelta).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(cancelReplace))
	return cancelReplace, err
}

func main() {
}
