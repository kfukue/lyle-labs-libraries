package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	SYSTEM_NAME         = "SYSTEM"
	BIG_INT_EXP         = 18
	COINGECKO_SOURCE_ID = 3
	USD_ID              = 34
	// structured value id
	LIVE_INTERVAL_STRUCTURED_VALUE_ID                        = 58
	DAILY_STRUCTURED_VALUE_ID                                = 5
	FIAT_STRUCTURED_VALUE_ID                                 = 9
	ASSET_PAIR_STRUCTURED_VALUE_ID                           = 12
	LIVE_DATA_MARKET_DATA_TYPE_STRUCTURED_VALUE_ID           = 7
	END_OF_DAY_MARKET_DATA_TYPE_STRUCTURED_VALUE_ID          = 8
	SUCCESS_STRUCTURED_VALUE_ID                              = 52
	RUNNING_STRUCTURED_VALUE_ID                              = 53
	WARNING_STRUCTURED_VALUE_ID                              = 54
	FAILED_STRUCTURED_VALUE_ID                               = 55
	LIVE_JOB_CATEGORY_STRUCTURED_VALUE_ID                    = 56
	EOD_JOB_CATEGORY_STRUCTURED_VALUE_ID                     = 57
	HARVESTING_FREQUENCY_STRUCTURED_VALUE_ID                 = 73
	STRATEGY_MARKET_DATA_IMPORT_JOB_TYPE_STRUCTURED_VALUE_ID = 74
	STRATEGY_SNAPSHOT_MARKET_DATA_TYPE_STRUCTURED_VALUE_ID   = 75
	STRATEGY_EXEUCTION_JOB_TYPE_STRUCTURED_VALUE_ID          = 76
	CALCULATION_ON_CHAIN_ANALYTICS_STRUCTURED_VALUE_ID       = 89
	EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID                     = 84
	CREATE_AUDIT_TYPE_STRUCTURED_VALUE_ID                    = 86
	UPDATE_AUDIT_TYPE_STRUCTURED_VALUE_ID                    = 87
	GETH_RELATED_JOB_IMPORT_JOB_TYPE_STRUCTURED_VALUE_ID     = 88
	IMPOPRT_BALANCE_AND_TRANSFERS_STRUCTURED_VALUE_ID        = 81
	IMPOPRT_SWAPS_STRUCTURED_VALUE_ID                        = 83
	GETH_RELATED_JOB_DELETE_JOB_TYPE_STRUCTURED_VALUE_ID     = 90
	SMART_CONTRACT_TAX_STRUCTURED_VALUE_ID                   = 95
	FIXED_STRUCTURED_VALUE_ID                                = 94
	PERCENTAGE_STRUCTURED_VALUE_ID                           = 93
	// structured value type ids
	JOB_STATUS_STRUCTURED_VALUE_TYPE_ID            = 14
	JOB_CATEGORY_STRUCTURED_VALUE_TYPE_ID          = 15
	STRATEGY_STRUCTURED_VALUE_TYPE_ID              = 16
	FREQUENCY_STRUCTURED_VALUE_TYPE_ID             = 17
	STEP_ACTION_TYPE_STRUCTURED_VALUE_TYPE_ID      = 18
	LIQUITIDY_POOL_TYPE_STRUCTURED_VALUE_TYPE_ID   = 20
	ADDRESS_TYPE_STRUCTURED_VALUE_TYPE_ID          = 22
	CONTRACT_ADDRESS_TYPE_STRUCTURED_VALUE_TYPE_ID = 85
	AUDIT_TYPE_STRUCTURED_VALUE_TYPE_ID            = 23
	RATE_TYPE_STRUCTURED_VALUE_TYPE_ID             = 24
	TAX_TYPE_STRUCTURED_VALUE_TYPE_ID              = 25
)

func GetEnv() string {
	return os.Getenv("APP_ENV")
}

func LookupEnv(k string) string {
	var envValue string
	if GetEnv() == "production" {
		envValue = MustGetenv(k)
	} else {
		envValue = GoDotEnvVariable(k)
	}
	return envValue
}

func MustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.", k)
	}
	return v
}

// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
