package source

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

func GetSource(dbConnPgx utils.PgxIface, sourceID *int) (*Source, error) {
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
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM sources 
	WHERE id = $1`, *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	source, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Source])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &source, nil
}

func RemoveSource(dbConnPgx utils.PgxIface, sourceID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemovePositionJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM sources WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *sourceID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetSourceList(dbConnPgx utils.PgxIface, ids []int) ([]Source, error) {
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
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM sources`
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
	sources, err := pgx.CollectRows(results, pgx.RowToStructByName[Source])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sources, nil
}

func UpdateSource(dbConnPgx utils.PgxIface, source *Source) error {
	// if the source id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if source.ID == nil || *source.ID == 0 {
		return errors.New("source has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdatePositionJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE sources SET 
		name=$1,  
		alternate_name=$2, 
		url=$3,
		ticker=$4,
		description=$5,
		updated_by=$6, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$7`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		source.Name,          //1
		source.AlternateName, //2
		source.URL,           //3
		source.Ticker,        //4
		source.Description,   //5
		source.UpdatedBy,     //6
		source.ID,            //7
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertSource(dbConnPgx utils.PgxIface, source *Source) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertMarketDataJob DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO sources  
	(
		uuid,
		name, 
		alternate_name, 
		url,
		ticker,
		description,
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
			current_timestamp at time zone 'UTC', $7, 
			current_timestamp at time zone 'UTC')
		RETURNING id`,
		source.Name,          //1
		source.AlternateName, //2
		source.URL,           //3
		source.Ticker,        //4
		source.Description,   //5
		source.CreatedBy,     //6
		source.CreatedBy,     //7
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

func InsertSources(dbConnPgx utils.PgxIface, sources []Source) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range sources {
		source := sources[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(source.UUID)
		row := []interface{}{
			uuidString,           //1
			source.Name,          //2
			source.AlternateName, //3
			source.URL,           //4
			source.Ticker,        //5
			source.Description,   //6
			source.CreatedBy,     //7
			&source.CreatedAt,    //8
			source.CreatedBy,     //9
			&now,                 //10
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"sources"},
		[]string{
			"uuid",           //1
			"name",           //2
			"alternate_name", //3
			"url",            //4
			"ticker",         //5
			"description",    //6
			"created_by",     //7
			"created_at",     //8
			"updated_by",     //9
			"updated_at",     //10
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertSources: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetSourceListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Source, error) {
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
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM sources 
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
	sources, err := pgx.CollectRows(results, pgx.RowToStructByName[Source])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sources, nil
}

func GetTotalSourcesCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM sources
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
