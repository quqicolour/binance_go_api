package main

import (
	"binance_connector"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
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
	fromId    *int
	startTime *uint64
	endTime   *uint64
	limit     *int
}

// 一个或多个token
type inputTokens struct {
	symbol  *string
	symbols *[]string
}

// 新订单
type NewOrder struct {
	timeInForce             *string
	quantity                *float64
	quoteOrderQty           *float64
	price                   *float64
	newClientOrderId        *string
	strategyId              *int
	strategyType            *int
	stopPrice               *float64
	trailingDelta           *int
	icebergQty              *float64
	newOrderRespType        *string
	selfTradePreventionMode *string
	recvWindow              *int
}

// 取消订单
type CancelOrder struct {
	orderId            *int64
	origClientOrderId  *string
	newClientOrderId   *string
	cancelRestrictions *string
	recvWindow         *int
}

// 所有订单
type AllOrders struct {
	orderId    *int64
	startTime  *uint64
	endTime    *uint64
	limit      *int
	recvWindow *int
}

// 当前打开的某个token所有未成交订单
type CurrentTokenAllOpenOrders struct {
	symbol     string
	recvWindow int
}

// 获取账户下的订单
type GetMyTrades struct {
	startTime *uint64
	endTime   *uint64
	fromId    *int64
	limit     *int
	orderId   *int64
}

// 检查一个订单状态
type QueryOrder struct {
	orderId           *int64
	origClientOrderId *string
	recvWindow        *int
}

// 替代
type CancelReplace struct {
	timeInForce             *string
	quantity                *float64
	quoteOrderQty           *float64
	price                   *float64
	cancelNewClientOrderId  *string
	cancelOrigClientOrderId *string
	cancelOrderId           *int64
	newClientOrderId        *string
	strategyId              *int32
	strategyType            *int32
	stopPrice               *float64
	trailingDelta           *int64
	icebergQty              *float64
	newOrderRespType        *string
	selfTradePreventionMode *string
	cancelRestrictions      *string
	recvWindow              *int
}

func initClient(thisApiKey, thisSecretKey string, proxyURL string) (*binance_connector.Client, error) {
	timeout := time.Second * 10
	client := binance_connector.NewClient(thisApiKey, thisSecretKey, baseURL)
	if client == nil {
		return nil, nil
	}

	// 如果代理 URL 不为空，则设置代理
	if proxyURL != "" {
		// 解析代理 URL 字符串为 *url.URL 类型
		proxyParsed, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		proxy := &http.Transport{
			Proxy: http.ProxyURL(proxyParsed),
		}
		client.HTTPClient = &http.Client{
			Timeout:   timeout,
			Transport: proxy,
		}
	}

	return client, nil
}

// 得到账户信息
func getAccountInformation(
	thisApiKey,
	thisSecretKey string,
	timestamp int64,
	ai AccountInformation,
	proxyURL string,
) (*binance_connector.AccountResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
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
	proxyURL string,
) error {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil
	}
	// NewPingService
	ping := reClient.NewPingService().Do(context.Background())
	fmt.Println(binance_connector.PrettyPrint(ping))
	return ping
}

