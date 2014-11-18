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

func TestFind(t *testing.T) {
	array := reflect.ValueOf([]string{"a", "b", "c"})
	a := reflect.ValueOf("a")
	e := reflect.ValueOf("e")

	// item in array
	index := find(array, a)
	if index != 0 {
		t.Error("Should return item index if it's exist")
	}

	// item not in array
	index = find(array, e)
	if index != -1 {
		t.Error("Should return -1 if item don't exist")
	}
}

func TestRemove(t *testing.T) {
	array := reflect.ValueOf([]string{"a", "b", "c"})
	length := array.Len()
	index1 := 0
	index2 := 3

	// remove item
	v := remove(array, index1)
	if v.Len() != length-1 {
		t.Error("Should remove item with index")
	}

	// output range
	v = remove(array, index2)
	if v.Len() != length {
		t.Error("Should not remove item if index out of range")
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
