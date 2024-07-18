package thirdparty

import (
	"errors"
	"github.com/tidwall/gjson"
	"yufuture-gpt/common/utils"
)

func getAccessToken(url, appid, secret string) (string, error) {
	url = url + "?grant_type=client_credential" + "&appid=" + appid + "&secret=" + secret
	var responseData interface{}
	err := utils.HTTPGetAndParseJSON(url, &responseData)
	if err != nil {
		return "", err
	}
	responseDataString, err := utils.AnyToString(responseData)
	if err != nil {
		return "", err
	}
	accessToken := gjson.Get(responseDataString, "access_token").String()
	if accessToken == "" {
		return "", errors.New("获取access_token失败: " + responseDataString)
	}
	return accessToken, nil
}

type Scene struct {
	SceneID int `json:"scene_id"`
}
type ActionInfo struct {
	Scene Scene `json:"scene"`
}
type Param struct {
	ExpireSeconds int        `json:"expire_seconds"`
	ActionName    string     `json:"action_name"`
	ActionInfo    ActionInfo `json:"action_info"`
}

func getTicket(url, accessToken string) (string, error) {
	url = url + "?access_token=" + accessToken
	requestData := Param{
		ExpireSeconds: 120,
		ActionName:    "QR_SCENE",
		ActionInfo: ActionInfo{
			Scene: Scene{
				SceneID: 123,
			},
		},
	}
	var responseData interface{}
	err := utils.HTTPPostAndParseJSON(url, requestData, &responseData)
	if err != nil {
		return "", err
	}
	responseDataString, err := utils.AnyToString(responseData)
	if err != nil {
		return "", err
	}
	ticket := gjson.Get(responseDataString, "ticket").String()
	if ticket == "" {
		return "", errors.New("获取ticket失败: " + responseDataString)
	}
	return ticket, nil
}

// GetWechatLoginQrCode 获取微信登录二维码
func GetWechatLoginQrCode(accessTokenUrl, appid, secret, ticketUrl, qrCodeUrl string) (ticketQrCodeUrl string, ticket string, err error) {
	// 获取access_token
	token, err := getAccessToken(accessTokenUrl, appid, secret)
	if err != nil {
		return "", "", err
	}
	// 获取ticket
	ticket, err = getTicket(ticketUrl, token)
	if err != nil {
		return "", "", err
	}
	// 获取登录二维码
	ticketQrCodeUrl = qrCodeUrl + "?ticket=" + ticket
	return ticketQrCodeUrl, ticket, nil
}
