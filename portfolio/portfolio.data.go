package portfolio

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

func GetPortfolio(dbConnPgx utils.PgxIface, portfolioID *int) (*Portfolio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	user_email,
	description,
	base_asset_id,
	portfolio_type_id,
	parent_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM portfolios 
	WHERE id = $1`, *portfolioID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	portfolio, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Portfolio])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &portfolio, nil
}

func RemovePortfolio(dbConnPgx utils.PgxIface, portfolioID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM portfolios WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *portfolioID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetPortfolios(dbConnPgx utils.PgxIface, ids []int) ([]Portfolio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	user_email,
	description,
	base_asset_id,
	portfolio_type_id,
	parent_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM portfolios
	`
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
	portfolios, err := pgx.CollectRows(results, pgx.RowToStructByName[Portfolio])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return portfolios, nil
}

func UpdatePortfolio(dbConnPgx utils.PgxIface, portfolio *Portfolio) error {
	// if the portfolio id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if portfolio.ID == nil || *portfolio.ID == 0 {
		return errors.New("portfolio has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateMarketDataJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE portfolios SET 
		name=$1,  
		alternate_name=$2, 
		start_date =$3,
		end_date =$4,
		user_email=$5,
		description=$6,
		base_asset_id=$7,
		portfolio_type_id=$8,
		parent_id=$9,
		updated_by=$10, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$11`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		portfolio.Name,            //1
		portfolio.AlternateName,   //2
		portfolio.StartDate,       //3
		portfolio.EndDate,         //4
		portfolio.UserEmail,       //5
		portfolio.Description,     //6
		portfolio.BaseAssetID,     //7
		portfolio.PortfolioTypeID, //8
		portfolio.ParentID,        //9
		portfolio.UpdatedBy,       //10
		portfolio.ID,              //11
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertPortfolio(dbConnPgx utils.PgxIface, portfolio *Portfolio) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertMarketDataJob DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var insertID int
	var insertUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO portfolios 
	(
		name,  
		uuid,
		alternate_name, 
		start_date,
		end_date,
		user_email,
		description,
		base_asset_id,
		portfolio_type_id,
		parent_id,
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
			current_timestamp at time zone 'UTC',
			$10,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id, uuid`,
		portfolio.Name,            //1
		portfolio.AlternateName,   //2
		portfolio.StartDate,       //3
		portfolio.EndDate,         //4
		portfolio.UserEmail,       //5
		portfolio.Description,     //6
		portfolio.BaseAssetID,     //7
		portfolio.PortfolioTypeID, //8
		portfolio.ParentID,        //9
		portfolio.CreatedBy,       //10
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

func InsertPortfolios(dbConnPgx utils.PgxIface, portfolios []Portfolio) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range portfolios {
		portfolio := portfolios[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(portfolio.UUID)
		row := []interface{}{
			uuidString,                //1
			portfolio.Name,            //2
			portfolio.AlternateName,   //3
			portfolio.StartDate,       //4
			portfolio.EndDate,         //5
			portfolio.UserEmail,       //6
			portfolio.Description,     //7
			portfolio.BaseAssetID,     //8
			portfolio.PortfolioTypeID, //9
			portfolio.ParentID,        //10
			portfolio.CreatedBy,       //11
			&now,                      //12
			portfolio.CreatedBy,       //13
			&now,                      //14
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"portfolios"},
		[]string{
			"uuid",              //1
			"name",              //2
			"alternate_name",    //3
			"start_date",        //4
			"end_date",          //5
			"user_email",        //6
			"description",       //7
			"base_asset_id",     //8
			"portfolio_type_id", //9
			"parent_id",         //10
			"created_by",        //11
			"created_at",        //12
			"updated_by",        //13
			"updated_at",        //14
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertPortfolios: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetPortfolioListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Portfolio, error) {
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
		user_email,
		description,
		base_asset_id,
		portfolio_type_id,
		parent_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM portfolios 
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
	portfolios, err := pgx.CollectRows(results, pgx.RowToStructByName[Portfolio])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return portfolios, nil
}

func GetTotalPortfoliosCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM portfolios
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
