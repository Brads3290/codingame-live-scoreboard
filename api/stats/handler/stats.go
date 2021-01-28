package handler

import (
	"codingame-live-scoreboard/apishared"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/dbschema"
	"codingame-live-scoreboard/schema/errors"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"strings"
)

type statsResponse struct {
	NumberOfRounds int                             `json:"number_of_rounds"`
	PlayerStats    map[string]*playerStatsResponse `json:"player_stats"`
}

type playerStatsResponse struct {
	Username     string `json:"username"`
	RoundsPlayed int    `json:"rounds_played"`
}

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return apishared.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {
		httpMethod := request.RequestContext.HTTP.Method

		switch strings.ToUpper(httpMethod) {
		case "GET":
			return getStats(ctx, request)
		default:
			return 405, nil, errors.New("unsupported HTTP method: " + httpMethod)
		}
	})
}

func getStats(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, interface{}, error) {

	// Get the event ID from the query string and sanity check for existence and validity
	eventIdStr, ok := request.PathParameters["event_id"]
	if !ok {
		return 400, nil, errors.New("missing event_id url param")
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		return 400, nil, err
	}

	// Use the event ID to fetch all the rounds from the DB
	rounds, err := ddb.GetAllRoundsForEvent(eventId.String())
	if err != nil {
		return 500, nil, err
	}

	// Convert to round ID list for fetching the results
	roundsIds := make([]string, len(rounds))
	for i, r := range rounds {
		roundsIds[i] = r.RoundId
	}

	// Fetch the individual results
	results, err := ddb.GetAllResultsForRounds(roundsIds...)
	if err != nil {
		return 500, nil, err
	}

	// Also fetch the players for the next part, where we will need the player name
	players, err := ddb.GetAllPlayersInEvent(eventId.String())
	if err != nil {
		return 500, nil, err
	}

	// Create a map of player ids to player objects
	playerMap := make(map[string]dbschema.PlayerModel)
	for _, p := range players {
		if _, ok := playerMap[p.PlayerId]; !ok {
			playerMap[p.PlayerId] = p
		}
	}

	// Create the individual player stats object
	playerStats := make(map[string]*playerStatsResponse)
	for _, r := range results {
		if _, ok := playerStats[r.PlayerId]; !ok {
			playerStats[r.PlayerId] = &playerStatsResponse{
				Username:     playerMap[r.PlayerId].Name,
				RoundsPlayed: 0,
			}
		}

		playerStats[r.PlayerId].RoundsPlayed += 1
	}

	// Create the data and return it
	data := statsResponse{
		NumberOfRounds: len(rounds),
		PlayerStats:    playerStats,
	}

	return 200, data, nil
}
