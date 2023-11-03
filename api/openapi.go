package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/chasonnchen/wechat_bot/service"
	"github.com/chasonnchen/wechat_bot/service/openapi"

	"github.com/gin-gonic/gin"
)

type CommonParams struct {
	AppId int32  `form:"appid" binding:"required"`
	Ts    int64  `form:"ts" binding:"required"`
	Once  string `form:"once" binding:"required"`
	Sign  string `form:"sign" binding:"required"`
}

type CommonResponse struct {
	Status int32       `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func GetOriJson(ctx *gin.Context) string {
	oriJson := "{}"
	data, err := ctx.GetRawData()
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	if err != nil {
		log.Printf("get ori post json fail, err:%v", err)
		return oriJson
	}
	log.Printf("get ori post json is %v\n", string(data))

	if len(string(data)) > 0 {
		return string(data)
	}

	return oriJson
}

func CheckSign(ctx *gin.Context, param CommonParams) error {
	// 判断时间戳是否过期

	appKey := openapi.NewUserService().GetAppKeyByAppId(param.AppId)
	strSignInput := appKey + "appid=" + strconv.FormatInt(int64(param.AppId), 10) + "&ts=" + strconv.FormatInt(param.Ts, 10) + "&once=" + param.Once + "&data=" + GetOriJson(ctx)
	log.Printf("sign input  is %s\n", strSignInput)
	h := md5.New()
	h.Write([]byte(strSignInput))
	sign := hex.EncodeToString(h.Sum(nil))
	log.Printf("sign is %s\n", sign)

	if sign != param.Sign {
		return fmt.Errorf("bad sign")
	}

	return nil
}

func RoomGetAll(ctx *gin.Context) {
	// 解析参数并简单检查
	var commonParams CommonParams
	if err := ctx.ShouldBindQuery(&commonParams); err != nil {
		ctx.JSON(http.StatusOK, CommonResponse{
			Status: 10001,
			Msg:    "公参有误，请查看接口文档",
			Data:   "",
		})
		return
	}
	// 检查签名
	err := CheckSign(ctx, commonParams)
	if err != nil {
		ctx.JSON(http.StatusOK, CommonResponse{
			Status: 10003,
			Msg:    "签名错误，请查看接口文档",
			Data:   "",
		})
		return
	}

	// 根据appid查询查询所有群
	roomList := service.NewContactService().GetByAppId(commonParams.AppId)
	type Room struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	var rooms []Room
	for _, item := range roomList {
		rooms = append(rooms, Room{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	// 返回结果
	ctx.JSON(http.StatusOK, CommonResponse{
		Status: 20000,
		Msg:    "success",
		Data:   rooms,
	})
}

func MsgSend(ctx *gin.Context) {
	// 解析参数并简单检查
	var commonParams CommonParams
	if err := ctx.ShouldBindQuery(&commonParams); err != nil {
		ctx.JSON(http.StatusOK, CommonResponse{
			Status: 10001,
			Msg:    "公参有误，请查看接口文档",
			Data:   "",
		})
		return
	}
	// 检查签名
	err := CheckSign(ctx, commonParams)
	if err != nil {
		ctx.JSON(http.StatusOK, CommonResponse{
			Status: 10003,
			Msg:    "签名错误，请查看接口文档",
			Data:   "",
		})
		return
	}

	// 获取POST参数
	type MsgSendRequest struct {
		Id  string `form:"id" binding:"required"`
		Msg string `form:"msg" binding:"required"`
	}
	var msgSendRequest MsgSendRequest
	if err := ctx.ShouldBindJSON(&msgSendRequest); err != nil {
		ctx.JSON(http.StatusOK, CommonResponse{
			Status: 10002,
			Msg:    "POST JSON参数有误，请查看接口文档",
			Data:   "",
		})
		return
	}

	// TODO 检查room是否属于当前appid
	// 根据id发送消息
	service.NewContactService().SayTextToContact(msgSendRequest.Id, msgSendRequest.Msg)

	// 返回结果
	ctx.JSON(http.StatusOK, CommonResponse{
		Status: 20000,
		Msg:    "success",
		Data:   "",
	})
}

/*func VipMsgSend(ctx *gin.Context) {
	// 获取GET参数
	type MsgSendRequest struct {
		Id  string `form:"id" binding:"required"`
		Msg string `form:"msg" binding:"required"`
	}
	var msgSendRequest MsgSendRequest
	if err := ctx.ShouldBindQuery(&msgSendRequest); err != nil {
		ctx.JSON(http.StatusOK, CommonResponse{
			Status: 10002,
			Msg:    "POST JSON参数有误，请查看接口文档",
			Data:   "",
		})
		return
	}

	// 根据id发送消息
	service.NewContactService().SayTextToContact(msgSendRequest.Id, msgSendRequest.Msg)

	// 返回结果
	ctx.JSON(http.StatusOK, CommonResponse{
		Status: 20000,
		Msg:    "success",
		Data:   "",
	})
}*/
