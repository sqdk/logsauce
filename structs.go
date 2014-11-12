package logsauce

import (
	"gopkg.in/mgo.v2/bson"
)

type LogLine struct {
	Id              bson.ObjectId `bson:"_id,omitempty" json:",omitempty"`
	Line            string
	Timestamp       int64
	DataTimestamp   int64
	HostId          bson.ObjectId `bson:",omitempty" json:",omitempty"`
	TypeId          bson.ObjectId `bson:",omitempty" json:",omitempty"`
	Filename        string
	Filepath        string
	Alias           string            `bson:",omitempty" json:",omitempty"`
	TokenizedObject map[string]string `bson:",omitempty" json:",omitempty"`
}

type LogType struct {
	Id            bson.ObjectId `bson:"_id,omitempty"`
	Name          string        `bson:",omitempty" json:",omitempty"`
	Description   string        `bson:",omitempty" json:",omitempty"`
	Pattern       string        `bson:",omitempty" json:",omitempty"`
	DateFormat    string        `bson:",omitempty" json:",omitempty"`
	DateFieldname string        `bson:",omitempty" json:",omitempty"`
}

type Host struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Token string
	Name  string
}

type ComputeRequest struct {
	Operation   string `json:"op,omitempty"`
	Parameter1  string `json:"p1,omitempty"`
	Parameter2  string `json:"p2,omitempty"`
	Parameter3  string `json:"p3,omitempty"`
	TimeStart   int64  `json:"t0,omitempty"`
	TimeEnd     int64  `json:"t1,omitempty"`
	Filename    string `json:"f,omitempty"`
	Host        string `json:"h,omitempty"`
	LogtypeName string `json:"l,omitempty"`
}

type ComputeResponse struct {
	Values       map[string]int `json:"v,omitempty"`
	TimeStart    int64          `json:"t,omitempty"`
	IntervalSize int64          `json:"is,omitempty"`
	FieldName    string         `json:"fn,omitempty"`
	Filename     string         `json:"f,omitempty"`
	HostId       bson.ObjectId  `bson:",omitempty" json:"h,omitempty"`
	LogtypeName  string         `json:"l,omitempty"`
}
