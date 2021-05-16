package infrastructures

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	dbr "github.com/gocraft/dbr/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ISQLConnection is
type ISQLConnection interface {
	EsignRead() *dbr.Session
	EsignWrite() *dbr.Session
}

// SQLConnection define sql connection.
type SQLConnection struct{}

var (
	dbEsignRead, dbEsignWrite *dbr.Connection
	err                       error
)

// Connection open a new database connection
func Connection(db, descriptor string, maxIdle, maxConns int) *dbr.Connection {
	conn, err := dbr.Open(db, descriptor, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"action": "connection for " + db,
			"event":  db + "_error_connection",
		}).Error(err)
		os.Exit(0)
	}

	conn.SetMaxOpenConns(maxConns)
	conn.SetMaxIdleConns(maxIdle)

	return conn
}

// EsignRead create a new bareksa_marketdata session
func (s *SQLConnection) EsignRead() *dbr.Session {
	if dbEsignRead == nil {
		dbEsignRead = Connection(
			viper.GetString("database.client.driver"),
			viper.GetString("database.client.read"),
			viper.GetInt("database.client.max_idle"),
			viper.GetInt("database.client.max_cons"),
		)
	}

	sess := dbEsignRead.NewSession(nil)
	return sess
}

// EsignWrite create a new bareksa_marketdata session
func (s *SQLConnection) EsignWrite() *dbr.Session {
	if dbEsignWrite == nil {
		dbEsignWrite = Connection(
			viper.GetString("database.client.driver"),
			viper.GetString("database.client.write"),
			viper.GetInt("database.client.max_idle"),
			viper.GetInt("database.client.max_cons"),
		)
	}

	sess := dbEsignWrite.NewSession(nil)
	return sess
}
