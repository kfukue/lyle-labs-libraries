package position

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

func GetPosition(dbConnPgx utils.PgxIface, positionID *int) (*Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid, 
		name, 
		alternate_name, 
		account_id,
		portfolio_id,
		frequency_id,
		start_date,
		end_date,
		base_asset_id,
		quote_asset_id,
		quantity,
		cost_basis,
		profit,
		total_amount,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM positions 
	WHERE id = $1`, *positionID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	position, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Position])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &position, nil
}

func GetPositionByDatesAndAccount(dbConnPgx utils.PgxIface, startDate, endDate time.Time, frequencyID, baseAssetID, quoteAssetID, accountID *int) (*Position, error) {
	layoutPostgres := utils.LayoutPostgres
	// log.Println(fmt.Sprintf("using start : %s, end : %s", startDate.Format(layoutPostgres), endDate.Format(layoutPostgres)))
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid, 
		name, 
		alternate_name, 
		account_id,
		portfolio_id,
		frequency_id,
		start_date,
		end_date,
		base_asset_id,
		quote_asset_id,
		quantity,
		cost_basis,
		profit,
		total_amount,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		FROM positions
		WHERE 
		start_date = $1
		AND end_date = $2
		AND frequency_id = $3
		AND base_asset_id = $4
		AND quote_asset_id = $5
		AND account_id = $6
		`, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), *frequencyID, *baseAssetID, *quoteAssetID, *accountID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	position, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Position])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &position, nil
}

func GetPositionByDatesAccountsForAllTradeableAssets(dbConnPgx utils.PgxIface, startDate, endDate time.Time, frequencyID, quoteAssetID, accountID *int) ([]Position, error) {
	layoutPostgres := utils.LayoutPostgres
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		uuid, 
		name, 
		alternate_name, 
		account_id,
		portfolio_id,
		frequency_id,
		start_date,
		end_date,
		base_asset_id,
		quote_asset_id,
		quantity,
		cost_basis,
		profit,
		total_amount,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		FROM positions
		WHERE 
		start_date = $1
		AND end_date = $2
		AND frequency_id = $3
		AND quote_asset_id = $4
		AND account_id = $5
		AND base_asset_id IN (SELECT id FROM public.get_current_assets)
		ORDER BY start_date asc
		`, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), *frequencyID, *quoteAssetID, *accountID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	positions, err := pgx.CollectRows(results, pgx.RowToStructByName[Position])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return positions, nil
}

func GetPositionBetweenDatesAndAccountAllCurrentAssets(dbConnPgx utils.PgxIface, startDate, endDate time.Time, frequencyID, quoteAssetID, accountID *int) ([]Position, error) {
	layoutPostgres := utils.LayoutPostgres
	// log.Println(fmt.Sprintf("using start : %s, end : %s", startDate.Format(layoutPostgres), endDate.Format(layoutPostgres)))
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		uuid, 
		name, 
		alternate_name, 
		account_id,
		portfolio_id,
		frequency_id,
		start_date,
		end_date,
		base_asset_id,
		quote_asset_id,
		quantity,
		cost_basis,
		profit,
		total_amount,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		FROM positions
		WHERE 
		start_date >= $1
		AND end_date <= $2
		AND frequency_id = $3
		AND base_asset_id IN (SELECT id FROM public.get_current_assets)
		AND quote_asset_id = $4
		AND account_id = $5
		ORDER BY start_date asc
		`, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), *frequencyID, *quoteAssetID, *accountID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	positions, err := pgx.CollectRows(results, pgx.RowToStructByName[Position])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return positions, nil
}

