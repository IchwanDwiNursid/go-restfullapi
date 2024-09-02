package main

import (
	"database/sql"
	"log"

	"github.com/IchwanDwiNursid/go_restfullapi/cmd/api"
	"github.com/IchwanDwiNursid/go_restfullapi/config"
	"github.com/IchwanDwiNursid/go_restfullapi/db"
	"github.com/go-sql-driver/mysql"
)



func main(){

	db , err := db.NewMySQLStorage(mysql.Config{
		User: config.Envs.DBUser,
		Passwd: config.Envs.DBPassword,
		Addr: config.Envs.DBAddress,
		DBName: config.Envs.DBName,
		AllowNativePasswords: true,
		ParseTime: true,
	})	

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewApiServer(":8080",db)
	
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB){
	err := db.Ping()
	if err != nil{
		log.Fatal(err)
	}

	log.Println("DB : Successfully connected!")
}