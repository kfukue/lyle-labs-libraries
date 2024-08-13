package gethlyleaddresses

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
)

func GetGethAddress(dbConnPgx utils.PgxIface, gethAddressID *int) (*GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	uuid, 
	name,
	alternate_name,
	description,
	address_str,
  	address_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_addresses 
	WHERE id = $1
	`, *gethAddressID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	gethAddress, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethAddress])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethAddress, nil
}

func GetGethAddressByAddressStr(dbConnPgx utils.PgxIface, addressStr string) (*GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	uuid, 
	name,
	alternate_name,
	description,
	address_str,
  	address_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_addresses 
	WHERE address_str = $1
	`, addressStr)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	gethAddress, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethAddress])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethAddress, nil
}

func GetGethAddressList(dbConnPgx utils.PgxIface) ([]GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	address_str,
  	address_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_addresses`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethAddresses, err := pgx.CollectRows(results, pgx.RowToStructByName[GethAddress])
	return gethAddresses, nil
}

func GetGethAddressListByAddressStr(dbConnPgx utils.PgxIface, addressStrList []string) ([]GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	uuid, 
	name,
	alternate_name,
	description,
	address_str,
  	address_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_addresses 
	WHERE address_str = ANY($1)
	`, pq.Array(addressStrList))

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethAddresses, err := pgx.CollectRows(results, pgx.RowToStructByName[GethAddress])
	return gethAddresses, nil
}

func GetGethAddressListByIds(dbConnPgx utils.PgxIface, addressIDs []int) ([]GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	uuid, 
	name,
	alternate_name,
	description,
	address_str,
  	address_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_addresses 
	WHERE id = ANY($1)
	`, pq.Array(addressIDs))

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethAddresses, err := pgx.CollectRows(results, pgx.RowToStructByName[GethAddress])
	return gethAddresses, nil
}

func RemoveGethAddress(dbConnPgx utils.PgxIface, gethAddressID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethAddress DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_addresses WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *gethAddressID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func UpdateGethAddress(dbConnPgx utils.PgxIface, gethAddress *GethAddress) error {
	// if the gethAddress id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethAddress.ID == nil || *gethAddress.ID == 0 {
		return errors.New("gethAddress has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethAddress DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_addresses SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		address_str=$4,
		address_type_id=$5,
		updated_by=$6, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$7 `
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		gethAddress.Name,          //1
		gethAddress.AlternateName, //2
		gethAddress.Description,   //3
		gethAddress.AddressStr,    //4
		gethAddress.AddressTypeID, //5
		gethAddress.UpdatedBy,     //6
		gethAddress.ID,            //7
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethAddress(dbConnPgx utils.PgxIface, gethAddress *GethAddress) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethAddress DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var ID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_addresses  
	(
		uuid, 
		name,
		alternate_name,
		description,
		address_str,
  		address_type_id,
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
			current_timestamp at time zone 'UTC',
			$6,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		gethAddress.Name,          //1
		gethAddress.AlternateName, //2
		gethAddress.Description,   //3
		gethAddress.AddressStr,    //4
		gethAddress.AddressTypeID, //5
		gethAddress.CreatedBy,     //6
	).Scan(&ID)
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
	return int(ID), nil
}

func InsertGethAddressList(dbConnPgx utils.PgxIface, gethAddressList []GethAddress) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethAddressList {
		gethAddress := gethAddressList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(gethAddress.UUID)
		row := []interface{}{
			uuidString,                //1
			gethAddress.Name,          //2
			gethAddress.AlternateName, //3
			gethAddress.Description,   //4
			gethAddress.AddressStr,    //5
			gethAddress.AddressTypeID, //6
			gethAddress.CreatedBy,     //7
			gethAddress.CreatedAt,     //8
			gethAddress.CreatedBy,     //9
			&now,                      //10
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_addresses"},
		[]string{
			"uuid",            //1
			"name",            //2
			"alternate_name",  //3
			"description",     //4
			"address_str",     //5
			"address_type_id", //6
			"created_by",      //7
			"created_at",      //8
			"updated_by",      //9
			"updated_at",      //10
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethAddressList: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetGethAddressListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,  
		uuid, 
		name,
		alternate_name,
		description,
		address_str,
		address_type_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_addresses 
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
	gethAddressList, err := pgx.CollectRows(results, pgx.RowToStructByName[GethAddress])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethAddressList, nil
}

func GetTotalGethAddressCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM geth_addresses`)
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
