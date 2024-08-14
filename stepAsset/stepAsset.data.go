package stepasset

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

func GetStepAsset(dbConnPgx utils.PgxIface, stepAssetID *int) (*StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_parameter,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM step_assets 
	WHERE id = $1`, *stepAssetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	stepAsset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[StepAsset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &stepAsset, nil
}

func RemoveStepAsset(dbConnPgx utils.PgxIface, stepAssetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveStepAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM step_assets WHERE id = $1`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *stepAssetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetStepAssets(dbConnPgx utils.PgxIface, ids []int) ([]StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_parameter,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM step_assets`
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
	stepAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[StepAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return stepAssets, nil
}

func GetStepAssetsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_parameter,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM step_assets
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	stepAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[StepAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return stepAssets, nil
}

func GetStartAndEndDateDiffStepAssets(dbConnPgx utils.PgxIface, diffInDate *int) ([]StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_parameter,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM step_assets
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	stepAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[StepAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return stepAssets, nil
}

func UpdateStepAsset(dbConnPgx utils.PgxIface, stepAsset *StepAsset) error {
	// if the stepAsset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if stepAsset.ID == nil || *stepAsset.ID == 0 {
		return errors.New("stepAsset has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateStepAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE step_assets SET 
		step_id=$1,
		asset_id=$2,
		swap_asset_id=$3,
		target_pool_id=$4,
		uuid=$5,
		name=$6,
		alternate_name=$7,
		start_date=$8,
		end_date=$9,
		description=$10,
		action_parameter=$11,
		updated_by=$12, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$13`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		stepAsset.StepID,          //1
		stepAsset.AssetID,         //2
		stepAsset.SwapAssetID,     //3
		stepAsset.TargetPoolID,    //4
		stepAsset.UUID,            //5
		stepAsset.Name,            //6
		stepAsset.AlternateName,   //7
		stepAsset.StartDate,       //8
		stepAsset.EndDate,         //9
		stepAsset.Description,     //10
		stepAsset.ActionParameter, //11
		stepAsset.UpdatedBy,       //12
		stepAsset.ID,              //13
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertStepAsset(dbConnPgx utils.PgxIface, stepAsset *StepAsset) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertStepAsset DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO step_assets 
	(
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_parameter,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
		) VALUES (
			$1,
			$2, 
			$3, 
			$4, 
			uuid_generate_v4(), 
			$5, 
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			current_timestamp at time zone 'UTC',
			$11,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		stepAsset.StepID,       //1
		stepAsset.AssetID,      //2
		stepAsset.SwapAssetID,  //3
		stepAsset.TargetPoolID, //4
		// &stepAsset.UUID,
		stepAsset.Name,            //5
		stepAsset.AlternateName,   //6
		stepAsset.StartDate,       //7
		stepAsset.EndDate,         //8
		stepAsset.Description,     //9
		stepAsset.ActionParameter, //10
		stepAsset.CreatedBy,       //11
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
func InsertStepAssets(dbConnPgx utils.PgxIface, stepAssets []StepAsset) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range stepAssets {
		stepAsset := stepAssets[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(stepAsset.UUID)
		row := []interface{}{
			stepAsset.StepID,          //1
			stepAsset.AssetID,         //2
			stepAsset.SwapAssetID,     //3
			stepAsset.TargetPoolID,    //4
			uuidString,                //5
			stepAsset.Name,            //6
			stepAsset.AlternateName,   //7
			&stepAsset.StartDate,      //8
			&stepAsset.EndDate,        //9
			stepAsset.Description,     //10
			stepAsset.ActionParameter, //11
			stepAsset.CreatedBy,       //12
			&now,                      //13
			stepAsset.CreatedBy,       //14
			&now,                      //15
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"step_assets"},
		[]string{
			"step_id",          //1
			"asset_id",         //2
			"swap_asset_id",    //3
			"target_pool_id",   //4
			"uuid",             //5
			"name",             //6
			"alternate_name",   //7
			"start_date",       //8
			"end_date",         //9
			"description",      //10
			"action_parameter", //11
			"created_by",       //12
			"created_at",       //13
			"updated_by",       //14
			"updated_at",       //15
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertStepAssets: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetStepAssetListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_parameter,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM step_assets 
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
	stepAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[StepAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return stepAssets, nil
}

func GetTotalStepAssetsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM step_assets
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
