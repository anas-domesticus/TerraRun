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
	Replace       int
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
		var createCount, updateCount, destroyCount, noopCount, replaceCount int
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
			if change.Change.Actions.Replace() {
				replaceCount++
			}

			var before []byte
			var err error
			if change.Change.Before != nil {
				before, err = json.MarshalIndent(change.Change.Before, "", " ")
				if err != nil {
					return nil
				}
			}

			mergedAfter := mergeAfterMap(change.Change.After, change.Change.AfterUnknown)
			var after []byte
			if mergedAfter != nil {
				after, err = json.MarshalIndent(mergedAfter, "", " ")
				if err != nil {
					return nil
				}
			}

			changeDetails = append(changeDetails, ChangeDetail{
				ResourceName: change.Address,
				Before:       string(before),
				After:        string(after),
			})
		}
		outSlice = append(outSlice, HTMLTableData{
			StackName:     v.Stack.Path,
			Create:        createCount,
			Update:        updateCount,
			Replace:       replaceCount,
			Noop:          noopCount,
			Destroy:       destroyCount,
			ChangeDetails: changeDetails,
		})
	}
	return outSlice
}

func mergeAfterMap(known, unknown interface{}) map[string]interface{} {
	var toReturn = make(map[string]interface{})
	switch uk := unknown.(type) {
	case map[string]interface{}:
		switch k := known.(type) {
		case map[string]interface{}:
			toReturn = k
			for i := range uk {
				toReturn[i] = "{{unknown}}"
			}
		default:
			return nil
		}
	default:
		return nil
	}

	return toReturn
}
