package assettax

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

func GetAllAssetTaxesByTaxType(dbConnPgx utils.PgxIface, taxTypeID *int) ([]AssetTax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	asset_taxes.tax_id,
	asset_taxes.asset_id,
	asset_taxes.uuid, 
	asset_taxes.name, 
	asset_taxes.alternate_name, 
	asset_taxes.tax_rate_override,
	asset_taxes.tax_rate_type_id,
	asset_taxes.description,
	asset_taxes.created_by, 
	asset_taxes.created_at, 
	asset_taxes.updated_by, 
	asset_taxes.updated_at 
	FROM asset_taxes 
	LEFT JOIN taxes
	ON taxes.tax_id = taxes.id 
	WHERE 
	taxes.tax_type_id = $1
	`, *taxTypeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetTaxes, err := pgx.CollectRows(results, pgx.RowToStructByName[AssetTax])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetTaxes, nil
}

func GetAssetTax(dbConnPgx utils.PgxIface, taxID, assetID *int) (*AssetTax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	tax_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	tax_rate_override,
	tax_rate_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM asset_taxes 
	WHERE tax_id = $1
	AND asset_id = $2
	`, *taxID, *assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetTax, err := pgx.CollectOneRow(results, pgx.RowToStructByName[AssetTax])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &assetTax, nil
}

func RemoveAssetTax(dbConnPgx utils.PgxIface, taxID, assetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveAssetTax DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM asset_taxes WHERE tax_id = $1 AND asset_id =$2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *taxID, *assetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetAssetTaxList(dbConnPgx utils.PgxIface, assetIds []int, taxIds []int) ([]AssetTax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	tax_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	tax_rate_override,
	tax_rate_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM asset_taxes`
	if len(assetIds) > 0 || len(taxIds) > 0 {
		additionalQuery := ` WHERE`
		if len(assetIds) > 0 {
			assetStrIds := utils.SplitToString(assetIds, ",")
			additionalQuery += fmt.Sprintf(`asset_id IN (%s)`, assetStrIds)
		}
		if len(taxIds) > 0 {
			if len(assetIds) > 0 {
				additionalQuery += `AND `
			}
			taxStrIds := utils.SplitToString(taxIds, ",")
			additionalQuery += fmt.Sprintf(`tax_id IN (%s)`, taxStrIds)
		}
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetTaxes, err := pgx.CollectRows(results, pgx.RowToStructByName[AssetTax])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetTaxes, nil
}

func UpdateAssetTax(dbConnPgx utils.PgxIface, assetTax *AssetTax) error {
	// if the assetTax id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (assetTax.TaxID == nil || *assetTax.TaxID == 0) || (assetTax.AssetID == nil || *assetTax.AssetID == 0) {
		return errors.New("assetTax has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE asset_taxes SET 
		name=$1,  
		alternate_name=$2, 
		tax_rate_override=$3,
		tax_rate_type_id=$4,
		description=$5,
		updated_by=$6, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE tax_id=$7 AND asset_id=$8`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		assetTax.Name,            //1
		assetTax.AlternateName,   //2
		assetTax.TaxRateOverride, //3
		assetTax.TaxRateTypeID,   //4
		assetTax.Description,     //5
		assetTax.UpdatedBy,       //6
		assetTax.TaxID,           //7
		assetTax.AssetID,         //8
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertAssetTax(dbConnPgx utils.PgxIface, assetTax *AssetTax) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertAsset DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var TaxID int
	var AssetID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO asset_taxes  
	(
		tax_id,
		asset_id,
		uuid,	 
		name, 
		alternate_name,  
		tax_rate_override,
		tax_rate_type_id,
		description,
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
		RETURNING tax_id, asset_id`,
		assetTax.TaxID,           //1
		assetTax.AssetID,         //2
		assetTax.Name,            //3
		assetTax.AlternateName,   //4
		assetTax.TaxRateOverride, //5
		assetTax.TaxRateTypeID,   //6
		assetTax.Description,     //7
		assetTax.CreatedBy,       //8
	).Scan(&TaxID, &AssetID)
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
	return int(TaxID), int(AssetID), nil
}

func InsertAssetTaxes(dbConnPgx utils.PgxIface, assetTaxes []AssetTax) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for _, assetTax := range assetTaxes {
		uuidString := &pgtype.UUID{}
		uuidString.Set(assetTax.UUID)
		row := []interface{}{
			assetTax.TaxID,           //1
			assetTax.AssetID,         //2
			uuidString,               //3
			assetTax.Name,            //4
			assetTax.AlternateName,   //5
			assetTax.TaxRateOverride, //6
			assetTax.TaxRateTypeID,   //7
			assetTax.Description,     //8
			assetTax.CreatedBy,       //9
			&now,                     //10
			assetTax.CreatedBy,       //11
			&now,                     //12
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"asset_taxes"},
		[]string{
			"tax_id",            //1
			"asset_id",          //2
			"uuid",              //3
			"name",              //4
			"alternate_name",    //5
			"tax_rate_override", //6
			"tax_rate_type_id",  //7
			"description",       //8
			"created_by",        //9
			"created_at",        //10
			"updated_by",        //11
			"updated_at",        //12
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertAssetTaxes: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetAssetTaxListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]AssetTax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		tax_id,
		asset_id,
		uuid, 
		name, 
		alternate_name, 
		tax_rate_override,
		tax_rate_type_id,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM asset_taxes 
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
	assetTaxes, err := pgx.CollectRows(results, pgx.RowToStructByName[AssetTax])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetTaxes, nil
}

func GetTotalAssetTaxCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM asset_taxes`)
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
