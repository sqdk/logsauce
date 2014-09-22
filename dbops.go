package logsauce

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var dbMap *gorp.DbMap

func InitializeDB(config Configuration) {
	connectionString := fmt.Sprintf("%v:%v@tcp(%v)/%v",
		config.ServerConfiguration.DbUsername,
		config.ServerConfiguration.DbPassword,
		config.ServerConfiguration.DbAddress,
		config.ServerConfiguration.DbName)

	log.Println(connectionString)

	dbConnection, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Panicln(err)
	}

	dbM := &gorp.DbMap{Db: dbConnection, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8"}}

	dbM.AddTableWithName(Host{}, "hosts").SetKeys(true, "Id")
	dbM.AddTableWithName(LogLine{}, "loglines").SetKeys(true, "Id")

	err = dbM.CreateTablesIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	dbMap = dbM
}

func getDbMap() *gorp.DbMap {
	return dbMap
}

func getAllHosts() ([]Host, error) {
	dbMap := getDbMap()

	var hosts []Host
	_, err := dbMap.Select(&hosts, "SELECT * FROM hosts")
	if err != nil {
		log.Println(err)
		return hosts, err
	}

	return hosts, nil
}

func insertLogline(logline LogLine) error {
	dbMap := getDbMap()
	err := dbMap.Insert(&logline)
	return err
}

func getLoglinesForHost(hostId int64) {
	var lines []LogLine
	dbMap := getDbMap()
	_, err := dbMap.Select(&lines, "SELECT * FROM loglines WHERE HostId = ?", hostId)
	if err != nil {
		log.Println(err)
		return lines, err
	}

	return lines, nil
}

func getLoglinesForFileForHost(hostId int64, filepath string) {
	var lines []LogLine
	dbMap := getDbMap()
	_, err := dbMap.Select(&lines, "SELECT * FROM loglines WHERE HostId = ? AND Filepath = ?", hostId, filepath)
	if err != nil {
		log.Println(err)
		return lines, err
	}

	return lines, nil
}
