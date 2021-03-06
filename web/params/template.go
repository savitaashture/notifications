package params

import (
	"encoding/json"
	"io"
	"text/template"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/valiant"
)

type Template struct {
	Name     string          `json:"name" validate-required:"true"`
	Text     string          `json:"text"`
	HTML     string          `json:"html" validate-required:"true"`
	Subject  string          `json:"subject"`
	Metadata json.RawMessage `json:"metadata"`
}

type TemplateCreateError struct{}

func (err TemplateCreateError) Error() string {
	return "Failed to create Template in the database"
}

func NewTemplate(body io.Reader) (Template, error) {
	var template Template
	validator := valiant.NewValidator(body)

	err := validator.Validate(&template)
	if err != nil {
		switch err.(type) {
		case valiant.RequiredFieldError:
			return template, ValidationError([]string{err.Error()})
		default:
			return template, ParseError{}
		}
	}

	if template.Metadata == nil {
		template.Metadata = json.RawMessage("{}")
	}

	err = template.validateSyntax()
	if err != nil {
		return Template{}, err
	}

	template.setDefaults()

	return template, nil
}

func (t Template) validateSyntax() error {
	toValidate := map[string]string{
		"Subject": t.Subject,
		"Text":    t.Text,
		"HTML":    t.HTML,
	}

	for field, contents := range toValidate {
		_, err := template.New("test").Parse(contents)
		if err != nil {
			return ValidationError([]string{field + " syntax is malformed please check your braces"})
		}
	}

	return nil
}

func (t Template) ToModel() models.Template {
	return models.Template{
		Name:     t.Name,
		Text:     t.Text,
		HTML:     t.HTML,
		Subject:  t.Subject,
		Metadata: string(t.Metadata),
	}
}

func (t *Template) setDefaults() {
	if t.Subject == "" {
		t.Subject = "{{.Subject}}"
	}
}
