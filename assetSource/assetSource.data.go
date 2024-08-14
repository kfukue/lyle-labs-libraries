package assetsource

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
)

func GetAllAssetSourceBySourceAndAssetType(dbConnPgx utils.PgxIface, sourceID, assetTypeID *int) ([]AssetSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
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
	`, *sourceID, *assetTypeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetSources, err := pgx.CollectRows(results, pgx.RowToStructByName[AssetSource])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetSources, nil
}

func GetAssetSource(dbConnPgx utils.PgxIface, sourceID, assetID *int) (*AssetSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
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
	`, *sourceID, *assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	assetSource, err := pgx.CollectOneRow(row, pgx.RowToStructByName[AssetSource])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &assetSource, nil
}

func GetAssetSourceByTicker(dbConnPgx utils.PgxIface, sourceID *int, sourceIdentifier string) (*AssetSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
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
	`, *sourceID, sourceIdentifier)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	assetSource, err := pgx.CollectOneRow(row, pgx.RowToStructByName[AssetSource])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &assetSource, nil
}

func RemoveAssetSource(dbConnPgx utils.PgxIface, sourceID, assetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveAssetSource DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM asset_sources WHERE source_id = $1 AND asset_id =$2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *sourceID, *assetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetAssetSourceList(dbConnPgx utils.PgxIface, assetIds []int, sourceIds []int) ([]AssetSource, error) {
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
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	defer results.Close()
	assetSources, err := pgx.CollectRows(results, pgx.RowToStructByName[AssetSource])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetSources, nil
}

func UpdateAssetSource(dbConnPgx utils.PgxIface, assetSource *AssetSource) error {
	// if the assetSource id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (assetSource.SourceID == nil || *assetSource.SourceID == 0) || (assetSource.AssetID == nil || *assetSource.AssetID == 0) {
		return errors.New("assetSource has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE asset_sources SET 
		name=$1,  
		alternate_name=$2, 
		source_identifier=$3,
		description=$4,
		source_data= $5,
		updated_by=$6, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE source_id=$7 AND asset_id=$8`

	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		assetSource.Name,             //1
		assetSource.AlternateName,    //2
		assetSource.SourceIdentifier, //3
		assetSource.Description,      //4
		assetSource.SourceData,       //5
		assetSource.UpdatedBy,        //6
		assetSource.SourceID,         //7
		assetSource.AssetID,          //8
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertAssetSource(dbConnPgx utils.PgxIface, assetSource *AssetSource) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertAsset DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var SourceID int
	var AssetID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO asset_sources  
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
		assetSource.SourceID,         //1
		assetSource.AssetID,          //2
		assetSource.Name,             //3
		assetSource.AlternateName,    //4
		assetSource.SourceIdentifier, //5
		assetSource.Description,      //6
		assetSource.SourceData,       //7
		assetSource.CreatedBy,        //8
	).Scan(&SourceID, &AssetID)

	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, -1, err
	}
	return int(SourceID), int(AssetID), nil
}

func InsertAssetSources(dbConnPgx utils.PgxIface, assetSources []AssetSource) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range assetSources {
		assetSource := assetSources[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(assetSource.UUID)
		row := []interface{}{
			assetSource.SourceID,         //1
			assetSource.AssetID,          //2
			uuidString,                   //3
			assetSource.Name,             //4
			assetSource.AlternateName,    //5
			assetSource.SourceIdentifier, //6
			assetSource.Description,      //7
			assetSource.SourceData,       //8
			assetSource.CreatedBy,        //9
			&now,                         //10
			assetSource.CreatedBy,        //11
			&now,                         //12
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"asset_sources"},
		[]string{
			"source_id",         //1
			"asset_id",          //2
			"uuid",              //3
			"name",              //4
			"alternate_name",    //5
			"source_identifier", //6
			"description",       //7
			"source_data",       //8
			"created_by",        //9
			"created_at",        //10
			"updated_by",        //11
			"updated_at",        //12
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertAssetSources: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetAssetSourceListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]AssetSource, error) {
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
	FROM asset_sources 
	`
	if len(_filters) > 0 {
		sql += "WHERE "
		for i, filter := range _filters {
			sql += filter
			if i < len(_filters)-1 {
				sql += " OR "
			}
		}
	}
	if _order != "" && _sort != "" {
		sql += fmt.Sprintf(" ORDER BY %s %s ", _sort, _order)
	}
	if (_start != nil && *_start > 0) && (_end != nil && *_end > 0) {
		pageSize := *_end - *_start
		sql += fmt.Sprintf(" OFFSET %d LIMIT %d ", *_start, pageSize)
	}

	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetSources, err := pgx.CollectRows(results, pgx.RowToStructByName[AssetSource])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetSources, nil
}

func GetTotalAssetSourceCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM asset_sources`)
	totalCount := 0
	err := row.Scan(
		&totalCount,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &totalCount, nil
}
