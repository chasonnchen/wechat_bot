package service

import (
)

func InitService() {
	// 初始化webot基础配置
	NewRoomQaService().init()
}
