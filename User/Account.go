package User

import (
	"database/sql"
	"fmt"
)

const (
	// AccountTableName is the name of the table for the Account model
	AccountTableName = "Account"
	// AccountIdNumberCol is the column name of the model's first name
	AccountIdNumberCol = "IdNumber"
	// AccountEmailCol is the column name of the model's last name
	AccountEmailCol = "Email"
	// AccountDeviceIdCol is the column name of the model's DeviceId
	AccountDeviceIdCol = "DeviceId"
)

// Account is the database model for a Account
type Account struct {
	IdNumber uint
	Email    string
	DeviceId uint
}

// CreateAccountTable uses db to create a new table for Account models, and returns the result
func CreateAccountTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf("CREATE TABLE %s (%s varchar(255), %s varchar(255), %s int)",
			AccountTableName,
			AccountIdNumberCol,
			AccountEmailCol,
			AccountDeviceIdCol,
		),
	)
}

// InsertAccount inserts Account into db
func InsertAccount(db *sql.DB, Account Account) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf("INSERT INTO %s VALUES(?, ?, ?)", AccountTableName),
		Account.IdNumber,
		Account.Email,
		Account.DeviceId,
	)
}

// SelectAccount selects a Account with the given first & last names and DeviceId. On success, writes the result into result and on failure, returns a non-nil error and makes no modifications to result
func SelectAccount(db *sql.DB, IdNumber, Email string, DeviceId uint, result *Account) error {
	row := db.QueryRow(
		fmt.Sprintf(
			"SELECT * FROM %s WHERE %s=? AND %s=? AND %s=?",
			AccountTableName,
			AccountIdNumberCol,
			AccountEmailCol,
			AccountDeviceIdCol,
		),
		IdNumber,
		Email,
		DeviceId,
	)
	var retEmail string
	var retIdNumber, retDeviceId uint
	if err := row.Scan(&retIdNumber, &retEmail, &retDeviceId); err != nil {
		return err
	}
	result.IdNumber = retIdNumber
	result.Email = retEmail
	result.DeviceId = retDeviceId
	return nil
}

// UpdateAccount updates the Account with the given first & last names and DeviceId with newAccount. Returns a non-nil error if the update failed, and nil if the update succeeded
func UpdateAccount(db *sql.DB, IdNumber, Email string, DeviceId uint, newAccount Account) error {
	_, err := db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET %s=?,%s=?,%s=? WHERE %s=? AND %s=? AND %s=?",
			AccountTableName,
			AccountIdNumberCol,
			AccountEmailCol,
			AccountDeviceIdCol,
			AccountIdNumberCol,
			AccountEmailCol,
			AccountDeviceIdCol,
		),
		newAccount.IdNumber,
		newAccount.Email,
		newAccount.DeviceId,
		IdNumber,
		Email,
		DeviceId,
	)
	return err
}

// DeleteAccount deletes the Account with the given first & last names and DeviceId. Returns a non-nil error if the delete failed, and nil if the delete succeeded
func DeleteAccount(db *sql.DB, IdNumber, Email string, DeviceId uint) error {
	_, err := db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE %s=? AND %s=? AND %s=?",
			AccountTableName,
			AccountIdNumberCol,
			AccountEmailCol,
			AccountDeviceIdCol,
		),
		IdNumber,
		Email,
		DeviceId,
	)
	return err
}
