package User

import (
	"database/sql"
	"fmt"
)

const (
	// AccountTableName is the name of the table for the Account model
	AccountTableName = "Account"
	// AccountFirstNameCol is the column name of the model's first name
	AccountFirstNameCol = "first_name"
	// AccountLastNameCol is the column name of the model's last name
	AccountLastNameCol = "last_name"
	// AccountAgeCol is the column name of the model's age
	AccountAgeCol = "age"
)

// Account is the database model for a Account
type Account struct {
	FirstName string
	LastName  string
	Age       uint
}

// CreateAccountTable uses db to create a new table for Account models, and returns the result
func CreateAccountTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf("CREATE TABLE %s (%s varchar(255), %s varchar(255), %s int)",
			AccountTableName,
			AccountFirstNameCol,
			AccountLastNameCol,
			AccountAgeCol,
		),
	)
}

// InsertAccount inserts Account into db
func InsertAccount(db *sql.DB, Account Account) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf("INSERT INTO %s VALUES(?, ?, ?)", AccountTableName),
		Account.FirstName,
		Account.LastName,
		Account.Age,
	)
}

// SelectAccount selects a Account with the given first & last names and age. On success, writes the result into result and on failure, returns a non-nil error and makes no modifications to result
func SelectAccount(db *sql.DB, firstName, lastName string, age uint, result *Account) error {
	row := db.QueryRow(
		fmt.Sprintf(
			"SELECT * FROM %s WHERE %s=? AND %s=? AND %s=?",
			AccountTableName,
			AccountFirstNameCol,
			AccountLastNameCol,
			AccountAgeCol,
		),
		firstName,
		lastName,
		age,
	)
	var retFirstName, retLastName string
	var retAge uint
	if err := row.Scan(&retFirstName, &retLastName, &retAge); err != nil {
		return err
	}
	result.FirstName = retFirstName
	result.LastName = retLastName
	result.Age = retAge
	return nil
}

// UpdateAccount updates the Account with the given first & last names and age with newAccount. Returns a non-nil error if the update failed, and nil if the update succeeded
func UpdateAccount(db *sql.DB, firstName, lastName string, age uint, newAccount Account) error {
	_, err := db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET %s=?,%s=?,%s=? WHERE %s=? AND %s=? AND %s=?",
			AccountTableName,
			AccountFirstNameCol,
			AccountLastNameCol,
			AccountAgeCol,
			AccountFirstNameCol,
			AccountLastNameCol,
			AccountAgeCol,
		),
		newAccount.FirstName,
		newAccount.LastName,
		newAccount.Age,
		firstName,
		lastName,
		age,
	)
	return err
}

// DeleteAccount deletes the Account with the given first & last names and age. Returns a non-nil error if the delete failed, and nil if the delete succeeded
func DeleteAccount(db *sql.DB, firstName, lastName string, age uint) error {
	_, err := db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE %s=? AND %s=? AND %s=?",
			AccountTableName,
			AccountFirstNameCol,
			AccountLastNameCol,
			AccountAgeCol,
		),
		firstName,
		lastName,
		age,
	)
	return err
}
