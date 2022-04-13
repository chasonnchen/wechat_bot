package service

func InitService() {
	// 初始化webot基础配置
	NewCacheService().init()

	NewContactService().init()
	NewQaService().init()
	NewForwardService().init()
	NewForwardMediaService().init()
	NewRoomService().init()
	NewGroupService().init()
}
