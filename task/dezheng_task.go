package task

import (
	"log"
	"strconv"
	"time"

	"github.com/chasonnchen/wechat_bot/service"
)

var (
	dezhengTask = &DezhengTask{}
)

type DezhengTask struct {
}

func NewDezhengTask() *DezhengTask {
	return dezhengTask
}

func (d *DezhengTask) Start() {
	d.work()
	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				d.work()
			}
		}
	}()
}

func (d *DezhengTask) work() {

	// 晚上10点半到早上8点半 不推送
	layout := "1504"
	timeStr, _ := strconv.Atoi(time.Now().Format(layout))
	// 提醒发朋友圈
	if timeStr == 930 {
		contactIdList := service.NewGroupService().GetContactIdListByGroupId(12)
		log.Printf("dezheng task id list is %#v", contactIdList)
		msg := "大家记得发个【朋友圈】，真诚介绍【靠谱新房项目】或者【优质二手房】 :)"
		for _, contactId := range contactIdList {
			service.NewContactService().SayTextToContact(contactId, msg)
		}
	}

	// 提醒看视频
	if timeStr == 2005 {
		contactIdList := service.NewGroupService().GetContactIdListByGroupId(12)
		log.Printf("dezheng task id list is %#v", contactIdList)
		msg := "看20点抖音直播，熟悉新房项目，储备新房项目知识。下面是传送门->"
		// 抖音传送门
		csm := "8h:/ 4【乐居邢台（广缘）的个人主页】长按复制此条消息，长按复制打开抖音搜索，查看TA的更多作品##wOgpp7PYjr8##[抖音口令]"
		for _, contactId := range contactIdList {
			service.NewContactService().SayTextToContact(contactId, msg)
			time.Sleep(1 * time.Second)
			service.NewContactService().SayTextToContact(contactId, csm)
			/*contact := service.NewContactService().GetById(contactId)
			            if contact.Type == 2 {
			                // 群里回答并at all
			                currRoom := service.NewGlobleService().GetBot().Room().Load(contact.Id)
			                atContact, err := currRoom.MemberAll(nil)
			                if err != nil {
			                    log.Printf("dezheng task get at list err. %v", err)
			                    continue
			                }
			                log.Printf("AAA list all is %v", atContact)

			                // todo 过滤自己
			                newAtList :=  make([]_interface.IContact, 0)
			                for _, v := range atContact {
			                    log.Printf("BBB list v is %v", v)
			                    if v.Self() == false {
			                        newAtList = append(newAtList, v)
			                    }
			                }
			                currRoom.Say(strings.Trim(msg, "\n"), newAtList...)
			            } else {
						    service.NewContactService().SayTextToContact(contactId, msg)
			            }*/
		}
	}
}
