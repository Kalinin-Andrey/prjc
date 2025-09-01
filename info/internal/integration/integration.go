package integration

import (
	"errors"
	"go.uber.org/zap"
	"info/internal/integration/cmc_api"
	"info/internal/integration/cmc_pro_api"
	"info/internal/integration/oracul_analytics_api"
)

type AppConfig struct {
	NameSpace   string
	Subsystem   string
	Service     string
	Environment string
}

type Integration struct {
	CmcAPI             *cmc_api.CmcApiClient
	CmcProAPI          *cmc_pro_api.CmcApiClient
	OraculAnalyticsAPI *oracul_analytics_api.OraculAnalyticsAPIClient
}

func New(appConfig *AppConfig, cfg *Config, logger *zap.Logger) (*Integration, error) {
	integration := &Integration{}

	if cfg.CmcAPI != nil {
		integration.CmcAPI = cmc_api.New(&cmc_api.AppConfig{
			NameSpace: appConfig.NameSpace,
			Subsystem: appConfig.Subsystem,
			Service:   appConfig.Service,
		}, cfg.CmcAPI, logger)
	}

	if cfg.CmcProAPI != nil {
		integration.CmcProAPI = cmc_pro_api.New(&cmc_pro_api.AppConfig{
			NameSpace: appConfig.NameSpace,
			Subsystem: appConfig.Subsystem,
			Service:   appConfig.Service,
		}, cfg.CmcProAPI, logger)
	}

	if cfg.OraculAnalyticsAPI != nil {
		integration.OraculAnalyticsAPI = oracul_analytics_api.New(&oracul_analytics_api.AppConfig{
			NameSpace: appConfig.NameSpace,
			Subsystem: appConfig.Subsystem,
			Service:   appConfig.Service,
		}, cfg.OraculAnalyticsAPI, logger)
	}

	return integration, nil
}

func (intgr *Integration) Close() error {
	return errors.Join()
}
