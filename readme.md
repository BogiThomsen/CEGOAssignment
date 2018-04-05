# CEGO jobinterview assignment
Given the choice between Go, PHP, JS, and Bash, I have chosen Go as the language, as I for quite some time have wanted learn Go, but never got started. 
The database choice fell on MariaDB, as I like MySQL, but after Oracle acquired MySQL they might have a conflict of interrest in developing MySQL and I made the switch. That said, this solution should also work for MySQL, but this hasn't been tested. 

## The solution
The solution is a console application that goes through the following steps:
	1. Connect to the database and run the SQL statement presented in the TOML config file.
	2. Ask for confirmation of deletion of data, so the user can catch any errors in their SQL statement before data is deleted.
	3. Writes results of query to file.
	4. Reads file and, per row, queries the database on id and compares data in file against data in database.
	5. It the data in file is correct, it deletes the data from the database.

## How to use
Fill out the config.toml file, build and run the application.
To change SQL statement or database, just change it in the config file and run again. No need to rebuild.


## Testing
Still to come

## Concerns and future work
This solution will run any SQL statement passed from the config file. As the user also has to supply the username/password for the database, this hasn't been a focus to protect against. 
The solution could be created as a package instead of an application.

