package metrics

import (
	"fmt"
	"reflect"
	"sync"
)

type Metrics struct {
	// Handled by the network layer

	TotalWsConnections    int64 `metrics:"total_ws_connections"`
	TotalWsErrors         int64 `metrics:"total_ws_errors"`
	TotalMessagesSent     int64 `metrics:"total_messages_sent"`
	TotalUnknownCommands  int64 `metrics:"total_unknown_commands"`
	TotalCommandsExecuted int64 `metrics:"total_commands_executed"`
	TotalCommandsFailed   int64 `metrics:"total_commands_failed"`

	// Handled by the game repo layer

	TotalGames      int64 `metrics:"total_games"`
	TotalUsers      int64 `metrics:"total_users"`
	UsersInGames    int64 `metrics:"users_in_games"`
	UsersConnected  int64 `metrics:"users_connected"`
	GamesInProgress int64 `metrics:"games_in_progress"`

	lock sync.Mutex
}

var metrics Metrics

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

func AddCommandExecuted() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.TotalCommandsExecuted++
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

func AddUserInGame() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.UsersInGames++
}

func RemoveUserInGame() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.UsersInGames--
}

func AddUserConnected() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.UsersConnected++
}

func RemoveUserConnected() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.UsersConnected--
}

func AddGameInProgress() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.GamesInProgress++
}

func RemoveGameInProgress() {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.GamesInProgress--
}

func getMetrics() Metrics {
	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	return metrics
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
