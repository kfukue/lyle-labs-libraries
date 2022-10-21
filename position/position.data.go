package position

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetPosition(positionID int) (*Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	WHERE id = $1`, positionID)

	position := &Position{}
	err := row.Scan(
		&position.ID,
		&position.UUID,
		&position.Name,
		&position.AlternateName,
		&position.AccountID,
		&position.PortfolioID,
		&position.FrequnecyID,
		&position.StartDate,
		&position.EndDate,
		&position.BaseAssetID,
		&position.QuoteAssetID,
		&position.Quantity,
		&position.CostBasis,
		&position.Profit,
		&position.TotalAmount,
		&position.Description,
		&position.CreatedBy,
		&position.CreatedAt,
		&position.UpdatedBy,
		&position.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return position, nil
}

func GetPositionByDates(startDate, endDate time.Time, frequencyID int, baseAssetID int, quoteAssetID int) (*Position, error) {
	layoutPostgres := utils.LayoutPostgres
	// log.Println(fmt.Sprintf("using start : %s, end : %s", startDate.Format(layoutPostgres), endDate.Format(layoutPostgres)))
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
		`, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), frequencyID, baseAssetID, quoteAssetID)
	position := &Position{}
	err := row.Scan(
		&position.ID,
		&position.UUID,
		&position.Name,
		&position.AlternateName,
		&position.AccountID,
		&position.PortfolioID,
		&position.FrequnecyID,
		&position.StartDate,
		&position.EndDate,
		&position.BaseAssetID,
		&position.QuoteAssetID,
		&position.Quantity,
		&position.CostBasis,
		&position.Profit,
		&position.TotalAmount,
		&position.Description,
		&position.CreatedBy,
		&position.CreatedAt,
		&position.UpdatedBy,
		&position.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return position, nil
}

func GetPositionByDatesAndAccount(startDate, endDate time.Time, frequencyID int, baseAssetID int, quoteAssetID int, accountID *int) (*Position, error) {
	layoutPostgres := utils.LayoutPostgres
	// log.Println(fmt.Sprintf("using start : %s, end : %s", startDate.Format(layoutPostgres), endDate.Format(layoutPostgres)))
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
		`, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), frequencyID, baseAssetID, quoteAssetID, *accountID)
	position := &Position{}
	err := row.Scan(
		&position.ID,
		&position.UUID,
		&position.Name,
		&position.AlternateName,
		&position.AccountID,
		&position.PortfolioID,
		&position.FrequnecyID,
		&position.StartDate,
		&position.EndDate,
		&position.BaseAssetID,
		&position.QuoteAssetID,
		&position.Quantity,
		&position.CostBasis,
		&position.Profit,
		&position.TotalAmount,
		&position.Description,
		&position.CreatedBy,
		&position.CreatedAt,
		&position.UpdatedBy,
		&position.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return position, nil
}

func GetPositionByDatesAccountsForAllTradeableAssets(startDate, endDate time.Time, frequencyID int, quoteAssetID int, accountID *int) ([]Position, error) {
	layoutPostgres := utils.LayoutPostgres
	// log.Println(fmt.Sprintf("using start : %s, end : %s", startDate.Format(layoutPostgres), endDate.Format(layoutPostgres)))
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
		`, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), frequencyID, quoteAssetID, *accountID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	positions := make([]Position, 0)
	for results.Next() {
		var position Position
		results.Scan(
			&position.ID,
			&position.UUID,
			&position.Name,
			&position.AlternateName,
			&position.AccountID,
			&position.PortfolioID,
			&position.FrequnecyID,
			&position.StartDate,
			&position.EndDate,
			&position.BaseAssetID,
			&position.QuoteAssetID,
			&position.Quantity,
			&position.CostBasis,
			&position.Profit,
			&position.TotalAmount,
			&position.Description,
			&position.CreatedBy,
			&position.CreatedAt,
			&position.UpdatedBy,
			&position.UpdatedAt,
		)

		positions = append(positions, position)
	}
	return positions, nil
}

