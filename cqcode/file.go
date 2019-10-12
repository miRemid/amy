package cqcode

import (
	"fmt"
)

// Image 图片cq码
// uri支持本地图片或url链接，同时支持base64编码，需base64://前缀
// cache为是否使用缓存
func Image(uri string, cache bool) string {
	if cache{
		return fmt.Sprintf("[CQ:image,file=%s]", uri)
	}
	return fmt.Sprintf("[CQ:image,cache=0,file=%s]", uri)
}

// Recode 语音cq码
// uri支持本地图片或url链接，同时支持base64编码，需base64://前缀
// cache为是否使用缓存
func Recode(uri string, cache bool) string {
	if cache{
		return fmt.Sprintf("[CQ:recode,file=%s]", uri)
	}
	return fmt.Sprintf("[CQ:recode,cache=0,file=%s]", uri)
}