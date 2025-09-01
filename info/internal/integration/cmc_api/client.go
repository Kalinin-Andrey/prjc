package cmc_api

import (
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/minipkg/httpclient"
	prometheus_utils "github.com/minipkg/prometheus-utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"info/internal/domain/concentration"
	"info/internal/domain/currency"
	"info/internal/domain/portfolio_item"
	"info/internal/domain/price_and_cap"
	"info/internal/pkg/apperror"
	"info/internal/pkg/log_key"
	"strconv"
)

type httpClient interface {
	Get(ctx context.Context, path string, opts ...httpclient.RequestOption) ([]byte, int, error)
	Post(ctx context.Context, path string, reqObj interface{}, opts ...httpclient.RequestOption) ([]byte, int, error)
}

type AppConfig struct {
	NameSpace string
	Subsystem string
	Service   string
}

type Config struct {
	Httpconfig httpclient.Config
	Cookie     string
}

type CmcApiClient struct {
	config     *Config
	httpClient httpClient
	logger     *zap.Logger
}

const (
	Name                  = "CmcApiClient"
	ContentType           = "application/json; charset=utf-8"
	HeaderParam_RequestId = "X-Request-Id"
	HeaderParam_Cookie    = "Cookie"
	HeaderParam_Platform  = "platform"
	HeaderVal_Web         = "web"

	ErrorMessage_Success = "SUCCESS"

	ChartRange_1D  = "1D"
	ChartRange_7D  = "7D"
	ChartRange_1M  = "1M"
	ChartRange_1Y  = "1Y"
	ChartRange_All = "All"

	AnalyticsRange_1M  = "month1"
	AnalyticsRange_1Y  = "year1"
	AnalyticsRange_All = "all"

	URI_GetDetailChart      string = "/data-api/v3/cryptocurrency/detail/chart"
	URI_GetAnalytics        string = "/data-api/v3/cryptocurrency/info/get-analytics"
	URI_GetCurrencySimple   string = "/data-api/v3/cryptocurrency/market-pairs/latest"
	URI_GetPortfolioSummary string = "/asset/v3/portfolio/query-summary"
)

var ChartRangeList = []interface{}{
	ChartRange_1D,
	ChartRange_7D,
	ChartRange_1M,
	ChartRange_1Y,
	ChartRange_All,
}
var AnalyticsRangeList = []interface{}{
	AnalyticsRange_1M,
	AnalyticsRange_1Y,
	AnalyticsRange_All,
}
var ChartRange2AnalyticsRangeMap = map[string]string{
	ChartRange_1M:  AnalyticsRange_1M,
	ChartRange_1Y:  AnalyticsRange_1Y,
	ChartRange_All: AnalyticsRange_All,
}

func ChartRangeValidate(s string) error {
	return validation.Validate(s, validation.Required, validation.In(ChartRangeList...))
}

func AnalyticsRangeValidate(s string) error {
	return validation.Validate(s, validation.Required, validation.In(AnalyticsRangeList...))
}

func ChartRange2AnalyticsRange(s string) (string, error) {
	res, ok := ChartRange2AnalyticsRangeMap[s]
	if !ok {
		return "", apperror.ErrNotFound
	}
	return res, nil
}

func New(appConfig *AppConfig, conf *Config, logger *zap.Logger) *CmcApiClient {
	client := httpclient.New(conf.Httpconfig, prometheus_utils.NewHttpClientMetrics(appConfig.NameSpace, appConfig.Subsystem, appConfig.Service, conf.Httpconfig.Name).SetCuttingPathOpts(&prometheus_utils.CuttingPathOpts{IsNeedToRemoveQueryInPath: true}))
	return &CmcApiClient{
		config:     conf,
		httpClient: client,
		logger:     logger,
	}
}

func (c *CmcApiClient) getDefaultRequestOptions() (requestId string, options []httpclient.RequestOption) {
	requestId = uuid.NewV4().String()
	return requestId, []httpclient.RequestOption{
		httpclient.WithContentType(ContentType),
		httpclient.WithHeader(HeaderParam_RequestId, requestId),
		httpclient.WithHeader(HeaderParam_Platform, HeaderVal_Web),
	}
}

