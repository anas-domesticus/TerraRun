package internal

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed templates/report.tpl.html
var reportTemplate string

type ShowOutputSet []ShowOutput

func (sos *ShowOutputSet) GenerateHTMLReport() []byte {
	t, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return nil
	}
	var rendered bytes.Buffer
	err = t.Execute(&rendered, sos)
	if err != nil {
		return nil
	}

	return rendered.Bytes()
}
