package gethlylejobs

import (
	"fmt"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/utils"
)

type GethProcessVlogJob struct {
	ID               *int `json:"id"`
	GethProcessJobID *int `json:"gethProcessJobId"`

	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	AlternateName string    `json:"alternateName"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	Description   string    `json:"description"`
	StatusID      *int      `json:"statusId"`
	JobCategoryID *int      `json:"jobCategoryId"`

	AssetID     *int     `json:"assetId"`
	ChainID     *int     `json:"chainId"`
	TxnHash     string   `json:"txnHash"`
	AddressID   *int     `json:"addressId"`
	BlockNumber *uint64  `json:"blockNumber"`
	IndexNumber *uint    `json:"indexNumber"`
	TopicsStr   []string `json:"topicsStr"`

	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedBy string    `json:"updatedBy"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func UpdateFailedGethProcessVlogJob(gethProcessVlogJob *GethProcessVlogJob, msg string, doUpdate bool) error {
	failedStatus := utils.FAILED_STRUCTURED_VALUE_ID
	gethProcessVlogJob.StatusID = &failedStatus
	gethProcessVlogJob.Description = fmt.Sprintf("%s \n %s", gethProcessVlogJob.Description, msg)
	if doUpdate {
		err := UpdateGethProcessVlogJob(gethProcessVlogJob)
		if err != nil {
			msg = fmt.Sprintf("Failed in UpdateGethProcessVlogJob, err %v", err)
			log.Fatal(msg)
			return err
		}
	}
	return nil
}
