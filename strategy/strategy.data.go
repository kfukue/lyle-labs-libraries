package strategy

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

func GetStrategy(dbConnPgx utils.PgxIface, strategyID *int) (*Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid, 
		name, 
		alternate_name, 
		start_date,
		end_date,
		description,
		strategy_type_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM strategies 
	WHERE id = $1`, *strategyID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	strategy, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Strategy])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &strategy, nil
}

func RemoveStrategy(dbConnPgx utils.PgxIface, strategyID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveStrategy DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM strategies WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *strategyID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetStrategies(dbConnPgx utils.PgxIface, ids []int) ([]Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		uuid, 
		name, 
		alternate_name, 
		start_date,
		end_date,
		description,
		strategy_type_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM strategies`
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
	strategies, err := pgx.CollectRows(results, pgx.RowToStructByName[Strategy])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategies, nil
}

func GetStrategiesByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]Strategy, error) {
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
	strategy_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM strategies
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategies, err := pgx.CollectRows(results, pgx.RowToStructByName[Strategy])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategies, nil
}

func GetStartAndEndDateDiffStrategies(dbConnPgx utils.PgxIface, diffInDate *int) ([]Strategy, error) {
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
	strategy_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM strategies
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategies, err := pgx.CollectRows(results, pgx.RowToStructByName[Strategy])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategies, nil
}

func UpdateStrategy(dbConnPgx utils.PgxIface, strategy *Strategy) error {
	// if the strategy id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if strategy.ID == nil || *strategy.ID == 0 {
		return errors.New("strategy has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateStrategy DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE strategies SET 
		name=$1,  
		alternate_name=$2, 
		start_date =$3,
		end_date =$4,
		description=$5, 
		strategy_type_id=$6, 
		updated_by=$7, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$8`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		strategy.Name,           //1
		strategy.AlternateName,  //2
		strategy.StartDate,      //3
		strategy.EndDate,        //4
		strategy.Description,    //5
		strategy.StrategyTypeID, //6
		strategy.UpdatedBy,      //7
		strategy.ID,             //8
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertStrategy(dbConnPgx utils.PgxIface, strategy *Strategy) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertStrategy DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO strategies 
	(
		uuid,
		name,  
		alternate_name, 
		start_date,
		end_date,
		description,
		strategy_type_id,
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
			current_timestamp at time zone 'UTC',
			$7,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		strategy.Name,           //1
		strategy.AlternateName,  //2
		strategy.StartDate,      //3
		strategy.EndDate,        //4
		strategy.Description,    //5
		strategy.StrategyTypeID, //6
		strategy.CreatedBy,      //7
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
func InsertStrategies(dbConnPgx utils.PgxIface, strategies []Strategy) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range strategies {
		strategy := strategies[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(strategy.UUID)
		row := []interface{}{
			uuidString,              //1
			strategy.Name,           //2
			strategy.AlternateName,  //3
			&strategy.StartDate,     //4
			&strategy.EndDate,       //5
			strategy.Description,    //6
			strategy.StrategyTypeID, //7
			strategy.CreatedBy,      //8
			&now,                    //9
			strategy.CreatedBy,      //10
			&now,                    //11
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"strategies"},
		[]string{
			"uuid",             //1
			"name",             //2
			"alternate_name",   //3
			"start_date",       //4
			"end_date",         //5
			"description",      //6
			"strategy_type_id", //7
			"created_by",       //8
			"created_at",       //9
			"updated_by",       //10
			"updated_at",       //11
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertStrategies: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetStrategyListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid, 
		name, 
		alternate_name, 
		start_date,
		end_date,
		description,
		strategy_type_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM strategies 
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
	strategyList, err := pgx.CollectRows(results, pgx.RowToStructByName[Strategy])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyList, nil
}

func GetTotalStrategiesCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM strategies
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
