package tax

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetTax(dbConnPgx utils.PgxIface, taxID *int) (*Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid, 
		name, 
		alternate_name, 
		start_date,
		end_date,
		start_block,
		end_block,
		tax_rate,
		tax_rate_type_id,
		contract_address_str,
		contract_address_id,
		tax_type_id,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM taxes 
	WHERE id = $1`, *taxID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	tax, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Tax])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &tax, nil
}

func RemoveTax(dbConnPgx utils.PgxIface, taxID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveTax DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM taxes WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *taxID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetTaxes(dbConnPgx utils.PgxIface, ids []int) ([]Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	start_block,
	end_block,
	tax_rate,
	tax_rate_type_id,
	contract_address_str,
	contract_address_id,
	tax_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM taxes`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	taxes, err := pgx.CollectRows(results, pgx.RowToStructByName[Tax])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return taxes, nil
}

func GetTaxesByAssetID(dbConnPgx utils.PgxIface, assetID *int) ([]Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	taxes.id,
	taxes.uuid, 
	taxes.name, 
	taxes.alternate_name, 
	taxes.start_date,
	taxes.end_date,
	taxes.start_block,
	taxes.end_block,
	taxes.tax_rate,
	taxes.tax_rate_type_id,
	taxes.contract_address_str,
	taxes.contract_address_id,
	taxes.tax_type_id,
	taxes.description,
	taxes.created_by, 
	taxes.created_at, 
	taxes.updated_by, 
	taxes.updated_at 
	FROM taxes 
	JOIN asset_taxes ON taxes.id = asset_taxes.tax_id
	WHERE asset_taxes.asset_id = $1
	`, *assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	taxes, err := pgx.CollectRows(results, pgx.RowToStructByName[Tax])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return taxes, nil
}

func GetTaxesByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	start_block,
	end_block,
	tax_rate,
	tax_rate_type_id,
	contract_address_str,
	contract_address_id,
	tax_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM taxes
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	taxes, err := pgx.CollectRows(results, pgx.RowToStructByName[Tax])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return taxes, nil
}

func UpdateTax(dbConnPgx utils.PgxIface, tax *Tax) error {
	// if the tax id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if tax.ID == nil || *tax.ID == 0 {
		return errors.New("tax has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateTax DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE taxes SET 
		name=$1,  
		alternate_name=$2, 
		start_date =$3,
		end_date =$4,
		start_block=$5,
		end_block=$6,
		tax_rate =$7,
		tax_rate_type_id=$8,
		contract_address_str=$9,
		contract_address_id=$10,
		tax_type_id=$11,
		description=$12, 
		updated_by=$13, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$14`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		tax.Name,               //1
		tax.AlternateName,      //2
		tax.StartDate,          //3
		tax.EndDate,            //4
		tax.StartBlock,         //5
		tax.EndBlock,           //6
		tax.TaxRate,            //7
		tax.TaxRateTypeID,      //8
		tax.ContractAddressStr, //9
		tax.ContractAddressID,  //10
		tax.TaxTypeID,          //11
		tax.Description,        //12
		tax.UpdatedBy,          //13
		tax.ID,                 //14
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertTax(dbConnPgx utils.PgxIface, tax *Tax) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertTax DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO taxes 
	(
		uuid,
		name,  
		alternate_name, 
		start_date,
		end_date,
		start_block,
		end_block,
		tax_rate,
		tax_rate_type_id,
		contract_address_str,
		contract_address_id,
		tax_type_id,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES (
			uuid_generate_v4(), 
			$1,
			$2, 
			$3, 
			$4, 
			$5, 
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			current_timestamp at time zone 'UTC',
			$13,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		tax.Name,               //1
		tax.AlternateName,      //2
		tax.StartDate,          //3
		tax.EndDate,            //4
		tax.StartBlock,         //5
		tax.EndBlock,           //6
		tax.TaxRate,            //7
		tax.TaxRateTypeID,      //8
		tax.ContractAddressStr, //9
		tax.ContractAddressID,  //10
		tax.TaxTypeID,          //11
		tax.Description,        //12
		tax.CreatedBy,          //13
	).Scan(&insertID)

	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	return int(insertID), nil
}
func InsertTaxes(dbConnPgx utils.PgxIface, taxes []Tax) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range taxes {
		tax := taxes[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(tax.UUID)
		row := []interface{}{
			uuidString,             //1
			tax.Name,               //2
			tax.AlternateName,      //3
			tax.StartDate,          //4
			tax.EndDate,            //5
			tax.StartBlock,         //6
			tax.EndBlock,           //7
			tax.TaxRate,            //8
			tax.TaxRateTypeID,      //9
			tax.ContractAddressStr, //10
			tax.ContractAddressID,  //11
			tax.TaxTypeID,          //12
			tax.Description,        //13
			tax.CreatedBy,          //14
			&now,                   //15
			tax.CreatedBy,          //16
			&now,                   //17
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"taxes"},
		[]string{
			"uuid",                 //1
			"name",                 //2
			"alternate_name",       //3
			"start_date",           //4
			"end_date",             //5
			"start_block",          //6
			"end_block",            //7
			"tax_rate",             //8
			"tax_rate_type_id",     //9
			"contract_address_str", //10
			"contract_address_id",  //11
			"tax_type_id",          //12
			"description",          //13
			"created_by",           //14
			"created_at",           //15
			"updated_by",           //16
			"updated_at",           //17
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertTaxes: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func GetTaxListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
		id,
		uuid, 
		name, 
		alternate_name, 
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM taxes
	`
	if len(_filters) > 0 {
		sql += "WHERE "
		for i, filter := range _filters {
			sql += filter
			if i < len(_filters)-1 {
				sql += " AND "
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
	taxes, err := pgx.CollectRows(results, pgx.RowToStructByName[Tax])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return taxes, nil
}

func GetTotalTaxCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM taxes
	`)
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
