package bottypes

import (
	"encoding/json"
	"fmt"
)

//正则式
const (
	//添加好友请求
	DefaultRequestFriendPattern = `.*(?:"request_type":"friend").*`
	//入群邀请请求
	DefaultRequestInviteGroupPattern = `.*(?:"request_type":"group").*`
	//私聊消息
	DefaultPrivateMsgPattern = `.*(?:"message_type":"private").*(?:"post_type":"message")`
	//群聊消息
	DefaultGroupMsgPattern = `.*(?:"message_type":"group").*(?:"post_type":"message")`

	//进群通知
	DefaultInGroupPattern = `.*(?:"notice_type":"group_increase").*`
	//退群通知
	DefaultOutGroupPattern = `.*(?:"notice_type":"group_decrease").*`
)

//类型
const (
	_ = iota
	DefaultPrivateMsg
	DefaultGroupMsg

	DefaultInGroup
	DefaultOutGroup

	DefualtRequestFriend
	DefualtRequestInviteGroup
)

//---------LoginInfo------------
type LoginInfo struct {
	NickName string `json:"nickname"`
	UserID   int64  `json:"user_id"`
}

//---------GroupInfo------------
type GroupInfo struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}

//---------RequestFrined---------
type RequestFrined struct {
	Comment     string `json:"comment"`
	Flag        string `json:"flag"`
	PostType    string `json:"post_type"`
	RequestType string `json:"request_type"`
	SelfID      int64  `json:"self_id"`
	Time        int    `json:"time"`
	UserID      int64  `json:"user_id"`
}

//---------PrivateMsg-----------
type PrivateMsg struct {
	Font        int    `json:"font"`
	Message     string `json:"message"`
	MessageID   int64  `json:"message_id"`
	MessageType string `json:"message_type"`
	PostType    string `json:"post_type"`
	RawMessage  string `json:"raw_message"`
	SelfID      int64  `json:"self_id"`
	Sender      struct {
		Age      int    `json:"age"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserID   int64  `json:"user_id"`
	} `json:"sender"`
	SubType string `json:"sub_type"`
	Time    int    `json:"time"`
	UserID  int64  `json:"user_id"`
}

type PrivateMsgRespone struct {
	Reply      string `json:"reply"`
	AutoEscape bool   `json:"auto_escape"`
}

//---------GroupMsg------------------
type GroupMsg struct {
	Anonymous   interface{} `json:"anonymous"`
	Font        int         `json:"font"`
	GroupID     int64       `json:"group_id"`
	Message     string      `json:"message"`
	MessageID   int         `json:"message_id"`
	MessageType string      `json:"message_type"`
	PostType    string      `json:"post_type"`
	RawMessage  string      `json:"raw_message"`
	SelfID      int64       `json:"self_id"`
	Sender      struct {
		Age      int    `json:"age"`
		Area     string `json:"area"`
		Card     string `json:"card"`
		Level    string `json:"level"`
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Sex      string `json:"sex"`
		Title    string `json:"title"`
		UserID   int64  `json:"user_id"`
	} `json:"sender"`
	SubType string `json:"sub_type"`
	Time    int    `json:"time"`
	UserID  int64  `json:"user_id"`
}

//---------GroupMemberInfo-----------
type GroupMemberInfo struct {
	Data struct {
		Age             int    `json:"age"`
		Area            string `json:"area"`
		Card            string `json:"card"`
		CardChangeable  bool   `json:"card_changeable"`
		GroupID         int    `json:"group_id"`
		JoinTime        int    `json:"join_time"`
		LastSentTime    int    `json:"last_sent_time"`
		Level           string `json:"level"`
		Nickname        string `json:"nickname"`
		Role            string `json:"role"`
		Sex             string `json:"sex"`
		Title           string `json:"title"`
		TitleExpireTime int    `json:"title_expire_time"`
		Unfriendly      bool   `json:"unfriendly"`
		UserID          int    `json:"user_id"`
	} `json:"data"`
	Retcode int    `json:"retcode"`
	Status  string `json:"status"`
}

//-----------help functions---------------

func ConvertToMsgStruct(data *[]byte, msg interface{}) error {
	err := json.Unmarshal(*data, msg)
	if err != nil {
		return err
	}
	return nil
}

func MakePrivateMsgFastRespone(msg string, autoescap bool) []byte {
	msgRes := PrivateMsgRespone{Reply: msg, AutoEscape: autoescap}
	rs, err := json.Marshal(&msgRes)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return rs
}
