package customtest

import (
	"reflect"
	"strings"
	"testing"
)

func assertContains(t *testing.T, expected, actual interface{}) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		tLog(t, "value type not match %T (expected)\n\n\t != %T (actual)", expected, actual)
	}
	switch e := expected.(type) {
	case string:
		a := actual.(string)
		if !strings.Contains(a, e) {
			tLog(t, "   %v (expected)\n\n\t!= %v (actual)",
				expected, actual)
			t.FailNow()
		}
	default:
		tLog(t, "unsupported type %T(expected)", expected)
		t.FailNow()
	}
}
