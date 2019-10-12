package cqcode

import (
	"strconv"
)

// ClientMusic 发送平台音乐
// client为音乐平台，目前支持qq、163、xiami
// id为对应音乐平台的数字音乐id
func ClientMusic(client string, id int) string {
	return CqCode("music", CQParams{
		"type": client,
		"id": strconv.Itoa(id),
	})
}

// Music 发送自定义音乐分享
// url: 分享链接，点击分享后进入的音乐页面
// audio: 音频链接
// title: 标题
// content: 简介
// image: 封面链接，为空显示默认图片
func Music(url, audio, title, content, image string) string {
	return CqCode("music", CQParams{
		"type": "custom",
		"url": url,
		"audio": audio,
		"title": title,
		"content": content,
		"image": image,
	})
}

// Share 发送链接分享
// url: 分享的链接
// title: 标题
// content: 简介
// image: 图片链接，为空显示默认图片
func Share(url, title, content, image string) string {
	return CqCode("share", CQParams{
		"url": url,
		"title": title,
		"content": content,
		"image": image,
	})
}