func RemovePosition(dbConnPgx utils.PgxIface, positionID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemovePosition DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM positions WHERE id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *positionID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemovePositionByDateRangeAndAccount(dbConnPgx utils.PgxIface, startDate, endDate time.Time, accountID *int) error {
	layoutPostgres := utils.LayoutPostgres
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemovePositionByDateRangeAndAccount DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM positions WHERE 
	start_date BETWEEN $1 AND $2
	AND	account_id = $3`

	if _, err := dbConnPgx.Exec(ctx, sql, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), *accountID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveAllPositionByAccount(dbConnPgx utils.PgxIface, accountID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveAllPositionByAccount DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM positions WHERE account_id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *accountID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetPositions(dbConnPgx utils.PgxIface, ids []int) ([]Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	account_id,
	portfolio_id,
	frequency_id,
	start_date,
	end_date,
	base_asset_id,
	quote_asset_id,
	quantity,
	cost_basis,
	profit,
	total_amount,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM positions`
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

	positions, err := pgx.CollectRows(results, pgx.RowToStructByName[Position])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return positions, nil
}

func QueryPositions(dbConnPgx utils.PgxIface, positionQuery *PositionQuery) ([]Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	account_id,
	portfolio_id,
	frequency_id,
	start_date,
	end_date,
	base_asset_id,
	quote_asset_id,
	quantity,
	cost_basis,
	profit,
	total_amount,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM positions`
	additionalQuery := ""
	needAnd := false
	if positionQuery.AccountID != nil {
		additionalQuery += fmt.Sprintf(` account_id = %d`, *positionQuery.AccountID)
		needAnd = true
	}
	if positionQuery.PortfolioID != nil {
		if needAnd == true {
			additionalQuery += ` AND`
		}
		additionalQuery += fmt.Sprintf(` portfolio_id = %d`, *positionQuery.PortfolioID)
		needAnd = true
	}
	if positionQuery.FrequencyID != nil {
		if needAnd == true {
			additionalQuery += ` AND`
		}
		additionalQuery += fmt.Sprintf(` frequency_id = %d`, *positionQuery.FrequencyID)
		needAnd = true
	}
	if positionQuery.BaseAssetID != nil {
		if needAnd == true {
			additionalQuery += ` AND`
		}
		additionalQuery += fmt.Sprintf(` base_asset_id = %d`, *positionQuery.BaseAssetID)
		needAnd = true
	}
	if positionQuery.StartDate != nil {
		if needAnd == true {
			additionalQuery += ` AND`
		}
		additionalQuery += fmt.Sprintf(` start_date >= '%s'`, positionQuery.StartDate.Format(utils.LayoutPostgres))
		needAnd = true
	}
	if positionQuery.EndDate != nil {
		if needAnd == true {
			additionalQuery += ` AND`
		}
		additionalQuery += fmt.Sprintf(` end_date <= '%s'`, positionQuery.EndDate.Format(utils.LayoutPostgres))
		needAnd = true
	}
	if additionalQuery != "" {
		sql += ` WHERE` + additionalQuery
	}
	sql += ` ORDER BY start_date asc`
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	positions, err := pgx.CollectRows(results, pgx.RowToStructByName[Position])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return positions, nil
}

func UpdatePosition(dbConnPgx utils.PgxIface, position *Position) error {
	layoutPostgres := utils.LayoutPostgres
	// if the position id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if position.ID == nil || *position.ID == 0 {
		return errors.New("position has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdatePosition DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE positions SET 
		name=$1,  
		alternate_name=$2, 
		account_id=$3,
		portfolio_id=$4,
		frequency_id=$5, 
		start_date=$6, 
		end_date=$7, 
		base_asset_id=$8,
		quote_asset_id=$9, 
		quantity=$10, 
		cost_basis=$11, 
		profit=$12, 
		total_amount=$13, 
		description=$14, 
		updated_by=$15, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$16`

	if _, err := dbConnPgx.Exec(ctx, sql,
		position.Name,                             //1
		position.AlternateName,                    //2
		position.AccountID,                        //3
		position.PortfolioID,                      //4
		position.FrequnecyID,                      //5
		position.StartDate.Format(layoutPostgres), //6
		position.EndDate.Format(layoutPostgres),   //7
		position.BaseAssetID,                      //8
		position.QuoteAssetID,                     //9
		position.Quantity,                         //10
		position.CostBasis,                        //11
		position.Profit,                           //12
		position.TotalAmount,                      //13
		position.Description,                      //14
		position.UpdatedBy,                        //15
		position.ID,                               //16
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)

}

func InsertPosition(dbConnPgx utils.PgxIface, position *Position) (int, string, error) {
	layoutPostgres := utils.LayoutPostgres
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertPosition DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var insertID int
	var insertUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO positions 
	(
		name,  
		uuid,
		alternate_name, 
		account_id,
		portfolio_id,
		frequency_id,
		start_date,
		end_date,
		base_asset_id,
		quote_asset_id,
		quantity,
		cost_basis,
		profit,
		total_amount,
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
			$12,
			$13,
			$14,
			$15,
			current_timestamp at time zone 'UTC',
			$15,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id,uuid`,
		position.Name,                             //1
		position.AlternateName,                    //2
		position.AccountID,                        //3
		position.PortfolioID,                      //4
		position.FrequnecyID,                      //5
		position.StartDate.Format(layoutPostgres), //6
		position.EndDate.Format(layoutPostgres),   //7
		position.BaseAssetID,                      //8
		position.QuoteAssetID,                     //9
		position.Quantity,                         //10
		position.CostBasis,                        //11
		position.Profit,                           //12
		position.TotalAmount,                      //13
		position.Description,                      //14
		position.CreatedBy,                        //15
	).Scan(&insertID, &insertUUID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	return int(insertID), insertUUID, nil
}

func InsertPositions(dbConnPgx utils.PgxIface, positions []Position) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range positions {
		position := positions[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(position.UUID)
		row := []interface{}{
			uuidString,             //1
			position.Name,          //2
			position.AlternateName, //3
			position.AccountID,     //4
			position.PortfolioID,   //5
			position.FrequnecyID,   //6
			position.StartDate,     //7
			position.EndDate,       //8
			position.BaseAssetID,   //9
			position.QuoteAssetID,  //10
			position.Quantity,      //11
			position.CostBasis,     //12
			position.Profit,        //13
			position.TotalAmount,   //14
			position.Description,   //15
			position.CreatedBy,     //16
			&now,                   //17
			position.CreatedBy,     //18
			&now,                   //19
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"positions"},
		[]string{
			"uuid",           //1
			"name",           //2
			"alternate_name", //3
			"account_id",     //4
			"portfolio_id",   //5
			"frequency_id",   //6
			"start_date",     //7
			"end_date",       //8
			"base_asset_id",  //9
			"quote_asset_id", //10
			"quantity",       //11
			"cost_basis",     //12
			"profit",         //13
			"total_amount",   //14
			"description",    //15
			"created_by",     //16
			"created_at",     //17
			"updated_by",     //18
			"updated_at",     //19
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertPositions: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetPositionListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid, 
		name, 
		alternate_name, 
		account_id,
		portfolio_id,
		frequency_id,
		start_date,
		end_date,
		base_asset_id,
		quote_asset_id,
		quantity,
		cost_basis,
		profit,
		total_amount,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM positions 
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

	positions, err := pgx.CollectRows(results, pgx.RowToStructByName[Position])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return positions, nil
}

func GetTotalPositionsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM positions
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
