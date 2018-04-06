package main

import (
	"testing"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

)

func TestDeleteUserByIdTransactionPass(t *testing.T){
	//Create mock connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening mock database connection\n", err)
	}
	defer db.Close()

	//Set expectations
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM users").WillReturnResult(sqlmock.NewResult(0,1))
	mock.ExpectCommit()

	//Execute function
	err = DeleteUsersByIdTransaction(db)
	if err != nil {
		t.Errorf("Error was not expected when deleting users: %s\n", err)
	}

	//Make sure expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expectations were not met: %s\n", err)
	}

}

func TestDeleteUserByIdTransactionExpectRollbackOnError(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening mock database connection\n", err)
	}
	defer db.Close()

	//Set expectations
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM users").WillReturnError(fmt.Errorf("Deletion error"))
	mock.ExpectRollback()

	//Execute function
	err = DeleteUsersByIdTransaction(db)
	if err == nil {
		t.Errorf("Error was expected when deleting users, found no error.\n")
	}

	//Make sure expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expectations were not met: %s\n", err)
	}
}

func TestCheckDbAgainstFilePass(t *testing.T){
	//Create mock connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening mock database connection\n", err)
	}
	defer db.Close()

	//Create user array
	var users []User
	users = append(users, NewUser("1", "testFirstName", "testLastName", "testEmail"))
	
	//Create columns
	columns := []string{"id", "firstName", "lastName", "email"}

	//Set expectations
	mock.ExpectQuery("SELECT id, firstName, lastName, email FROM users").WithArgs("1").WillReturnRows(sqlmock.NewRows(columns).AddRow("1", "testFirstName", "testLastName", "testEmail"))
	
	//Execute function
	actualResult, err := CheckDbAgainstFile(db, users)
	if err != nil {
		t.Fatalf("An error '%s' was not expected when querying database\n", err)
	}

	var expectedResult = true

	if actualResult != expectedResult {
		t.Errorf("Expected %t but got %t\n", expectedResult, actualResult)
	}

	//Make sure expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expectations were not met: %s\n", err)
	}
}

func TestCheckDbAgainstFileFailComparison(t *testing.T){
	//Create mock connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening mock database connection\n", err)
	}
	defer db.Close()

	//Create user array
	var users []User
	users = append(users, NewUser("1", "testFirstName", "testLastName", "testEmail"))
	//Create columns
	columns := []string{"id", "firstName", "lastName", "email"}

	//Set expectations: returning wrong first name 
	mock.ExpectQuery("SELECT id, firstName, lastName, email FROM users").WithArgs("1").WillReturnRows(sqlmock.NewRows(columns).AddRow("1", "testWrongFirstName", "testLastName", "testEmail"))
	
	//Execute function
	actualResult, err := CheckDbAgainstFile(db, users)

	//Expecting to return false
	var expectedResult = false

	//Testing expected vs actual
	if actualResult != expectedResult {
		t.Errorf("Expected %t but got %t\n", expectedResult, actualResult)
	}

	//Make sure expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expectations were not met: %s\n", err)
	}
}

func TestCheckDbAgainstFileFailErrorQuery(t *testing.T){
	//Create mock connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening mock database connection\n", err)
	}
	defer db.Close()

	//Create user array
	var users []User
	users = append(users, NewUser("1", "testFirstName", "testLastName", "testEmail"))

	//Set expectations: query returns an error
	mock.ExpectQuery("SELECT id, firstName, lastName, email FROM users").WithArgs("1").WillReturnError(fmt.Errorf("Query error"))
	
	//Execute function
	_, err = CheckDbAgainstFile(db, users)
	if err == nil {
		t.Errorf("Error was expected when querying database, found no error.\n")
	}

	//Make sure expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expectations were not met: %s\n", err)
	}
}

func TestCreateTempTableForDeletePass(t *testing.T){
	//Create mock connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening mock database connection\n", err)
	}
	defer db.Close()

	//Set expectations
	mock.ExpectExec("CREATE TEMPORARY TABLE IF NOT EXISTS idForDelete").WillReturnResult(sqlmock.NewResult(0,1))

	//Execute function
	err = CreateTempTableForDelete(db)
	if err != nil {
		t.Fatalf("An error '%s' was not expected when creating temporary table\n", err)
	}

	//Make sure expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expectations were not met: %s\n", err)
	}

}

func TestCreateTempTableForDeleteErrorOnCreateTable(t *testing.T){
	//Create mock connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening mock database connection\n", err)
	}
	defer db.Close()

	//Set expectations
	mock.ExpectExec("CREATE TEMPORARY TABLE IF NOT EXISTS idForDelete").WillReturnError(fmt.Errorf("Create temporary table error"))

	//Execute function
	err = CreateTempTableForDelete(db)
	if err == nil {
		t.Errorf("Error was expected when querying database, found no error.\n")
	}

	//Make sure expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expectations were not met: %s\n", err)
	}

}

func TestProcessStatementPass(t *testing.T){
	//Good statement
	stmt := "SELECT id, firstName, lastName, email FROM users WHERE firstName LIKE 'R%';"
	expectedResult := stmt

	actualResult, err := ProcessStatement(stmt)
	if err != nil {
		t.Errorf("An error '%s' was not expected when matching statement.\n", err)
	}

	if expectedResult != actualResult {
		t.Errorf("Expected %s but got %s\n", expectedResult, actualResult)
	}
}

func TestProcessStatementPassAdditionalStatement(t *testing.T){
	//Additional statement
	stmt := "SELECT id, firstName, lastName, email FROM users; DROP TABLE users;"
	//Expecting to drop additional statement
	expectedResult := "SELECT id, firstName, lastName, email FROM users;"

	actualResult, err := ProcessStatement(stmt)
	if err != nil {
		t.Errorf("An error '%s' was not expected when matching statement.\n", err)
	}

	if expectedResult != actualResult {
		t.Errorf("Expected %s but got %s\n", expectedResult, actualResult)
	}
}



func TestProcessStatementFailBadStatement(t *testing.T){
	//Good statement
	stmt := "SELECT id FROM users;"

	_, err := ProcessStatement(stmt)
	if err == nil {
		t.Errorf("Expected an error when matching bad statement, found no error\n")
	}
}