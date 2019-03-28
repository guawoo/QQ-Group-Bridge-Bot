package observer

import (
	"errors"
	"io/ioutil"

	// "bottypes"
	"fmt"
	"net/http"
	"regexp"
)

type (
	msgType int
	handle  func(*[]byte) []byte
)

var (
	broadCastTypes = make(map[string]msgType)
	typeHandles    = make(map[msgType]handle)
)

func registerType(pattern string, t msgType) {
	broadCastTypes[pattern] = t
}

func registerTypeHandle(t msgType, h handle) {
	typeHandles[t] = h
}

func RegisterMsgHandle(pattern string, t msgType, h handle) {
	registerType(pattern, t)
	registerTypeHandle(t, h)
}

func getMsgTypeFromRawData(data *[]byte) (msgType, error) {
	for pattern, t := range broadCastTypes {
		rs, _ := regexp.Match(pattern, *data)
		if rs {
			return t, nil
		}
	}

	err := fmt.Sprintf("没有找到匹配的消息类型。==>\n%s\n", *data)
	return -1, errors.New(err)
}

func entry(rs http.ResponseWriter, rq *http.Request) {

	body, _ := ioutil.ReadAll(rq.Body)

	t, err := getMsgTypeFromRawData(&body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fastresp := typeHandles[t](&body)

	if fastresp != nil {
		rs.Write(fastresp)
		return
	}

	rs.WriteHeader(http.StatusNoContent)

}

func RunHttpObserver() {
	http.HandleFunc("/", entry)
	http.ListenAndServe(":8081", nil)

}
