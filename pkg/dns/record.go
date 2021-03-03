package dns

import (
	"fmt"
	"log"
)

// Record struct containging single record with corresponding action
type Record struct {
	action int
	url    string
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
	actionString := "unknown"

	switch r.action {
	case ActionBlock:
		actionString = "BLOCK"
	case ActionBlockRegex:
		actionString = "BLOCK_REGEX"
	case ActionLog:
		actionString = "LOG"
	}

	return fmt.Sprintf("rule: %v, url: %v", actionString, r.url)
}

// CreateRecordBlock create record structure with type ActionBlock
func CreateRecordBlock(url string) *Record {
	rec := Record{action: ActionBlock, url: url}
	return &rec
}

// CreateRecord create record structure
func CreateRecord(url string, action int) *Record {
	if action != ActionBlock && action != ActionBlockRegex && action != ActionLog {
		log.Printf("Wrong action provided: %v", action)
		return nil
	}
	rec := Record{action: action, url: url}
	return &rec
}
