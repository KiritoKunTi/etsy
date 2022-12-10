package db

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type DBConfig struct {
	DBUsername string
	DBPassword string
	DBName string
	DBHost string
	DBPort int
}

var config DBConfig

var DB *sql.DB

func init(){
	parsingDBConfig()
	initDB()
}

func parsingDBConfig(){
	config = DBConfig{
		DBUsername: *flag.String("dbusername", "", "database username of our website"),
		DBPassword: *flag.String("dbpassword", "", "database password of our website"),
		DBName: *flag.String("dbname", "", "database name of our website"),
		DBHost: *flag.String("dbhost", "", "database host of our website"),
		DBPort: *flag.Int("dbport", 5432, "database port of our website"),
	}
}

func initDB(){
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUsername, config.DBPassword, config.DBName,
	)
	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("We have problems with connection database", err)
		os.Exit(1)
	}
	//configDB()
}

func configDB(){
	st, ioErr := ioutil.ReadFile("db/setup.sql")
	if ioErr != nil {
		fmt.Println("Cannot read data/setup.sql")
		os.Exit(1)
	}
	DB.Exec(string(st))
}
