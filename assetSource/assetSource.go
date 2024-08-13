package assetsource

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// Asset
type AssetSource struct {
	SourceID         *int           `json:"sourceId" db:"source_id"`                 //1
	AssetID          *int           `json:"assetId" db:"asset_id"`                   //2
	UUID             string         `json:"uuid" db:"uuid"`                          //3
	Name             string         `json:"name" db:"name"`                          //4
	AlternateName    string         `json:"alternateName" db:"alternate_name"`       //5
	SourceIdentifier string         `json:"sourceIdentifier" db:"source_identifier"` //6
	Description      string         `json:"description" db:"description"`            //7
	SourceData       types.JSONText `json:"sourceData" db:"source_data"`             //8
	CreatedBy        string         `json:"createdBy" db:"created_by"`               //9
	CreatedAt        time.Time      `json:"createdAt" db:"created_at"`               //10
	UpdatedBy        string         `json:"updatedBy" db:"updated_by"`               //11
	UpdatedAt        time.Time      `json:"updatedAt" db:"updated_at"`               //12
}
