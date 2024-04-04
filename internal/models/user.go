package models

import (
	"fmt"
	"time"
	"wolfdog/internal/consts"
)

type Users struct {
	Id     int64     `json:"id" `
	Name   string    `json:"name"`
	Passwd string    `json:"passwd"`
	Email  string    `json:"email"`
	Mobile string    `json:"mobile"`
	Status int       `json:"status"`
	Ext    string    `json:"ext"`
	Salt   string    `json:"salt"`
	Ctime  int       `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}
type UserRow struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

var UsersStatusOk = 1
var UsersStatusDel = 10
var UsersStatusDef = 0

var usersTable = "users"

func (u *Users) GetRow() bool {
	tx := consts.GormDB.Table(usersTable).First(&u)
	if tx.Error != nil {
		fmt.Printf("查询用户失败: %v\n", tx.Error)
		return false
	} else if tx.RowsAffected == 0 {
		fmt.Println("未查询到用户数据")
		return false
	}
	return true
}
func (u *Users) GetAll() ([]Users, error) {
	users := make([]Users, 0)
	tx := consts.GormDB.Table(usersTable).Find(&users)
	if tx.Error != nil {
		fmt.Printf("查询用户失败: %v\n", tx.Error)
		return users, tx.Error
	} else if tx.RowsAffected == 0 {
		fmt.Println("未查询到用户数据")
		return users, nil
	}
	return users, nil
}

//func (u *Users) Add(trace *Trace, device *Device) (int64, error) {
//	session := mEngine.NewSession()
//	defer session.Close()
//	// add Begin() before any action
//	if err := session.Begin(); err != nil {
//		return 0, err
//	}
//	_, err := session.Insert(u)
//	if err != nil {
//		return 0, err
//	}
//
//	trace.Uid = u.Id
//	_, err = session.Insert(trace)
//	if err != nil {
//		return 0, err
//	}
//	device.Uid = u.Id
//	_, err = session.Insert(device)
//	if err != nil {
//		return 0, err
//	}
//	return u.Id, session.Commit()
//}

func (u *Users) Add(trace *Trace, device *Device) (int64, error) {
	// 使用事务处理
	tx := consts.GormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 插入用户信息
	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 设置关联的用户ID并插入轨迹信息
	trace.Uid = u.Id
	if err := tx.Create(trace).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 设置关联的用户ID并插入设备信息
	device.Uid = u.Id
	if err := tx.Create(device).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return 0, err
	}

	return u.Id, nil
}
func IsExistsMobile(mobile string) bool {
	model := Users{Mobile: mobile}
	return model.GetRow()
}
func (u *Users) GetRowById() (UserRow, error) {
	var userRow UserRow
	tx := consts.GormDB.Table(usersTable).First(&userRow)
	if tx.Error != nil {
		fmt.Printf("查询用户失败: %v\n", tx.Error)
		return UserRow{}, tx.Error
	} else if tx.RowsAffected == 0 {
		fmt.Println("未查询到用户数据")
		return UserRow{}, nil
	}
	return userRow, nil
}
