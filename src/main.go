package main

import (
	"bottypes"
	"broadcast"
	"config"
	"event"
	"fmt"
	"observer"
)

var BOT_LOGIN_INFO bottypes.LoginInfo

func init() {

	BOT_LOGIN_INFO, _ = broadcast.GetLoginInfo()
	fmt.Printf("%+v\n", BOT_LOGIN_INFO)

	registerInit()
}

func main() {
	broadcast.RunStation()
	fmt.Println("station running..")
	observer.RunHttpObserver()
	fmt.Println("observer running..")

}

func registerInit() {
	//加好友请求
	observer.RegisterMsgHandle(
		bottypes.DefaultRequestFriendPattern,
		bottypes.DefualtRequestFriend,
		func(b *[]byte) []byte {

			var rqfriend bottypes.RequestFrined

			err := bottypes.ConvertToMsgStruct(b, &rqfriend)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if "湖南游戏" == rqfriend.Comment {
				fmt.Println(fmt.Sprintf("同意%d的加好友请求。", rqfriend.UserID))
				return []byte(`{"approve":true,"remark":""}`)
			}
			return nil
		})

	//邀请入群请求
	observer.RegisterMsgHandle(
		bottypes.DefaultRequestInviteGroupPattern,
		bottypes.DefualtRequestInviteGroup,
		func(b *[]byte) []byte {

			return []byte(`{"approve":true,"reason":"None"}`)
		})

	//------------Bot入群通知处理---------------
	observer.RegisterMsgHandle(
		bottypes.DefaultInGroupPattern,
		bottypes.DefaultInGroup,
		func(b *[]byte) []byte {
			fmt.Println("进入了一个群。")

			inGroupMsg := struct {
				GroupID    int64  `json:"group_id"`
				NoticeType string `json:"notice_type"`
				OperatorID int    `json:"operator_id"`
				PostType   string `json:"post_type"`
				SelfID     int64  `json:"self_id"`
				SubType    string `json:"sub_type"`
				Time       int    `json:"time"`
				UserID     int64  `json:"user_id"`
			}{}

			err := bottypes.ConvertToMsgStruct(b, &inGroupMsg)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if inGroupMsg.UserID != BOT_LOGIN_INFO.UserID {
				return nil
			}

			broadcast.DoAction(func() {
				grouplist, err := broadcast.GetGroupList()
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, val := range grouplist {
					if inGroupMsg.GroupID == val.GroupID {
						config.AddToJoinedGroups(val)
						break
					}
				}
			})

			broadcast.DoAction(func() {
				broadcast.SendGroupMsg("大家好，我是QQ群桥接机器人。\r\n输入#help查看具体用法。\r\n本Q没付费，不支持传送图片和自定义表情。", inGroupMsg.GroupID)
			})

			return nil
		})
	//------------退群通知处理------------

	observer.RegisterMsgHandle(
		bottypes.DefaultOutGroupPattern,
		bottypes.DefaultOutGroup,
		func(b *[]byte) []byte {
			fmt.Println("退出了一个群。")

			outGroupMsg := struct {
				GroupID    int64  `json:"group_id"`
				NoticeType string `json:"notice_type"`
				OperatorID int    `json:"operator_id"`
				PostType   string `json:"post_type"`
				SelfID     int64  `json:"self_id"`
				SubType    string `json:"sub_type"`
				Time       int    `json:"time"`
				UserID     int64  `json:"user_id"`
			}{}

			err := bottypes.ConvertToMsgStruct(b, &outGroupMsg)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if outGroupMsg.UserID != BOT_LOGIN_INFO.UserID {
				return nil
			}

			broadcast.DoAction(func() {
				config.RemoveFromeJoinedGroups(outGroupMsg.GroupID)
			})

			return nil
		})
	//------------私聊消息处理------------
	observer.RegisterMsgHandle(
		bottypes.DefaultPrivateMsgPattern,
		bottypes.DefaultPrivateMsg,
		func(b *[]byte) []byte {

			fmt.Printf("%s\n", *b)
			var msg bottypes.PrivateMsg
			err := bottypes.ConvertToMsgStruct(b, &msg)
			if err != nil {
				fmt.Print(err)
				return nil
			}
			fmt.Printf("%+v\n", msg)

			// frs := bottypes.MakePrivateMsgFastRespone("hello", true)

			// if frs != nil {
			// 	return frs
			// }
			broadcast.DoAction(func() {
				broadcast.SendPrivateMsg(msg.Message, msg.Sender.UserID)
			})

			return nil
		})

	//------------群聊消息处理------------
	observer.RegisterMsgHandle(
		bottypes.DefaultGroupMsgPattern,
		bottypes.DefaultGroupMsg,
		func(i *[]byte) []byte {

			fmt.Printf("%s\n", i)

			var msg bottypes.GroupMsg
			err := bottypes.ConvertToMsgStruct(i, &msg)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if msg.Message[0] == '#' {
				fmt.Println("system event")
				event.SystemEventHandle(msg)
			} else if msg.Message[0] == '<' && msg.Message[1] == '<' {
				fmt.Println("broadcast event")
				event.BroadCastEventHanlde(msg)
			}

			return nil
		})
}
