package gethlylejobstopics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetGethProcessJobTopic(dbConnPgx utils.PgxIface, gethProcessJobTopicID *int) (*GethProcessJobTopic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	geth_process_job_id,
	uuid, 
	name,
	alternate_name,
	description,
	status_id,
	topic_str,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_process_job_topics 
	WHERE id = $1
	`, *gethProcessJobTopicID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	gethProcessJobTopic, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethProcessJobTopic])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethProcessJobTopic, nil
}

func GetGethProcessJobTopicList(dbConnPgx utils.PgxIface) ([]GethProcessJobTopic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	geth_process_job_id,
	uuid, 
	name,
	alternate_name,
	description,
	status_id,
	topic_str,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_process_job_topics`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethProcessJobTopics, err := pgx.CollectRows(results, pgx.RowToStructByName[GethProcessJobTopic])

	return gethProcessJobTopics, nil
}

func RemoveGethProcessJobTopic(dbConnPgx utils.PgxIface, gethProcessJobTopicID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethProcessJobTopic DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_process_job_topics WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *gethProcessJobTopicID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func UpdateGethProcessJobTopic(dbConnPgx utils.PgxIface, gethProcessJobTopic *GethProcessJobTopic) error {
	// if the gethProcessJobTopic id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethProcessJobTopic.ID == nil || *gethProcessJobTopic.ID == 0 {
		return errors.New("gethProcessJobTopic has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethAddress DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_process_job_topics SET
		geth_process_job_id = $1,
		name=$2,
		alternate_name=$3,
		description=$4,
		status_id=$5,
		topic_str=$6,
		updated_by=$7, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$8 `
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		gethProcessJobTopic.GethProcessJobID, //1
		gethProcessJobTopic.Name,             //2
		gethProcessJobTopic.AlternateName,    //3
		gethProcessJobTopic.Description,      //4
		gethProcessJobTopic.StatusID,         //5
		gethProcessJobTopic.TopicStr,         //6
		gethProcessJobTopic.UpdatedBy,        //7
		gethProcessJobTopic.ID,               //8
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethProcessJobTopic(dbConnPgx utils.PgxIface, gethProcessJobTopic *GethProcessJobTopic) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethAddress DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var ID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_process_job_topics  
	(
		geth_process_job_id,
		uuid, 
		name,
		alternate_name,
		description,
		status_id,
		topic_str,
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
		gethProcessJobTopic.GethProcessJobID, //1
		gethProcessJobTopic.Name,             //2
		gethProcessJobTopic.AlternateName,    //3
		gethProcessJobTopic.Description,      //4
		gethProcessJobTopic.StatusID,         //5
		gethProcessJobTopic.TopicStr,         //6
		gethProcessJobTopic.CreatedBy,        //7
	).Scan(&ID)
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
	return int(ID), nil
}

func InsertGethProcessJobTopicList(dbConnPgx utils.PgxIface, gethProcessJobTopicList []GethProcessJobTopic) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethProcessJobTopicList {
		gethProcessJobTopic := gethProcessJobTopicList[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(gethProcessJobTopic.UUID)
		row := []interface{}{
			gethProcessJobTopic.GethProcessJobID, //1
			uuidString,                           //2
			gethProcessJobTopic.Name,             //3
			gethProcessJobTopic.AlternateName,    //4
			gethProcessJobTopic.Description,      //5
			gethProcessJobTopic.StatusID,         //6
			gethProcessJobTopic.TopicStr,         //7
			gethProcessJobTopic.CreatedBy,        //8
			&gethProcessJobTopic.CreatedAt,       //9
			gethProcessJobTopic.CreatedBy,        //10
			&now,                                 //11
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_process_job_topics"},
		[]string{
			"geth_process_job_id", //1
			"uuid",                //2
			"name",                //3
			"alternate_name",      //4
			"description",         //5
			"status_id",           //6
			"topic_str",           //7
			"created_by",          //8
			"created_at",          //9
			"updated_by",          //10
			"updated_at",          //11
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethProcessJobTopicList: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetGethProcessJobTopicListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethProcessJobTopic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,  
		geth_process_job_id,
		uuid, 
		name,
		alternate_name,
		description,
		status_id,
		topic_str,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		FROM geth_process_job_topics 
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
	gethProcessJobList, err := pgx.CollectRows(results, pgx.RowToStructByName[GethProcessJobTopic])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethProcessJobList, nil
}

func GetTotalGethProcessJobTopicCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM geth_process_job_topics`)
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
