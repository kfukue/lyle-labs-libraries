package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/kfukue/lyle-labs-libraries/utils"
	_ "github.com/lib/pq"
)

const (
	port        = 5432
	dbname      = "assetdb"
	sslmode     = "verify-ca"
	sslrootcert = `server-ca.pem`
	sslcert     = `client-cert.pem`
	sslkey      = `client-key.pem`
)

var DbConn *sql.DB
var DbConnPgx *pgx.Conn

func init() {
	DbConn, _ = SetupDatabase()
}

// SetupDatabase
func SetupDatabase() (*sql.DB, error) {
	godotenv.Load()
	var err error
	dbURI := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s  sslrootcert=%s sslcert=%s sslkey=%s"
	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}
	if utils.GetEnv() == "production" {
		log.Println("In production!")
		dbSecretPath := utils.MustGetenv("DB_SECRET_PATH")
		// Set each value dynamically w/ Sprintf
		password, err := utils.AccessSecretVersion(dbSecretPath)
		if err != nil {
			log.Fatal(err)
		}
		userSecretPath := utils.MustGetenv("USER_SECRET_PATH")
		user, err := utils.AccessSecretVersion(userSecretPath)
		if err != nil {
			log.Fatal(err)
		}
		instanceSecretPath := utils.MustGetenv("INSTANCE_SECRET_PATH")
		instanceConnectionName, err := utils.AccessSecretVersion(instanceSecretPath)
		if err != nil {
			log.Fatal(err)
		}
		dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", user, password, dbname, socketDir, instanceConnectionName)
		if err != nil {
			log.Fatal(fmt.Errorf("sql.Open: %v", err))
		}
	} else {
		dir := utils.GoDotEnvVariable("SSL_CERT_FILE_PATH")
		filepath.Abs(dir)
		fmt.Println("in base : " + dir)
		hostSecretPath := utils.GoDotEnvVariable("HOST_SECRET_PATH")
		host, err := utils.AccessSecretVersion(hostSecretPath)
		if err != nil {
			log.Fatal(err)
		}
		user := utils.GoDotEnvVariable("DB_USER")
		password := utils.GoDotEnvVariable("DB_PASS")
		dbURI = fmt.Sprintf(dbURI, host, port, user, password, dbname, sslmode, filepath.Join(dir, sslrootcert), filepath.Join(dir, sslcert), filepath.Join(dir, sslkey))
	}
	DbConn, err := sql.Open("pgx", dbURI)
	// DbConn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	DbConnPgx, err = pgx.Connect(context.Background(), dbURI)

	if err != nil {
		log.Fatal(err)
	}
	err = DbConn.Ping()
	if err != nil {
		panic(err)
	}
	DbConn.SetMaxOpenConns(25)
	DbConn.SetMaxIdleConns(25)
	DbConn.SetConnMaxLifetime(600 * time.Second)
	DbConnPgx.ConnInfo().RegisterDataType(pgtype.DataType{
		Value: &shopspring.Numeric{},
		Name:  "numeric",
		OID:   pgtype.NumericOID,
	})
	DbConnPgx.ConnInfo().RegisterDataType(pgtype.DataType{
		Value: &pgtypeuuid.UUID{},
		Name:  "uuid",
		OID:   pgtype.UUIDOID,
	})
	return DbConn, nil
}
