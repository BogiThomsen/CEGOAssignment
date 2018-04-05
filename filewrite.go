package main

import (
    "fmt"
    "time"
    "log"
    "io"
    "os"
    //"strings"
    "path/filepath"
    "encoding/csv"
)

var fileName string

//Writes user data to CSV file with the delimiter ;
func WriteUsersToFile(users *[]User) (err error){
	//Create subfolder for data if it doesnt exist
	err = os.MkdirAll(filepath.Join(".","DeletedUserData"), os.ModePerm)
	if err != nil {
        err = fmt.Errorf("Error creating folder: %s", err)
        return
    }
	//Creates file
	fileName = BuildCSVFileName()
	file, err := os.Create(fileName)
    if err != nil {
        err = fmt.Errorf("Cannot create file: %s", err)
        return
    } else {
    	fmt.Printf("Writing %d users to file.\n", len(*users))
    }
    //close the file when done
    defer file.Close()

    //Writes CSV header
    _, err = fmt.Fprintln(file, "id;firstName;lastName;email")
    if err != nil{
    	return
    }
    //Writes all users to file
    for _, user := range *users {
		_, err = fmt.Fprintln(file, user.ID + ";" + user.firstName + ";" + user.lastName + ";" + user.email)
		if err != nil{
    		return
    	}
	}
	return
}

//Reads users from file created in current session
func ReadUsersFromFile() (users []User, err error){
	fmt.Println("Reading users from file.")
	//open file that was created this session
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	//create reader
	reader := csv.NewReader(f)
	reader.Comma = ';'
	//handle header
	_, err = reader.Read()
	if err != nil {
		return
	}
	//iterate over lines and create users
	for {   
			//read line, break if end of file.
            line, err := reader.Read()
            if err != nil {
                if err == io.EOF {
                    break
                }
                log.Fatal(err)
            }
            //create user
        	user := NewUser(line[0], line[1], line[2], line[3])
			
			//append to array of users
			users = append(users, user)    
        }
	
	//Set error if no users in file
	if len(users)==0{
		err = fmt.Errorf("No users in file.\n")
	}
	return
}

//Builds filename based on execution directory and time of creation
func BuildCSVFileName() string {
	ex, err := os.Executable()
    if err != nil {
        log.Fatal(err)
    }
    exPath := filepath.Dir(ex)
	return filepath.Join(exPath, "DeletedUserData", "DeletedUsers_" + time.Now().Format("2006-01-02T150405") + ".csv")
}

func DeleteFile(){
	err := os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
}
