package gohut

import "github.com/tidwall/gjson"

func gjsonStringArray(g gjson.Result) []string {
	arr := g.Array()
	slice := make([]string, len(arr))

	for i, item := range g.Array() {
		slice[i] = item.String()
	}

	return slice
}
