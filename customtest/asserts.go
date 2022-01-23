package customtest

import (
	"reflect"
	"testing"
)

func assertContains(t *testing.T, expected, actual interface{}) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		tLog(t, "value type not match %T (expected)\n\n\t != %T (actual)", expected, actual)

	}
}
