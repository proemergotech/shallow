package shallow

import (
	"encoding/json"
	"testing"

	"github.com/kr/pretty"
	"github.com/proemergotech/errors/v2"
)

func TestMerge(t *testing.T) {

	for name, data := range map[string]struct {
		current         test
		incoming        string
		want            test
		wantChangedKeys []string
	}{
		"string_change": {
			current:  testData(nil),
			incoming: `{"string":"test2"}`,
			want: testData(func(t test) test {
				t.String = "test2"
				return t
			}),
			wantChangedKeys: []string{"string"},
		},
		"string_null": {
			current:  testData(nil),
			incoming: `{"string":null}`,
			want: testData(func(t test) test {
				t.String = ""
				return t
			}),
			wantChangedKeys: []string{"string"},
		},
		"string_same": {
			current:         testData(nil),
			incoming:        `{"string":"string_val"}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"string_ptr_change": {
			current:  testData(nil),
			incoming: `{"string_ptr":"test2"}`,
			want: testData(func(t test) test {
				t.StringPtr = stringPtr("test2")
				return t
			}),
			wantChangedKeys: []string{"string_ptr"},
		},
		"string_ptr_null": {
			current:  testData(nil),
			incoming: `{"string_ptr":null}`,
			want: testData(func(t test) test {
				t.StringPtr = nil
				return t
			}),
			wantChangedKeys: []string{"string_ptr"},
		},
		"string_ptr_same": {
			current:         testData(nil),
			incoming:        `{"string_ptr":"string_ptr_val"}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"bool_change": {
			current:  testData(nil),
			incoming: `{"bool":false}`,
			want: testData(func(t test) test {
				t.Bool = false
				return t
			}),
			wantChangedKeys: []string{"bool"},
		},
		"bool_null": {
			current:  testData(nil),
			incoming: `{"bool":null}`,
			want: testData(func(t test) test {
				t.Bool = false
				return t
			}),
			wantChangedKeys: []string{"bool"},
		},
		"bool_same": {
			current:         testData(nil),
			incoming:        `{"bool":true}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"bool_ptr_change": {
			current:  testData(nil),
			incoming: `{"bool_ptr":false}`,
			want: testData(func(t test) test {
				t.BoolPtr = boolPtr(false)
				return t
			}),
			wantChangedKeys: []string{"bool_ptr"},
		},
		"bool_ptr_null": {
			current:  testData(nil),
			incoming: `{"bool_ptr":null}`,
			want: testData(func(t test) test {
				t.BoolPtr = nil
				return t
			}),
			wantChangedKeys: []string{"bool_ptr"},
		},
		"bool_ptr_same": {
			current:         testData(nil),
			incoming:        `{"bool_ptr":true}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"nested_change": {
			current:  testData(nil),
			incoming: `{"nested":{"string":"test2"}}`,
			want: testData(func(t test) test {
				t.Nested = Nested{String: "test2"}
				return t
			}),
			wantChangedKeys: []string{"nested"},
		},
		"nested_null": {
			current:  testData(nil),
			incoming: `{"nested":null}`,
			want: testData(func(t test) test {
				t.Nested = Nested{}
				return t
			}),
			wantChangedKeys: []string{"nested"},
		},
		"nested_same": {
			current:         testData(nil),
			incoming:        `{"nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"nested_ptr_change": {
			current:  testData(nil),
			incoming: `{"nested_ptr":{"string":"test2"}}`,
			want: testData(func(t test) test {
				t.NestedPtr = &Nested{String: "test2"}
				return t
			}),
			wantChangedKeys: []string{"nested_ptr"},
		},
		"nested_ptr_null": {
			current:  testData(nil),
			incoming: `{"nested_ptr":null}`,
			want: testData(func(t test) test {
				t.NestedPtr = nil
				return t
			}),
			wantChangedKeys: []string{"nested_ptr"},
		},
		"nested_ptr_same": {
			current:         testData(nil),
			incoming:        `{"nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_string_change": {
			current:  testData(nil),
			incoming: `{"anonym_string":"test2"}`,
			want: testData(func(t test) test {
				t.AnonymString = "test2"
				return t
			}),
			wantChangedKeys: []string{"anonym_string"},
		},
		"anonym_string_null": {
			current:  testData(nil),
			incoming: `{"anonym_string":null}`,
			want: testData(func(t test) test {
				t.AnonymString = ""
				return t
			}),
			wantChangedKeys: []string{"anonym_string"},
		},
		"anonym_string_same": {
			current:         testData(nil),
			incoming:        `{"anonym_string":"anonym_string_val"}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_nested_change": {
			current:  testData(nil),
			incoming: `{"anonym_nested":{"string":"test2"}}`,
			want: testData(func(t test) test {
				t.AnonymNested = Nested{String: "test2"}
				return t
			}),
			wantChangedKeys: []string{"anonym_nested"},
		},
		"anonym_nested_null": {
			current:  testData(nil),
			incoming: `{"anonym_nested":null}`,
			want: testData(func(t test) test {
				t.AnonymNested = Nested{}
				return t
			}),
			wantChangedKeys: []string{"anonym_nested"},
		},
		"anonym_nested_same": {
			current:         testData(nil),
			incoming:        `{"anonym_nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_nested_ptr_change": {
			current:  testData(nil),
			incoming: `{"anonym_nested_ptr":{"string":"test2"}}`,
			want: testData(func(t test) test {
				t.AnonymNestedPtr = &Nested{String: "test2"}
				return t
			}),
			wantChangedKeys: []string{"anonym_nested_ptr"},
		},
		"anonym_nested_ptr_null": {
			current:  testData(nil),
			incoming: `{"anonym_nested_ptr":null}`,
			want: testData(func(t test) test {
				t.AnonymNestedPtr = nil
				return t
			}),
			wantChangedKeys: []string{"anonym_nested_ptr"},
		},
		"anonym_nested_ptr_same": {
			current:         testData(nil),
			incoming:        `{"anonym_nested_ptr":{"string":"nested_ptr_string_val","string_ptr":"nested_ptr_string_ptr_val","bool":true,"bool_ptr":true}}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_ptr_string_change": {
			current:  testData(nil),
			incoming: `{"anonym_ptr_string":"test2"}`,
			want: testData(func(t test) test {
				t.AnonymPtrString = "test2"
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr_string"},
		},
		"anonym_ptr_string_null": {
			current:  testData(nil),
			incoming: `{"anonym_ptr_string":null}`,
			want: testData(func(t test) test {
				t.AnonymPtrString = ""
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr_string"},
		},
		"anonym_ptr_string_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_string":"anonym_ptr_string_val"}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_ptr_nested_change": {
			current:  testData(nil),
			incoming: `{"anonym_ptr_nested":{"string":"test2"}}`,
			want: testData(func(t test) test {
				t.AnonymPtrNested = Nested{String: "test2"}
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr_nested"},
		},
		"anonym_ptr_nested_null": {
			current:  testData(nil),
			incoming: `{"anonym_ptr_nested":null}`,
			want: testData(func(t test) test {
				t.AnonymPtrNested = Nested{}
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr_nested"},
		},
		"anonym_ptr_nested_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_ptr_nested_ptr_change": {
			current:  testData(nil),
			incoming: `{"anonym_ptr_nested_ptr":{"string":"test2"}}`,
			want: testData(func(t test) test {
				t.AnonymPtrNestedPtr = &Nested{String: "test2"}
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr_nested_ptr"},
		},
		"anonym_ptr_nested_ptr_null": {
			current:  testData(nil),
			incoming: `{"anonym_ptr_nested_ptr":null}`,
			want: testData(func(t test) test {
				t.AnonymPtrNestedPtr = nil
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr_nested_ptr"},
		},
		"anonym_ptr_nested_ptr_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr_nested_ptr":{"string":"nested_ptr_string_val","string_ptr":"nested_ptr_string_ptr_val","bool":true,"bool_ptr":true}}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_ptr2_string_change": {
			current:  testData(nil),
			incoming: `{"anonym_ptr2_string":"test2"}`,
			want: testData(func(t test) test {
				t.AnonymPtr2String = "test2"
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr2_string"},
		},
		"anonym_ptr2_string_null": {
			current:  testData(nil),
			incoming: `{"anonym_ptr2_string":null}`,
			want: testData(func(t test) test {
				t.AnonymPtr2String = ""
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr2_string"},
		},
		"anonym_ptr2_string_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_string":"anonym_ptr_string_val"}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_ptr2_nested_change": {
			current:  testData(nil),
			incoming: `{"anonym_ptr2_nested":{"string":"test2"}}`,
			want: testData(func(t test) test {
				t.AnonymPtr2Nested = Nested{String: "test2"}
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr2_nested"},
		},
		"anonym_ptr2_nested_null": {
			current:  testData(nil),
			incoming: `{"anonym_ptr2_nested":null}`,
			want: testData(func(t test) test {
				t.AnonymPtr2Nested = Nested{}
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr2_nested"},
		},
		"anonym_ptr2_nested_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_nested":{"string":"nested_string_val","string_ptr":"nested_string_ptr_val","bool":true,"bool_ptr":true}}`,
			want:            testData(nil),
			wantChangedKeys: []string{},
		},
		"anonym_ptr2_nested_ptr_change": {
			current:  testData(nil),
			incoming: `{"anonym_ptr2_nested_ptr":{"string":"test2"}}`,
			want: testData(func(t test) test {
				t.AnonymPtr2NestedPtr = &Nested{String: "test2"}
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr2_nested_ptr"},
		},
		"anonym_ptr2_nested_ptr_null": {
			current:  testData(nil),
			incoming: `{"anonym_ptr2_nested_ptr":null}`,
			want: testData(func(t test) test {
				t.AnonymPtr2NestedPtr = nil
				return t
			}),
			wantChangedKeys: []string{"anonym_ptr2_nested_ptr"},
		},
		"anonym_ptr2_nested_ptr_same": {
			current:         testData(nil),
			incoming:        `{"anonym_ptr2_nested_ptr":{"string":"nested_ptr_string_val","string_ptr":"nested_ptr_string_ptr_val","bool":true,"bool_ptr":true}}`,
			want:            testData(nil),
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
		gotChangedKeys, err := Merge(&orig, &update, keys)
		if err != nil {
			t.Fatalf("%+v", errors.WithStack(err))
		}
		got := orig

		if diff := pretty.Diff(data.want, got); len(diff) > 0 {
			t.Errorf("%v: diffs (want/got): %v", name, pretty.Diff(data.want, got))
		}

		if diff := pretty.Diff(data.wantChangedKeys, gotChangedKeys); len(diff) > 0 {
			t.Errorf("%v changedKeys: diffs (want/got): %v", name, pretty.Diff(data.wantChangedKeys, gotChangedKeys))
		}
	}
}
