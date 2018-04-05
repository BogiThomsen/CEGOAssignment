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
	err := connectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	//Close the database when done
	defer closeDatabase() 	

	//Run select statement
	err = runSelectStatement(&users)
	if len(users) == 0 {
		err = fmt.Errorf("No users returned by statement: \"%s\"\n", config.Stmt)
	}
	if err != nil{
		log.Fatal(err)
	} else {
		fmt.Printf("%d users retrieved from database:\n", len(users))
	}

	//Prints users to console for user confirmation
	printUsers(users)
	
	fmt.Printf("Do you wish to delete these %d users from the database? (y/n)\n", len(users))
	//Ask user for confirmation before deleting users
	conf := userConfirmation()
	if conf {
		//Write users to file
		err = writeUsersToFile(&users)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("%d users written to file.\n", len(users))
		}


		//Checking file data vs database data
		fileUsers, err := readUsersFromFile()
		if err != nil {
			log.Fatal(err)
		}

		//Check data in database against data in file
		correct := checkDbAgainstFile(fileUsers)
		if correct {
			fmt.Println("Comparison finished. Data is correct")
			//Deletes users based on ID
			err = deleteUsersByIdTransaction()
			if err != nil {
				fmt.Println("Delete failed. Now deleting created file.")
				//Delete file
				deleteFile()
				err = fmt.Errorf("File deleted due to: %s\n", err)
				log.Fatal(err)
			}
		} else {
			err = fmt.Errorf("Data in file is not correct. Delete aborted. Now deleting file.")
			deleteFile()
			log.Fatal(err)
		} 
	}
}

//Asks user for y/n confirmation returns bool if answer is y, loops on wrong input
func userConfirmation() bool {
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
		return userConfirmation()
	}
}


//prints user data to console
func printUsers(users []User){
	//iterate over users in array
	for _, user := range users {
		//print to console
		fmt.Println("id: " + user.ID + ", firstName: " + user.firstName + ", lastName: " + user.lastName + ", email: " + user.email)
	}
}