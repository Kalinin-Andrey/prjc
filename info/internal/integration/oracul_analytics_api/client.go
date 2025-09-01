package oracul_analytics_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/minipkg/httpclient"
	prometheus_utils "github.com/minipkg/prometheus-utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"info/internal/domain/oracul_analytics"
	"info/internal/pkg/apperror"
	"info/internal/pkg/log_key"
	"strconv"
	"time"
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
}

type OraculAnalyticsAPIClient struct {
	config     *Config
	httpClient httpClient
	logger     *zap.Logger
}

const (
	Name                  = "OraculAnalyticsAPI"
	ContentType           = "application/json; charset=utf-8"
	HeaderParam_RequestId = "X-Request-Id"
	HeaderParam_Cookie    = "Cookie"

	URI_GetHoldersStats string = "/api/holders_stats"
)

var _ oracul_analytics.OraculAnalyticsAPIClient = (*OraculAnalyticsAPIClient)(nil)

func New(appConfig *AppConfig, conf *Config, logger *zap.Logger) *OraculAnalyticsAPIClient {
	client := httpclient.New(conf.Httpconfig, prometheus_utils.NewHttpClientMetrics(appConfig.NameSpace, appConfig.Subsystem, appConfig.Service, conf.Httpconfig.Name).SetCuttingPathOpts(&prometheus_utils.CuttingPathOpts{IsNeedToRemoveQueryInPath: true}))
	return &OraculAnalyticsAPIClient{
		config:     conf,
		httpClient: client,
		logger:     logger,
	}
}

func (c *OraculAnalyticsAPIClient) getDefaultRequestOptions() (requestId string, options []httpclient.RequestOption) {
	requestId = uuid.NewV4().String()
	return requestId, []httpclient.RequestOption{
		httpclient.WithContentType(ContentType),
		httpclient.WithHeader(HeaderParam_RequestId, requestId),
	}
}

func (c *OraculAnalyticsAPIClient) GetHoldersStats(ctx context.Context, currencyID uint, blockchain string, coinAddress string) (*oracul_analytics.ImportData, error) {
	if blockchain == "" || coinAddress == "" {
		return nil, apperror.ErrNotFound
	}

	const funcName = "GetHoldersStats"
	resp := &GetHoldersStatsResponse{}
	requestId, options := c.getDefaultRequestOptions()

	ts := time.Now().UTC()
	uri := URI_GetHoldersStats + "?coin_address=" + coinAddress + "&blockchain=" + blockchain + "&start_at=" + ts.Add(-1*time.Hour*24*365).Format(time.DateOnly) + "&end_at=" + ts.Format(time.DateOnly) + "&total_candles=27"

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

	res, err := resp.ImportData(currencyID, ts)
	if err != nil {
		c.logger.Error("error while convertation result", zap.String(log_key.ApiClient, Name), zap.String(log_key.Func, funcName), zap.Error(err))
		return nil, fmt.Errorf(funcName+" [%w] error while convertation result; requestId: %s; uri: %s; response: %s; error: %w;", apperror.ErrInternal, requestId, uri, string(data), err)
	}

	return res, nil
}
