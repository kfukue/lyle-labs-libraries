package tax

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetTax(taxID int) (*Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
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
	WHERE id = $1`, taxID)

	tax := &Tax{}
	err := row.Scan(
		&tax.ID,
		&tax.UUID,
		&tax.Name,
		&tax.AlternateName,
		&tax.StartDate,
		&tax.EndDate,
		&tax.TaxRate,
		&tax.TaxRateTypeID,
		&tax.ContractAddressStr,
		&tax.ContractAddressID,
		&tax.TaxTypeID,
		&tax.Description,
		&tax.CreatedBy,
		&tax.CreatedAt,
		&tax.UpdatedBy,
		&tax.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return tax, nil
}

func GetTopTenTaxes() ([]Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
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
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	taxes := make([]Tax, 0)
	for results.Next() {
		var tax Tax
		results.Scan(
			&tax.ID,
			&tax.UUID,
			&tax.Name,
			&tax.AlternateName,
			&tax.StartDate,
			&tax.EndDate,
			&tax.TaxRate,
			&tax.TaxRateTypeID,
			&tax.ContractAddressStr,
			&tax.ContractAddressID,
			&tax.TaxTypeID,
			&tax.Description,
			&tax.CreatedBy,
			&tax.CreatedAt,
			&tax.UpdatedBy,
			&tax.UpdatedAt,
		)

		taxes = append(taxes, tax)
	}
	return taxes, nil
}

func RemoveTax(taxID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM taxes WHERE id = $1`, taxID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetTaxes(ids []int) ([]Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
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
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	taxes := make([]Tax, 0)
	for results.Next() {
		var tax Tax
		results.Scan(
			&tax.ID,
			&tax.UUID,
			&tax.Name,
			&tax.AlternateName,
			&tax.StartDate,
			&tax.EndDate,
			&tax.TaxRate,
			&tax.TaxRateTypeID,
			&tax.ContractAddressStr,
			&tax.ContractAddressID,
			&tax.TaxTypeID,
			&tax.Description,
			&tax.CreatedBy,
			&tax.CreatedAt,
			&tax.UpdatedBy,
			&tax.UpdatedAt,
		)

		taxes = append(taxes, tax)
	}
	return taxes, nil
}

func GetTaxesByUUIDs(UUIDList []string) ([]Tax, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
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
	taxes := make([]Tax, 0)
	for results.Next() {
		var tax Tax
		results.Scan(
			&tax.ID,
			&tax.UUID,
			&tax.Name,
			&tax.AlternateName,
			&tax.StartDate,
			&tax.EndDate,
			&tax.TaxRate,
			&tax.TaxRateTypeID,
			&tax.ContractAddressStr,
			&tax.ContractAddressID,
			&tax.TaxTypeID,
			&tax.Description,
			&tax.CreatedBy,
			&tax.CreatedAt,
			&tax.UpdatedBy,
			&tax.UpdatedAt,
		)

		taxes = append(taxes, tax)
	}
	return taxes, nil
}

func UpdateTax(tax Tax) error {
	// if the tax id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if tax.ID == nil || *tax.ID == 0 {
		return errors.New("tax has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE taxes SET 
		name=$1,  
		alternate_name=$2, 
		start_date =$3,
		end_date =$4,
		tax_rate =$5,
		tax_rate_type_id=$6,
		contract_address_str=$7,
		contract_address_id=$8,
		tax_type_id=$9,
		description=$10, 
		updated_by=$11, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$12`,
		tax.Name,               //1
		tax.AlternateName,      //2
		tax.StartDate,          //3
		tax.EndDate,            //4
		tax.TaxRate,            //5
		tax.TaxRateTypeID,      //6
		tax.ContractAddressStr, //7
		tax.ContractAddressID,  //8
		tax.TaxTypeID,          //9
		tax.Description,        //10
		tax.UpdatedBy,          //11
		tax.ID)                 //12
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertTax(tax Tax) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO taxes 
	(
		name,  
		uuid,
		alternate_name, 
		start_date,
		end_date,
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
			$1,
			uuid_generate_v4(), 
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
			current_timestamp at time zone 'UTC',
			$11,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		tax.Name,               //1
		tax.AlternateName,      //2
		tax.StartDate,          //3
		tax.EndDate,            //4
		tax.TaxRate,            //5
		tax.TaxRateTypeID,      //6
		tax.ContractAddressStr, //7
		tax.ContractAddressID,  //8
		tax.Description,        //9
		tax.TaxTypeID,          //10
		tax.CreatedBy,          //11
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertTaxes(taxes []Tax) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range taxes {
		tax := taxes[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(tax.UUID)
		row := []interface{}{
			tax.Name,               //1
			uuidString,             //2
			tax.AlternateName,      //3
			tax.StartDate,          //4
			tax.EndDate,            //5
			tax.TaxRate,            //6
			tax.TaxRateTypeID,      //7
			tax.ContractAddressStr, //8
			tax.ContractAddressID,  //9
			tax.TaxTypeID,          //10
			tax.Description,        //11
			tax.CreatedBy,          //12
			&now,                   //13
			tax.CreatedBy,          //14
			&now,                   //15
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"taxes"},
		[]string{
			"name",                 //1
			"uuid",                 //2
			"alternate_name",       //3
			"start_date",           //4
			"end_date",             //5
			"tax_rate",             //6
			"tax_rate_type_id",     //7
			"contract_address_str", //8
			"contract_address_id",  //9
			"tax_type_id",          //10
			"description",          //11
			"created_by",           //12
			"created_at",           //13
			"updated_by",           //14
			"updated_at",           //15
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}

	return nil
}
