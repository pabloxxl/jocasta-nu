package dns

import (
	"fmt"
	"log"

	"github.com/pabloxxl/jocasta-nu/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/dns/dnsmessage"
)

// Record struct containging single record with corresponding action
type Record struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Action int                `bson:"action"`
	URL    string             `bson:"url"`
	Type   dnsmessage.Type    `bson:"type"`
}

const (
	// ActionUknown from error or not found record1
	ActionNo = -1
	// ActionBlock block given url
	ActionBlock = iota
	// ActionLog log given url
	ActionLog = iota
)

func (r *Record) string() string {
	actionString := ActionToString(r.Action)
	return fmt.Sprintf("rule: %v, url: %v", actionString, r.URL)
}

func RecordToString(record Record) string {
	return record.string()
}

func ActionToString(action int) string {
	actionString := "NO ACTION"

	switch action {
	case ActionBlock:
		actionString = "BLOCK"
	case ActionLog:
		actionString = "LOG"
	}

	return actionString
}

func StringToAction(action string) int {
	actionInt := ActionNo

	switch action {
	case "BLOCK":
		actionInt = ActionBlock
	case "LOG":
		actionInt = ActionLog
	}

	return actionInt
}

// CreateRecordBlock create record structure with type ActionBlock
func CreateRecordBlock(url string) *Record {
	rec := Record{Action: ActionBlock, URL: url}
	return &rec
}

// CreateRecord create record structure
func CreateRecord(url string, action int, recordType dnsmessage.Type) *Record {
	rec := Record{Action: action, URL: url, Type: recordType}
	return &rec
}

func IsRecordEmpty(record Record) bool {
	if record.URL == "" && record.Action == ActionNo {
		return true
	}
	return false
}

func CreateManyRecordsFromDB(client *mongo.Client, key string, value interface{}) *[]Record {
	var records []Record
	var record Record

	recordsFromDB := db.GetAny(client, "records", "", nil)

	for _, value := range recordsFromDB {
		data, err := bson.Marshal(value)
		if err != nil {
			log.Fatal("Failed to marshal data")
		}
		err = bson.Unmarshal(data, &record)
		if err != nil {
			log.Fatal("Failed to unmarshal data")
		}

		records = append(records, record)
	}

	return &records
}

func CreateAllRecordsFromDB(client *mongo.Client) *[]Record {
	return CreateManyRecordsFromDB(client, "", nil)
}

func GetOneRecordFromDB(url string, recordType dnsmessage.Type) Record {
	record := Record{URL: "", Action: ActionNo}
	client := db.CreateClient()
	filter := map[string]interface{}{"url": url, "type": recordType}
	recordFromDB := db.GetOne(client, "records", filter)
	if recordFromDB.Err() != nil {
		return record
	}
	err := recordFromDB.Decode(&record)
	if err != nil {
		log.Fatal(err)
	}
	return record
}

func GetRecord(client *mongo.Client, url string, messageType dnsmessage.Type) int {
	record := GetOneRecordFromDB(url, messageType)

	if !IsRecordEmpty(record) {
		log.Printf("%s %s is marked from database query: %s", url, TypeToString(messageType), ActionToString(record.Action))
	}
	return record.Action
}