func (c *CmcApiClient) getRequestOptionsWithCookie() (requestId string, options []httpclient.RequestOption) {
	requestId, options = c.getDefaultRequestOptions()
	return requestId, append(options, httpclient.WithHeader(HeaderParam_Cookie, c.config.Cookie))
}

func (c *CmcApiClient) GetDetailChart(ctx context.Context, currencyID uint, tRange string) (*price_and_cap.PriceAndCapList, error) {
	if err := ChartRangeValidate(tRange); err != nil {
		return nil, err
	}

	const funcName = "GetDetailChart"
	resp := &DetailChartResponse{}
	requestId, options := c.getDefaultRequestOptions()
	uri := URI_GetDetailChart + "?id=" + strconv.FormatUint(uint64(currencyID), 10) + "&range=" + tRange

	data, code, err := c.httpClient.Get(ctx, uri, options...)
	if err != nil {
		c.logger.Error("httpClient.Get error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(Name+"."+funcName+" [%w] http error: %s; requestId: %s; uri: %s", apperror.ErrInternal, err.Error(), requestId, uri)
	}
	if code != 200 {
		c.logger.Error("httpClient.Get error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err), zap.Int(log_key.Code, code))
		return nil, fmt.Errorf(funcName+" [%w] http response error code: "+strconv.Itoa(code)+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	if err = json.Unmarshal(data, resp); err != nil {
		c.logger.Error("json.Unmarshal error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(funcName+" [%w] json.Unmarshal error: %s; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, err.Error(), requestId, uri, string(data))
	}

	if resp.Status.ErrorCode != "0" || resp.Status.ErrorMessage != ErrorMessage_Success {
		c.logger.Error("response with error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.String(log_key.ErrorCode, resp.Status.ErrorCode), zap.String(log_key.ErrorMessage, resp.Status.ErrorMessage))
		return nil, fmt.Errorf(funcName+" [%w] response with error; code: "+resp.Status.ErrorCode+"; error message: "+resp.Status.ErrorMessage+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	resp.Data.CurrencyID = currencyID
	res, err := resp.Data.PriceAndCapList()
	if err != nil {
		c.logger.Error("error while convertation result", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(funcName+" [%w] error while convertation result; requestId: %s; uri: %s; response: %s; error: %w;", apperror.ErrInternal, requestId, uri, string(data), err)
	}

	return res, nil
}

func (c *CmcApiClient) GetAnalytics(ctx context.Context, currencyID uint, tRange string) (*concentration.ConcentrationList, error) {
	var err error
	if err := ChartRangeValidate(tRange); err != nil {
		return nil, err
	}
	tRange, err = ChartRange2AnalyticsRange(tRange)
	if err != nil {
		return nil, err
	}

	const funcName = "GetAnalytics"
	resp := &GetAnalyticsResponse{}
	requestId, options := c.getDefaultRequestOptions()
	uri := URI_GetAnalytics + "?cryptoId=" + strconv.FormatUint(uint64(currencyID), 10) + "&timeRangeType=" + tRange

	data, code, err := c.httpClient.Get(ctx, uri, options...)
	if err != nil {
		c.logger.Error("httpClient.Get error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(Name+"."+funcName+" [%w] http error: %s; requestId: %s; uri: %s", apperror.ErrInternal, err.Error(), requestId, uri)
	}
	if code != 200 {
		c.logger.Error("httpClient.Get error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err), zap.Int(log_key.Code, code))
		return nil, fmt.Errorf(funcName+" [%w] http response error code: "+strconv.Itoa(code)+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	if err = json.Unmarshal(data, resp); err != nil {
		c.logger.Error("json.Unmarshal error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(funcName+" [%w] json.Unmarshal error: %s; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, err.Error(), requestId, uri, string(data))
	}

	if resp.Status.ErrorCode != "0" || resp.Status.ErrorMessage != ErrorMessage_Success {
		c.logger.Error("response with error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.String(log_key.ErrorCode, resp.Status.ErrorCode), zap.String(log_key.ErrorMessage, resp.Status.ErrorMessage))
		return nil, fmt.Errorf(funcName+" [%w] response with error; code: "+resp.Status.ErrorCode+"; error message: "+resp.Status.ErrorMessage+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	resp.Data.CurrencyID = currencyID
	res, err := resp.Data.ConcentrationList()
	if err != nil {
		c.logger.Error("error while convertation result", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(funcName+" [%w] error while convertation result; requestId: %s; uri: %s; response: %s; error: %w;", apperror.ErrInternal, requestId, uri, string(data), err)
	}

	return res, nil
}

func (c *CmcApiClient) GetCurrency(ctx context.Context, currencySlug string) (*currency.Currency, error) {
	var err error
	const funcName = "GetCurrency"
	resp := &GetCurrencyResponse{}
	requestId, options := c.getDefaultRequestOptions()
	uri := URI_GetCurrencySimple + "?start=1&limit=10&category=spot&slug=" + currencySlug

	data, code, err := c.httpClient.Get(ctx, uri, options...)
	if err != nil {
		c.logger.Error("httpClient.Get error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(Name+"."+funcName+" [%w] http error: %s; requestId: %s; uri: %s", apperror.ErrInternal, err.Error(), requestId, uri)
	}
	if code != 200 {
		c.logger.Error("httpClient.Get error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err), zap.Int(log_key.Code, code))
		return nil, fmt.Errorf(funcName+" [%w] http response error code: "+strconv.Itoa(code)+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	if err = json.Unmarshal(data, resp); err != nil {
		c.logger.Error("json.Unmarshal error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(funcName+" [%w] json.Unmarshal error: %s; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, err.Error(), requestId, uri, string(data))
	}

	if resp.Status.ErrorCode != "0" || resp.Status.ErrorMessage != ErrorMessage_Success {
		c.logger.Error("response with error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.String(log_key.ErrorCode, resp.Status.ErrorCode), zap.String(log_key.ErrorMessage, resp.Status.ErrorMessage))
		return nil, fmt.Errorf(funcName+" [%w] response with error; code: "+resp.Status.ErrorCode+"; error message: "+resp.Status.ErrorMessage+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	resp.Data.Slug = currencySlug
	res := resp.Data.Currency()

	return res, nil
}

func (c *CmcApiClient) getPortfolioSummaryRequest(portfolioSourceId string) *GetPortfolioSummaryRequest {
	return &GetPortfolioSummaryRequest{
		PortfolioSourceId: portfolioSourceId,
		PortfolioType:     "manual",
		CryptoUnit:        2781,
		CurrentPage:       1,
		PageSize:          1000,
	}
}

func (c *CmcApiClient) GetPortfolioSummary(ctx context.Context, portfolioSourceId string) (*portfolio_item.PortfolioItemList, error) {
	var err error
	const funcName = "GetPortfolioSummary"
	resp := &GetPortfolioSummaryResponse{}
	requestId, options := c.getRequestOptionsWithCookie()
	uri := URI_GetPortfolioSummary

	data, code, err := c.httpClient.Post(ctx, uri, c.getPortfolioSummaryRequest(portfolioSourceId), options...)
	if err != nil {
		c.logger.Error("httpClient.Get error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(Name+"."+funcName+" [%w] http error: %s; requestId: %s; uri: %s", apperror.ErrInternal, err.Error(), requestId, uri)
	}
	if code != 200 {
		c.logger.Error("httpClient.Get error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err), zap.Int(log_key.Code, code))
		return nil, fmt.Errorf(funcName+" [%w] http response error code: "+strconv.Itoa(code)+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	if err = json.Unmarshal(data, resp); err != nil {
		c.logger.Error("json.Unmarshal error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(funcName+" [%w] json.Unmarshal error: %s; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, err.Error(), requestId, uri, string(data))
	}

	if resp.Status.ErrorCode != "0" || resp.Status.ErrorMessage != ErrorMessage_Success {
		c.logger.Error("response with error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.String(log_key.ErrorCode, resp.Status.ErrorCode), zap.String(log_key.ErrorMessage, resp.Status.ErrorMessage))
		return nil, fmt.Errorf(funcName+" [%w] response with error; code: "+resp.Status.ErrorCode+"; error message: "+resp.Status.ErrorMessage+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	if resp.Data == nil || len(resp.Data.ManualSummary) == 0 {
		c.logger.Error("response with empty data", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName))
		return nil, fmt.Errorf(funcName+" [%w] response with empty data; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, requestId, uri, string(data))
	}

	return resp.Data.ManualSummary[0].List.SetPortfolioSourceId(portfolioSourceId).PortfolioItemList(), nil
}
