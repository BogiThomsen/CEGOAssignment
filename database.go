package main

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/BurntSushi/toml"
	"fmt"
	"strings"
	"regexp"
)

//Config struct for toml config file
type tomlConfig struct {
	Title   string
	Stmt	string
	DB      database `toml:"database"`
}

//Database struct for toml onfig file
type database struct {
	User  		string
	Password   	string
	Dbname 		string 
	Host 		string
	Port 		string
}


var config tomlConfig
var processedStmt string

//Reads config file
//Processes config statement
//Creates connection to database
func ConnectToDatabase() (db *sql.DB, err error) {

	//Read config file into config struct
	_, err = toml.DecodeFile("config.toml", &config)
	if err != nil {
		return
	}
	processedStmt, err = ProcessStatement(config.Stmt)
	if err != nil {
		return
	}
	//Create DSN string from config struct
	DSN := (config.DB.User + ":" + config.DB.Password + "@("+config.DB.Host+":"+config.DB.Port+")/" + config.DB.Dbname)

	//Opens connection to database
	db, err = sql.Open("mysql", DSN)
	return
}

//Closes database connection
func CloseDatabase(db *sql.DB){
	db.Close()
}

//Runs a select statement and adds users to user array
//Can only run "SELECT * FROM users" statements.
//Can still handle any WHERE clause 
func RunSelectStatement(db *sql.DB, users *[]User) (err error){
	//create user struct
	var id, firstName, lastName, email string

	// Run select stmt
	fmt.Printf("Running Select Statement: \"%s\"\n", processedStmt)
	rows, err := db.Query(processedStmt)
	if err != nil {
		return
	}
	defer rows.Close()

	//Iterate over rows
	for rows.Next() {
		//scan all columns into user struct
		err := rows.Scan(&id, &firstName, &lastName, &email)
		if err != nil {
			log.Fatal(err)
		}
		user := NewUser(id, firstName, lastName, email)
		//add current user to user struct array
		*users = append(*users, user)
	}
	if err != nil{
		return
	}
	return
}

//Creates temporary table with select statement.
//This is used for passing into an IN statement for deletion of users
func CreateTempTableForDelete(db *sql.DB) (err error){
	tableStatement := fmt.Sprintf("CREATE TEMPORARY TABLE IF NOT EXISTS idForDelete AS (%s);", strings.Trim(processedStmt, ";"))
	_, err = db.Exec(tableStatement)
	if err != nil {
		return
	}
	return
}


//Begins a transaction, deletes users based on ID, Commits if no errors, else rolls back transaction.
func DeleteUsersByIdTransaction (db *sql.DB) (err error){
	//Open transaction connection
	tx, err := db.Begin()
    if err != nil {
        return
    }

    //Defer function to rollback or commit depending on errors
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
        err = tx.Commit()
    }()
     
    fmt.Println("Deleting users from database.")
    //Transaction runs delete statement    
    result, err := tx.Exec("DELETE FROM users WHERE id IN (SELECT id FROM idForDelete);")
    if err != nil{
    	return
    } else {
    	rowsAffected, _ := result.RowsAffected()
    	fmt.Printf("%d rows deleted from database.\n", rowsAffected)
    }
    return
}

func CheckDbAgainstFile(db *sql.DB, users []User) (b bool, err error){
	//Get user from DB
	var id, firstName, lastName, email string
	fmt.Println("Comparing file to database.")
	stmt := "SELECT id, firstName, lastName, email FROM users WHERE id = ?"
	for _, user := range users{
		err := db.QueryRow(stmt, user.ID).Scan(&id, &firstName, &lastName, &email)
		if err != nil {
			return false, err
		}
		if (id != user.ID || firstName != user.firstName || lastName != user.lastName || email != user.email){
			return false, err
		}
	}
	return true, err
}

//Makes sure statement selects all rows from user
//Removes additional statements after first ";""
func ProcessStatement(stmt string) (result string, err error) {
	//Create match string
	match := "^(SELECT id, firstName, lastName, email FROM users)[^;]*"
	matched, err := regexp.MatchString(match, stmt)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("Statement doesnt match \"%s\".\n", match)
	} 
	//Removes additional statements
	result = fmt.Sprintf("%s%s", strings.Split(stmt, ";")[0], ";")
	return result, err
}


