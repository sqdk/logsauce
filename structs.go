package logsauce

import (
	"github.com/coopernurse/gorp"
	"gopkg.in/mgo.v2/bson"
)

type LogLine struct {
	Id            bson.ObjectId `bson:"_id,omitempty"`
	Line          string
	Timestamp     int64
	HostId        bson.ObjectId `bson:",omitempty"`
	TypeId        bson.ObjectId `bson:",omitempty"`
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
	Id          bson.ObjectId `bson:"_id"`
	Name        string
	Description string
	Pattern     string
	DateFormat  string
}

type Host struct {
	Id    bson.ObjectId `bson:"_id"`
	Token string
	Name  string
}
