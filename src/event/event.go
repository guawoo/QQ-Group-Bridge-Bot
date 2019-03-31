package event

import (
	"bottypes"
	"broadcast"
	"bytes"
	"config"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	_ = iota
	SYSTEM_JOIN_SQUARE
	SYSTEM_LEAVE_SQUARE
	SYSTEM_LIST_SQUARE
	SYSTEM_HELP
	SYSTEM_BAN_USER_ID
	SYSTEM_BAN_GROUP_ID
	SYSTEM_LIST_BAN

	BROADCAST_SEND_MSG
)

func GetSystemEventType(msg *string) int {
	if len(*msg) >= 5 && "#join" == (*msg)[:5] {
		return SYSTEM_JOIN_SQUARE
	} else if len(*msg) >= 6 && "#leave" == (*msg)[:6] {
		return SYSTEM_LEAVE_SQUARE
	} else if len(*msg) >= 5 && "#list" == (*msg)[:5] {
		return SYSTEM_LIST_SQUARE
	} else if len(*msg) >= 5 && "#help" == (*msg)[:5] {
		return SYSTEM_HELP
	} else if len(*msg) >= 7 && "#banuid" == (*msg)[:7] {
		return SYSTEM_BAN_USER_ID
	} else if len(*msg) == 8 && "#banlist" == (*msg)[:8] {
		return SYSTEM_LIST_BAN
	}

	return -1
}

func GetBroadCastEventType(msg *string) int {
	return BROADCAST_SEND_MSG
}

func SystemEventHandle(msg interface{}) {
	typeOfMsg := reflect.TypeOf(msg)

	//处理群聊中的系统事件
	if typeOfMsg.Name() == "GroupMsg" {
		fmt.Println("group system msg..")

		groupMsg := msg.(bottypes.GroupMsg)

		eventType := GetSystemEventType(&groupMsg.Message)

		switch eventType {
		//加入广场事件
		case SYSTEM_JOIN_SQUARE:
			systemEventJoinSquare(&groupMsg)
		//离开广场事件
		case SYSTEM_LEAVE_SQUARE:
			systemEventLeaveSquare(&groupMsg)
		//显示广场群列表事件
		case SYSTEM_LIST_SQUARE:
			systemEventListSquare(&groupMsg)
		case SYSTEM_HELP:
			systemEventHelp(&groupMsg)
		case SYSTEM_BAN_USER_ID:
			systemEventBanUserID(&groupMsg)
		case SYSTEM_LIST_BAN:
			fmt.Println("banlist")
			systemEventListBan(&groupMsg)
		}

	}
}

func BroadCastEventHanlde(msg interface{}) {
	typeOfMsg := reflect.TypeOf(msg)

	if typeOfMsg.Name() == "GroupMsg" {
		fmt.Println("group broadcast msg..")

		groupMsg := msg.(bottypes.GroupMsg)

		eventType := GetBroadCastEventType(&groupMsg.Message)

		switch eventType {
		case BROADCAST_SEND_MSG:
			broadCastEvnetSendMsg(&groupMsg)
		}
	}
}

func systemEventBanUserID(groupMsg *bottypes.GroupMsg) {
	if groupMsg.Sender.Role == "owner" || groupMsg.Sender.Role == "admin" {
		fmt.Println("屏蔽QQID操作。")
		id, _ := strconv.ParseInt(
			strings.Replace(groupMsg.Message[7:], " ", "", -1),
			10, 64)

		fmt.Printf("屏蔽id:%d", id)

		config.BanUserID(id, groupMsg.GroupID)
	}
}

func systemEventListSquare(groupMsg *bottypes.GroupMsg) {
	broadcast.DoAction(func() {
		var strbuf bytes.Buffer
		strbuf.WriteString("当前加入广场的群有：\r\n")
		strbuf.WriteString("------------------------\r\n")

		list, err := config.ListSquareGroups()

		if err != nil {
			broadcast.SendGroupMsg(err.Error(), groupMsg.GroupID)
			return
		}

		for _, val := range list {
			strbuf.WriteString(fmt.Sprintf("%s (%d)\r\n", val.GroupName, val.GroupID))
		}

		broadcast.SendGroupMsg(strbuf.String(), groupMsg.GroupID)

	})
}

