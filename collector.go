package tagcollector

import (
	"fmt"
	"reflect"
)

type FieldTagEntry struct {
	FieldName  string
	FieldValue interface{}
	TagMap     map[string]string
}

func Collect(inputStruct interface{}, tagNames []string) (fieldTagEntries []*FieldTagEntry, err error) {
	v := reflect.ValueOf(inputStruct).Elem()
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not struct")
	}

	for i := 0; i < v.NumField(); i++ {
		fieldType := v.Type().Field(i)
		fieldName := fieldType.Name
		fieldValue := v.Field(i)

		tagNameValueMap := make(map[string]string)
		for _, tagName := range tagNames {
			tagValue := fieldType.Tag.Get(tagName)
			tagNameValueMap[tagName] = tagValue
		}

		fieldTagEntries = append(fieldTagEntries, &FieldTagEntry{
			FieldName:  fieldName,
			FieldValue: fieldValue.Interface(),
			TagMap:     tagNameValueMap,
		})
	}

	return fieldTagEntries, nil
}
