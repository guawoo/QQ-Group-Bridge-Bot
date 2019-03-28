package test

import (
	"bottypes"
	// "encoding/json"
	"fmt"
	"testing"
)

func TestMatchType(t *testing.T) {
	src := `{"anonymous":null,"font":52363088,"group_id":838699508,"message":"sdf","message_id":264,"message_type":"group","post_type":"message","raw_message":"sdf","self_id":2165685084,"sender":{"age":38,"area":"长沙","card":"","level":"冒泡","nickname":"No.142","role":"owner","sex":"male","title":""
,"user_id":1044715076},"sub_type":"normal","time":1552630169,"user_id":1044715076}`

	pattern := `.*(?:"message_type":"group").*(?:"post_type":"message")`

	if MatchType(pattern, src) != true {
		t.Error("sorry..")
	}
}

func TestUnJson(t *testing.T) {
	data := `{"font":136421120,"message":"dddd","message_id":52,"message_type":"private","post_type":"message","raw_message":"dddd","self_id":2165685084,"sender":{"age":38,"nickname":"No.142","sex":"male","user_id":1044715076},"sub_type":"friend","time":1551068033,"user_id":1044715076}`
	var msg bottypes.PrivateMsg
	// err := json.Unmarshal([]byte(data), &msg)
	databytes := []byte(data)
	bottypes.ConvertToMsgStruct(&databytes, &msg)

	rs := fmt.Sprintf("%+v", msg)
	t.Log(rs)
}

func TestStringSlice(t *testing.T) {
	str := "123456"
	func(s *string) {
		k := (*s)[:3]
		fmt.Println(k)
	}(&str)

	fmt.Println(str)
}
