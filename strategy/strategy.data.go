package strategy

import (
	"context"
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

func GetStrategy(strategyID int) (*Strategy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	WHERE id = $1`, strategyID)

	strategy := &Strategy{}
	err := row.Scan(
		&strategy.ID,
		&strategy.UUID,
		&strategy.Name,
		&strategy.AlternateName,
		&strategy.StartDate,
		&strategy.EndDate,
		&strategy.Description,
		&strategy.StrategyTypeID,
		&strategy.CreatedBy,
		&strategy.CreatedAt,
		&strategy.UpdatedBy,
		&strategy.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategy, nil
}

func GetTopTenStrategies() ([]Strategy, error) {
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
	strategy_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM strategies 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategies := make([]Strategy, 0)
	for results.Next() {
		var strategy Strategy
		results.Scan(
			&strategy.ID,
			&strategy.UUID,
			&strategy.Name,
			&strategy.AlternateName,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.Description,
			&strategy.StrategyTypeID,
			&strategy.CreatedBy,
			&strategy.CreatedAt,
			&strategy.UpdatedBy,
			&strategy.UpdatedAt,
		)

		strategies = append(strategies, strategy)
	}
	return strategies, nil
}

func RemoveStrategy(strategyID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM strategies WHERE id = $1`, strategyID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetStrategies(ids []int) ([]Strategy, error) {
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
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategies := make([]Strategy, 0)
	for results.Next() {
		var strategy Strategy
		results.Scan(
			&strategy.ID,
			&strategy.UUID,
			&strategy.Name,
			&strategy.AlternateName,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.Description,
			&strategy.StrategyTypeID,
			&strategy.CreatedBy,
			&strategy.CreatedAt,
			&strategy.UpdatedBy,
			&strategy.UpdatedAt,
		)

		strategies = append(strategies, strategy)
	}
	return strategies, nil
}

func GetStrategiesByUUIDs(UUIDList []string) ([]Strategy, error) {
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
	strategies := make([]Strategy, 0)
	for results.Next() {
		var strategy Strategy
		results.Scan(
			&strategy.ID,
			&strategy.UUID,
			&strategy.Name,
			&strategy.AlternateName,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.Description,
			&strategy.StrategyTypeID,
			&strategy.CreatedBy,
			&strategy.CreatedAt,
			&strategy.UpdatedBy,
			&strategy.UpdatedAt,
		)

		strategies = append(strategies, strategy)
	}
	return strategies, nil
}

func GetStartAndEndDateDiffStrategies(diffInDate int) ([]Strategy, error) {
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
	strategy_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM strategies
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategies := make([]Strategy, 0)
	for results.Next() {
		var strategy Strategy
		results.Scan(
			&strategy.ID,
			&strategy.UUID,
			&strategy.Name,
			&strategy.AlternateName,
			&strategy.StartDate,
			&strategy.EndDate,
			&strategy.Description,
			&strategy.StrategyTypeID,
			&strategy.CreatedBy,
			&strategy.CreatedAt,
			&strategy.UpdatedBy,
			&strategy.UpdatedAt,
		)

		strategies = append(strategies, strategy)
	}
	return strategies, nil
}

func UpdateStrategy(strategy Strategy) error {
	// if the strategy id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if strategy.ID == nil || *strategy.ID == 0 {
		return errors.New("strategy has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE strategies SET 
		name=$1,  
		alternate_name=$2, 
		start_date =$3,
		end_date =$4,
		description=$5, 
		strategy_type_id=$6, 
		updated_by=$7, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$8`,
		strategy.Name,           //1
		strategy.AlternateName,  //2
		strategy.StartDate,      //3
		strategy.EndDate,        //4
		strategy.Description,    //5
		strategy.StrategyTypeID, //6
		strategy.UpdatedBy,      //7
		strategy.ID)             //8
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertStrategy(strategy Strategy) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO strategies 
	(
		name,  
		uuid,
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
			$1,
			uuid_generate_v4(), 
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
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertStrategies(strategies []Strategy) error {
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
			strategy.Name,            //1
			uuidString,               //2
			strategy.AlternateName,   //3
			&strategy.StartDate,      //4
			&strategy.EndDate,        //5
			strategy.Description,     //6
			*strategy.StrategyTypeID, //7
			strategy.CreatedBy,       //8
			&now,                     //9
			strategy.CreatedBy,       //10
			&now,                     //11
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"strategies"},
		[]string{
			"name",             //1
			"uuid",             //2
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
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}

	return nil
}
