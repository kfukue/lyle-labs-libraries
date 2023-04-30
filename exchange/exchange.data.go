package exchange

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetExchange(exchangeID int) (*Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		url,
		start_date,
		end_date,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM exchanges 
	WHERE id = $1`, exchangeID)

	exchange := &Exchange{}
	err := row.Scan(
		&exchange.ID,
		&exchange.UUID,
		&exchange.Name,
		&exchange.AlternateName,
		&exchange.Url,
		&exchange.StartDate,
		&exchange.EndDate,
		&exchange.Description,
		&exchange.CreatedBy,
		&exchange.CreatedAt,
		&exchange.UpdatedBy,
		&exchange.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return exchange, nil
}

func RemoveExchange(exchangeID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM exchanges WHERE id = $1`, exchangeID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetExchanges() ([]Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		url,
		start_date,
		end_date,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM exchanges`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	exchanges := make([]Exchange, 0)
	for results.Next() {
		var exchange Exchange
		results.Scan(
			&exchange.ID,
			&exchange.UUID,
			&exchange.Name,
			&exchange.AlternateName,
			&exchange.Url,
			&exchange.StartDate,
			&exchange.EndDate,
			&exchange.Description,
			&exchange.CreatedBy,
			&exchange.CreatedAt,
			&exchange.UpdatedBy,
			&exchange.UpdatedAt,
		)

		exchanges = append(exchanges, exchange)
	}
	return exchanges, nil
}

func GetExchangeList(ids []int) ([]Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		url,
		start_date,
		end_date,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM exchanges`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := database.DbConn.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	exchanges := make([]Exchange, 0)
	for results.Next() {
		var exchange Exchange
		results.Scan(
			&exchange.ID,
			&exchange.UUID,
			&exchange.Name,
			&exchange.AlternateName,
			&exchange.Url,
			&exchange.StartDate,
			&exchange.EndDate,
			&exchange.Description,
			&exchange.CreatedBy,
			&exchange.CreatedAt,
			&exchange.UpdatedBy,
			&exchange.UpdatedAt,
		)
		exchanges = append(exchanges, exchange)
	}
	return exchanges, nil
}

func GetExchangesByUUIDs(UUIDList []string) ([]Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		url,
		start_date,
		end_date,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM exchanges
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	exchanges := make([]Exchange, 0)
	for results.Next() {
		var exchange Exchange
		results.Scan(
			&exchange.ID,
			&exchange.UUID,
			&exchange.Name,
			&exchange.AlternateName,
			&exchange.Url,
			&exchange.StartDate,
			&exchange.EndDate,
			&exchange.Description,
			&exchange.CreatedBy,
			&exchange.CreatedAt,
			&exchange.UpdatedBy,
			&exchange.UpdatedAt,
		)

		exchanges = append(exchanges, exchange)
	}
	return exchanges, nil
}

func GetStartAndEndDateDiffExchanges(diffInDate int) ([]Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		url,
		start_date,
		end_date,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM exchanges
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	exchanges := make([]Exchange, 0)
	for results.Next() {
		var exchange Exchange
		results.Scan(
			&exchange.ID,
			&exchange.UUID,
			&exchange.Name,
			&exchange.AlternateName,
			&exchange.Url,
			&exchange.StartDate,
			&exchange.EndDate,
			&exchange.Description,
			&exchange.CreatedBy,
			&exchange.CreatedAt,
			&exchange.UpdatedBy,
			&exchange.UpdatedAt,
		)

		exchanges = append(exchanges, exchange)
	}
	return exchanges, nil
}

func UpdateExchange(exchange Exchange) error {
	// if the exchange id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if exchange.ID == nil || *exchange.ID == 0 {
		return errors.New("exchange has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE exchanges SET 
		uuid=$1,
		name=$2,
		alternate_name=$3,
		url=$4,
		start_date=$5,
		end_date=$6,
		description=$7,
		updated_by=$8, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$9`,
		exchange.UUID,          //1
		exchange.Name,          //2
		exchange.AlternateName, //3
		exchange.Url,           //4
		exchange.StartDate,     //5
		exchange.EndDate,       //6
		exchange.Description,   //7
		exchange.CreatedBy,     //8
		exchange.CreatedAt,
		exchange.UpdatedBy, //8
		exchange.UpdatedAt,
		exchange.ID) //9
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertExchange(exchange Exchange) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO exchanges 
	(
		uuid,
		name,
		alternate_name,
		url,
		start_date,
		end_date,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
		) VALUES (
			$1,
			$2, 
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
		RETURNING id`,
		&exchange.UUID,          //1
		&exchange.Name,          //2
		&exchange.AlternateName, //3
		&exchange.Url,           //4
		&exchange.StartDate,     //5
		&exchange.EndDate,       //6
		&exchange.Description,   //7
		&exchange.CreatedBy,     //8
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertExchanges(exchanges []Exchange) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range exchanges {
		exchange := exchanges[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(exchange.UUID)
		row := []interface{}{
			uuidString,             //1
			exchange.Name,          //2
			exchange.AlternateName, //3
			exchange.Url,           //4
			&exchange.StartDate,    //5
			&exchange.EndDate,      //6
			exchange.Description,   //7
			exchange.CreatedBy,     //8
			&now,                   //9
			exchange.CreatedBy,     //10
			&now,                   //11
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"exchanges"},
		[]string{
			"uuid",           //1
			"name",           //2
			"alternate_name", //3
			"url",            //4
			"start_date",     //5
			"end_date",       //6
			"description",    //7
			"created_by",     //8
			"created_at",     //9
			"updated_by",     //10
			"updated_at",     //11
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

// exchange chain methods

func UpdateExchangeChainByUUID(exchangeChain ExchangeChain) error {
	// if the exchange id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if exchangeChain.ExchangeID == nil || *exchangeChain.ExchangeID == 0 || exchangeChain.ChainID == nil || *exchangeChain.ChainID == 0 {
		return errors.New("exchangeChain has invalid IDs")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE exchange_chains SET 
		
		exchange_id=$1,
		chain_id=$2,
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE 
		uuid=$5,`,
		exchangeChain.ExchangeID,  //1
		exchangeChain.ChainID,     //2
		exchangeChain.Description, //3
		exchangeChain.UpdatedBy,   //4
		exchangeChain.UUID)        //5
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertExchangeChain(exchangeChain ExchangeChain) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO exchange_chains 
	(
		uuid,
		exchange_id,
		chain_id,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
		) VALUES (
			$1,
			$2, 
			$3, 
			$4, 
			$5, 
			current_timestamp at time zone 'UTC',
			$5,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		&exchangeChain.UUID,        //1
		&exchangeChain.ExchangeID,  //2
		&exchangeChain.ChainID,     //3
		&exchangeChain.Description, //4
		&exchangeChain.CreatedBy,   //5
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertExchangeChains(exchangeChains []ExchangeChain) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range exchangeChains {
		exchangeChain := exchangeChains[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(exchangeChain.UUID)
		row := []interface{}{
			uuidString,                //1
			exchangeChain.ExchangeID,  //2
			exchangeChain.ChainID,     //3
			exchangeChain.Description, //4
			exchangeChain.CreatedBy,   //5
			&now,                      //6
			exchangeChain.CreatedBy,   //7
			&now,                      //8
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"exchange_chains"},
		[]string{
			"uuid",        //1
			"exchange_id", //2
			"chain_id",    //3
			"description", //4
			"created_by",  //5
			"created_at",  //6
			"updated_by",  //7
			"updated_at",  //8
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