func GetPositionBetweenDatesAndAccountAllCurrentAssets(startDate, endDate time.Time, frequencyID int, quoteAssetID int, accountID *int) ([]Position, error) {
	layoutPostgres := utils.LayoutPostgres
	// log.Println(fmt.Sprintf("using start : %s, end : %s", startDate.Format(layoutPostgres), endDate.Format(layoutPostgres)))
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
		`, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), frequencyID, quoteAssetID, *accountID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	positions := make([]Position, 0)
	for results.Next() {
		var position Position
		results.Scan(
			&position.ID,
			&position.UUID,
			&position.Name,
			&position.AlternateName,
			&position.AccountID,
			&position.PortfolioID,
			&position.FrequnecyID,
			&position.StartDate,
			&position.EndDate,
			&position.BaseAssetID,
			&position.QuoteAssetID,
			&position.Quantity,
			&position.CostBasis,
			&position.Profit,
			&position.TotalAmount,
			&position.Description,
			&position.CreatedBy,
			&position.CreatedAt,
			&position.UpdatedBy,
			&position.UpdatedAt,
		)

		positions = append(positions, position)
	}
	return positions, nil
}

func RemovePosition(positionID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM positions WHERE id = $1`, positionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemovePositionByDateRangeAndAccount(startDate, endDate time.Time, accountID int) error {
	layoutPostgres := utils.LayoutPostgres
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM positions WHERE 
	start_date BETWEEN $1 AND $2
	AND	account_id = $3`, startDate.Format(layoutPostgres), endDate.Format(layoutPostgres), accountID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveAllPositionByAccount(accountID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM positions WHERE 
	account_id = $1`, accountID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetPositions(ids []int) ([]Position, error) {
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
	results, err := database.DbConn.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	positions := make([]Position, 0)
	for results.Next() {
		var position Position
		results.Scan(
			&position.ID,
			&position.UUID,
			&position.Name,
			&position.AlternateName,
			&position.AccountID,
			&position.PortfolioID,
			&position.FrequnecyID,
			&position.StartDate,
			&position.EndDate,
			&position.BaseAssetID,
			&position.QuoteAssetID,
			&position.Quantity,
			&position.CostBasis,
			&position.Profit,
			&position.TotalAmount,
			&position.Description,
			&position.CreatedBy,
			&position.CreatedAt,
			&position.UpdatedBy,
			&position.UpdatedAt,
		)

		positions = append(positions, position)
	}
	return positions, nil
}

func QueryPositions(positionQuery *PositionQuery) ([]Position, error) {
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
	results, err := database.DbConn.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	positions := make([]Position, 0)
	for results.Next() {
		var position Position
		results.Scan(
			&position.ID,
			&position.UUID,
			&position.Name,
			&position.AlternateName,
			&position.AccountID,
			&position.PortfolioID,
			&position.FrequnecyID,
			&position.StartDate,
			&position.EndDate,
			&position.BaseAssetID,
			&position.QuoteAssetID,
			&position.Quantity,
			&position.CostBasis,
			&position.Profit,
			&position.TotalAmount,
			&position.Description,
			&position.CreatedBy,
			&position.CreatedAt,
			&position.UpdatedBy,
			&position.UpdatedAt,
		)

		positions = append(positions, position)
	}
	return positions, nil
}

func UpdatePosition(position Position) error {
	layoutPostgres := utils.LayoutPostgres
	// if the position id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if position.ID == nil || *position.ID == 0 {
		return errors.New("position has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE positions SET 
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
		WHERE id=$16`,
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
		position.ID)                               //16
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertPosition(position Position) (int, error) {
	layoutPostgres := utils.LayoutPostgres
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO positions 
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
		RETURNING id`,
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
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func InsertPositions(positions []Position) error {
	txn, err := database.DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	layoutPostgres := utils.LayoutPostgres
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"positions",
		"name",
		"uuid",
		"alternate_name",
		"account_id",
		"portfolio_id",
		"frequency_id",
		"start_date",
		"end_date",
		"base_asset_id",
		"quote_asset_id",
		"quantity",
		"cost_basis",
		"profit",
		"total_amount",
		"description",
		"created_by",
		"created_at",
		"updated_by",
		"updated_at",
	))
	if err != nil {
		log.Fatal(err)
	}
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	for _, position := range positions {
		_, err = stmt.Exec(
			position.Name, //1
			// `uuid_generate_v4()`,
			uuid.New(),
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
			now.Format(layoutPostgres),
			position.CreatedBy, //15
			now.Format(layoutPostgres),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
