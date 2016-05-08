package mdl // import "kego.io/editor/client_old/mdl"

// ke: {"package": {"notest": true}}

import (
	"bytes"
	"math/rand"

	"honnef.co/go/js/dom"
)

func randomId() string {
	randInt := func(min int, max int) int {
		return min + rand.Intn(max-min)
	}
	var result bytes.Buffer
	var temp string
	for i := 0; i < 20; {
		if string(randInt(65, 90)) != temp {
			temp = string(randInt(65, 90))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}

func get(tag string) dom.Element {
	return dom.GetWindow().Document().CreateElement(tag)
}