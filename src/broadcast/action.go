package broadcast

import (
	"bottypes"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	IP   = "127.0.0.1"
	PORT = 5700
)

const (
	API_SEND_PRIVATE_MSG      = "/send_private_msg"
	API_SEND_GROUP_MSG        = "/send_group_msg"
	API_GET_GROUP_MEMBER_INFO = "/get_group_member_info"
	API_GET_GROUP_LIST        = "/get_group_list"
	API_GET_LOGIN_INFO        = "/get_login_info"
)

// func Test() {
// 	action := func() {
// 		fmt.Println("hi every one..")
// 	}

// 	job := _Job{Action: action}

// 	_JobQueue <- job
// }

func SendPrivateMsg(msg string, QQ int64) {
	url := fmt.Sprintf("http://%s:%d%s", IP, PORT, API_SEND_PRIVATE_MSG)

	fmt.Printf("私聊API地址：%s\n", url)
	fmt.Printf("QQ: %d\n", QQ)

	j, err := json.Marshal(struct {
		UserID     int64  `json:"user_id"`
		Message    string `json:"message"`
		AutoEscape bool   `json:"auto_escape"`
	}{
		QQ, msg, false,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("发送json: %s\n", j)

	dataReader := bytes.NewBuffer(j)

	req, err := http.NewRequest("POST", url, dataReader)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	resdata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("返回内容：%s\n", resdata)

}

func SendGroupMsg(msg string, groupID int64) {
	url := fmt.Sprintf("http://%s:%d%s", IP, PORT, API_SEND_GROUP_MSG)

	fmt.Printf("群聊API地址：%s\n", url)
	fmt.Printf("groupID: %d\n", groupID)

	j, err := json.Marshal(struct {
		GroupID    int64  `json:"group_id"`
		Message    string `json:"message"`
		AutoEscape bool   `json:"auto_escape"`
	}{
		groupID, msg, false,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("发送json: %s\n", j)

	dataReader := bytes.NewBuffer(j)

	req, err := http.NewRequest("POST", url, dataReader)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	resdata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("返回内容：%s\n", resdata)
}

func GetGroupMemberInfo(groupid int64, QQ int64) (bottypes.GroupMemberInfo, error) {
	url := fmt.Sprintf("http://%s:%d%s", IP, PORT, API_GET_GROUP_MEMBER_INFO)

	fmt.Printf("获得群成员信息API地址：%s\n", url)
	fmt.Printf("QQ: %d, GroupID: %d \n", QQ, groupid)

	j, err := json.Marshal(struct {
		UserID  int64 `json:"user_id"`
		GroupID int64 `json:"group_id"`
		NoCache bool  `json:"no_cache"`
	}{
		QQ, groupid, false,
	})

	if err != nil {
		fmt.Println(err)
		return bottypes.GroupMemberInfo{}, err
	}

	fmt.Printf("发送json: %s\n", j)

	dataReader := bytes.NewBuffer(j)

	req, err := http.NewRequest("POST", url, dataReader)
	if err != nil {
		fmt.Println(err)
		return bottypes.GroupMemberInfo{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return bottypes.GroupMemberInfo{}, err
	}

	defer res.Body.Close()

	resdata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return bottypes.GroupMemberInfo{}, err
	}

	info := bottypes.GroupMemberInfo{}

	err = bottypes.ConvertToMsgStruct(&resdata, &info)
	if err != nil {
		return bottypes.GroupMemberInfo{}, err
	}

	return info, nil
}

func GetGroupList() ([]bottypes.GroupInfo, error) {
	url := fmt.Sprintf("http://%s:%d%s", IP, PORT, API_GET_GROUP_LIST)

	fmt.Printf("获得群列表API地址：%s\n", url)

	null := []byte(`{"data":0}`)

	fmt.Printf("发送json: %s\n", null)

	dataReader := bytes.NewBuffer(null)

	req, err := http.NewRequest("POST", url, dataReader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	resdata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rsdatastruct := struct {
		Data    []bottypes.GroupInfo
		Retcode int    `json:"retcode"`
		Status  string `json:"status"`
	}{}

	err = bottypes.ConvertToMsgStruct(&resdata, &rsdatastruct)
	if err != nil {
		return nil, err
	}

	return rsdatastruct.Data, nil
}

func GetLoginInfo() (bottypes.LoginInfo, error) {
	url := fmt.Sprintf("http://%s:%d%s", IP, PORT, API_GET_LOGIN_INFO)

	fmt.Printf("获得登录信息API地址：%s\n", url)

	null := []byte(`{"data":0}`)

	fmt.Printf("发送json: %s\n", null)

	dataReader := bytes.NewBuffer(null)

	req, err := http.NewRequest("POST", url, dataReader)
	if err != nil {
		fmt.Println(err)
		return bottypes.LoginInfo{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return bottypes.LoginInfo{}, err
	}

	defer res.Body.Close()

	resdata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return bottypes.LoginInfo{}, err
	}

	var info = struct {
		Data struct {
			Nickname string `json:"nickname"`
			UserID   int64  `json:"user_id"`
		} `json:"data"`
		Retcode int    `json:"retcode"`
		Status  string `json:"status"`
	}{}

	err = bottypes.ConvertToMsgStruct(&resdata, &info)
	if err != nil {
		return bottypes.LoginInfo{}, err
	}

	return bottypes.LoginInfo{NickName: info.Data.Nickname, UserID: info.Data.UserID}, nil
}

func DoAction(action func()) {
	job := _Job{Action: action}
	_JobQueue <- job
}
