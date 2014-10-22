package logsauce

import (
	_ "github.com/go-sql-driver/mysql"
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
		Key:    []string{"Name"},
		Unique: true,
		Sparse: true,
	}

	//Configure hosts collection
	/*err = session.DB("logsauce").C("hosts").EnsureIndex(idIndex)
	if err != nil {
		log.Panic(err)
	}*/

	err = session.DB("logsauce").C("hosts").EnsureIndex(nameIndex)
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < len(config.ServerConfiguration.Hosts); i++ {
		err := session.DB("logsauce").C("hosts").Insert(config.ServerConfiguration.Hosts[i])
		if err != nil {
			log.Println(err)
		}
	}

	//Configure logs collection
	/*err = session.DB("logsauce").C("logs").EnsureIndex(idIndex)
	if err != nil {
		log.Panic(err)
	}*/

	//Configure logtypes collection
	/*err = session.DB("logsauce").C("logtypes").EnsureIndex(idIndex)
	if err != nil {
		log.Panic(err)
	}*/

	err = session.DB("logsauce").C("logtypes").EnsureIndex(nameIndex)
	if err != nil {
		log.Panic(err)
	}

	dbSession = session
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

	log.Println(logline)

	err := logCollection.Insert(&logline)
	if err != nil {
		log.Println(err)
	}
}

func getLoglinesForPeriodForHostnameAndFilepath(hostname, filepath string, startUnix, endUnix int64) ([]LogLine, error) {
	host, err := getHostWithName(hostname)
	if err != nil {
		return []LogLine{}, err
	}

	logCollection := getLogCollection()
	var loglines []LogLine
	//err = logCollection.Find(bson.M{"HostId": host.Id, "Filepath": filepath, "Timestamp": bson.M{"$gte": startUnix, "$lte": endUnix}}).All(&loglines)
	err = logCollection.Find(bson.M{"hostid": host.Id, "filepath": filepath, "timestamp": bson.M{"$gte": startUnix, "$lte": endUnix}}).All(&loglines)
	if err != nil {
		return []LogLine{}, err
	}

	return loglines, nil
}

func getHostWithToken(token string) (Host, error) {
	hostCollection := getHostsCollection()

	var host Host

	err := hostCollection.Find(bson.M{"Token": token}).One(&host)
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

	err := hostCollection.Find(bson.M{"Name": name}).One(&host)
	if err != nil {
		return Host{}, err
	}

	return host, nil
}
