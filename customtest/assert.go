package customtest

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tLog(t, "   %v (expected)\n\n\t!= %v (actual)",
			expected, actual)
		t.FailNow()
	}
}

func NotEqual(t *testing.T, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		tLog(t, "value should not equal %#v", actual)
		t.FailNow()
	}
}

func Nil(t *testing.T, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		tLog(t, "value should ne equal %#v", actual)
	}
}

func notNil(t *testing.T, object interface{}) {
	if isNil(object) {
		tLog(t, "Expected value not to be nil", object)
		t.FailNow()
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

func tLog(t *testing.T, fmt string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	targs := make([]interface{}, len(args)+2)
	targs[0] = file
	targs[1] = line
	copy(targs[2:], args)
	t.Logf("\033[31m%s:%d:\n\n\t"+fmt+"\033[39m\n\n", targs...)
}
