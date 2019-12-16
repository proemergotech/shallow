package shallow

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

type test struct {
	String    string  `json:"string"`
	StringPtr *string `json:"string_ptr,omitempty"`
	Bool      bool    `json:"bool,omitempty"`
	BoolPtr   *bool   `json:"bool_ptr"`
	Nested    Nested  `json:"nested"`
	NestedPtr *Nested `json:"nested_ptr"`
	Anonym
	*AnonymPtr
}

type Nested struct {
	String    string  `json:"string,omitempty"`
	StringPtr *string `json:"string_ptr"`
	Bool      bool    `json:"bool"`
	BoolPtr   *bool   `json:"bool_ptr"`
}

type Anonym struct {
	AnonymString    string  `json:"anonym_string,omitempty"`
	AnonymStringPtr *string `json:"anonym_string_ptr"`
	AnonymBool      bool    `json:"anonym_bool"`
	AnonymBoolPtr   *bool   `json:"anonym_bool_ptr,omitempty"`
	AnonymNested    Nested  `json:"anonym_nested"`
	AnonymNestedPtr *Nested `json:"anonym_nested_ptr"`
}

type AnonymPtr struct {
	AnonymPtrString    string  `json:"anonym_ptr_string"`
	AnonymPtrStringPtr *string `json:"anonym_ptr_string_ptr"`
	AnonymPtrBool      bool    `json:"anonym_ptr_bool"`
	AnonymPtrBoolPtr   *bool   `json:"anonym_ptr_bool_ptr"`
	AnonymPtrNested    Nested  `json:"anonym_ptr_nested"`
	AnonymPtrNestedPtr *Nested `json:"anonym_ptr_nested_ptr"`
	*AnonymPtr2
}

type AnonymPtr2 struct {
	AnonymPtr2String    string  `json:"anonym_ptr2_string"`
	AnonymPtr2StringPtr *string `json:"anonym_ptr2_string_ptr"`
	AnonymPtr2Bool      bool    `json:"anonym_ptr2_bool"`
	AnonymPtr2BoolPtr   *bool   `json:"anonym_ptr2_bool_ptr"`
	AnonymPtr2Nested    Nested  `json:"anonym_ptr2_nested"`
	AnonymPtr2NestedPtr *Nested `json:"anonym_ptr2_nested_ptr"`
}

func boolPtr(b bool) *bool {
	return &b
}
func stringPtr(str string) *string {
	return &str
}

func testData(modify func(test) test) test {
	data := test{
		String:    "string_val",
		StringPtr: stringPtr("string_ptr_val"),
		Bool:      true,
		BoolPtr:   boolPtr(true),
		Nested: Nested{
			String:    "nested_string_val",
			StringPtr: stringPtr("nested_string_ptr_val"),
			Bool:      true,
			BoolPtr:   boolPtr(true),
		},
		NestedPtr: &Nested{
			String:    "nested_ptr_string_val",
			StringPtr: stringPtr("nested_ptr_string_ptr_val"),
			Bool:      true,
			BoolPtr:   boolPtr(true),
		},
		Anonym: Anonym{
			AnonymString:    "anonym_string_val",
			AnonymStringPtr: stringPtr("anonym_string_ptr_val"),
			AnonymBool:      true,
			AnonymBoolPtr:   boolPtr(true),
			AnonymNested: Nested{
				String:    "nested_string_val",
				StringPtr: stringPtr("nested_string_ptr_val"),
				Bool:      true,
				BoolPtr:   boolPtr(true),
			},
			AnonymNestedPtr: &Nested{
				String:    "nested_ptr_string_val",
				StringPtr: stringPtr("nested_ptr_string_ptr_val"),
				Bool:      true,
				BoolPtr:   boolPtr(true),
			},
		},
		AnonymPtr: &AnonymPtr{
			AnonymPtrString:    "anonym_ptr_string_val",
			AnonymPtrStringPtr: stringPtr("anonym_ptr_string_ptr_val"),
			AnonymPtrBool:      true,
			AnonymPtrBoolPtr:   boolPtr(true),
			AnonymPtrNested: Nested{
				String:    "nested_string_val",
				StringPtr: stringPtr("nested_string_ptr_val"),
				Bool:      true,
				BoolPtr:   boolPtr(true),
			},
			AnonymPtrNestedPtr: &Nested{
				String:    "nested_ptr_string_val",
				StringPtr: stringPtr("nested_ptr_string_ptr_val"),
				Bool:      true,
				BoolPtr:   boolPtr(true),
			},
			AnonymPtr2: &AnonymPtr2{
				AnonymPtr2String:    "anonym_ptr_string_val",
				AnonymPtr2StringPtr: stringPtr("anonym_ptr_string_ptr_val"),
				AnonymPtr2Bool:      true,
				AnonymPtr2BoolPtr:   boolPtr(true),
				AnonymPtr2Nested: Nested{
					String:    "nested_string_val",
					StringPtr: stringPtr("nested_string_ptr_val"),
					Bool:      true,
					BoolPtr:   boolPtr(true),
				},
				AnonymPtr2NestedPtr: &Nested{
					String:    "nested_ptr_string_val",
					StringPtr: stringPtr("nested_ptr_string_ptr_val"),
					Bool:      true,
					BoolPtr:   boolPtr(true),
				},
			},
		},
	}

	if modify != nil {
		data = modify(data)
	}

	return data
}

