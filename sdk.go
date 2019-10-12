package amy

import (
	cqmsg "github.com/miRemid/amy/message"

	"net/http"
	"io/ioutil"
	"log"
	"fmt"
	"errors"
	// "bytes"
	"encoding/json"
)

const (
	// SendReqFailed 发送请求失败
	SendReqFailed	= iota	
	// JSONMarshalFailed 消息格式化错误
	JSONMarshalFailed
)

// URL 生成发送的URL
func (api *API) URL(router string, asyc bool) string {
	url := fmt.Sprintf("http://%s:%d/%s", api.CQHTTP, api.CQPORT, router)
	if asyc{
		url = url + "_async"
	}	
	log.Printf("Url: %s\n", url)
	return url
}

// Send 发送消息
// msg 发送的消息
// url 发送的api
// cr 返回的消息管道
func (api *API) Send(msg cqmsg.CQMessage, url string, cr CResponse) {
	if cr == nil {
		cr = make(CResponse, 1)
	}
	var response Response
	// 获取数据	
	reader, err := msg.Value()
	if err != nil {
		response.Error = err
		cr <- response
		return
	}
	// 生成Request
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		response.Error = err
		cr <- response
		return
	}
	// 设置头消息
	req.Header.Add("Content-Type", "application/json")
	if api.Token != ""{
		req.Header.Add("Authorization", "Bearer " + api.Token)
	}
	if err != nil {		
		response.Error = err
		cr <- response		
		return
	}
	// 发送请求
	resp, err := api.Client.Do(req)
	if err != nil {
		response.Error = err
		cr <- response		
		return
	}	
	defer resp.Body.Close()
	// 读取回应消息
	// 如果正常Post
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		response.Reader = body
		response.Error = err
		cr <- response		
		return
	}
	response.Error = errors.New("http client do request failed")
	cr <- response			
}

// sendHandler 消息相应处理助手
func (api *API) sendHandler(res Response) (response cqmsg.CQMessageID, ok bool){
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s\n", SendReqFailed, res.Error.Error())
		return
	}
	if err := json.Unmarshal(res.Reader, &response); err != nil {
		log.Printf("Amy Error[%d]:%s\n", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// SendRaw 原生发送消息
func (api *API) SendRaw(msg cqmsg.CQRawMsg, async bool) (response cqmsg.CQMessageID, ok bool) {
	c := make(CResponse, 1)
	go api.Send(msg, api.URL("send_msg", async), c)	
	return api.sendHandler(<- c)
}

// SendPrivateMsg 发送私人消息
func (api *API) SendPrivateMsg(msg cqmsg.CQPrivateMsg, async bool) (response cqmsg.CQMessageID, ok bool) {
	c := make(CResponse, 1)
	go api.Send(msg, api.URL("send_private_msg", async), c)		
	return api.sendHandler(<- c)
}

// SendGroupMsg 发送群组消息
func (api *API) SendGroupMsg(msg cqmsg.CQGroupMsg, async bool) (response cqmsg.CQMessageID, ok bool) {
	c := make(CResponse, 1)
	go api.Send(msg, api.URL("send_group_msg", async), c)
	return api.sendHandler(<- c)
}

// SendDiscussMsg 发送讨论组消息
func (api *API) SendDiscussMsg(msg cqmsg.CQDiscussMsg, async bool) (response cqmsg.CQMessageID, ok bool) {
	c := make(CResponse, 1)
	go api.Send(msg, api.URL("send_discuss_msg", async), c)
	return api.sendHandler(<- c)
}




// noEchoHandler 无响应请求助手
func (api *API) noEchoHandler(msg cqmsg.CQJSON, router string, async bool) bool {	
	c := make(CResponse, 1)
	go api.Send(msg, api.URL(router, async), c)
	res := <- c
	return res.Error == nil	
}

// CQHTTPRestart 重启CQHTTP插件
func (api *API) CQHTTPRestart(delay int, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"delay": delay,
	}, "set_restart_plugin", async)
}

// CQHTTPCleanData 清理数据目录
func (api *API) CQHTTPCleanData(dir string, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"data_dir": dir,
	}, "clean_data_dir ", async)
}

// CQHTTPCleanLog 清理插件日志
func (api *API) CQHTTPCleanLog(async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{},"clean_plugin_log", async)
}

// DeleteMsg 撤回消息
func (api *API) DeleteMsg(id int, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"message_id": id,
	}, "delete_msg", async)	
}

// SendLike 发送好友赞
func (api *API) SendLike(userid uint, times int, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"user_id": userid,
		"times":	times,
	}, "send_like", async)
}

// GroupKick 群组踢人
func (api *API) GroupKick(groupid uint, userid uint, reject, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,
		"user_id":	userid,
		"reject_add_request": reject,
	}, "set_group_kick", async)
}

// GroupBan 群组禁言
func (api *API) GroupBan(groupid uint, userid uint, duration int, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,
		"user_id":	userid,
		"duration": duration,
	}, "set_group_ban", async)
}

// GroupAnoymousBan 匿名禁言
// Flag 参数为消息上报中的参数
func (api *API) GroupAnoymousBan(groupid uint, flag string, duration int, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,
		"flag":	flag,
		"duration": duration,
	}, "set_group_anonymous_ban", async)
}

// GroupWholeBan 全体禁言
func (api *API) GroupWholeBan(groupid uint, enable, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,
		"enable": enable,
	}, "set_group_whole_ban", async)
}

// SetGroupAdmin 设置管理员
func (api *API) SetGroupAdmin(groupid uint, userid uint, enable, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,
		"user_id": userid,
		"enable": enable,
	}, "set_group_admin", async)
}

