package exclude

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExclude(t *testing.T) {
	var input interface{}
	var rules interface{}
	var output interface{}

	tests := []string{"enum", "head-asterisk", "tail-asterisk", "tail-asterisk2"}

	for _, test := range tests {
		t.Logf("Executing %s", test)
		input = parse(t, "input.json")
		rules = parse(t, filepath.Join(test, "rule.json"))
		Exclude(input, rules)
		output = parse(t, filepath.Join(test, "output.json"))
		if !reflect.DeepEqual(output, input) {
			t.Error("Result & expected don't match")
			t.Logf("Result %s", stringify(t, input))
			t.Logf("Expected %s", stringify(t, output))
			t.Logf("Rules %s", stringify(t, rules))
		}
	}
}

func parse(t *testing.T, filename string) (data interface{}) {
	file, err := ioutil.ReadFile(filepath.Join("example", filename))
	if err != nil {
		t.Errorf("%s read error", filename, err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		t.Errorf("%s parse error", filename, err)
	}
	return
}

func stringify(t *testing.T, data interface{}) string {
	file, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Errorf("json stringify error", err)
	}
	return string(file)
}
