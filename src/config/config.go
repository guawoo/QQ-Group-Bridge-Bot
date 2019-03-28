package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"bottypes"
	"sync"
)

const CONFIG_FILE = "QQBotconfig.json"

var configOBj _Config
var lock *sync.RWMutex    //文件操作锁
var objLock *sync.RWMutex //对象操作锁

type _Config struct {
	//加入的群
	JoinedGroups []bottypes.GroupInfo `json:"joined_groups"`
	//加入聊天广场的群
	SquareGroups []bottypes.GroupInfo `json:"square_Groups"`

	BanUserID  map[int64][]int64 `json:"ban_userid"`
	BanGroupID map[int64][]int64 `json:"ban_groupid"`
}

func InitConfig() error {
	fmt.Println("init config file..")

	if _, err := os.Stat(CONFIG_FILE); err == nil || os.IsExist(err) {
		// 如果文件存在则读取
		fmt.Println("read config file..")
		data, err := ioutil.ReadFile(CONFIG_FILE)
		if err != nil {
			return errors.New("读取config文件失败。")
		}

		err = json.Unmarshal(data, &configOBj)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", configOBj)

	} else {
		//如果文件不存在,则创建
		fmt.Println("create config file..")
		_, err := os.Create(CONFIG_FILE)
		if err != nil {
			return errors.New("config文件创建失败。")
		}

		configOBj = _Config{}
		configOBj.BanGroupID = make(map[int64][]int64)
		configOBj.BanUserID = make(map[int64][]int64)
	}

	lock = new(sync.RWMutex)
	objLock = new(sync.RWMutex)

	return nil
}

func saveConfigToFile(obj _Config) {

	objLock.Lock()
	data, err := json.Marshal(obj)
	objLock.Unlock()

	if err != nil {
		fmt.Println(err)
		return
	}

	lock.Lock()

	err = ioutil.WriteFile(CONFIG_FILE, data, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}

	lock.Unlock()
}

func GetGroupInfoFromJoinedGroups(groupID int64) (bottypes.GroupInfo, error) {
	objLock.RLock()
	for index, val := range configOBj.JoinedGroups {
		if groupID == val.GroupID {

			gi := configOBj.JoinedGroups[index]
			objLock.RUnlock()

			return gi, nil
		}
	}
	objLock.RUnlock()
	return bottypes.GroupInfo{}, errors.New("没有找到相关群的信息。")
}

func AddToSquareGroups(info bottypes.GroupInfo) error {
	objLock.RLock()

	for _, i := range configOBj.SquareGroups {
		if i.GroupID == info.GroupID {

			objLock.RUnlock()

			return errors.New("已经加入")
		}
	}

	objLock.RUnlock()

	objLock.Lock()
	configOBj.SquareGroups = append(configOBj.SquareGroups, info)
	objLock.Unlock()

	saveConfigToFile(configOBj)

	return nil
}

func RemoveFromSquareGroups(groupID int64) error {
	objLock.RLock()
	for index, val := range configOBj.SquareGroups {
		if val.GroupID == groupID {

			objLock.RUnlock()

			objLock.Lock()
			configOBj.SquareGroups = append(configOBj.SquareGroups[:index], configOBj.SquareGroups[index+1:]...)
			objLock.Unlock()

			saveConfigToFile(configOBj)
			return nil
		}
	}
	objLock.RUnlock()
	return errors.New("没有找到group id,无法删除。")
}

func AddToJoinedGroups(info bottypes.GroupInfo) error {
	objLock.RLock()
	for _, i := range configOBj.JoinedGroups {
		if i.GroupID == info.GroupID {

			objLock.RUnlock()

			return errors.New("加入了已经加入群，逻辑错误，请检查。")
		}
	}
	objLock.RUnlock()

	objLock.Lock()
	configOBj.JoinedGroups = append(configOBj.JoinedGroups, info)
	objLock.Unlock()

	saveConfigToFile(configOBj)

	return nil
}

func RemoveFromeJoinedGroups(groupID int64) error {
	objLock.RLock()
	for index, val := range configOBj.JoinedGroups {
		if groupID == val.GroupID {
			objLock.RUnlock()

			objLock.Lock()
			configOBj.JoinedGroups = append(configOBj.JoinedGroups[:index], configOBj.JoinedGroups[index+1:]...)
			objLock.Unlock()

			RemoveFromSquareGroups(groupID)

			saveConfigToFile(configOBj)

			return nil
		}
	}
	objLock.RUnlock()

	return errors.New("没有找到group id,无法删除。")
}

func ListSquareGroups() ([]bottypes.GroupInfo, error) {
	if len(configOBj.SquareGroups) <= 0 {
		return nil, errors.New("广场中没有任何群。")
	}

	objLock.RLock()
	list := configOBj.SquareGroups
	objLock.RUnlock()

	return list, nil
}

func BanUserID(userid int64, owngroup int64) {
	objLock.Lock()
	configOBj.BanUserID[owngroup] = append(configOBj.BanUserID[owngroup], userid)
	objLock.Unlock()
	saveConfigToFile(configOBj)
}

func LiftUserID(userid int64, owngroup int64) error {
	objLock.RLock()
	for index, val := range configOBj.BanUserID[owngroup] {
		if userid == val {
			objLock.RUnlock()

			objLock.Lock()
			configOBj.BanUserID[owngroup] = append(configOBj.BanUserID[owngroup][:index], configOBj.BanUserID[owngroup][index+1:]...)
			objLock.Unlock()

			saveConfigToFile(configOBj)
			return nil
		}
	}
	objLock.RUnlock()
	return errors.New("无法删除禁言ID。")
}
