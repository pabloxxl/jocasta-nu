package dns

import (
	"log"
	"net"
	"time"

	"github.com/pabloxxl/jocasta-nu/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Stat struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Action    string             `bson:"action" json:"action"`
	URL       string             `bson:"url" json:"url"`
	Type      string             `bson:"type" json:"type"`
	IP        string             `bson:"ip" json:"ip"`
	Port      int                `bson:"port" json:"port"`
	Timestamp int64              `bson:"timestamp" json:"timestamp"`
}

type StatCollection struct {
	Number_of_records          int    `json:"number_of_records"`
	Number_of_requests         int    `json:"number_of_requests"`
	Number_of_blocked_requests int    `json:"number_of_blocked_requests"`
	Requests                   []Stat `json:"requests"`
}

func putStat(client *mongo.Client, actionString string, data MessageData, ip net.IP, port int) {
	for _, question := range data.Questions {
		stat := Stat{Action: actionString, URL: question.URL, Type: question.Type.String(), IP: ip.String(), Port: port, Timestamp: time.Now().Unix()}
		db.PutAny(client, "stats", stat)
	}
}

func getAllStats(client *mongo.Client) []Stat {
	var stats []Stat
	var stat Stat

	statsFromDB := db.GetAny(client, "stats", "", nil)

	for _, value := range statsFromDB {
		data, err := bson.Marshal(value)
		if err != nil {
			log.Fatal("Failed to marshal data")
		}
		err = bson.Unmarshal(data, &stat)
		if err != nil {
			log.Fatal("Failed to unmarshal data")
		}
		stats = append(stats, stat)
	}

	return stats
}

func GetStatsCollection(client *mongo.Client) *StatCollection {
	var statCollection StatCollection

	statCollection.Number_of_records = db.CountDocuments(client, "records", "", nil)
	countBlock := db.CountDocuments(client, "stats", "action", ActionToString(ActionBlock))
	countLog := db.CountDocuments(client, "stats", "action", ActionToString(ActionLog))

	statCollection.Number_of_blocked_requests = countBlock
	statCollection.Number_of_requests = countBlock + countLog
	statCollection.Requests = getAllStats(client)

	return &statCollection
}
