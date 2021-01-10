package schema

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
)

type ScoreData struct {
	EventId      string
	ActiveRounds []RoundData `json:"active_rounds"`
}

func (s *ScoreData) ToDynamoDbMap() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"Event_ID": &dynamodb.AttributeValue{S: &s.EventId},
	}
}

type RoundData struct {
	EventId string
	RoundId string            `json:"round_id"`
	Mode    string            `json:"mode"`
	Players []PlayerRoundData `json:"players"`
}

func (r *RoundData) ToDynamoDbMap() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"Event_ID": &dynamodb.AttributeValue{S: &r.EventId},
		"Round_ID": &dynamodb.AttributeValue{S: &r.RoundId},
		"Mode":     &dynamodb.AttributeValue{S: &r.Mode},
	}
}

type PlayerRoundData struct {
	Name          string `json:"name"`
	Rank          int    `json:"rank"`
	Score         int    `json:"score"`
	SessionStatus string `json:"session_status"`
}

type PlayerData struct {
	EventId  string
	PlayerId string
	Name     string
}

func (p *PlayerData) ToDynamoDbMap() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"Event_ID": &dynamodb.AttributeValue{S: &p.EventId},
		"Round_ID": &dynamodb.AttributeValue{S: &p.PlayerId},
		"Name":     &dynamodb.AttributeValue{S: &p.Name},
	}
}

type ResultData struct {
	RoundId  string
	PlayerId string
	Status   string
	Rank     int
	Score    int
}

func (r *ResultData) ToDynamoDbMap() map[string]*dynamodb.AttributeValue {
	rank := strconv.Itoa(r.Rank)
	score := strconv.Itoa(r.Score)

	return map[string]*dynamodb.AttributeValue{
		"Round_ID":  &dynamodb.AttributeValue{S: &r.RoundId},
		"Player_ID": &dynamodb.AttributeValue{S: &r.PlayerId},
		"Status":    &dynamodb.AttributeValue{S: &r.Status},
		"Rank":      &dynamodb.AttributeValue{S: &rank},
		"Score":     &dynamodb.AttributeValue{S: &score},
	}
}
