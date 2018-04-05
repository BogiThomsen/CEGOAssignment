package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

//User struct for storing data
type User struct{
	ID string
	firstName string
	lastName string
	email string
}

func main() {
	//create array of user structs.
	var users []User

	//Connect to databse
	db, err := ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	//Close the database when done
	defer CloseDatabase(db) 	

	//Run select statement
	err = RunSelectStatement(db, &users)
	if len(users) == 0 {
		err = fmt.Errorf("No users returned by statement: \"%s\"\n", processedStmt)
	}
	if err != nil{
		log.Fatal(err)
	} else {
		fmt.Printf("%d users retrieved from database:\n", len(users))
	}

	//Prints users to console for user confirmation
	PrintUsers(users)
	
	fmt.Printf("Do you wish to delete these %d users from the database? (y/n)\n", len(users))
	//Ask user for confirmation before deleting users
	conf := UserConfirmation()
	if conf {
		//Write users to file
		err = WriteUsersToFile(&users)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("%d users written to file.\n", len(users))
		}


		//Checking file data vs database data
		fileUsers, err := ReadUsersFromFile()
		if err != nil {
			log.Fatal(err)
		}

		//Check data in database against data in file
		correct, err := CheckDbAgainstFile(db, fileUsers)
		if err != nil {
			log.Fatal(err)
		}
		if correct {
			fmt.Println("Comparison finished. Data is correct")

			//Creating temporary table to delete users from database
			err = CreateTempTableForDelete(db)
			if err != nil {
				return
			}
			//Deletes users based on ID
			err = DeleteUsersByIdTransaction(db)
			if err != nil {
				fmt.Println("Delete failed. Now deleting created file.")
				//Delete file
				DeleteFile()
				err = fmt.Errorf("File deleted due to: %s\n", err)
				log.Fatal(err)
			}
		} else {
			err = fmt.Errorf("Data in file is not correct. Delete aborted. Now deleting file.")
			DeleteFile()
			log.Fatal(err)
		} 
	}
}

//Asks user for y/n confirmation returns bool if answer is y, loops on wrong input
func UserConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	//Trims whitespace off input for comparison
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "y" {
		return true
	} else if input == "n" {
		return false
	} else {
		//Recursion on wrong input to ask again
		fmt.Println("Please type y or n, then ENTER:")
		return UserConfirmation()
	}
}


//prints user data to console
func PrintUsers(users []User){
	//iterate over users in array
	for _, user := range users {
		//print to console
		fmt.Println("id: " + user.ID + ", firstName: " + user.firstName + ", lastName: " + user.lastName + ", email: " + user.email)
	}
}

func NewUser (id string, fName string, lName string, mail  string) User {
	return User {
				ID: id,
				firstName: fName,
				lastName: lName, 
				email: mail,
			}
}