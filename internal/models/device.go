package models

type Device struct {
	Client string `json:"client"`
	Ctime  int    `json:"ctime"`
	Ext    string `json:"ext" `
	Id     int64  `json:"id" `
	Ip     int    `json:"ip"`
	Model  string `json:"model"`
	Uid    int64  `json:"uid" `
}
