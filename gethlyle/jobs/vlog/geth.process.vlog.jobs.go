package gethlylevlogjobs

import (
	"fmt"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/utils"
)

type GethProcessVlogJob struct {
	ID               *int      `json:"id" db:"id"`                                //1
	GethProcessJobID *int      `json:"gethProcessJobId" db:"geth_process_job_id"` //2
	UUID             string    `json:"uuid" db:"uuid"`                            //3
	Name             string    `json:"name" db:"name"`                            //4
	AlternateName    string    `json:"alternateName" db:"alternate_name"`         //5
	StartDate        time.Time `json:"startDate" db:"start_date"`                 //6
	EndDate          time.Time `json:"endDate" db:"end_date"`                     //7
	Description      string    `json:"description" db:"description"`              //8
	StatusID         *int      `json:"statusId" db:"status_id"`                   //9
	JobCategoryID    *int      `json:"jobCategoryId" db:"job_category_id"`        //10
	AssetID          *int      `json:"assetId" db:"asset_id"`                     //11
	ChainID          *int      `json:"chainId" db:"chain_id"`                     //12
	TxnHash          string    `json:"txnHash" db:"txn_hash"`                     //13
	AddressID        *int      `json:"addressId" db:"address_id"`                 //14
	BlockNumber      *uint64   `json:"blockNumber" db:"block_number"`             //15
	IndexNumber      *uint     `json:"indexNumber" db:"index_number"`             //16
	TopicsStrArray   []string  `json:"topicsStr" db:"topics_str"`                 //17
	CreatedBy        string    `json:"createdBy" db:"created_by"`                 //18
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`                 //19
	UpdatedBy        string    `json:"updatedBy" db:"updated_by"`                 //20
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`                 //21

}

func UpdateFailedGethProcessVlogJob(dbConnPgx utils.PgxIface, gethProcessVlogJob *GethProcessVlogJob, msg string, doUpdate bool) error {
	failedStatus := utils.FAILED_STRUCTURED_VALUE_ID
	gethProcessVlogJob.StatusID = &failedStatus
	gethProcessVlogJob.Description = fmt.Sprintf("%s \n %s", gethProcessVlogJob.Description, msg)
	if doUpdate {
		err := UpdateGethProcessVlogJob(dbConnPgx, gethProcessVlogJob)
		if err != nil {
			msg = fmt.Sprintf("Failed in UpdateGethProcessVlogJob, err %v", err)
			log.Println(msg)
			return err
		}
	}
	return nil
}