// 得到binance系统时间
func getTime(
	thisApiKey,
	thisSecretKey string,
	proxyURL string,
) (*binance_connector.ServerTimeResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
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
	proxyURL string,
) (*binance_connector.ExchangeInfoResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
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
	limit *int,
	proxyURL string,
) (*binance_connector.OrderBookResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	service := reClient.NewOrderBookService().Symbol(symbol)
	if limit != nil {
		service = service.Limit(*limit)
	}
	// orderBook, err := reClient.NewOrderBookService().
	// 	Symbol(symbol).Limit(*limit).Do(context.Background())
	orderBook, err := service.Do(context.Background())
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
	limit *int,
	proxyURL string,
) ([]*binance_connector.RecentTradesListResponse, error) {

	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}

	service := reClient.NewRecentTradesListService().Symbol(symbol)
	if limit != nil {
		service = service.Limit(*limit)
	}

	// RecentTradesList
	// recentTradesList, err := reClient.NewRecentTradesListService().
	// 	Symbol(symbol).Limit(limit).Do(context.Background())
	recentTradesList, err := service.Do(context.Background())
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
	fromId *int64,
	limit *int64,
	proxyURL string,
) ([]*binance_connector.RecentTradesListResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	service := reClient.NewHistoricalTradeLookupService().Symbol(symbol)
	if fromId != nil {
		service = service.FromId(*fromId)
	}
	if limit != nil {
		service = service.FromId(*limit)
	}
	// historicalTradeLookup, err := reClient.NewHistoricalTradeLookupService().
	// 	Symbol(symbol).FromId(fromId).Limit(limit).Do(context.Background())
	historicalTradeLookup, err := service.Do(context.Background())
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
	proxyURL string,
) ([]*binance_connector.AggTradesListResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	// AggTradesList
	service := reClient.NewAggTradesListService().Symbol(symbol)
	if at.fromId != nil {
		service = service.FromId(*at.fromId)
	}
	if at.limit != nil {
		service = service.Limit(*at.limit)
	}
	if at.startTime != nil {
		service = service.StartTime(*at.startTime)
	}
	if at.endTime != nil {
		service = service.EndTime(*at.endTime)
	}
	// AggTradesList
	// aggTradesList, err := reClient.NewAggTradesListService().
	// 	Symbol(symbol).FromId(at.fromId).Limit(at.limit).StartTime(at.startTime).
	// 	EndTime(at.endTime).Do(context.Background())
	aggTradesList, err := service.Do(context.Background())
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
	symbol,
	tickerType,
	windowSize string,
	proxyURL string,
) (*binance_connector.TickerResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	// Ticker
	ticker, err := reClient.NewTickerService().
		Symbol(symbol).Type(tickerType).WindowSize(windowSize).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(ticker))
	return ticker, err
}

