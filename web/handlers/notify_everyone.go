package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type NotifyEveryone struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    strategies.StrategyInterface
	database    models.DatabaseInterface
}

func NewNotifyEveryone(notify NotifyInterface, errorWriter ErrorWriterInterface,
	strategy strategies.StrategyInterface, database models.DatabaseInterface) NotifyEveryone {
	return NotifyEveryone{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
		database:    database,
	}
}

func (handler NotifyEveryone) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := handler.database.Connection()
	err := handler.Execute(w, req, connection, context, handler.strategy)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}
}

func (handler NotifyEveryone) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface, context stack.Context,
	strategy strategies.StrategyInterface) error {

	output, err := handler.notify.Execute(connection, req, context, "", strategy, params.GUIDValidator{})
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return nil
}