func TestEmptyKeys(t *testing.T) {
	orig := testData(nil)
	orig.String = "hello world"
	update := testData(nil)
	update.String = "hello test"

	gotChangedKeys, err := Diff(&orig, &update, nil)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	wantedChangedKeys := []string{"string"}
	if diff := pretty.Diff(wantedChangedKeys, gotChangedKeys); len(diff) > 0 {
		t.Errorf("changedKeys: diffs (want/got): %v", pretty.Diff(wantedChangedKeys, gotChangedKeys))
	}
}

func TestUseTag(t *testing.T) {
	type testWithTags struct {
		String1 string `taga:"taga_string1" tagb:"tagb_string1"`
		String2 string `taga:"taga_string2"`
		String3 string `tagb:"tagb_string3"`
	}

	testData := func(modify func(testWithTags) testWithTags) testWithTags {
		data := testWithTags{
			String1: "string1",
			String2: "string2",
			String3: "string3",
		}

		if modify != nil {
			data = modify(data)
		}

		return data
	}

	for name, data := range map[string]struct {
		current         testWithTags
		update          testWithTags
		keys            map[string]interface{}
		tag             string
		want            testWithTags
		wantChangedKeys []string
	}{
		"taga": {
			current: testData(nil),
			update: testWithTags{
				String1: "test1",
				String2: "test2",
			},
			keys: map[string]interface{}{
				"taga_string1": nil,
				"taga_string2": nil,
			},
			tag:             "taga",
			wantChangedKeys: []string{"taga_string1", "taga_string2"},
			want: testData(func(t testWithTags) testWithTags {
				t.String1 = "test1"
				t.String2 = "test2"

				return t
			}),
		},
		"tagb": {
			current: testData(nil),
			update: testWithTags{
				String1: "test1",
				String3: "test3",
			},
			keys: map[string]interface{}{
				"tagb_string1": nil,
				"tagb_string3": nil,
			},
			tag:             "tagb",
			wantChangedKeys: []string{"tagb_string1", "tagb_string3"},
			want: testData(func(t testWithTags) testWithTags {
				t.String1 = "test1"
				t.String3 = "test3"

				return t
			}),
		},
	} {
		gotChangedKeys, err := Merge(&data.current, &data.update, data.keys, UseTag(data.tag))
		if err != nil {
			t.Fatalf("%+v", errors.WithStack(err))
		}
		got := data.current

		if diff := pretty.Diff(data.want, got); len(diff) > 0 {
			t.Errorf("%v: diffs (want/got): %v", name, pretty.Diff(data.want, got))
		}

		if diff := pretty.Diff(data.wantChangedKeys, gotChangedKeys); len(diff) > 0 {
			t.Errorf("%v changedKeys: diffs (want/got): %v", name, pretty.Diff(data.wantChangedKeys, gotChangedKeys))
		}
	}
}
