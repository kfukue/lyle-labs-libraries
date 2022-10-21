package stepasset

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

func GetStepAsset(stepAssetID int) (*StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		name,
		uuid,
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
	WHERE id = $1`, stepAssetID)

	stepAsset := &StepAsset{}
	err := row.Scan(
		&stepAsset.ID,
		&stepAsset.StepID,
		&stepAsset.AssetID,
		&stepAsset.SwapAssetID,
		&stepAsset.TargetPoolID,
		&stepAsset.Name,
		&stepAsset.UUID,
		&stepAsset.AlternateName,
		&stepAsset.StartDate,
		&stepAsset.EndDate,
		&stepAsset.Description,
		&stepAsset.ActionParameter,
		&stepAsset.CreatedBy,
		&stepAsset.CreatedAt,
		&stepAsset.UpdatedBy,
		&stepAsset.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return stepAsset, nil
}

func GetTopTenStrategies() ([]StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		name,
		uuid,
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
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	stepAssets := make([]StepAsset, 0)
	for results.Next() {
		var stepAsset StepAsset
		results.Scan(
			&stepAsset.ID,
			&stepAsset.StepID,
			&stepAsset.AssetID,
			&stepAsset.SwapAssetID,
			&stepAsset.TargetPoolID,
			&stepAsset.Name,
			&stepAsset.UUID,
			&stepAsset.AlternateName,
			&stepAsset.StartDate,
			&stepAsset.EndDate,
			&stepAsset.Description,
			&stepAsset.ActionParameter,
			&stepAsset.CreatedBy,
			&stepAsset.CreatedAt,
			&stepAsset.UpdatedBy,
			&stepAsset.UpdatedAt,
		)

		stepAssets = append(stepAssets, stepAsset)
	}
	return stepAssets, nil
}

func RemoveStepAsset(stepAssetID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM step_assets WHERE id = $1`, stepAssetID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetStepAssets(ids []int) ([]StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		name,
		uuid,
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
		additionalQuery := fmt.Sprintf(`WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := database.DbConn.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	stepAssets := make([]StepAsset, 0)
	for results.Next() {
		var stepAsset StepAsset
		results.Scan(
			&stepAsset.ID,
			&stepAsset.StepID,
			&stepAsset.AssetID,
			&stepAsset.SwapAssetID,
			&stepAsset.TargetPoolID,
			&stepAsset.Name,
			&stepAsset.UUID,
			&stepAsset.AlternateName,
			&stepAsset.StartDate,
			&stepAsset.EndDate,
			&stepAsset.Description,
			&stepAsset.ActionParameter,
			&stepAsset.CreatedBy,
			&stepAsset.CreatedAt,
			&stepAsset.UpdatedBy,
			&stepAsset.UpdatedAt,
		)

		stepAssets = append(stepAssets, stepAsset)
	}
	return stepAssets, nil
}

func GetStepAssetsByUUIDs(UUIDList []string) ([]StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		name,
		uuid,
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
	stepAssets := make([]StepAsset, 0)
	for results.Next() {
		var stepAsset StepAsset
		results.Scan(
			&stepAsset.ID,
			&stepAsset.StepID,
			&stepAsset.AssetID,
			&stepAsset.SwapAssetID,
			&stepAsset.TargetPoolID,
			&stepAsset.Name,
			&stepAsset.UUID,
			&stepAsset.AlternateName,
			&stepAsset.StartDate,
			&stepAsset.EndDate,
			&stepAsset.Description,
			&stepAsset.ActionParameter,
			&stepAsset.CreatedBy,
			&stepAsset.CreatedAt,
			&stepAsset.UpdatedBy,
			&stepAsset.UpdatedAt,
		)

		stepAssets = append(stepAssets, stepAsset)
	}
	return stepAssets, nil
}

func GetStartAndEndDateDiffStepAssets(diffInDate int) ([]StepAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
		id,
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		name,
		uuid,
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
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	stepAssets := make([]StepAsset, 0)
	for results.Next() {
		var stepAsset StepAsset
		results.Scan(
			&stepAsset.ID,
			&stepAsset.StepID,
			&stepAsset.AssetID,
			&stepAsset.SwapAssetID,
			&stepAsset.TargetPoolID,
			&stepAsset.Name,
			&stepAsset.UUID,
			&stepAsset.AlternateName,
			&stepAsset.StartDate,
			&stepAsset.EndDate,
			&stepAsset.Description,
			&stepAsset.ActionParameter,
			&stepAsset.CreatedBy,
			&stepAsset.CreatedAt,
			&stepAsset.UpdatedBy,
			&stepAsset.UpdatedAt,
		)

		stepAssets = append(stepAssets, stepAsset)
	}
	return stepAssets, nil
}

func UpdateStepAsset(stepAsset StepAsset) error {
	// if the stepAsset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if stepAsset.ID == nil || *stepAsset.ID == 0 {
		return errors.New("stepAsset has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE step_assets SET 
		step_id=$1,
		asset_id=$2,
		swap_asset_id=$3,
		target_pool_id=$4,
		name=$5,
		uuid=$6,
		alternate_name=$7,
		start_date=$8,
		end_date=$9,
		description=$10,
		action_parameter=$11,
		updated_by=$12, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$13`,
		stepAsset.StepID,          //1
		stepAsset.AssetID,         //2
		stepAsset.SwapAssetID,     //3
		stepAsset.TargetPoolID,    //4
		stepAsset.Name,            //5
		stepAsset.UUID,            //6
		stepAsset.AlternateName,   //7
		stepAsset.StartDate,       //8
		stepAsset.EndDate,         //9
		stepAsset.Description,     //10
		stepAsset.ActionParameter, //11
		stepAsset.UpdatedBy,       //12
		stepAsset.ID)              //13
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertStepAsset(stepAsset StepAsset) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO step_assets 
	(
		step_id,
		asset_id,
		swap_asset_id,
		target_pool_id,
		name,
		uuid,
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
			$5, 
			uuid_generate_v4(), 
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
		&stepAsset.StepID,       //1
		&stepAsset.AssetID,      //2
		&stepAsset.SwapAssetID,  //3
		&stepAsset.TargetPoolID, //4
		&stepAsset.Name,         //5
		// &stepAsset.UUID,
		&stepAsset.AlternateName,   //6
		&stepAsset.StartDate,       //7
		&stepAsset.EndDate,         //8
		&stepAsset.Description,     //9
		&stepAsset.ActionParameter, //10
		&stepAsset.CreatedBy,       //11
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertStepAssets(stepAssets []StepAsset) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range stepAssets {
		stepAsset := stepAssets[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(stepAsset.UUID)
		row := []interface{}{
			*stepAsset.StepID,          //1
			*stepAsset.AssetID,         //2
			*stepAsset.SwapAssetID,     //3
			*stepAsset.TargetPoolID,    //4
			stepAsset.Name,             //5
			uuidString,                 //6
			stepAsset.AlternateName,    //7
			&stepAsset.StartDate,       //8
			&stepAsset.EndDate,         //9
			stepAsset.Description,      //10
			*stepAsset.ActionParameter, //11
			stepAsset.CreatedBy,        //12
			&now,                       //13
			stepAsset.CreatedBy,        //14
			&now,                       //15
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"step_assets"},
		[]string{
			"step_id",          //1
			"asset_id",         //2
			"swap_asset_id",    //3
			"target_pool_id",   //4
			"name",             //5
			"uuid",             //6
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
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}

	return nil
}
