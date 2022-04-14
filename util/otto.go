package util

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

func OttoRun(s string) (otto.Value, error) {
	return otto.New().Run(s)
}

// JSObjectToJSON 将js对象转为json
func JSObjectToJSON(s string) ([]byte, error) {
	vm := otto.New()
	v, err := vm.Run(fmt.Sprintf(`
		cs = %s
		JSON.stringify(cs)
`, s))
	if err != nil {
		return nil, err
	}
	return StringToBytes(v.String()), nil
}
