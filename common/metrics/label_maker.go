package metrics

import (
	"fmt"
	"reflect"
)

// labelMaker encapsulates logic for creating labels for metrics.
type labelMaker struct {
	keys         []string
	templateType reflect.Type
}

// newLabelMaker creates a new labelMaker instance given a label template. The label template may be nil.
func newLabelMaker(labelTemplate *struct{}) *labelMaker {
	lm := &labelMaker{
		keys: make([]string, 0),
	}

	if labelTemplate == nil {
		return lm
	}

	v := reflect.ValueOf(labelTemplate)
	t := v.Type()
	lm.templateType = t
	for i := 0; i < t.NumField(); i++ {
		lm.keys = append(lm.keys, t.Field(i).Name)
	}

	return lm
}

// getKeys provides the keys for the label struct.
func (l *labelMaker) getKeys() []string {
	return []string{}
}

// extractValues extracts the values from the given label struct.
func (l *labelMaker) extractValues(label *struct{}) ([]string, error) {
	values := make([]string, 0)

	if label == nil {
		return values, nil
	}

	if l.templateType != reflect.TypeOf(label) {
		return nil, fmt.Errorf(
			"label type mismatch, expected %v, got %v", l.templateType, reflect.TypeOf(label))
	}

	for i := 0; i < l.templateType.NumField(); i++ {
		v := reflect.ValueOf(label)
		values = append(values, v.Field(i).String())
	}

	return values, nil
}