func systemEventJoinSquare(groupMsg *bottypes.GroupMsg) {
	if groupMsg.Sender.Role == "owner" || groupMsg.Sender.Role == "admin" {

		groupinfo, err := config.GetGroupInfoFromJoinedGroups(groupMsg.GroupID)

		if err != nil {
			fmt.Println(err)
			broadcast.DoAction(func() {
				broadcast.SendGroupMsg("相关群信息正在同步中，请稍后在试。", groupMsg.GroupID)
			})
			return
		}

		err = config.AddToSquareGroups(groupinfo)
		if err != nil {
			broadcast.DoAction(func() {
				broadcast.SendGroupMsg("本群已经加入广场，不可重复加入", groupMsg.GroupID)
			})
			return
		}

		broadcast.DoAction(func() {
			broadcast.SendGroupMsg("本群成功加入广场。", groupMsg.GroupID)
		})

		broadcast.DoAction(func() {
			var strbuf bytes.Buffer
			strbuf.WriteString("当前加入广场的群有：\r\n")
			strbuf.WriteString("------------------------\r\n")

			list, err := config.ListSquareGroups()

			if err != nil {
				broadcast.SendGroupMsg(err.Error(), groupMsg.GroupID)
				return
			}

			for _, val := range list {
				strbuf.WriteString(fmt.Sprintf("%s (%d)\r\n", val.GroupName, val.GroupID))
			}

			broadcast.SendGroupMsg(strbuf.String(), groupMsg.GroupID)

		})

		//通知其它广场群，有群加入
		broadcast.DoAction(func() {
			list, err := config.ListSquareGroups()
			if err != nil {
				return
			}

			var gi bottypes.GroupInfo

			for _, val := range list {
				if val.GroupID == groupMsg.GroupID {
					gi = val
					break
				}
			}

			for _, val := range list {
				if val.GroupID != groupMsg.GroupID {
					broadcast.SendGroupMsg(fmt.Sprintf("群【%s(%d)】加入了广场。", gi.GroupName, gi.GroupID), val.GroupID)
				}
			}

		})

		return
	}

	fmt.Println("没有权限。。")
	broadcast.DoAction(func() {
		broadcast.SendGroupMsg(fmt.Sprintf("[CQ:at,qq=%d] 你没有权限，只有管理员和群主能操作系统命令。", groupMsg.Sender.UserID), groupMsg.GroupID)
	})
}

func systemEventLeaveSquare(groupMsg *bottypes.GroupMsg) {
	if groupMsg.Sender.Role == "owner" || groupMsg.Sender.Role == "admin" {
		err := config.RemoveFromSquareGroups(groupMsg.GroupID)
		if err != nil {
			fmt.Println(err)
			broadcast.DoAction(func() {
				broadcast.SendGroupMsg("本群已经退出广场。", groupMsg.GroupID)
			})
			return
		}
		fmt.Println("退出成功")
		broadcast.DoAction(func() {
			broadcast.SendGroupMsg("本群成功退出广场。", groupMsg.GroupID)
		})
		return
	}
	fmt.Println("没有权限。。")
	broadcast.DoAction(func() {
		broadcast.SendGroupMsg(fmt.Sprintf("[CQ:at,qq=%d] 你没有权限，只有管理员和群主能操作系统命令。", groupMsg.Sender.UserID), groupMsg.GroupID)
	})
}

func systemEventHelp(groupMsg *bottypes.GroupMsg) {
	broadcast.DoAction(func() {
		var buf bytes.Buffer
		buf.WriteString("输入 #join 将本群加入广场。(管理员权限)\r\n")
		buf.WriteString("输入 #leave 将本群从广场中退出。(管理员权限)\r\n")
		buf.WriteString("输入 #list 查看广场中有哪些群。\r\n")
		buf.WriteString("输入 <<你想要广播的文字 广场中的群都能看见此文字。\r\n")
		broadcast.SendGroupMsg(buf.String(), groupMsg.GroupID)
	})
}

func systemEventListBan(groupMsg *bottypes.GroupMsg) {
	list, err := config.ListBanUserID(groupMsg.GroupID)
	if err != nil {
		broadcast.DoAction(func() {
			broadcast.SendGroupMsg("没有任何被禁言人QQ ID。", groupMsg.GroupID)
		})
		return
	}

	var buf bytes.Buffer
	buf.WriteString("被禁言的ID有:\r\n")
	for _, val := range list {
		buf.WriteString(strconv.FormatInt(val, 10))
		buf.WriteString("\r\n")
	}
	broadcast.DoAction(func() {
		broadcast.SendGroupMsg(buf.String(), groupMsg.GroupID)
	})
}

func broadCastEvnetSendMsg(groupMsg *bottypes.GroupMsg) {
	broadcast.DoAction(func() {
		sqaureList, err := config.ListSquareGroups()

		if err != nil {

			broadcast.SendGroupMsg("广场里没有任何群，广播失败。", groupMsg.GroupID)
			return
		}

		var gi bottypes.GroupInfo

		for _, val := range sqaureList {
			if val.GroupID == groupMsg.GroupID {
				gi = val
				break
			}
		}

		if (bottypes.GroupInfo{}) == gi {
			broadcast.SendGroupMsg("本群没有加入广场，不能广播消息。", groupMsg.GroupID)
			return
		}

	SqaureListLoop:
		for _, val := range sqaureList {
			if val.GroupID != groupMsg.GroupID {

				//处理禁言
				banidlist, err := config.ListBanUserID(val.GroupID)
				if err != nil {
					continue
				}
				if len(banidlist) > 0 {
					for _, id := range banidlist {
						if id == groupMsg.Sender.UserID {
							fmt.Println("had a ban id: " + strconv.FormatInt(id, 10))
							continue SqaureListLoop
						}
					}
				}

				var buf bytes.Buffer

				buf.WriteString(fmt.Sprintf("来自广场的【%s(%d)】说：\r\n", groupMsg.Sender.Nickname, groupMsg.Sender.UserID))
				buf.WriteString(groupMsg.Message[2:])
				buf.WriteString("\r\n")

				t := time.Now()
				timestr := fmt.Sprintf("%d-%d-%d %d:%d:%d\r\n", t.Year(),
					t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

				buf.WriteString(fmt.Sprintf("->%s->群【%s(%d)】\r\n", timestr, gi.GroupName, gi.GroupID))

				broadcast.SendGroupMsg(buf.String(), val.GroupID)
			}
		}
	})

}
