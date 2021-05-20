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
	err = t.Execute(&rendered, sos.buildTableData())
	if err != nil {
		return nil
	}

	return rendered.Bytes()
}

type HTMLTableData struct {
	StackName     string
	Create        int
	Update        int
	Noop          int
	Destroy       int
	ChangeDetails []ChangeDetail
}

type ChangeDetail struct {
	ResourceName string
	Before       string
	After        string
}

func (sos *ShowOutputSet) buildTableData() []HTMLTableData {
	var outSlice []HTMLTableData
	for _, v := range *sos {
		var createCount, updateCount, destroyCount, noopCount int
		var changeDetails []ChangeDetail
		for _, change := range v.Plan.ResourceChanges {
			if change.Change.Actions.Create() {
				createCount++
			}
			if change.Change.Actions.Delete() {
				destroyCount++
			}
			if change.Change.Actions.Update() {
				updateCount++
			}
			if change.Change.Actions.NoOp() {
				noopCount++
			}
			changeDetails = append(changeDetails, ChangeDetail{
				ResourceName: change.Name,
				Before:       "foo",
				After:        "bar",
			})
		}
		outSlice = append(outSlice, HTMLTableData{
			StackName:     v.Stack.Path,
			Create:        createCount,
			Update:        updateCount,
			Noop:          noopCount,
			Destroy:       destroyCount,
			ChangeDetails: changeDetails,
		})
	}
	return outSlice
}

func GetJSONString(input interface{}) string {
	data, err := json.MarshalIndent(input, "", "\t")
	if err != nil {
		return ""
	}
	return string(data)
}
