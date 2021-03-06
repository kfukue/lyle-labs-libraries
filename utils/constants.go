package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	SYSTEM_NAME                                              = "SYSTEM"
	LIVE_INTERVAL_STRUCTURED_VALUE_ID                        = 58
	DAILY_STRUCTURED_VALUE_ID                                = 5
	FIAT_STRUCTURED_VALUE_ID                                 = 9
	ASSET_PAIR_STRUCTURED_VALUE_ID                           = 12
	USD_ID                                                   = 34
	LIVE_DATA_MARKET_DATA_TYPE_STRUCTURED_VALUE_ID           = 7
	END_OF_DAY_MARKET_DATA_TYPE_STRUCTURED_VALUE_ID          = 8
	COINGECKO_SOURCE_ID                                      = 3
	JOB_STATUS_STRUCTURED_VALUE_TYPE_ID                      = 14
	SUCCESS_STRUCTURED_VALUE_ID                              = 52
	RUNNING_STRUCTURED_VALUE_ID                              = 53
	WARNING_STRUCTURED_VALUE_ID                              = 54
	FAILED_STRUCTURED_VALUE_ID                               = 55
	JOB_CATEGORY_STRUCTURED_VALUE_TYPE_ID                    = 15
	LIVE_JOB_CATEGORY_STRUCTURED_VALUE_ID                    = 56
	EOD_JOB_CATEGORY_STRUCTURED_VALUE_ID                     = 57
	BIG_INT_EXP                                              = 18
	STRATEGY_STRUCTURED_VALUE_TYPE_ID                        = 16
	FREQUENCY_STRUCTURED_VALUE_TYPE_ID                       = 17
	STEP_ACTION_TYPE_STRUCTURED_VALUE_TYPE_ID                = 18
	HARVESTING_FREQUENCY_STRUCTURED_VALUE_ID                 = 73
	STRATEGY_MARKET_DATA_IMPORT_JOB_TYPE_STRUCTURED_VALUE_ID = 74
	STRATEGY_SNAPSHOT_MARKET_DATA_TYPE_STRUCTURED_VALUE_ID   = 75
	STRATEGY_EXEUCTION_JOB_TYPE_STRUCTURED_VALUE_ID          = 76
)

func GetEnv() string {
	return os.Getenv("APP_ENV")
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
