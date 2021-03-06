package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type DeleteTemplates struct {
	deleter     services.TemplateDeleterInterface
	errorWriter ErrorWriterInterface
}

func NewDeleteTemplates(deleter services.TemplateDeleterInterface, errorWriter ErrorWriterInterface) DeleteTemplates {
	return DeleteTemplates{
		deleter:     deleter,
		errorWriter: errorWriter,
	}
}

func (handler DeleteTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, stack stack.Context) {
	templateID := strings.Split(req.URL.Path, "/templates/")[1]

	err := handler.deleter.Delete(templateID)
	if err != nil {
		handler.errorWriter.Write(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
