package constants

import "time"

// SETTINGS
const (
	MAX_EVENT_RECORD_AGE = 20 * time.Second
)

// DYNAMODB TABLES
const (
	DB_TABLE_EVENTS            = "codezone-codingame-live-scoreboard-events"
	DB_TABLE_ROUNDS            = "codezone-codingame-live-scoreboard-rounds"
	DB_TABLE_RESULTS           = "codezone-codingame-live-scoreboard-results"
	DB_TABLE_XREF_ROUND_PLAYER = "codezone-codingame-live-scoreboard-xref_round_player"
	DB_TABLE_PLAYERS           = "codezone-codingame-live-scoreboard-players"
	DB_TABLE_SETTINGS          = "codezone-codingame-live-scoreboard-settings"
)

// LOGGING
const (
	API_LOGGER_KEY = "ApiEndpoint"
)

// URLS
const (
	CODINGAME_BASE_URL         = "https://www.codingame.com"
	CODINGAME_CLASHREPORT_PATH = "/services/ClashOfCode/findClashReportInfoByHandle"
)