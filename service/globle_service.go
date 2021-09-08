package service

import (
	"github.com/wechaty/go-wechaty/wechaty"
)

var (
	globleServiceInstance = &GlobleService{}
)

func NewGlobleService() *GlobleService {
	return globleServiceInstance
}

type GlobleService struct {
	Bot *wechaty.Wechaty
}

func (g *GlobleService) SetBot(bot *wechaty.Wechaty) {
	g.Bot = bot
}

func (g *GlobleService) GetBot() *wechaty.Wechaty {
	return g.Bot
}
