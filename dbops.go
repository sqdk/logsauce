package logsauce

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sqdk/samurai"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

var dbSession *mgo.Session

func InitializeDB(config Configuration) {

	session, err := mgo.Dial(config.ServerConfiguration.DbAddress)
	if err != nil {
		log.Panic(err)
	}

	err = session.Login(&mgo.Credential{Username: config.ServerConfiguration.DbUsername, Password: config.ServerConfiguration.DbPassword})
	if err != nil {
		log.Panic(err)
	}

	/*idIndex := mgo.Index{
		Key:    []string{"Id"},
		Unique: true,
		Sparse: true,
	}*/

	nameIndex := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
		Sparse: true,
	}

	timestampIndex := mgo.Index{
		Key: []string{"timestamp"},
	}

	//Configure hosts collection
	/*err = session.DB("logsauce").C("hosts").EnsureIndex(idIndex)
	if err != nil {
		log.Panic(err)
	}*/

	err = session.DB("logsauce").C("hosts").EnsureIndex(nameIndex)
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < len(config.ServerConfiguration.Hosts); i++ {
		err := session.DB("logsauce").C("hosts").Insert(config.ServerConfiguration.Hosts[i])
		if err != nil {
			log.Println(err)
		}
	}

	//Configure logs collection
	err = session.DB("logsauce").C("logs").EnsureIndex(timestampIndex)
	if err != nil {
		log.Panic(err)
	}

	//Configure logtypes collection
	/*err = session.DB("logsauce").C("logtypes").EnsureIndex(idIndex)
	if err != nil {
		log.Panic(err)
	}*/

	err = session.DB("logsauce").C("logtypes").EnsureIndex(nameIndex)
	if err != nil {
		log.Println(err)
	}

	dbSession = session

	err = insertNewLogType(LogType{Name: "apache", Description: "desc", Pattern: " (ip,nil,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver),code,size)", DateFormat: "02/Jan/2006:15:04:05", DateFieldname: "date"})
	if err != nil {
		log.Println(err)
	}

}

func getHostsCollection() *mgo.Collection {
	if err := dbSession.Ping(); err != nil {
		log.Println(err)
		return nil
	}

	return dbSession.DB("logsauce").C("hosts")
}

func getLogCollection() *mgo.Collection {
	if err := dbSession.Ping(); err != nil {
		log.Println(err)
		return nil
	}

	return dbSession.DB("logsauce").C("logs")
}

func getLogTypeCollection() *mgo.Collection {
	if err := dbSession.Ping(); err != nil {
		log.Println(err)
		return nil
	}

	return dbSession.DB("logsauce").C("logtypes")
}

func insertLogline(logline LogLine) {
	logCollection := getLogCollection()

	logline.Timestamp = time.Now().Unix()

	err := logCollection.Insert(&logline)
	if err != nil {
		log.Println(err)
	}
}

func getLoglinesForPeriodForHostnameAndFilepath(hostname, filepath string, startUnix, endUnix int64) ([]LogLine, error) {
	if startUnix <= -1 || endUnix <= -1 { //Negative timestamp fetches by datatimestamp

	}

	host, err := getHostWithName(hostname)
	if err != nil {
		if err == mgo.ErrNotFound {
			return []LogLine{}, errors.New("Cannot find host in db")
		} else {
			log.Println(err)
		}
	}

	logCollection := getLogCollection()
	var loglines []LogLine
	//err = logCollection.Find(bson.M{"hostid": host.Id, "filepath": filepath, "timestamp": bson.M{"$gte": startUnix, "$lte": endUnix}}).Sort("$orderby : { timestamp : -1 }").All(&loglines)
	log.Println("Fetching loglines")
	ts := time.Now()
	err = logCollection.Find(bson.M{"hostid": host.Id, "filepath": filepath, "timestamp": bson.M{"$gte": startUnix, "$lte": endUnix}}).All(&loglines)
	log.Printf("Fetched %v records in %v", len(loglines), time.Now().Sub(ts))
	if err != nil {
		return []LogLine{}, err
	}

	return loglines, nil
}

func getTokenizedDataForHostnameAndFilepath(patternName, hostname, filepath string, startUnix, endUnix int64) ([]LogLine, error) {
	loglines, err := getLoglinesForPeriodForHostnameAndFilepath(hostname, filepath, startUnix, endUnix)
	logCollection := getLogCollection()
	if err != nil {
		log.Println(err)
		return []LogLine{}, err
	}

	lt, err := getLogTypeWithName(patternName)
	if err != nil {
		return []LogLine{}, err
	}

	if err = samurai.ValidatePattern(lt.Pattern); err != nil {
		return []LogLine{}, err
	}

	if lt.DateFieldname != "" {
		log.Println("Tokenizing and parsing dates")
	} else {
		log.Println("No date fieldname. Only tokenizing.")
	}
	ts := time.Now()
	for i := 0; i < len(loglines); i++ {
		changed := false
		if len(loglines[i].TokenizedObject) == 0 {
			data := samurai.TokenizeBlock(loglines[i].Line, lt.Pattern)
			if data != nil {
				loglines[i].TokenizedObject = data
				changed = true
			}
		}

		if loglines[i].DataTimestamp <= 0 {
			if lt.DateFieldname != "" {
				timestamp := loglines[i].TokenizedObject[lt.DateFieldname]
				timestampTime, err := time.Parse(lt.DateFormat, timestamp)
				if err == nil {
					loglines[i].DataTimestamp = timestampTime.Unix()
					changed = true
				} else {
					log.Println(lt.DateFormat, timestamp)
					log.Println(err)
				}
			}
		}

		if changed {
			err = logCollection.UpdateId(loglines[i].Id, loglines[i])
			if err != nil {
				log.Println(err)
			}
		}
	}
	log.Printf("Decoded %v records in %v", len(loglines), time.Now().Sub(ts))

	return loglines, nil
}

func getHostWithToken(token string) (Host, error) {
	hostCollection := getHostsCollection()

	var host Host

	err := hostCollection.Find(bson.M{"token": token}).One(&host)
	if err != nil {
		return Host{}, err
	}

	return host, nil
}

func getHostWithId(id bson.ObjectId) (Host, error) {
	hostCollection := getHostsCollection()

	var host Host

	err := hostCollection.FindId(id).One(&host)
	if err != nil {
		return Host{}, err
	}

	return host, nil
}

func getHostWithName(name string) (Host, error) {
	hostCollection := getHostsCollection()

	var host Host

	err := hostCollection.Find(bson.M{"name": name}).One(&host)
	if err != nil {
		return Host{}, err
	}

	return host, nil
}

func getLogTypeWithName(name string) (LogType, error) {
	logtypeCollection := getLogTypeCollection()

	var logtype LogType

	err := logtypeCollection.Find(bson.M{"name": name}).One(&logtype)
	if err != nil {
		return LogType{}, err
	}

	return logtype, nil
}

func insertNewLogType(logType LogType) error {
	logtypeCollection := getLogTypeCollection()
	err := logtypeCollection.Insert(&logType)
	return err
}
