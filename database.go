package main

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/BurntSushi/toml"
	"fmt"
	"strings"
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


var db *sql.DB
var config tomlConfig

//Reads config file,creates connection to database
func connectToDatabase() (err error) {

	//Read config file into config struct
	_, err = toml.DecodeFile("config.toml", &config)
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
func closeDatabase(){
	db.Close()
}

//Runs a select statement and adds users to user array
//Can only run "SELECT * FROM users" statements.
//Can still handle any WHERE clause 
func runSelectStatement(users *[]User) (err error){
	//create user struct
	var user User

	// Run select stmt
	fmt.Printf("Running Select Statement: \"%s\"\n", config.Stmt)
	rows, err := db.Query(config.Stmt)
	if err != nil {
		return
	}
	defer rows.Close()

	//Iterate over rows
	for rows.Next() {
		//scan all columns into user struct
		err := rows.Scan(&user.ID, &user.firstName, &user.lastName, &user.email)
		if err != nil {
			log.Fatal(err)
		}
		//add current user to user struct array
		*users = append(*users, user)
	}
	if err != nil{
		return
	}
	return
}

//Creates temporary table 
func createTempTableForDelete() (err error){
	tableStatement := fmt.Sprintf("CREATE TEMPORARY TABLE IF NOT EXISTS idForDelete AS (%s);", strings.Trim(config.Stmt, ";"))
	_, err = db.Exec(tableStatement)
	if err != nil {
		return
	}
	return
}


//Begins a transaction, deletes users based on ID, Commits if no errors, else rolls back transaction.
func deleteUsersByIdTransaction () (err error){
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
    //Creating temporary table to delete users from database
	err = createTempTableForDelete()
	if err != nil {
		return
	} 
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

func checkDbAgainstFile(users []User) bool{
	//Get user from DB
	var id, firstName, lastName, email string
	fmt.Println("Comparing file to database.")
	stmt := "SELECT id, firstName, lastName, email FROM users WHERE id = ?"
	for _, user := range users{
		err := db.QueryRow(stmt, user.ID).Scan(&id, &firstName, &lastName, &email)
		if err != nil {
			log.Fatal(err)
		}
		if (id != user.ID || firstName != user.firstName || lastName != user.lastName || email != user.email){
			return false
		}
	}
	return true
}