// 一个token的当前平均价格。
func getAvgPrice(
	thisApiKey,
	thisSecretKey,
	symbol string,
	proxyURL string,
) (*binance_connector.AvgPriceResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
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
	proxyURL string,
) (*binance_connector.Ticker24hrResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	// Ticker24hr
	service := reClient.NewTicker24hrService()
	if it.symbol != nil {
		service = service.Symbol(*it.symbol)
	}
	if it.symbols != nil {
		service = service.Symbols(*it.symbols)
	}
	// Ticker24hr
	// ticker24hr, err := reClient.NewTicker24hrService().
	// 	Symbol(it.symbol).Symbols(it.symbols).Do(context.Background())
	ticker24hr, err := service.Do(context.Background())
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
	proxyURL string,
) (*binance_connector.TickerPriceResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}

	service := reClient.NewTickerPriceService()
	if it.symbol != nil {
		service = service.Symbol(*it.symbol)
	}
	if it.symbols != nil {
		service = service.Symbols(*it.symbols)
	}

	// TickerPrice, err := reClient.NewTickerPriceService().
	// 	Symbol(it.symbol).Symbols(it.symbols).Do(context.Background())
	TickerPrice, err := service.Do(context.Background())
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
	proxyURL string,
) ([]*binance_connector.TickerBookTickerResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	service := reClient.NewTickerBookTickerService()
	if it.symbol != nil {
		service = service.Symbol(*it.symbol)
	}
	if it.symbols != nil {
		service = service.Symbols(*it.symbols)
	}

	// TickerBookTicker, err := reClient.NewTickerBookTickerService().
	// 	Symbol(it.symbol).Symbols(it.symbols).Do(context.Background())
	TickerBookTicker, err := service.Do(context.Background())
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
	proxyURL string,
) ([]*binance_connector.NewAllOrdersResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}

	service := reClient.NewGetAllOrdersService().Symbol(symbol)
	if ao.orderId != nil {
		service = service.OrderId(*ao.orderId)
	}
	if ao.startTime != nil {
		service = service.StartTime(*ao.startTime)
	}
	if ao.endTime != nil {
		service = service.EndTime(*ao.endTime)
	}
	if ao.limit != nil {
		service = service.Limit(*ao.limit)
	}
	// Binance Get all account orders; active, canceled, or filled - GET /api/v3/allOrders
	// getAllOrders, err := reClient.NewGetAllOrdersService().Symbol(symbol).
	// 	OrderId(ao.orderId).StartTime(ao.startTime).
	// 	EndTime(ao.endTime).Limit(ao.limit).Do(context.Background())
	getAllOrders, err := service.Do(context.Background())
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
	proxyURL string,
) ([]*binance_connector.NewOpenOrdersResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
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
	proxyURL string,
) ([]*binance_connector.AccountTradeListResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	service := reClient.NewGetMyTradesService().Symbol(symbol)
	if gmt.fromId != nil {
		service = service.FromId(*gmt.fromId)
	}
	if gmt.limit != nil {
		service = service.Limit(*gmt.limit)
	}
	if gmt.orderId != nil {
		service = service.OrderId(*gmt.orderId)
	}
	if gmt.startTime != nil {
		service = service.StartTime(*gmt.startTime)
	}
	if gmt.endTime != nil {
		service = service.EndTime(*gmt.endTime)
	}
	// Binance Get trades for a specific account and symbol (USER_DATA) - GET /api/v3/myTrades
	// getMyTradesService, err := reClient.NewGetMyTradesService().
	// 	Symbol(symbol).StartTime(gmt.startTime).EndTime(gmt.endTime).FromId(gmt.fromId).
	// 	Limit(gmt.limit).OrderId(gmt.orderId).Do(context.Background())
	getMyTradesService, err := service.Do(context.Background())
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
	proxyURL string,
) (*binance_connector.GetOrderResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	service := reClient.NewGetOrderService().Symbol(symbol)
	if qo.orderId != nil {
		service = service.OrderId(*qo.orderId)
	}
	if qo.origClientOrderId != nil {
		service = service.OrigClientOrderId(*qo.origClientOrderId)
	}

	// Binance Query Order (USER_DATA) - GET /api/v3/order
	// queryOrder, err := reClient.NewGetOrderService().Symbol(symbol).OrderId(qo.orderId).
	// 	OrigClientOrderId(qo.origClientOrderId).Do(context.Background())
	queryOrder, err := service.Do(context.Background())
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
	proxyURL string,
) ([]*binance_connector.QueryCurrentOrderCountUsageResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}

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
	proxyURL string,
) (interface{}, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	service := reClient.NewCreateOrderService().Symbol(symbol).Side(side).Type(orderType)

	if no.icebergQty != nil {
		service = service.IcebergQuantity(*no.icebergQty)
	}
	if no.newClientOrderId != nil {
		service = service.NewClientOrderId(*no.newClientOrderId)
	}
	if no.newOrderRespType != nil {
		service = service.NewOrderRespType(*no.newOrderRespType)
	}
	if no.price != nil {
		service = service.Price(*no.price)
	}
	if no.quantity != nil {
		service = service.Quantity(*no.quantity)
	}
	if no.quoteOrderQty != nil {
		service = service.QuoteOrderQty(*no.quoteOrderQty)
	}
	if no.selfTradePreventionMode != nil {
		service = service.SelfTradePreventionMode(*no.selfTradePreventionMode)
	}
	if no.stopPrice != nil {
		service = service.StopPrice(*no.stopPrice)
	}
	if no.strategyId != nil {
		service = service.StrategyId(*no.strategyId)
	}
	if no.strategyType != nil {
		service = service.StrategyType(*no.strategyType)
	}
	if no.timeInForce != nil {
		service = service.TimeInForce(*no.timeInForce)
	}
	if no.trailingDelta != nil {
		service = service.TrailingDelta(*no.trailingDelta)
	}

	// newOrder, err := reClient.NewCreateOrderService().Symbol(symbol).
	// 	Side(side).Type(orderType).IcebergQuantity(no.icebergQty).
	// 	NewClientOrderId(no.newClientOrderId).NewOrderRespType(no.newOrderRespType).
	// 	Price(no.price).Quantity(no.quantity).
	// 	QuoteOrderQty(no.quoteOrderQty).SelfTradePreventionMode(no.selfTradePreventionMode).
	// 	StopPrice(no.stopPrice).StrategyId(no.strategyId).StrategyType(no.strategyType).
	// 	TimeInForce(no.timeInForce).TrailingDelta(no.trailingDelta).Do(context.Background())
	newOrder, err := service.Do(context.Background())

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
	proxyURL string,
) (*binance_connector.CancelOrderResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	service := reClient.NewCancelOrderService().Symbol(symbol)
	if co.orderId != nil {
		service = service.OrderId(*co.orderId)
	}
	if co.origClientOrderId != nil {
		service = service.OrigClientOrderId(*co.origClientOrderId)
	}
	if co.newClientOrderId != nil {
		service = service.NewClientOrderId(*co.newClientOrderId)
	}
	if co.cancelRestrictions != nil {
		service = service.CancelRestrictions(*co.cancelRestrictions)
	}
	// cancelOrder, err := reClient.NewCancelOrderService().Symbol(symbol).
	// 	OrderId(co.orderId).OrigClientOrderId(co.origClientOrderId).
	// 	NewClientOrderId(co.newClientOrderId).CancelRestrictions(co.cancelRestrictions).Do(context.Background())
	cancelOrder, err := service.Do(context.Background())
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
	proxyURL string,
) ([]*binance_connector.CancelOrderResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
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
	proxyURL string,
) (*binance_connector.CancelReplaceResponse, error) {
	reClient, clientErr := initClient(thisApiKey, thisSecretKey, proxyURL)
	if clientErr != nil {
		return nil, clientErr
	}
	service := reClient.NewCancelReplaceService().
		Symbol(symbol).Side(side).OrderType(orderType).CancelReplaceMode(cancelReplaceMode)

	if cr.cancelRestrictions != nil {
		service = service.CancelRestrictions(*cr.cancelRestrictions)
	}
	if cr.cancelOrderId != nil {
		service = service.CancelOrderId(*cr.cancelOrderId)
	}
	if cr.cancelNewClientOrderId != nil {
		service = service.CancelNewClientOrderId(*cr.cancelNewClientOrderId)
	}
	if cr.cancelOrigClientOrderId != nil {
		service = service.CancelOrigClientOrderId(*cr.cancelOrigClientOrderId)
	}
	if cr.timeInForce != nil {
		service = service.TimeInForce(*cr.timeInForce)
	}
	if cr.icebergQty != nil {
		service = service.IcebergQty(*cr.icebergQty)
	}
	if cr.quantity != nil {
		service = service.Quantity(*cr.quantity)
	}
	if cr.price != nil {
		service = service.Price(*cr.price)
	}
	if cr.newOrderRespType != nil {
		service = service.NewOrderRespType(*cr.newOrderRespType)
	}
	if cr.newClientOrderId != nil {
		service = service.NewClientOrderId(*cr.newClientOrderId)
	}
	if cr.selfTradePreventionMode != nil {
		service = service.SelfTradePreventionMode(*cr.selfTradePreventionMode)
	}
	if cr.strategyId != nil {
		service = service.StrategyId(*cr.strategyId)
	}
	if cr.strategyType != nil {
		service = service.StrategyType(*cr.strategyType)
	}
	if cr.stopPrice != nil {
		service = service.StopPrice(*cr.stopPrice)
	}
	if cr.timeInForce != nil {
		service = service.TimeInForce(*cr.timeInForce)
	}
	if cr.trailingDelta != nil {
		service = service.TrailingDelta(*cr.trailingDelta)
	}

	// cancelReplace, err := reClient.NewCancelReplaceService().
	// 	Symbol(symbol).Side(side).OrderType(orderType).CancelRestrictions(cr.cancelRestrictions).
	// 	CancelReplaceMode(cancelReplaceMode).CancelOrderId(cr.cancelOrderId).
	// 	CancelNewClientOrderId(cr.cancelNewClientOrderId).CancelOrigClientOrderId(cr.cancelOrigClientOrderId).
	// 	TimeInForce(cr.timeInForce).IcebergQty(cr.icebergQty).Quantity(cr.quantity).
	// 	Price(cr.price).NewOrderRespType(cr.newOrderRespType).NewClientOrderId(cr.newClientOrderId).
	// 	SelfTradePreventionMode(cr.selfTradePreventionMode).StrategyId(cr.strategyId).StrategyType(cr.strategyType).
	// 	StopPrice(cr.stopPrice).TimeInForce(cr.timeInForce).TrailingDelta(cr.trailingDelta).Do(context.Background())
	cancelReplace, err := service.Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(binance_connector.PrettyPrint(cancelReplace))
	return cancelReplace, err
}

func main() {
}
