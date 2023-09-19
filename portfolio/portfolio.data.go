package portfolio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetPortfolio(portfolioID int) (*Portfolio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	WHERE id = $1`, portfolioID)

	portfolio := &Portfolio{}
	err := row.Scan(
		&portfolio.ID,
		&portfolio.UUID,
		&portfolio.Name,
		&portfolio.AlternateName,
		&portfolio.StartDate,
		&portfolio.EndDate,
		&portfolio.UserEmail,
		&portfolio.Description,
		&portfolio.BaseAssetID,
		&portfolio.PortfolioTypeID,
		&portfolio.ParentID,
		&portfolio.CreatedBy,
		&portfolio.CreatedAt,
		&portfolio.UpdatedBy,
		&portfolio.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return portfolio, nil
}

func GetTopTenPortfolios() ([]Portfolio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	portfolios := make([]Portfolio, 0)
	for results.Next() {
		var portfolio Portfolio
		results.Scan(
			&portfolio.ID,
			&portfolio.UUID,
			&portfolio.Name,
			&portfolio.AlternateName,
			&portfolio.StartDate,
			&portfolio.EndDate,
			&portfolio.UserEmail,
			&portfolio.Description,
			&portfolio.BaseAssetID,
			&portfolio.PortfolioTypeID,
			&portfolio.ParentID,
			&portfolio.CreatedBy,
			&portfolio.CreatedAt,
			&portfolio.UpdatedBy,
			&portfolio.UpdatedAt,
		)

		portfolios = append(portfolios, portfolio)
	}
	return portfolios, nil
}

func RemovePortfolio(portfolioID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM portfolios WHERE id = $1`, portfolioID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetPortfolios(ids []int) ([]Portfolio, error) {
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
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	portfolios := make([]Portfolio, 0)
	for results.Next() {
		var portfolio Portfolio
		results.Scan(
			&portfolio.ID,
			&portfolio.UUID,
			&portfolio.Name,
			&portfolio.AlternateName,
			&portfolio.StartDate,
			&portfolio.EndDate,
			&portfolio.UserEmail,
			&portfolio.Description,
			&portfolio.BaseAssetID,
			&portfolio.PortfolioTypeID,
			&portfolio.ParentID,
			&portfolio.CreatedBy,
			&portfolio.CreatedAt,
			&portfolio.UpdatedBy,
			&portfolio.UpdatedAt,
		)

		portfolios = append(portfolios, portfolio)
	}
	return portfolios, nil
}

func UpdatePortfolio(portfolio Portfolio) error {
	// if the portfolio id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if portfolio.ID == nil || *portfolio.ID == 0 {
		return errors.New("portfolio has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE portfolios SET 
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
		WHERE id=$11`,
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
		portfolio.ID)              //11
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertPortfolio(portfolio Portfolio) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO portfolios 
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
		RETURNING id`,
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
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
