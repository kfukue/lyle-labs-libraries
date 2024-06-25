package assetsource

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// Asset
type AssetSource struct {
	SourceID         *int           `json:"sourceId" db:"source_id"`
	AssetID          *int           `json:"assetId" db:"asset_id"`
	UUID             string         `json:"uuid" db:"uuid"`
	Name             string         `json:"name" db:"name"`
	AlternateName    string         `json:"alternateName" db:"alternate_name"`
	SourceIdentifier string         `json:"sourceIdentifier" db:"source_identifier"`
	Description      string         `json:"description" db:"description"`
	SourceData       types.JSONText `json:"sourceData" db:"source_data"`
	CreatedBy        string         `json:"createdBy" db:"created_by"`
	CreatedAt        time.Time      `json:"createdAt" db:"created_at"`
	UpdatedBy        string         `json:"updatedBy" db:"updated_by"`
	UpdatedAt        time.Time      `json:"updatedAt" db:"updated_at"`
}
