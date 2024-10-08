package test

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)


type NameStruct struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
}


func TestMongo(t *testing.T) {
	session,err := mgo.DialWithTimeout("mongodb://10.160.84.104:27017,10.160.84.105:27017,10.160.84.106:27017/JanusTest?replicaSet=testdb",
		time.Duration(300000000))

	if err != nil{
		return
	}

	ts := NameStruct{
		Id:bson.ObjectIdHex("604205b6ac363e197be1f3e9"),
		Name: "",
	}
	tsBytes,_ := json.Marshal(ts)
	t.Log(string(tsBytes))

	err = session.DB("JanusTest").C("junhu.li").Update(bson.M{"_id":bson.ObjectIdHex("604205b6ac363e197be1f3e9")},ts)

}
