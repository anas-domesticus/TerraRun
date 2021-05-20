package internal

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
)

//go:embed templates/report.tpl.html
var reportTemplate string

type ShowOutputSet []ShowOutput

func (sos *ShowOutputSet) GenerateHTMLReport() []byte {
	t, err := template.New("report").
		Funcs(template.FuncMap{"getJSONString": GetJSONString}).
		Parse(reportTemplate)
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

func GetJSONString(input interface{}) string {
	data, err := json.MarshalIndent(input, "", "\t")
	if err != nil {
		return ""
	}
	return string(data)
}
