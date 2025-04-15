package model

type ProcessActions string

const (
	ProcessActionNone   ProcessActions = "none"   // 	未处理
	ProcessActionPass   ProcessActions = "pass"   // 	通过
	ProcessActionReject ProcessActions = "reject" // 	拒绝
	ProcessActionModify ProcessActions = "modify" // 	修改
	ProcessActionDelete ProcessActions = "delete" // 	删除
)
