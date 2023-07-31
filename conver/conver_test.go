package conver

import (
	"reflect"
	"testing"
)

func TestConver(t *testing.T) {

	a := IntToString(55)
	t.Log(reflect.TypeOf(a))

}
