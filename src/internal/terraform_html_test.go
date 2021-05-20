package internal

import (
	"encoding/json"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowOutputSet_GenerateHTMLReport(t *testing.T) {
	testJson := `{"format_version":"0.1","terraform_version":"0.13.0","planned_values":{"root_module":{"resources":[{"address":"null_resource.this","mode":"managed","type":"null_resource","name":"this","provider_name":"registry.terraform.io/hashicorp/null","schema_version":0},{"address":"random_string.random","mode":"managed","type":"random_string","name":"random","provider_name":"registry.terraform.io/hashicorp/random","schema_version":1,"values":{"keepers":null,"length":16,"lower":true,"min_lower":0,"min_numeric":0,"min_special":0,"min_upper":0,"number":true,"override_special":"/@£$","special":true,"upper":true}}]}},"resource_changes":[{"address":"null_resource.this","mode":"managed","type":"null_resource","name":"this","provider_name":"registry.terraform.io/hashicorp/null","change":{"actions":["create"],"before":null,"after":{},"after_unknown":{"id":true,"triggers":true}}},{"address":"random_string.random","mode":"managed","type":"random_string","name":"random","provider_name":"registry.terraform.io/hashicorp/random","change":{"actions":["create"],"before":null,"after":{"keepers":null,"length":16,"lower":true,"min_lower":0,"min_numeric":0,"min_special":0,"min_upper":0,"number":true,"override_special":"/@£$","special":true,"upper":true},"after_unknown":{"id":true,"result":true}}}],"configuration":{"provider_config":{"null":{"name":"null"},"random":{"name":"random"}},"root_module":{"resources":[{"address":"null_resource.this","mode":"managed","type":"null_resource","name":"this","provider_config_key":"null","expressions":{"triggers":{"references":["random_string.random"]}},"schema_version":0},{"address":"random_string.random","mode":"managed","type":"random_string","name":"random","provider_config_key":"random","expressions":{"length":{"constant_value":16},"override_special":{"constant_value":"/@£$"},"special":{"constant_value":true}},"schema_version":1}]}}}`
	output := tfjson.Plan{}
	json.Unmarshal([]byte(testJson), &output)
	set := ShowOutputSet{}
	set = append(set, ShowOutput{Stack: TerraformStack{Path: "some_path"}, Plan: output})
	set = append(set, ShowOutput{Stack: TerraformStack{Path: "some_other_path"}, Plan: output})
	html := set.GenerateHTMLReport()
	assert.True(t, len(html) > 0)

}
