package logsauce

import (
	"database/sql"
	"github.com/coopernurse/gorp"
)

type LogLine struct {
	Id            int64
	Line          string
	Timestamp     int64
	HostId        int64
	Filename      string
	Filepath      string
	Alias         string
	TokenizedLine string
}

func (l *LogLine) UpdateTokenizedLineTransactional(transaction gorp.Transaction, string tokenizedLine) {
	l.TokenizedLine = tokenizedLine
	transaction.Update(&l)
}

type LogType struct {
	Id          int64
	Name        string
	Description string
	Pattern     string
	DateFormat  string
}

type Host struct {
	Id    int64
	Token string
	Name  string
}