// SetGroupAnoymous 设置群组匿名
func (api *API) SetGroupAnoymous(groupid uint, enable, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,
		"enable": enable,
	}, "set_group_anonymous", async)
}

// SetGroupCard 设置群名片
func (api *API) SetGroupCard(groupid uint, userid uint, card string, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,
		"user_id": userid,
		"card": card,
	}, "set_group_card", async)
}

// SetGroupLeave 退群
// isDismiss 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散
func (api *API) SetGroupLeave(groupid uint, isDismiss, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,		
		"is_dismiss": isDismiss,
	}, "set_group_leave", async)
}

// SetGroupSpecialTitle 设置专属头衔
func (api *API) SetGroupSpecialTitle(groupid uint, userid uint, title string, duration int, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"group_id": groupid,		
		"user_id": userid,
		"special_title": title,
		"duration": duration,
	}, "set_group_special_title", async)
}

// SetDiscussLeave 退讨论组
func (api *API) SetDiscussLeave(discussid uint, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"discuss_id": discussid,		
	}, "set_discuss_leave", async)
}

// SetFriendAddRequest 处理加好友请求
func (api *API) SetFriendAddRequest(flag, remark string, approve, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"flag": flag,
		"remark": remark,
		"approve": approve,
	}, "set_friend_add_request", async)
}

// SetGroupAddRequest 处理加群请求
func (api *API) SetGroupAddRequest(flag, subtype, reason string, approve, async bool) bool {
	return api.noEchoHandler(cqmsg.CQJSON{
		"flag": flag,
		"reason": reason,
		"approve": approve,
		"sub_type": subtype,
	}, "set_group_add_request ", async)
}







// GetHandler 获取信息助手
func (api *API) getHandler(msg cqmsg.CQJSON, router string, async bool) Response{
	c := make(CResponse, 1)
	go api.Send(msg, api.URL(router, async), c)
	return <- c
}

// GetLoginInfo 获取登陆账号信息
func (api *API) GetLoginInfo(async bool) (response cqmsg.CQLoginInfo, ok bool){
	res := api.getHandler(cqmsg.CQJSON{}, "get_login_info", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, &response)
	if err !=nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// GetStrangerInfo 获取他人信息
func (api *API) GetStrangerInfo(userid uint, cache, async bool) (response cqmsg.CQStrangerInfo, ok bool){
	res := api.getHandler(cqmsg.CQJSON{
		"user_id": userid,
		"no_cache": cache,
	}, "get_stranger_info", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, &response)
	if err !=nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// GetGroupList 获取群组列表
func (api *API) GetGroupList(async bool) (response cqmsg.CQGroupList, ok bool){
	res := api.getHandler(cqmsg.CQJSON{}, "get_group_list", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, &response)
	if err !=nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// GetGroupMemberInfo 获取群成员信息
func (api *API) GetGroupMemberInfo(groupid uint, userid uint, cache, async bool) (response cqmsg.CQGroupMember, ok bool){
	res := api.getHandler(cqmsg.CQJSON{
		"group_id": groupid,
		"user_id": userid,
		"cache": cache,
	}, "get_group_member_info", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, &response)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// GetGroupMemberList 获取群成员列表
func (api *API) GetGroupMemberList(groupid uint, async bool) (response cqmsg.CQGroupMemberList, ok bool){
	res := api.getHandler(cqmsg.CQJSON{
		"group_id": groupid,		
	}, "get_group_member_list", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, &response)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// GetCookies 获取cookies
func (api *API) GetCookies(domain string, async bool) (response cqmsg.CQCookiesInfo, ok bool) {
	res := api.getHandler(cqmsg.CQJSON{
		"domain": domain,
	}, "get_cookies", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, &response)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// GetCSRFToken 获取CSRF认证token
func (api *API) GetCSRFToken(async bool) (response cqmsg.CQCRFSTokenInfo, ok bool) {
	res := api.getHandler(cqmsg.CQJSON{}, "get_csrf_token", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, response)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}	
	return response, true
}

// GetRecord 获取语音
func (api *API) GetRecord(file, format string, fullPath, async bool) (response cqmsg.CQRecodeInfo, ok bool) {
	res := api.getHandler(cqmsg.CQJSON{
		"file": file,
		"out_format": format,
		"full_path": fullPath,
	}, "get_record", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}	
	err := json.Unmarshal(res.Reader, &response)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// GetImage 获取图片
func (api *API) GetImage(file string, async bool) (response cqmsg.CQImageInfo, ok bool) {
	res := api.getHandler(cqmsg.CQJSON{
		"file": file,		
	}, "get_image", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}	
	err := json.Unmarshal(res.Reader, &response)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// CanSendImage 检查是否能发送图片
func (api *API) CanSendImage(async bool) (ok bool) {
	res := api.getHandler(cqmsg.CQJSON{}, "can_send_image", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	var can cqmsg.CQResponse
	err := json.Unmarshal(res.Reader, &can)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return can.Data["yes"].(bool)
}

// GetStatus 检查插件运行状态
func (api *API) GetStatus(async bool) (response cqmsg.CQHTTPStatus, ok bool) {
	res := api.getHandler(cqmsg.CQJSON{}, "get_status", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, &response)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}

// GetVersion 获取版本信息
func (api *API) GetVersion(async bool) (response cqmsg.CQHTTPVersion, ok bool) {
	res := api.getHandler(cqmsg.CQJSON{}, "get_version_info", async)
	if res.Error != nil {
		log.Printf("Amy Error[%d]:%s", SendReqFailed, res.Error.Error())
		return
	}
	err := json.Unmarshal(res.Reader, &response)
	if err != nil {
		log.Printf("Amy Error[%d]:%s", JSONMarshalFailed, err.Error())
		return
	}
	return response, true
}