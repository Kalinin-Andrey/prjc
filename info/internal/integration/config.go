package integration

import (
	"info/internal/integration/cmc_api"
	"info/internal/integration/cmc_pro_api"
	"info/internal/integration/oracul_analytics_api"
)

type Config struct {
	CmcAPI             *cmc_api.Config
	CmcProAPI          *cmc_pro_api.Config
	OraculAnalyticsAPI *oracul_analytics_api.Config
}

type UsageConfig struct {
	CmcAPI             bool
	CmcProAPI          bool
	OraculAnalyticsAPI bool
}
