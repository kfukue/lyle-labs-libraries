package assettax

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

func GetAllAssetTaxesByTaxType(taxTypeID int) ([]AssetTax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	asset_taxes.tax_id,
	asset_taxes.asset_id,
	asset_taxes.uuid, 
	asset_taxes.name, 
	asset_taxes.alternate_name, 
	asset_taxes.tax_rate_override,
	asset_taxes.description,
	asset_taxes.created_by, 
	asset_taxes.created_at, 
	asset_taxes.updated_by, 
	asset_taxes.updated_at 
	FROM asset_taxes 
	LEFT JOIN taxes
	ON taxes.tax_id = taxes.id 
	WHERE 
	taxes.tax_type_id = $2
	`, taxTypeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetTaxes := make([]AssetTax, 0)
	for results.Next() {
		var assetTax AssetTax
		results.Scan(
			&assetTax.TaxID,
			&assetTax.AssetID,
			&assetTax.UUID,
			&assetTax.Name,
			&assetTax.AlternateName,
			&assetTax.TaxRateOverride,
			&assetTax.Description,
			&assetTax.CreatedBy,
			&assetTax.CreatedAt,
			&assetTax.UpdatedBy,
			&assetTax.UpdatedAt,
		)

		assetTaxes = append(assetTaxes, assetTax)
	}
	return assetTaxes, nil
}
func GetAssetTax(taxID int, assetID int) (*AssetTax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	tax_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	tax_rate_override,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM asset_taxes 
	WHERE tax_id = $1
	AND asset_id = $2
	`, taxID, assetID)

	assetTax := &AssetTax{}
	err := row.Scan(
		&assetTax.TaxID,
		&assetTax.AssetID,
		&assetTax.UUID,
		&assetTax.Name,
		&assetTax.AlternateName,
		&assetTax.TaxRateOverride,
		&assetTax.Description,
		&assetTax.CreatedBy,
		&assetTax.CreatedAt,
		&assetTax.UpdatedBy,
		&assetTax.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetTax, nil
}

func GetAssetTaxByTicker(taxID int, taxIdentifier string) (*AssetTax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	tax_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	tax_rate_override,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM asset_taxes 
	WHERE tax_id = $1
	AND tax_identifier = $2
	`, taxID, taxIdentifier)

	assetTax := &AssetTax{}
	err := row.Scan(
		&assetTax.TaxID,
		&assetTax.AssetID,
		&assetTax.UUID,
		&assetTax.Name,
		&assetTax.AlternateName,
		&assetTax.TaxRateOverride,
		&assetTax.Description,
		&assetTax.CreatedBy,
		&assetTax.CreatedAt,
		&assetTax.UpdatedBy,
		&assetTax.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetTax, nil
}

func RemoveAssetTax(taxID int, assetID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM asset_taxes WHERE 
	tax_id = $1 AND asset_id =$2`, taxID, assetID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetAssetTaxList(assetIds []int, taxIds []int) ([]AssetTax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	tax_id,
	asset_id,
	uuid, 
	name, 
	alternate_name, 
	tax_rate_override,
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
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetTaxes := make([]AssetTax, 0)
	for results.Next() {
		var assetTax AssetTax
		results.Scan(
			&assetTax.TaxID,
			&assetTax.AssetID,
			&assetTax.UUID,
			&assetTax.Name,
			&assetTax.AlternateName,
			&assetTax.TaxRateOverride,
			&assetTax.Description,
			&assetTax.CreatedBy,
			&assetTax.CreatedAt,
			&assetTax.UpdatedBy,
			&assetTax.UpdatedAt,
		)

		assetTaxes = append(assetTaxes, assetTax)
	}
	return assetTaxes, nil
}

func UpdateAssetTax(assetTax AssetTax) error {
	// if the assetTax id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (assetTax.TaxID == nil || *assetTax.TaxID == 0) || (assetTax.AssetID == nil || *assetTax.AssetID == 0) {
		return errors.New("assetTax has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE asset_taxes SET 
		name=$1,  
		alternate_name=$2, 
		tax_rate_override=$3,
		description=$4,
		updated_by=$5, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE tax_id=$6 AND asset_id=$7`,
		assetTax.Name,            //1
		assetTax.AlternateName,   //2
		assetTax.TaxRateOverride, //3
		assetTax.Description,     //4
		assetTax.UpdatedBy,       //5
		assetTax.TaxID,           //6
		assetTax.AssetID,         //7
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertAssetTax(assetTax AssetTax) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var TaxID int
	var AssetID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO asset_taxes  
	(
		tax_id,
		asset_id,
		uuid,	 
		name, 
		alternate_name,  
		tax_rate_override,
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
			current_timestamp at time zone 'UTC',
			$7,
			current_timestamp at time zone 'UTC'
		)
		RETURNING tax_id, asset_id`,
		assetTax.TaxID,           // 1
		assetTax.AssetID,         // 2
		assetTax.Name,            // 3
		assetTax.AlternateName,   //4
		assetTax.TaxRateOverride, //5
		assetTax.Description,     //6
		assetTax.CreatedBy,       //7
	).Scan(&TaxID, &AssetID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(TaxID), int(AssetID), nil
}
