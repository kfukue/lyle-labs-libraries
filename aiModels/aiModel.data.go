package aimodels

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

func GetAIModel(dbConnPgx utils.PgxIface, aiModelID *int) (*AIModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
	id,
	uuid, 
	name, 
	alternate_name, 
	url,
	ticker,
	description,
	ollama_name,
	params_size,
	quantiized_size,
	base_model_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM ai_models 
	WHERE id = $1`, *aiModelID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	aiModel, err := pgx.CollectOneRow(row, pgx.RowToStructByName[AIModel])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &aiModel, nil
}

func RemoveAIModel(dbConnPgx utils.PgxIface, aiModelID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemovePositionJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM ai_models WHERE id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *aiModelID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetAIModelList(dbConnPgx utils.PgxIface, ids []int) ([]AIModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	url,
	ticker,
	description,
	ollama_name,
	params_size,
	quantiized_size,
	base_model_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM ai_models`
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

	aiModels, err := pgx.CollectRows(results, pgx.RowToStructByName[AIModel])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return aiModels, nil
}

func UpdateAIModel(dbConnPgx utils.PgxIface, aiModel *AIModel) error {
	// if the aiModel id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if aiModel.ID == nil || *aiModel.ID == 0 {
		return errors.New("aiModel has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdatePositionJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE ai_models SET 
		name=$1,  
		alternate_name=$2, 
		url=$3,
		ticker=$4,
		description=$5,
		ollama_name=$6,
		params_size=$7,
		quantiized_size=$8,
		base_model_id=$9,
		updated_by=$10, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$11`

	if _, err := dbConnPgx.Exec(ctx, sql,
		aiModel.Name,          //1
		aiModel.AlternateName, //2
		aiModel.URL,           //3
		aiModel.Ticker,        //4
		aiModel.Description,   //5
		aiModel.OllamaName,    //6
		aiModel.ParamsSize,    //7
		aiModel.QuantizedSize, //8
		aiModel.BaseModelID,   //9
		aiModel.UpdatedBy,     //10
		aiModel.ID,            //11
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertAIModel(dbConnPgx utils.PgxIface, aiModel *AIModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertMarketDataJob DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO ai_models  
	(
		uuid,
		name, 
		alternate_name, 
		url,
		ticker,
		description,
		ollama_name,
		params_size,
		quantiized_size,
		base_model_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES (
		 	uuid_generate_v4(),
		 	$1,
			$2, $3, 
			$4, 
			$5, $6, 
			$7, $8, 
			$9, $10,
			current_timestamp at time zone 'UTC', $11, 
			current_timestamp at time zone 'UTC')
		RETURNING id`,
		aiModel.Name,          //1
		aiModel.AlternateName, //2
		aiModel.URL,           //3
		aiModel.Ticker,        //4
		aiModel.Description,   //5
		aiModel.OllamaName,    //6
		aiModel.ParamsSize,    //7
		aiModel.QuantizedSize, //8
		aiModel.BaseModelID,   //9
		aiModel.CreatedBy,     //10
		aiModel.CreatedBy,     //11
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

func InsertAIModels(dbConnPgx utils.PgxIface, aiModels []AIModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range aiModels {
		aiModel := aiModels[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(aiModel.UUID)
		row := []interface{}{
			uuidString,            //1
			aiModel.Name,          //2
			aiModel.AlternateName, //3
			aiModel.URL,           //4
			aiModel.Ticker,        //5
			aiModel.Description,   //6
			aiModel.OllamaName,    //7
			aiModel.ParamsSize,    //8
			aiModel.QuantizedSize, //9
			aiModel.BaseModelID,   //10
			aiModel.CreatedBy,     //11
			&aiModel.CreatedAt,    //12
			aiModel.CreatedBy,     //13
			&now,                  //14
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"ai_models"},
		[]string{
			"uuid",            //1
			"name",            //2
			"alternate_name",  //3
			"url",             //4
			"ticker",          //5
			"description",     //6
			"ollama_name",     //7
			"params_size",     //8
			"quantiized_size", //9
			"base_model_id",   //10
			"created_by",      //11
			"created_at",      //12
			"updated_by",      //13
			"updated_at",      //14
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertAIModels: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetAIModelListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]AIModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid, 
		name, 
		alternate_name, 
		url,
		ticker,
		description,
		ollama_name,
		params_size,
		quantiized_size,
		base_model_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM ai_models 
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

	aiModels, err := pgx.CollectRows(results, pgx.RowToStructByName[AIModel])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return aiModels, nil
}

func GetTotalAIModelsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM ai_models
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
