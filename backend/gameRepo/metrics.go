package gameRepo

import (
	"fmt"
	"reflect"
	"sync"
)

type Metrics struct {
	// Handled by the network layer

	TotalWsConnections    int `metrics:"total_ws_connections"`
	TotalWsErrors         int `metrics:"total_ws_errors"`
	TotalMessagesSent     int `metrics:"total_messages_sent"`
	TotalUnknownCommands  int `metrics:"total_unknown_commands"`
  TotalCommandDuraion int `metrics:"total_command_duration"`
	TotalCommandsExecuted int `metrics:"total_commands_executed"`
	TotalCommandsFailed   int `metrics:"total_commands_failed"`

	// Handled by the game repo layer

	TotalGames             int `metrics:"total_games"`
	TotalUsers             int `metrics:"total_users"`
	UsersInGames           int `metrics:"users_in_games"`
	UsersConnected         int `metrics:"users_connected"`
	GamesInProgress        int `metrics:"games_in_progress"`
	TotalGamePurges        int `metrics:"total_game_purges"`
	TotalGamePurgeDuration int `metrics:"total_game_purge_duration"`
	TotalGamesPurged       int `metrics:"total_games_purged"`

	lock sync.Mutex
}

var metrics Metrics

func AddGamePurgeData(duration int, gamesPurged int) {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()

	metrics.TotalGamePurges++
	metrics.TotalGamePurgeDuration += duration
	metrics.TotalGamesPurged += gamesPurged
}

func AddGame() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalGames++
}

func AddUser() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalUsers++
}

func AddWsConnection() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalWsConnections++
}

func RemoveWsConnection() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalWsConnections--
}

func AddWsError() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalWsErrors++
}

func AddCommandExecuted(duration int) {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalCommandsExecuted++
  metrics.TotalCommandDuraion += duration
}

func AddMessageSent() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalMessagesSent++
}

func AddUnknownCommand() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalUnknownCommands++
}

func AddCommandFailed() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalCommandsFailed++
}

func getMetrics() Metrics {
	// safe copy of the metrics
	metrics.lock.Lock()
	ret := metrics
	metrics.lock.Unlock()

	games := Repo.GetGames()
	metrics.GamesInProgress = len(games)
	metrics.UsersInGames = 0
	metrics.UsersConnected = 0

	for _, game := range games {
		gameMetrics := game.Metrics()
		metrics.UsersConnected += gameMetrics.PlayersConnected
		metrics.UsersInGames += gameMetrics.Players
	}

	return ret
}

func GetMetrics() string {
	metrics := getMetrics()
	ret := ""

	elem := reflect.TypeOf(metrics)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)

		tag := field.Tag.Get("metrics")
		if tag == "" {
			continue
		}
		value := reflect.ValueOf(metrics).FieldByName(field.Name).Int()
		ret += fmt.Sprintf("%s %d\n", tag, value)
	}

	return ret
}
