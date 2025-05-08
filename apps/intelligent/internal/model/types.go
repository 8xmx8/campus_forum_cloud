package model

type ProcessActions string

const (
	ProcessActionNone   ProcessActions = "none"   // 	未处理
	ProcessActionPass   ProcessActions = "pass"   // 	通过
	ProcessActionReject ProcessActions = "reject" // 	拒绝
	ProcessActionModify ProcessActions = "modify" // 	修改
	ProcessActionDelete ProcessActions = "delete" // 	删除
)

var processActionMap = map[int64]ProcessActions{
	0: ProcessActionNone,
	1: ProcessActionPass,
	2: ProcessActionModify,
	3: ProcessActionModify,
	4: ProcessActionModify,
	5: ProcessActionModify,
}

func GetProcessAction(level int64) ProcessActions {
	if level > 5 {
		return ProcessActionModify
	}
	if action, ok := processActionMap[level]; ok {
		return action
	}
	return ProcessActionNone
}
