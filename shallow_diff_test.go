package shallow

import (
	"encoding/json"
	"testing"

	"github.com/kr/pretty"
	"github.com/proemergotech/errors/v2"
)

func TestDiff(t *testing.T) {

	for name, data := range map[string]struct {
		current         test
		incoming        string
		wantChangedKeys []string
	}{
		"string_change": {
			current:         testData(nil),
			incoming:        `{"string":"test2"}`,
			wantChangedKeys: []string{"string"},
		},
		"string_null": {
			current:         testData(nil),
			incoming:        `{"string":null}`,
			wantChangedKeys: []string{"string"},
		},
		"string_same": {
			current:         testData(nil),
			incoming:        `{"string":"string_val"}`,
			wantChangedKeys: []string{},
		},
		"string_ptr_change": {
			current:         testData(nil),
			incoming:        `{"string_ptr":"test2"}`,
			wantChangedKeys: []string{"string_ptr"},
		},
		"string_ptr_null": {
			current:         testData(nil),
			incoming:        `{"string_ptr":null}`,
			wantChangedKeys: []string{"string_ptr"},
		},
		"string_ptr_same": {
			current:         testData(nil),
			incoming:        `{"string_ptr":"string_ptr_val"}`,
			wantChangedKeys: []string{},
		},
		"bool_change": {
			current:         testData(nil),
			incoming:        `{"bool":false}`,
			wantChangedKeys: []string{"bool"},
		},
		"bool_null": {
			current:         testData(nil),
			incoming:        `{"bool":null}`,
			wantChangedKeys: []string{"bool"},
		},
		"bool_same": {
			current:         testData(nil),
			incoming:        `{"bool":true}`,
			wantChangedKeys: []string{},
		},
		"bool_ptr_change": {
			current:         testData(nil),
			incoming:        `{"bool_ptr":false}`,
			wantChangedKeys: []string{"bool_ptr"},
		},
		"bool_ptr_null": {
			current:         testData(nil),
			incoming:        `{"bool_ptr":null}`,
			wantChangedKeys: []string{"bool_ptr"},
		},
		"bool_ptr_same": {
			current:         testData(nil),
			incoming:        `{"bool_ptr":true}`,
			wantChangedKeys: []string{},
		},
		"nested_change": {
			current:         testData(nil),
			incoming:        `{"nested":{"string":"test2"}}`,
			wantChangedKeys: []string{"nested"},
		},
		"nested_null": {
			current:         testData(nil),
			incoming:        `{"nested":null}`,
			wantChangedKeys: []string{"nested"},
		},
		"nested_same": {
			current:         testData(nil),
			incoming:        `{"nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			wantChangedKeys: []string{},
		},
		"nested_ptr_change": {
			current:         testData(nil),
			incoming:        `{"nested_ptr":{"string":"test2"}}`,
			wantChangedKeys: []string{"nested_ptr"},
		},
		"nested_ptr_null": {
			current:         testData(nil),
			incoming:        `{"nested_ptr":null}`,
			wantChangedKeys: []string{"nested_ptr"},
		},
		"nested_ptr_same": {
			current:         testData(nil),
			incoming:        `{"nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			wantChangedKeys: []string{},
		},
		"anonym_string_change": {
			current:         testData(nil),
			incoming:        `{"anonym_string":"test2"}`,
			wantChangedKeys: []string{"anonym_string"},
		},
		"anonym_string_null": {
			current:         testData(nil),
			incoming:        `{"anonym_string":null}`,
			wantChangedKeys: []string{"anonym_string"},
		},
		"anonym_string_same": {
			current:         testData(nil),
			incoming:        `{"anonym_string":"anonym_string_val"}`,
			wantChangedKeys: []string{},
		},
		"anonym_nested_change": {
			current:         testData(nil),
			incoming:        `{"anonym_nested":{"string":"test2"}}`,
			wantChangedKeys: []string{"anonym_nested"},
		},
		"anonym_nested_null": {
			current:         testData(nil),
			incoming:        `{"anonym_nested":null}`,
			wantChangedKeys: []string{"anonym_nested"},
		},
		"anonym_nested_same": {
			current:         testData(nil),
			incoming:        `{"anonym_nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			wantChangedKeys: []string{},
		},
		"anonym_nested_ptr_change": {
			current:         testData(nil),
			incoming:        `{"anonym_nested_ptr":{"string":"test2"}}`,
			wantChangedKeys: []string{"anonym_nested_ptr"},
		},
		"anonym_nested_ptr_null": {
			current:         testData(nil),
			incoming:        `{"anonym_nested_ptr":null}`,
			wantChangedKeys: []string{"anonym_nested_ptr"},
		},
		"anonym_nested_ptr_same": {
			current:         testData(nil),
			incoming:        `{"anonym_nested_ptr":{"string":"nested_ptr_string_val","string_ptr":"nested_ptr_string_ptr_val","bool":true,"bool_ptr":true}}`,
			wantChangedKeys: []string{},
		},
		"anonym_ptr_string_change": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_string":"test2"}`,
			wantChangedKeys: []string{"anonym_ptr_string"},
		},
		"anonym_ptr_string_null": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_string":null}`,
			wantChangedKeys: []string{"anonym_ptr_string"},
		},
		"anonym_ptr_string_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_string":"anonym_ptr_string_val"}`,
			wantChangedKeys: []string{},
		},
		"anonym_ptr_nested_change": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_nested":{"string":"test2"}}`,
			wantChangedKeys: []string{"anonym_ptr_nested"},
		},
		"anonym_ptr_nested_null": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_nested":null}`,
			wantChangedKeys: []string{"anonym_ptr_nested"},
		},
		"anonym_ptr_nested_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			wantChangedKeys: []string{},
		},
		"anonym_ptr_nested_ptr_change": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_nested_ptr":{"string":"test2"}}`,
			wantChangedKeys: []string{"anonym_ptr_nested_ptr"},
		},
		"anonym_ptr_nested_ptr_null": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_nested_ptr":null}`,
			wantChangedKeys: []string{"anonym_ptr_nested_ptr"},
		},
		"anonym_ptr_nested_ptr_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_nested_ptr":{"string":"nested_ptr_string_val","string_ptr":"nested_ptr_string_ptr_val","bool":true,"bool_ptr":true}}`,
			wantChangedKeys: []string{},
		},
		"anonym_ptr2_string_change": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_string":"test2"}`,
			wantChangedKeys: []string{"anonym_ptr2_string"},
		},
		"anonym_ptr2_string_null": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_string":null}`,
			wantChangedKeys: []string{"anonym_ptr2_string"},
		},
		"anonym_ptr2_string_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_string":"anonym_ptr_string_val"}`,
			wantChangedKeys: []string{},
		},
		"anonym_ptr2_nested_change": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_nested":{"string":"test2"}}`,
			wantChangedKeys: []string{"anonym_ptr2_nested"},
		},
		"anonym_ptr2_nested_null": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_nested":null}`,
			wantChangedKeys: []string{"anonym_ptr2_nested"},
		},
		"anonym_ptr2_nested_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			wantChangedKeys: []string{},
		},
		"anonym_ptr2_nested_ptr_change": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_nested_ptr":{"string":"test2"}}`,
			wantChangedKeys: []string{"anonym_ptr2_nested_ptr"},
		},
		"anonym_ptr2_nested_ptr_null": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_nested_ptr":null}`,
			wantChangedKeys: []string{"anonym_ptr2_nested_ptr"},
		},
		"anonym_ptr2_nested_ptr_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_nested_ptr":{"string":"nested_ptr_string_val","string_ptr":"nested_ptr_string_ptr_val","bool":true,"bool_ptr":true}}`,
			wantChangedKeys: []string{},
		},
	} {
		orig := data.current
		update := test{}
		err := json.Unmarshal([]byte(data.incoming), &update)
		if err != nil {
			t.Fatalf("%+v", errors.WithStack(err))
		}
		var keys map[string]interface{}
		err = json.Unmarshal([]byte(data.incoming), &keys)
		if err != nil {
			t.Fatalf("%+v", errors.WithStack(err))
		}
		gotChangedKeys, err := Diff(&orig, &update, keys)
		if err != nil {
			t.Fatalf("%+v", errors.WithStack(err))
		}

		if diff := pretty.Diff(data.wantChangedKeys, gotChangedKeys); len(diff) > 0 {
			t.Errorf("%v changedKeys: diffs (want/got): %v", name, pretty.Diff(data.wantChangedKeys, gotChangedKeys))
		}
	}
}
