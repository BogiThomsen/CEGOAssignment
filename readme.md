# CEGO jobinterview assignment
Given the choice between Go, PHP, JS, and Bash, I have chosen Go as the language, as I for quite some time have wanted learn Go, but never got started. <br />
The database choice fell on MariaDB as I chose to learn a new language for this task, and haven't worked with CockroachDB or SQLite, I chose something familiar. This solution should also work for MySQL, but this hasn't been tested. <br />

## The solution
The solution is a console application that goes through the following steps: <br />
	1. Connect to the database and processes and runs the SQL statement presented in the TOML config file. <br />
	2. Ask for confirmation of deletion of data, so the user can catch any errors in their SQL statement before data is deleted. <br />
	3. Writes results of query to file. <br />
	4. Reads file and, per row, queries the database on id and compares data in file against data in database. <br />
	5. If the data in file is correct, it deletes the data from the database. <br /><br />
This solution will match the SQL statement in the config file against: <br />
`"SELECT id, firstName, lastName, email FROM users"`, and additional statements after first `;` are removed before the statement is run.

## How to use
First get needed packages: <br />
`go get gopkg.in/DATA-DOG/go-sqlmock.v1` <br />
`go get github.com/BurntSushi/toml` <br />
`go get github.com/go-sql-driver/mysql` <br />

Fill out the config.toml file, build and run the application. <br />
To change SQL statement or database, just change it in the config file and run again. No need to rebuild. <br />


## Testing
Run tests with `go test -race`

## Concerns
The solution might be voulnerable to SQL injection<br />

## Future work
The solution could be created as a package instead of an application. <br />
Update to take DSN as input instead of TOML config, and use https://godoc.org/github.com/go-sql-driver/mysql#Config

