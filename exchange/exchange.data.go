package exchange

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

func GetExchange(dbConnPgx utils.PgxIface, exchangeID *int) (*Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		exchange_type_id,
		url,
		start_date,
		end_date,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM exchanges 
	WHERE id = $1`, *exchangeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	exchange, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Exchange])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &exchange, nil
}

func RemoveExchange(dbConnPgx utils.PgxIface, exchangeID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveExchange DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM exchanges WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *exchangeID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
func GetExchangeList(dbConnPgx utils.PgxIface, ids []int) ([]Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		exchange_type_id,
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
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	exchanges, err := pgx.CollectRows(results, pgx.RowToStructByName[Exchange])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return exchanges, nil
}

func GetExchangesByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		exchange_type_id,
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
	exchanges, err := pgx.CollectRows(results, pgx.RowToStructByName[Exchange])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return exchanges, nil
}

func GetStartAndEndDateDiffExchanges(dbConnPgx utils.PgxIface, diffInDate *int) ([]Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		exchange_type_id,
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
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	exchanges, err := pgx.CollectRows(results, pgx.RowToStructByName[Exchange])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return exchanges, nil
}

func UpdateExchange(dbConnPgx utils.PgxIface, exchange *Exchange) error {
	// if the exchange id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if exchange.ID == nil || *exchange.ID == 0 {
		return errors.New("exchange has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE exchanges SET 
		name=$1,
		alternate_name=$2,
		exchange_type_id =$3,
		url=$4,
		start_date=$5,
		end_date=$6,
		description=$7,
		updated_by=$8, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$9`

	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		exchange.Name,           //1
		exchange.AlternateName,  //2
		exchange.ExchangeTypeID, //3
		exchange.Url,            //4
		exchange.StartDate,      //5
		exchange.EndDate,        //6
		exchange.Description,    //7
		exchange.UpdatedBy,      //8
		exchange.ID,             //9
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)

}

func InsertExchange(dbConnPgx utils.PgxIface, exchange *Exchange) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertAsset DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO exchanges 
	(
		uuid,
		name,
		alternate_name,
		exchange_type_id,
		url,
		start_date,
		end_date,
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
			current_timestamp at time zone 'UTC',
			$8,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		exchange.Name,           //1
		exchange.AlternateName,  //2
		exchange.ExchangeTypeID, //3
		exchange.Url,            //4
		exchange.StartDate,      //5
		exchange.EndDate,        //6
		exchange.Description,    //7
		exchange.CreatedBy,      //8
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
func InsertExchanges(dbConnPgx utils.PgxIface, exchanges []Exchange) error {
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
			uuidString,               //1
			exchange.Name,            //2
			exchange.AlternateName,   //3
			*exchange.ExchangeTypeID, //4
			exchange.Url,             //5
			exchange.StartDate,       //6
			exchange.EndDate,         //7
			exchange.Description,     //8
			exchange.CreatedBy,       //9
			&now,                     //10
			exchange.CreatedBy,       //11
			&now,                     //12
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"exchanges"},
		[]string{
			"uuid",             //1
			"name",             //2
			"alternate_name",   //3
			"exchange_type_id", //4
			"url",              //5
			"start_date",       //6
			"end_date",         //7
			"description",      //8
			"created_by",       //9
			"created_at",       //10
			"updated_by",       //11
			"updated_at",       //12
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertExchanges: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// exchange chain methods

func UpdateExchangeChainByUUID(dbConnPgx utils.PgxIface, exchangeChain *ExchangeChain) error {
	// if the exchange id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if exchangeChain.ExchangeID == nil || *exchangeChain.ExchangeID == 0 || exchangeChain.ChainID == nil || *exchangeChain.ChainID == 0 {
		return errors.New("exchangeChain has invalid IDs")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE exchange_chains SET 
		exchange_id=$1,
		chain_id=$2,
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE 
		uuid=$5,`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		exchangeChain.ExchangeID,  //1
		exchangeChain.ChainID,     //2
		exchangeChain.Description, //3
		exchangeChain.UpdatedBy,   //4
		exchangeChain.UUID,        //5
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertExchangeChain(dbConnPgx utils.PgxIface, exchangeChain *ExchangeChain) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertAsset DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO exchange_chains 
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
		exchangeChain.UUID,        //1
		exchangeChain.ExchangeID,  //2
		exchangeChain.ChainID,     //3
		exchangeChain.Description, //4
		exchangeChain.CreatedBy,   //5
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
func InsertExchangeChains(dbConnPgx utils.PgxIface, exchangeChains []ExchangeChain) error {
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
	copyCount, err := dbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("InsertExchangeChains copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetExchangeListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		exchange_type_id,
		url,
		start_date,
		end_date,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM exchanges
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
	exchanges, err := pgx.CollectRows(results, pgx.RowToStructByName[Exchange])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return exchanges, nil
}

func GetTotalExchangeCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM exchanges`)
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
