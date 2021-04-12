package dns

import (
	"fmt"
	"log"

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
	actionString := "unknown"

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
	if action != ActionBlock && action != ActionBlockRegex && action != ActionLog {
		log.Printf("Wrong action provided: %v", action)
		return nil
	}
	rec := Record{Action: action, URL: url}
	return &rec
}
