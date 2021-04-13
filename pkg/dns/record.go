package dns

import (
	"fmt"
	"log"

	"github.com/pabloxxl/jocasta-nu/pkg/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Record struct containging single record with corresponding action
type Record struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Action int                `bson:"action"`
	URL    string             `bson:"url"`
}

const (
	// ActionBlock block given url
	ActionBlock = iota
	// ActionBlockRegex block given url match
	ActionBlockRegex = iota
	// ActionLog log given url
	ActionLog = iota
)

func (r *Record) string() string {
	actionString := ActionToString(r.Action)
	return fmt.Sprintf("rule: %v, url: %v", actionString, r.URL)
}

func ActionToString(action int) string {
	actionString := "UNKNOWN"

	switch action {
	case ActionBlock:
		actionString = "BLOCK"
	case ActionBlockRegex:
		actionString = "BLOCK_REGEX"
	case ActionLog:
		actionString = "LOG"
	}

	return actionString
}

// CreateRecordBlock create record structure with type ActionBlock
func CreateRecordBlock(url string) *Record {
	rec := Record{Action: ActionBlock, URL: url}
	return &rec
}

// CreateRecord create record structure
func CreateRecord(url string, action int) *Record {
	rec := Record{Action: action, URL: url}
	return &rec
}

func CreateManyRecordsFromDB(key string, value interface{}) *[]Record {

	var records []Record
	client := db.CreateClient()

	recordsFromDB := db.GetAny(client, "records", "", nil)

	for _, value := range recordsFromDB {
		// This probably will not be needed in final version and clean database
		if _, ok := value["action"]; !ok {
			log.Printf("No action found for %s; discarding record", value["url"])
			continue
		}

		actionInt := int(value["action"].(int32))

		records = append(records, *CreateRecord(value["url"].(string), actionInt))
	}

	return &records
}

func CreateAllRecordsFromDB() *[]Record {
	return CreateManyRecordsFromDB("", nil)
}

// TODO add CreateOneRecordFromDB
