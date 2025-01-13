package test

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v any) {
	jsonText, _ := json.MarshalIndent(v, "", "\t")
	fmt.Println(string(jsonText))
}
