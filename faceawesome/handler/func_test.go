package handler_test

import (
	"fmt"
	"github.com/facedamon/faceawesome/handler"
	"testing"
)

// selfInfo
func selfInfo(k, v interface{}) {
	fmt.Printf("大家好,我叫%s,今年%d岁\n", k, v)
}

func TestHandlerFunc_Do(t *testing.T) {
	m := make(map[interface{}]interface{}, 3)
	m["tom"] = 11
	m["sunny"] = 23
	m["jon"] = 24

	handler.EachFunc(m, selfInfo)
}
