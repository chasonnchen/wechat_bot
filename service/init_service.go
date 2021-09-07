package service

func InitService() {
	// 初始化webot基础配置
	NewContactService().init()
	NewQaService().init()
}
