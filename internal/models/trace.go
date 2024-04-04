package models

type Trace struct {
	Ctime int    `json:"ctime"`
	Ext   string `json:"ext" `
	Id    int64  `json:"id"`
	Ip    int    `json:"ip"`
	Type  int    `json:"type"`
	Uid   int64  `json:"uid" `
}

var TraceTypeReg = 0
var TraceTypeLogin = 1
var TraceTypeOut = 2
var TraceTypeEdit = 3
var TraceTypeDel = 4
