package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type NotifyUAAScope struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    strategies.StrategyInterface
	database    models.DatabaseInterface
}

func NewNotifyUAAScope(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy strategies.StrategyInterface, database models.DatabaseInterface) NotifyUAAScope {
	return NotifyUAAScope{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
		database:    database,
	}
}

func (handler NotifyUAAScope) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := handler.database.Connection()
	err := handler.Execute(w, req, connection, context, handler.strategy)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}
}

func (handler NotifyUAAScope) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface, context stack.Context, strategy strategies.StrategyInterface) error {
	scope := strings.TrimPrefix(req.URL.Path, "/uaa_scopes/")

	output, err := handler.notify.Execute(connection, req, context, scope, strategy, params.GUIDValidator{})
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)

	return nil
}
