package gethlylejobs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
)

func GetGethProcessJobTopic(gethProcessJobTopicID int) (*GethProcessJobTopic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	`, gethProcessJobTopicID)

	gethProcessJobTopic := &GethProcessJobTopic{}
	err := row.Scan(
		&gethProcessJobTopic.ID,
		&gethProcessJobTopic.GethProcessJobID,
		&gethProcessJobTopic.UUID,
		&gethProcessJobTopic.Name,
		&gethProcessJobTopic.AlternateName,
		&gethProcessJobTopic.Description,
		&gethProcessJobTopic.StatusID,
		&gethProcessJobTopic.TopicStr,
		&gethProcessJobTopic.CreatedBy,
		&gethProcessJobTopic.CreatedAt,
		&gethProcessJobTopic.UpdatedBy,
		&gethProcessJobTopic.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethProcessJobTopic, nil
}

func GetGethProcessJobTopicList() ([]GethProcessJobTopic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	gethProcessJobTopics := make([]GethProcessJobTopic, 0)
	for results.Next() {
		var gethProcessJobTopic GethProcessJobTopic
		results.Scan(
			&gethProcessJobTopic.ID,
			&gethProcessJobTopic.GethProcessJobID,
			&gethProcessJobTopic.UUID,
			&gethProcessJobTopic.Name,
			&gethProcessJobTopic.AlternateName,
			&gethProcessJobTopic.Description,
			&gethProcessJobTopic.StatusID,
			&gethProcessJobTopic.TopicStr,
			&gethProcessJobTopic.CreatedBy,
			&gethProcessJobTopic.CreatedAt,
			&gethProcessJobTopic.UpdatedBy,
			&gethProcessJobTopic.UpdatedAt,
		)

		gethProcessJobTopics = append(gethProcessJobTopics, gethProcessJobTopic)
	}
	return gethProcessJobTopics, nil
}

func RemoveGethProcessJobTopic(gethProcessJobTopicID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_process_job_topics WHERE 
	id = $1`, *gethProcessJobTopicID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateGethProcessJobTopic(gethProcessJobTopic *GethProcessJobTopic) error {
	// if the gethProcessJobTopic id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethProcessJobTopic.ID == nil || *gethProcessJobTopic.ID == 0 {
		return errors.New("gethProcessJobTopic has invalid ID")
	}
	_, err := database.DbConnPgx.Exec(ctx, `UPDATE geth_process_job_topics SET 
		geth_process_job_id = $1,
		name=$2,
		alternate_name=$3,
		description=$4,
		status_id=$5,
		topic_str=$6,
		updated_by=$7, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$8 `,
		gethProcessJobTopic.GethProcessJobID, //1
		gethProcessJobTopic.Name,             //2
		gethProcessJobTopic.AlternateName,    //3
		gethProcessJobTopic.Description,      //4
		gethProcessJobTopic.StatusID,         //5
		gethProcessJobTopic.TopicStr,         //6
		gethProcessJobTopic.UpdatedBy,        //7
		gethProcessJobTopic.ID,               //8
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethProcessJobTopic(gethProcessJobTopic GethProcessJobTopic) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var ID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_process_job_topics  
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
		&gethProcessJobTopic.GethProcessJobID, //1
		gethProcessJobTopic.Name,              //2
		gethProcessJobTopic.AlternateName,     //3
		gethProcessJobTopic.Description,       //4
		gethProcessJobTopic.StatusID,          //5
		gethProcessJobTopic.TopicStr,          //6
		gethProcessJobTopic.CreatedBy,         //7
	).Scan(&ID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(ID), nil
}

func InsertGethProcessJobTopicList(gethProcessJobTopicList []GethProcessJobTopic) error {
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
			*gethProcessJobTopic.ID,               //1
			*gethProcessJobTopic.GethProcessJobID, //2
			uuidString,                            //3
			gethProcessJobTopic.Name,              //4
			gethProcessJobTopic.AlternateName,     //5
			gethProcessJobTopic.Description,       //6
			*gethProcessJobTopic.StatusID,         //7
			gethProcessJobTopic.TopicStr,          //8
			gethProcessJobTopic.CreatedBy,         //9
			&gethProcessJobTopic.CreatedAt,        //10
			gethProcessJobTopic.CreatedBy,         //11
			&now,                                  //12
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_process_job_topics"},
		[]string{
			"id",                  //1
			"geth_process_job_id", //2
			"uuid",                //3
			"name",                //4
			"alternate_name",      //5
			"description",         //6
			"status_id",           //7
			"topic_str",           //8
			"created_by",          //9
			"created_at",          //10
			"updated_by",          //11
			"updated_at",          //12
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
