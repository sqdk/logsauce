package logsauce

import (
	"github.com/coopernurse/gorp"
	"gopkg.in/mgo.v2/bson"
)

type LogLine struct {
	Id            bson.ObjectId
	Line          string
	Timestamp     int64
	HostId        bson.ObjectId
	TypeId        bson.ObjectId
	Filename      string
	Filepath      string
	Alias         string
	TokenizedLine string
}

func (l *LogLine) UpdateTokenizedLineTransactional(transaction gorp.Transaction, tokenizedLine string) {
	l.TokenizedLine = tokenizedLine
	transaction.Update(&l)
}

type LogType struct {
	Id          bson.ObjectId
	Name        string
	Description string
	Pattern     string
	DateFormat  string
}

type Host struct {
	Id    bson.ObjectId
	Token string
	Name  string
}
