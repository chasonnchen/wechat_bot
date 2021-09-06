package task

import (
	"github.com/wechaty/go-wechaty/wechaty"
)

func InitTask(bot *wechaty.Wechaty) {
	// 初始化webot基础配置
	NewSina7x24Task(bot).Start()
}
