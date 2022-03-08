package shallow

import (
	"reflect"
	"strings"

	"github.com/proemergotech/errors/v2"
)

type options struct {
	tag string
}

// Diff compare structs based on the following rule: for every field of the first struct,
// if the field's tag (specified by tag option, default "json") can be found in the keys map,
// compare the corresponding value in the first struct to the value in the second struct.
// If the keys map is nil, all field will be compared.
//
// First and second must be a pointer to a non-nil struct of the same type.
//
// Returns with a list of diff keys. This list can include elements that are NOT actually different if the first struct
// and the second struct had the same value for the given key, and the keys map contained this key.
//
// Traverses anonym fields with struct or struct pointer type, even if they are nested (anonym structs
// within anonym structs). Other types of anonym fields are not supported and will raise an error.
//
// Does NOT check nested fields other than anonym.
func Diff(first interface{}, second interface{}, keys map[string]interface{}, opts ...Option) (diffKeys []string, err error) {
	return process(first, second, keys, false, opts...)
}

// Merge the update struct into the dest struct based on the following rule: for every field of the update struct,
// if the field's tag (specified by tag option, default "json") can be found in the keys map,
// set the corresponding value in the dest struct to the value in the update struct.
//
// Dest and update must be a pointer to a non-nil struct of the same type.
//
// Returns with a list of updated keys. This list can include elements that are NOT actually changed if the dest struct
// and the update struct had the same value for the given key, and the keys map contained this key.
//
// Traverses anonym fields with struct or struct pointer type, even if they are nested (anonym structs
// within anonym structs). Other types of anonym fields are not supported and will raise an error.
//
// Does NOT merge nested fields other than anonym. For these, either the dest value
// is kept intact, or the update value is used, but the two are never merged.
func Merge(dest interface{}, update interface{}, keys map[string]interface{}, opts ...Option) (updatedKeys []string, err error) {
	return process(dest, update, keys, true, opts...)
}

func process(target interface{}, source interface{}, keys map[string]interface{}, merge bool, opts ...Option) (processedKeys []string, err error) {
	targetV := reflect.ValueOf(target)
	sourceV := reflect.ValueOf(source)
	if targetV.Kind() != reflect.Ptr || targetV.Elem().Kind() != reflect.Struct {
		return nil, errors.New("target and source must be a non-nil pointer to a struct with the same type")
	}
	if targetV.Type() != sourceV.Type() {
		return nil, errors.New("target and source must be a non-nil pointer to a struct with the same type")
	}
	if targetV.IsNil() || sourceV.IsNil() {
		return nil, errors.New("target and source must be a non-nil pointer to a struct with the same type")
	}

	o := &options{
		tag: "json",
	}
	for _, opt := range opts {
		opt(o)
	}

	processedKeys = make([]string, 0)
	err = processStructs(targetV.Elem(), sourceV.Elem(), o.tag, keys, &processedKeys, merge)
	if err != nil {
		return nil, err
	}

	return processedKeys, nil
}

func processStructs(targetV reflect.Value, sourceV reflect.Value, tag string, keys map[string]interface{}, processedKeys *[]string, merge bool) error {
	for i := 0; i < sourceV.NumField(); i++ {
		ft := sourceV.Type().Field(i)
		if ft.Anonymous {
			destAVal := targetV.Field(i)
			upAVal := sourceV.Field(i)
			if destAVal.Kind() == reflect.Struct {
				err := processStructs(destAVal, upAVal, tag, keys, processedKeys, merge)
				if err != nil {
					return err
				}
			} else if destAVal.Kind() == reflect.Ptr && destAVal.Elem().Kind() == reflect.Struct {
				if upAVal.IsNil() {
					continue
				}

				if destAVal.IsNil() {
					destAVal.Set(reflect.New(destAVal.Type().Elem()))
				}

				err := processStructs(destAVal.Elem(), upAVal.Elem(), tag, keys, processedKeys, merge)
				if err != nil {
					return err
				}
			} else {
				return errors.New("this method only handles anonym fields of kind struct or pointer to struct")
			}

			continue
		}

		tagVal := ft.Tag.Get(tag)
		tagVal = strings.SplitN(tagVal, ",", 2)[0]
		if tagVal == "" {
			continue
		}
		if keys != nil {
			if _, ok := keys[tagVal]; !ok {
				continue
			}
		}

		if reflect.DeepEqual(targetV.Field(i).Interface(), sourceV.Field(i).Interface()) {
			continue
		}

		*processedKeys = append(*processedKeys, tagVal)
		if merge {
			targetV.Field(i).Set(sourceV.Field(i))
		}
	}

	return nil
}

type Option func(*options)

// UseTag can be used to use struct tags other than json.
func UseTag(tag string) Option {
	return func(o *options) {
		o.tag = tag
	}
}
