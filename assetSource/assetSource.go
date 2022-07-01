package assetsource

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

// Asset
type AssetSource struct {
	SourceID         *int           `json:"sourceId"`
	AssetID          *int           `json:"assetId"`
	UUID             string         `json:"uuid"`
	Name             string         `json:"name"`
	AlternateName    string         `json:"alternateName"`
	SourceIdentifier string         `json:"sourceIdentifier"`
	Description      string         `json:"description"`
	SourceData       types.JSONText `json:"sourceData"`
	CreatedBy        string         `json:"createdBy"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedBy        string         `json:"updatedBy"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}
