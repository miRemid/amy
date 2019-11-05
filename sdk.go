package amy

import (
	cqmsg "github.com/miRemid/amy/message"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "bytes"
	"encoding/json"
)

const (
	// SendReqFailed 发送请求失败
	SendReqFailed = iota
	// JSONMarshalFailed 消息格式化错误
	JSONMarshalFailed
)

func (api *API) url(router string, asyc bool) string {
	url := fmt.Sprintf("http://%s:%d/%s", api.CQHTTP, api.CQPORT, router)
	if asyc {
		url = url + "_async"
	}
	log.Printf("Url: %s\n", url)
	return url
}

func (api *API) send(msg cqmsg.CQMessage, url string, cr CResponse) {
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
	if api.Token != "" {
		req.Header.Add("Authorization", "Bearer "+api.Token)
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
	response.Error = fmt.Errorf("请求API失败, url=%s, 状态码: %d", url, resp.StatusCode)
	cr <- response
}

// sendHandler 消息相应处理助手
func (api *API) sendHandler(res Response) (response cqmsg.CQMessageID, err error) {
	if res.Error != nil {
		err = res.Error
		return
	}
	if e := json.Unmarshal(res.Reader, &response); e != nil {
		err = e
		return
	}
	return response, nil
}

// Send 发送消息，根据flag值判断发送类型
// to 目标方id， msg
func (api *API) Send(to int, msg interface{}, auto, async bool, flag uint) (response cqmsg.CQMessageID, err error) {
	switch flag {
	case Private:
		msg := cqmsg.DefaultBuilder.PrivateMsg(to, msg, auto)
		return api.SendPrivateMsg(msg, async)
	case Group:
		msg := cqmsg.DefaultBuilder.GroupMsg(to, msg, auto)
		return api.SendGroupMsg(msg, async)
	case Discuss:
		msg := cqmsg.DefaultBuilder.DiscussMsg(to, msg, auto)
		return api.SendDiscussMsg(msg, async)
	default:
		return response, fmt.Errorf("wrong flag")
	}
}

// SendRaw 原生发送消息
func (api *API) SendRaw(msg cqmsg.CQRawMsg, async bool) (response cqmsg.CQMessageID, err error) {
	c := make(CResponse, 1)
	go api.send(msg, api.url("send_msg", async), c)
	return api.sendHandler(<-c)
}

// SendPrivateMsg 发送私人消息
func (api *API) SendPrivateMsg(msg cqmsg.CQPrivateMsg, async bool) (response cqmsg.CQMessageID, err error) {
	c := make(CResponse, 1)
	go api.send(msg, api.url("send_private_msg", async), c)
	return api.sendHandler(<-c)
}

// SendGroupMsg 发送群组消息
func (api *API) SendGroupMsg(msg cqmsg.CQGroupMsg, async bool) (response cqmsg.CQMessageID, err error) {
	c := make(CResponse, 1)
	go api.send(msg, api.url("send_group_msg", async), c)
	return api.sendHandler(<-c)
}

// SendDiscussMsg 发送讨论组消息
func (api *API) SendDiscussMsg(msg cqmsg.CQDiscussMsg, async bool) (response cqmsg.CQMessageID, err error) {
	c := make(CResponse, 1)
	go api.send(msg, api.url("send_discuss_msg", async), c)
	return api.sendHandler(<-c)
}

// noEchoHandler 无响应请求助手
func (api *API) noEchoHandler(msg cqmsg.CQMAP, router string, async bool) error {
	c := make(CResponse, 1)
	go api.send(msg, api.url(router, async), c)
	res := <-c
	return res.Error
}

// CQHTTPRestart 重启CQHTTP插件
func (api *API) CQHTTPRestart(delay int, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"delay": delay,
	}, "set_restart_plugin", async)
}

// CQHTTPCleanData 清理数据目录
func (api *API) CQHTTPCleanData(dir string, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"data_dir": dir,
	}, "clean_data_dir ", async)
}

// CQHTTPCleanLog 清理插件日志
func (api *API) CQHTTPCleanLog(async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{}, "clean_plugin_log", async)
}

// DeleteMsg 撤回消息
func (api *API) DeleteMsg(id int, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"message_id": id,
	}, "delete_msg", async)
}

// SendLike 发送好友赞
func (api *API) SendLike(userid uint, times int, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"user_id": userid,
		"times":   times,
	}, "send_like", async)
}

// GroupKick 群组踢人
func (api *API) GroupKick(groupid uint, userid uint, reject, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id":           groupid,
		"user_id":            userid,
		"reject_add_request": reject,
	}, "set_group_kick", async)
}

// GroupBan 群组禁言
func (api *API) GroupBan(groupid uint, userid uint, duration int, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id": groupid,
		"user_id":  userid,
		"duration": duration,
	}, "set_group_ban", async)
}

// GroupAnoymousBan 匿名禁言
// Flag 参数为消息上报中的参数
func (api *API) GroupAnoymousBan(groupid uint, flag string, duration int, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id": groupid,
		"flag":     flag,
		"duration": duration,
	}, "set_group_anonymous_ban", async)
}

// GroupWholeBan 全体禁言
func (api *API) GroupWholeBan(groupid uint, enable, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id": groupid,
		"enable":   enable,
	}, "set_group_whole_ban", async)
}

// SetGroupAdmin 设置管理员
func (api *API) SetGroupAdmin(groupid uint, userid uint, enable, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id": groupid,
		"user_id":  userid,
		"enable":   enable,
	}, "set_group_admin", async)
}

// SetGroupAnoymous 设置群组匿名
func (api *API) SetGroupAnoymous(groupid uint, enable, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id": groupid,
		"enable":   enable,
	}, "set_group_anonymous", async)
}

