package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *sql.DB
func DBConnect() *sql.DB{

	if os.Getenv("CLEARDB_DATABASE_URL") != "" {
		dbDriver:="mysql"
		dbUser:=os.Getenv("DB_USERNAME")
		dbPass:=os.Getenv("DB_PASSWORD")
		dbName:=os.Getenv("DB_NAME")
		Dbhostname:=os.Getenv("DB_HOSTNAME")
		dboption:="?parseTime=true"
		_db,err:=sql.Open(dbDriver,dbUser+":"+dbPass+"@tcp("+Dbhostname+":3306)/"+dbName+dboption)
		if err != nil {
			log.Fatal(err)
		}
		if err= db.Ping();err==nil{
			log.Println("1success")
		}else{
			log.Println("fail")
		}
		db=_db
	}else {
		//local
		godotenv.Load(".env")
		dbDriver:="mysql"
		dbUser:=os.Getenv("DB_USERNAME")
		dbPass:=os.Getenv("DB_PASSWORD")
		dbName:=os.Getenv("DB_NAME")
		dbOption:="?parseTime=true&loc=Asia%2FTokyo"
		dataSource:=dbUser+":"+dbPass+"@tcp(us-cdbr-east-04.cleardb.com:3306)/"+dbName+dbOption
		_db,err:=sql.Open(dbDriver,dataSource)
		if err != nil {
			log.Fatal(err)
		}
		if err= db.Ping();err==nil{
			log.Println("2success")
		}else{
			log.Println("fail")
		}
		db=_db
	}
	return db
}