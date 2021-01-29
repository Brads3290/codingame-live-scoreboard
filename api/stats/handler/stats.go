package handler

import (
	"codingame-live-scoreboard/api/stats/schema"
	"codingame-live-scoreboard/apishared"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/dbschema"
	"codingame-live-scoreboard/schema/errors"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"sort"
	"strings"
)

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

	// Take in a query string parameter to list what stats we want to return
	fetchListStr, ok := request.QueryStringParameters["fetch"]
	if !ok {
		return 400, nil, errors.New("Missing fetch parameter")
	}

	fetchList := strings.Split(fetchListStr, ",")

	// Dedupe fetch list so that we don't have multiple threads accessing the same field on the
	// StatsResponse object
	fetchList = dedupeStringList(fetchList)

	// Fetch some data that all stats methods have in common
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

	// Request URL:
	// /api/stats/{event_id}?fetch=round,language

	data := schema.StatsResponse{}

	// Fetch stats on different threads
	chans := make([]chan error, 0)
	for _, fetchStr := range fetchList {
		ch := make(chan error)

		switch fetchStr {
		case "round":
			go getRoundStats(eventId, &data, rounds, results, ch)
		case "language":
			go getLangStats(eventId, &data, rounds, results, ch)
		default:
			continue
		}

		chans = append(chans, ch)
	}

	// Join the threads and check for errors
	errs := make([]error, 0)
	for _, c := range chans {
		err := <-c
		if err != nil {
			errs = append(errs, err)
		}
	}

	// If we have any errors, we create a composite error and return that
	if len(errs) > 0 {
		return 500, nil, errors.NewComposite(errs...)
	}

	return 200, data, nil
}

func getRoundStats(eventId uuid.UUID, response *schema.StatsResponse, rounds []dbschema.RoundModel, results []dbschema.ResultModel, chErr chan error) {

	// Also fetch the players for the next part, where we will need the player name
	players, err := ddb.GetAllPlayersInEvent(eventId.String())
	if err != nil {
		chErr <- err
		return
	}

	// Create a map of player ids to player objects
	playerMap := make(map[string]dbschema.PlayerModel)
	for _, p := range players {
		if _, ok := playerMap[p.PlayerId]; !ok {
			playerMap[p.PlayerId] = p
		}
	}

	// Create the individual player stats object
	playerStats := make(map[string]*schema.PlayerRoundStats)
	for _, r := range results {
		if _, ok := playerStats[r.PlayerId]; !ok {
			playerStats[r.PlayerId] = &schema.PlayerRoundStats{
				Username:     playerMap[r.PlayerId].Name,
				RoundsPlayed: 0,
			}
		}

		playerStats[r.PlayerId].RoundsPlayed += 1
	}

	// Create the data and return it
	roundStats := schema.RoundStats{
		NumberOfRounds: len(rounds),
		PlayerStats:    playerStats,
	}

	response.Round = &roundStats
	chErr <- nil
}

func getLangStats(eventId uuid.UUID, response *schema.StatsResponse, rounds []dbschema.RoundModel, results []dbschema.ResultModel, chErr chan error) {
	langStats := schema.LanguageStats{
		Popularity: make([]schema.LanguagePopularityStat, 0),
	}

	languageCounts := make(map[string]int)
	for _, result := range results {
		if result.LanguageUsed == "" {
			continue
		}

		if _, ok := languageCounts[result.LanguageUsed]; !ok {
			languageCounts[result.LanguageUsed] = 0
		}

		languageCounts[result.LanguageUsed] += 1
	}

	type languageEntry struct {
		languageName string
		count        int
	}

	languages := make([]languageEntry, 0)
	for k, v := range languageCounts {
		languages = append(languages, languageEntry{
			languageName: k,
			count:        v,
		})
	}

	sort.Slice(languages, func(i, j int) bool {
		return languages[i].count < languages[j].count
	})

	for i, l := range languages {
		langStats.Popularity = append(langStats.Popularity, schema.LanguagePopularityStat{
			Name: l.languageName,
			Rank: i,
			Uses: l.count,
		})
	}

	response.Language = &langStats
	chErr <- nil
}

func dedupeStringList(ls []string) []string {
	newLs := make([]string, 0)

	for _, v := range ls {
		if !stringListContains(newLs, v) {
			newLs = append(newLs, v)
		}
	}

	return newLs
}

func stringListContains(ls []string, s string) bool {
	for _, v := range ls {
		if v == s {
			return true
		}
	}

	return false
}
