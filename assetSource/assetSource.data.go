package assetsource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetAllAssetSourceBySourceAndAssetType(sourceID int, assetTypeID int) ([]AssetSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	asset_sources.source_id,
	asset_sources.asset_id,
	asset_sources.uuid, 
	asset_sources.name, 
	asset_sources.alternate_name, 
	asset_sources.source_identifier,
	asset_sources.description,
	asset_sources.source_data,
	asset_sources.created_by, 
	asset_sources.created_at, 
	asset_sources.updated_by, 
	asset_sources.updated_at 
	FROM asset_sources 
	LEFT JOIN assets
	ON asset_sources.asset_id = assets.id 
	WHERE 
	asset_sources.source_id = $1
	AND assets.asset_type_id = $2
	`, sourceID, assetTypeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetSources := make([]AssetSource, 0)
	for results.Next() {
		var assetSource AssetSource
		results.Scan(
			&assetSource.SourceID,
			&assetSource.AssetID,
			&assetSource.UUID,
			&assetSource.Name,
			&assetSource.AlternateName,
			&assetSource.SourceIdentifier,
			&assetSource.Description,
			&assetSource.SourceData,
			&assetSource.CreatedBy,
			&assetSource.CreatedAt,
			&assetSource.UpdatedBy,
			&assetSource.UpdatedAt,
		)

		assetSources = append(assetSources, assetSource)
	}
	return assetSources, nil
}
func GetAssetSource(sourceID int, assetID int) (*AssetSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	source_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	source_identifier,
	description,
	source_data,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM asset_sources 
	WHERE source_id = $1
	AND asset_id = $2
	`, sourceID, assetID)

	assetSource := &AssetSource{}
	err := row.Scan(
		&assetSource.SourceID,
		&assetSource.AssetID,
		&assetSource.UUID,
		&assetSource.Name,
		&assetSource.AlternateName,
		&assetSource.SourceIdentifier,
		&assetSource.Description,
		&assetSource.SourceData,
		&assetSource.CreatedBy,
		&assetSource.CreatedAt,
		&assetSource.UpdatedBy,
		&assetSource.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetSource, nil
}

func GetAssetSourceByTicker(sourceID int, sourceIdentifier string) (*AssetSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	source_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	source_identifier,
	description,
	source_data,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM asset_sources 
	WHERE source_id = $1
	AND source_identifier = $2
	`, sourceID, sourceIdentifier)

	assetSource := &AssetSource{}
	err := row.Scan(
		&assetSource.SourceID,
		&assetSource.AssetID,
		&assetSource.UUID,
		&assetSource.Name,
		&assetSource.AlternateName,
		&assetSource.SourceIdentifier,
		&assetSource.Description,
		&assetSource.SourceData,
		&assetSource.CreatedBy,
		&assetSource.CreatedAt,
		&assetSource.UpdatedBy,
		&assetSource.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetSource, nil
}

func GetTopTenAssetSources() ([]AssetSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	source_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	source_identifier,
	description,
	source_data,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM asset_sources 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetSources := make([]AssetSource, 0)
	for results.Next() {
		var assetSource AssetSource
		results.Scan(
			&assetSource.SourceID,
			&assetSource.AssetID,
			&assetSource.UUID,
			&assetSource.Name,
			&assetSource.AlternateName,
			&assetSource.SourceIdentifier,
			&assetSource.Description,
			&assetSource.SourceData,
			&assetSource.CreatedBy,
			&assetSource.CreatedAt,
			&assetSource.UpdatedBy,
			&assetSource.UpdatedAt,
		)

		assetSources = append(assetSources, assetSource)
	}
	return assetSources, nil
}

func RemoveAssetSource(sourceID int, assetID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM asset_sources WHERE 
	source_id = $1 AND asset_id =$2`, sourceID, assetID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetAssetSourceList(assetIds []int, sourceIds []int) ([]AssetSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	source_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	source_identifier,
	description,
	source_data,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM asset_sources`
	if len(assetIds) > 0 || len(sourceIds) > 0 {
		additionalQuery := ` WHERE`
		if len(assetIds) > 0 {
			assetStrIds := utils.SplitToString(assetIds, ",")
			additionalQuery += fmt.Sprintf(`asset_id IN (%s)`, assetStrIds)
		}
		if len(sourceIds) > 0 {
			if len(assetIds) > 0 {
				additionalQuery += `AND `
			}
			sourceStrIds := utils.SplitToString(sourceIds, ",")
			additionalQuery += fmt.Sprintf(`source_id IN (%s)`, sourceStrIds)
		}
		sql += additionalQuery
	}
	results, err := database.DbConn.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetSources := make([]AssetSource, 0)
	for results.Next() {
		var assetSource AssetSource
		results.Scan(
			&assetSource.SourceID,
			&assetSource.AssetID,
			&assetSource.UUID,
			&assetSource.Name,
			&assetSource.AlternateName,
			&assetSource.SourceIdentifier,
			&assetSource.Description,
			&assetSource.SourceData,
			&assetSource.CreatedBy,
			&assetSource.CreatedAt,
			&assetSource.UpdatedBy,
			&assetSource.UpdatedAt,
		)

		assetSources = append(assetSources, assetSource)
	}
	return assetSources, nil
}

func UpdateAssetSource(assetSource AssetSource) error {
	// if the assetSource id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (assetSource.SourceID == nil || *assetSource.SourceID == 0) || (assetSource.AssetID == nil || *assetSource.AssetID == 0) {
		return errors.New("assetSource has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE asset_sources SET 
		name=$1,  
		alternate_name=$2, 
		source_identifier=$3,
		description=$4,
		source_data= $5,
		updated_by=$6, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE source_id=$7 AND asset_id=$8`,
		assetSource.Name,
		assetSource.AlternateName,
		assetSource.SourceIdentifier,
		assetSource.Description,
		assetSource.SourceData,
		assetSource.UpdatedBy,
		assetSource.SourceID,
		assetSource.AssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertAssetSource(assetSource AssetSource) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var SourceID int
	var AssetID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO asset_sources  
	(
		source_id,
		asset_id,
		uuid,	 
		name, 
		alternate_name,  
		source_identifier,
		description,
		source_data,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,
			$2,
			uuid_generate_v4(),
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			current_timestamp at time zone 'UTC',
			$8,
			current_timestamp at time zone 'UTC'
		)
		RETURNING source_id, asset_id`,
		assetSource.SourceID,         // 1
		assetSource.AssetID,          // 2
		assetSource.Name,             // 3
		assetSource.AlternateName,    //4
		assetSource.SourceIdentifier, //5
		assetSource.Description,      //6
		assetSource.SourceData,       //7
		assetSource.CreatedBy,        //8
	).Scan(&SourceID, &AssetID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(SourceID), int(AssetID), nil
}
