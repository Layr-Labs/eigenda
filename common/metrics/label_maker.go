package metrics

import (
	"fmt"
	"reflect"
)

// labelMaker encapsulates logic for creating labels for metrics.
type labelMaker struct {
	keys         []string
	emptyValues  []string
	templateType reflect.Type
	labelCount   int
}

// newLabelMaker creates a new labelMaker instance given a label template. The label template may be nil.
func newLabelMaker(labelTemplate any) (*labelMaker, error) {
	labeler := &labelMaker{
		keys: make([]string, 0),
	}

	if labelTemplate == nil {
		return labeler, nil
	}

	v := reflect.ValueOf(labelTemplate)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("label template must be a struct")
	}

	t := v.Type()
	labeler.templateType = t
	for i := 0; i < t.NumField(); i++ {

		fieldType := t.Field(i).Type
		if fieldType.Kind() != reflect.String {
			return nil, fmt.Errorf(
				"field %s has type %v, only string fields are supported", t.Field(i).Name, fieldType)
		}

		labeler.keys = append(labeler.keys, t.Field(i).Name)
	}

	labeler.emptyValues = make([]string, len(labeler.keys))
	labeler.labelCount = len(labeler.keys)

	return labeler, nil
}

// getKeys provides the keys for the label struct.
func (l *labelMaker) getKeys() []string {
	return l.keys
}

// extractValues extracts the values from the given label struct.
func (l *labelMaker) extractValues(label any) ([]string, error) {
	if l.templateType == nil || label == nil {
		return l.emptyValues, nil
	}

	if l.templateType != reflect.TypeOf(label) {
		return nil, fmt.Errorf(
			"label type mismatch, expected %v, got %v", l.templateType, reflect.TypeOf(label))
	}

	values := make([]string, 0, l.labelCount)
	for i := 0; i < l.labelCount; i++ {
		v := reflect.ValueOf(label)
		values = append(values, v.Field(i).String())
	}

	return values, nil
}