// SetGroupCard 设置群名片
func (api *API) SetGroupCard(groupid uint, userid uint, card string, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id": groupid,
		"user_id":  userid,
		"card":     card,
	}, "set_group_card", async)
}

// SetGroupLeave 退群
// isDismiss 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散
func (api *API) SetGroupLeave(groupid uint, isDismiss, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id":   groupid,
		"is_dismiss": isDismiss,
	}, "set_group_leave", async)
}

// SetGroupSpecialTitle 设置专属头衔
func (api *API) SetGroupSpecialTitle(groupid uint, userid uint, title string, duration int, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"group_id":      groupid,
		"user_id":       userid,
		"special_title": title,
		"duration":      duration,
	}, "set_group_special_title", async)
}

// SetDiscussLeave 退讨论组
func (api *API) SetDiscussLeave(discussid uint, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"discuss_id": discussid,
	}, "set_discuss_leave", async)
}

// SetFriendAddRequest 处理加好友请求
func (api *API) SetFriendAddRequest(flag, remark string, approve, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"flag":    flag,
		"remark":  remark,
		"approve": approve,
	}, "set_friend_add_request", async)
}

// SetGroupAddRequest 处理加群请求
func (api *API) SetGroupAddRequest(flag, subtype, reason string, approve, async bool) error {
	return api.noEchoHandler(cqmsg.CQMAP{
		"flag":     flag,
		"reason":   reason,
		"approve":  approve,
		"sub_type": subtype,
	}, "set_group_add_request ", async)
}

// GetHandler 获取信息助手
func (api *API) getHandler(msg cqmsg.CQMAP, router string, async bool) Response {
	c := make(CResponse, 1)
	go api.send(msg, api.url(router, async), c)
	return <-c
}

func gethandler(r Response, m interface{}) error {
	if r.Error != nil {
		return r.Error
	}
	err := json.Unmarshal(r.Reader, &m)
	if err != nil {
		return err
	}
	return nil
}

// GetLoginInfo 获取登陆账号信息
func (api *API) GetLoginInfo(async bool) (cqmsg.CQLoginInfo, error) {
	var response cqmsg.CQLoginInfo
	res := api.getHandler(cqmsg.CQMAP{}, "get_login_info", async)
	err := gethandler(res, &response)
	return response, err
}

// GetStrangerInfo 获取他人信息
func (api *API) GetStrangerInfo(userid uint, cache, async bool) (response cqmsg.CQStrangerInfo, err error) {
	res := api.getHandler(cqmsg.CQMAP{
		"user_id":  userid,
		"no_cache": cache,
	}, "get_stranger_info", async)
	err = gethandler(res, &response)
	return
}

// GetGroupList 获取群组列表
func (api *API) GetGroupList(async bool) (response cqmsg.CQGroupList, err error) {
	res := api.getHandler(cqmsg.CQMAP{}, "get_group_list", async)
	err = gethandler(res, &response)
	return
}

// GetGroupMemberInfo 获取群成员信息
func (api *API) GetGroupMemberInfo(groupid uint, userid uint, cache, async bool) (response cqmsg.CQGroupMember, err error) {
	res := api.getHandler(cqmsg.CQMAP{
		"group_id": groupid,
		"user_id":  userid,
		"cache":    cache,
	}, "get_group_member_info", async)
	err = gethandler(res, &response)
	return
}

// GetGroupMemberList 获取群成员列表
func (api *API) GetGroupMemberList(groupid uint, async bool) (response cqmsg.CQGroupMemberList, err error) {
	res := api.getHandler(cqmsg.CQMAP{
		"group_id": groupid,
	}, "get_group_member_list", async)
	err = gethandler(res, &response)
	return
}

// GetCookies 获取cookies
func (api *API) GetCookies(domain string, async bool) (response cqmsg.CQCookiesInfo, err error) {
	res := api.getHandler(cqmsg.CQMAP{
		"domain": domain,
	}, "get_coerries", async)
	err = gethandler(res, &response)
	return
}

// GetCSRFToken 获取CSRF认证token
func (api *API) GetCSRFToken(async bool) (response cqmsg.CQCRFSTokenInfo, err error) {
	res := api.getHandler(cqmsg.CQMAP{}, "get_csrf_terren", async)
	err = gethandler(res, &response)
	return
}

// GetRecord 获取语音
func (api *API) GetRecord(file, format string, fullPath, async bool) (response cqmsg.CQRecodeInfo, err error) {
	res := api.getHandler(cqmsg.CQMAP{
		"file":       file,
		"out_format": format,
		"full_path":  fullPath,
	}, "get_record", async)
	err = gethandler(res, &response)
	return
}

// GetImage 获取图片
func (api *API) GetImage(file string, async bool) (response cqmsg.CQImageInfo, err error) {
	res := api.getHandler(cqmsg.CQMAP{
		"file": file,
	}, "get_image", async)
	err = gethandler(res, &response)
	return
}

// CanSendImage 检查是否能发送图片
func (api *API) CanSendImage(async bool) (ok bool) {
	res := api.getHandler(cqmsg.CQMAP{}, "can_send_image", async)
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
func (api *API) GetStatus(async bool) (response cqmsg.CQHTTPStatus, err error) {
	res := api.getHandler(cqmsg.CQMAP{}, "get_status", async)
	err = gethandler(res, &response)
	return
}

// GetVersion 获取版本信息
func (api *API) GetVersion(async bool) (response cqmsg.CQHTTPVersion, err error) {
	res := api.getHandler(cqmsg.CQMAP{}, "get_version_info", async)
	err = gethandler(res, &response)
	return
}
