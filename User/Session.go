package User

import (
	"database/sql"
	"fmt"
)

//Session Table

const (
	TableName = "Session"

	SessionIdCol = "SessionId"

	UserIdCol = "UserId"
)

///Users/martinsinnuss/go/pkg/mod/github.com/arschles/go-in-5-minutes@v0.0.0-20200709150023-eb8196a64257/episode22/models/models.go

// Session is the database model for a Session
type Session struct {
	SessionId string
	UserId    uint
}

// CreateSessionTable uses db to create a new table for Session, and returns the result
func CreateSessionTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf("CREATE TABLE %s (%s varchar(255), %s int)",
			TableName,
			SessionIdCol,
			UserIdCol,
		),
	)
}

// InsertSession inserts Session into db
func InsertSession(db *sql.DB, session Session) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf("INSERT INTO %s VALUES($1, $2)", TableName),
		session.SessionId,
		session.UserId,
	)
}

// SelectSession selects a Session with the given id and email and DeviceId. On success, writes the result into result and on failure, returns a non-nil error and makes no modifications to result
func SelectSession(db *sql.DB, SessionId string, UserId uint, result *Session) error {
	row := db.QueryRow(
		fmt.Sprintf(
			"SELECT * FROM %s WHERE %s=$1 AND %s=$2",
			TableName,
			SessionIdCol,
			UserIdCol,
		),
		SessionId,
		UserId,
	)

	var retSessionId string
	var retUserId uint
	if err := row.Scan(&retSessionId, &retUserId); err != nil {
		return err
	}
	result.SessionId = retSessionId
	result.UserId = retUserId
	return nil
}

// UpdateSession updates the Session with the id, email and DeviceId with newSession. Returns a non-nil error if the update failed, and nil if the update succeeded
func UpdateSession(db *sql.DB, SessionId string, UserId uint, newSession Session) error {
	_, err := db.Exec(
		fmt.Sprintf(
			//"UPDATE %s SET %s,%s,%s WHERE %s AND %s AND %s",
			"UPDATE %s SET %s=$1,%s=$2 WHERE %s=$4 AND %s=$5",
			TableName,
			SessionIdCol,
			UserIdCol,
			SessionIdCol,
			UserIdCol,
		),
		newSession.SessionId,
		newSession.UserId,
		SessionId,
		UserId,
	)
	return err
}

// DeleteSession deletes the Session with the given id, email and DeviceId. Returns a non-nil error if the delete failed, and nil if the delete succeeded
func DeleteSession(db *sql.DB, SessionId string, UserId uint) error {
	_, err := db.Exec(
		fmt.Sprintf(
			//"DELETE FROM $1 WHERE $2 AND $3 AND $4",
			"DELETE FROM %s WHERE %s=$1 AND %s=$2",
			TableName,
			SessionIdCol,
			UserIdCol,
		),
		SessionId,
		UserId,
	)
	return err
}

func GetSession(db *sql.DB) error {
	_, err := db.Exec(

		"SELECT UserId FROM Session",
	)
	return err
}
