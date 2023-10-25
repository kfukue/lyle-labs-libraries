package gethlyleaddresses

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
	"github.com/lib/pq"
)

func GetGethAddress(gethAddressID int) (*GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	`, gethAddressID)

	gethAddress := &GethAddress{}
	err := row.Scan(
		&gethAddress.ID,
		&gethAddress.UUID,
		&gethAddress.Name,
		&gethAddress.AlternateName,
		&gethAddress.Description,
		&gethAddress.AddressStr,
		&gethAddress.AddressTypeID,
		&gethAddress.CreatedBy,
		&gethAddress.CreatedAt,
		&gethAddress.UpdatedBy,
		&gethAddress.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethAddress, nil
}

func GetGethAddressByAddressStr(addressStr string) (*GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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

	gethAddress := &GethAddress{}
	err := row.Scan(
		&gethAddress.ID,
		&gethAddress.UUID,
		&gethAddress.Name,
		&gethAddress.AlternateName,
		&gethAddress.Description,
		&gethAddress.AddressStr,
		&gethAddress.AddressTypeID,
		&gethAddress.CreatedBy,
		&gethAddress.CreatedAt,
		&gethAddress.UpdatedBy,
		&gethAddress.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethAddress, nil
}

func GetGethAddressList() ([]GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	gethAddresses := make([]GethAddress, 0)
	for results.Next() {
		var gethAddress GethAddress
		results.Scan(
			&gethAddress.ID,
			&gethAddress.UUID,
			&gethAddress.Name,
			&gethAddress.AlternateName,
			&gethAddress.Description,
			&gethAddress.AddressStr,
			&gethAddress.AddressTypeID,
			&gethAddress.CreatedBy,
			&gethAddress.CreatedAt,
			&gethAddress.UpdatedBy,
			&gethAddress.UpdatedAt,
		)

		gethAddresses = append(gethAddresses, gethAddress)
	}
	return gethAddresses, nil
}

func GetGethAddressListByAddressStr(addressStrList []string) (*GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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

	gethAddress := &GethAddress{}
	err := row.Scan(
		&gethAddress.ID,
		&gethAddress.UUID,
		&gethAddress.Name,
		&gethAddress.AlternateName,
		&gethAddress.Description,
		&gethAddress.AddressStr,
		&gethAddress.AddressTypeID,
		&gethAddress.CreatedBy,
		&gethAddress.CreatedAt,
		&gethAddress.UpdatedBy,
		&gethAddress.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethAddress, nil
}

func GetGethAddressListByIds(addressIDs []int) (*GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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

	gethAddress := &GethAddress{}
	err := row.Scan(
		&gethAddress.ID,
		&gethAddress.UUID,
		&gethAddress.Name,
		&gethAddress.AlternateName,
		&gethAddress.Description,
		&gethAddress.AddressStr,
		&gethAddress.AddressTypeID,
		&gethAddress.CreatedBy,
		&gethAddress.CreatedAt,
		&gethAddress.UpdatedBy,
		&gethAddress.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethAddress, nil
}

func RemoveGethAddress(gethAddressID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_addresses WHERE 
	id = $1`, *gethAddressID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateGethAddress(gethAddress GethAddress) error {
	// if the gethAddress id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethAddress.ID == nil || *gethAddress.ID == 0 {
		return errors.New("gethAddress has invalid ID")
	}
	_, err := database.DbConnPgx.Exec(ctx, `UPDATE geth_addresses SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		address_str=$4,
		address_type_id=$5,
		updated_by=$6, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$7 `,
		gethAddress.Name,          //1
		gethAddress.AlternateName, //2
		gethAddress.Description,   //3
		gethAddress.AddressStr,    //4
		gethAddress.AddressTypeID, //5
		gethAddress.UpdatedBy,     //6
		gethAddress.ID,            //7
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethAddress(gethAddress GethAddress) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var ID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_addresses  
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
		log.Println(err.Error())
		return 0, err
	}
	return int(ID), nil
}

func InsertGethAddressList(gethAddressList []*GethAddress) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range gethAddressList {
		gethAddress := gethAddressList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(gethAddress.UUID)
		row := []interface{}{
			uuidString,                 //2
			gethAddress.Name,           //3
			gethAddress.AlternateName,  //4
			gethAddress.Description,    //5
			gethAddress.AddressStr,     //6
			*gethAddress.AddressTypeID, //7
			gethAddress.CreatedBy,      //8
			&gethAddress.CreatedAt,     //9
			gethAddress.CreatedBy,      //10
			&now,                       //11
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_addresses"},
		[]string{
			"uuid",            //2
			"name",            //3
			"alternate_name",  //4
			"description",     //5
			"address_str",     //6
			"address_type_id", //7
			"created_by",      //8
			"created_at",      //9
			"updated_by",      //10
			"updated_at",      //11
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
