package test

import (
	"bottypes"
	"encoding/json"
	"fmt"
	"regexp"
)

func MatchType(pattern string, src string) bool {
	rs, err := regexp.Match(pattern, []byte(src))
	if err != nil {
		panic("Error!")
	}

	return rs
}

func UnJson(data string, m *bottypes.PrivateMsg) {
	err := json.Unmarshal([]byte(data), m)
	if err != nil {
		fmt.Println(err)
	}
}
