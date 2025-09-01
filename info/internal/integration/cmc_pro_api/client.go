package cmc_pro_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/minipkg/httpclient"
	prometheus_utils "github.com/minipkg/prometheus-utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"info/internal/domain/currency"
	"info/internal/pkg/apperror"
	"info/internal/pkg/fasthttp_tools"
	"info/internal/pkg/log_key"
	"strconv"
	"strings"
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
	Token      string
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
	HeaderParam_APIKey    = "X-CMC_PRO_API_KEY"

	ErrorMessage_Success = "SUCCESS"

	URI_GetCurrencies string = "/v2/cryptocurrency/quotes/latest"
)

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
	}
}

func (c *CmcApiClient) getRequestOptions() (requestId string, options []httpclient.RequestOption) {
	requestId, options = c.getDefaultRequestOptions()
	return requestId, append(options, httpclient.WithHeader(HeaderParam_APIKey, c.config.Token))
}

func (c *CmcApiClient) GetCurrenciesByIDs(ctx context.Context, currencyIDs *[]uint) (currencyMap currency.CurrencyMap, err error) {
	if currencyIDs == nil || len(*currencyIDs) == 0 {
		return nil, nil
	}

	return c.getCurrencies(ctx, "id="+fasthttp_tools.Uints2Str(currencyIDs, nil))
}

func (c *CmcApiClient) GetCurrenciesBySlugs(ctx context.Context, slugs *[]string) (currencyMap currency.CurrencyMap, err error) {
	if slugs == nil || len(*slugs) == 0 {
		return nil, nil
	}

	return c.getCurrencies(ctx, "slug="+strings.Join(*slugs, ","))
}

func (c *CmcApiClient) getCurrencies(ctx context.Context, params string) (currencyMap currency.CurrencyMap, err error) {
	const funcName = "getCurrencies"
	resp := &CurrencyQuotesResponse{}
	requestId, options := c.getRequestOptions()

	uri := URI_GetCurrencies + "?" + params

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

	if resp.Status.ErrorCode != 0 || (resp.Status.ErrorMessage != ErrorMessage_Success && resp.Status.ErrorMessage != "") {
		c.logger.Error("response with error", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Int(log_key.ErrorCode, resp.Status.ErrorCode), zap.String(log_key.ErrorMessage, resp.Status.ErrorMessage))
		return nil, fmt.Errorf(funcName+" [%w] response with error; code: %d; error message: "+resp.Status.ErrorMessage+"; requestId: %s; uri: %s; response: %s", apperror.ErrInternal, resp.Status.ErrorCode, requestId, uri, string(data))
	}

	if currencyMap, err = resp.Data.CurrencyMap(); err != nil {
		c.logger.Error("error while convertation result", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(funcName+" [%w] error while convertation result; requestId: %s; uri: %s; response: %s; error: %w;", apperror.ErrInternal, requestId, uri, string(data), err)
	}

	return currencyMap, nil
}
