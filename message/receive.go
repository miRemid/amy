package message

// CQResponse 相应基础
type CQResponse struct{
	Retcode	int							`json:"retcode"`
	Status	string						`json:"status"`
	Data	map[string]interface{}		`json:"data"`
}

// CQLoginInfo 登陆信息
type CQLoginInfo struct {
	CQResponse
	Data	CQLoginInfoData	`json:"data"`
}

// CQLoginInfoData 登陆数据
type CQLoginInfoData struct {
	UserID		int		`json:"user_id"`
	NickName	string	`json:"nickname"`
}

// CQCookiesInfo ck信息
type CQCookiesInfo struct {
	CQResponse
	Data	CookiesData	`json:"data"`
}

// CookiesData Data
type CookiesData struct {
	Cookies string	`json:"cookies"`
}

// CQCRFSTokenInfo crfsToken信息
type CQCRFSTokenInfo struct {
	CQResponse
	Data	CRFSToken	`json:"data"`
}

// CRFSToken token
type CRFSToken struct {
	Token int `json:"token"`
}

// CQRecodeInfo 语音消息
type CQRecodeInfo struct {
	CQResponse
	Data Recode `json:"data"`
}

// Recode 语音路径
type Recode struct {
	FilePath	string 	`json:"file"`
}

// CQImageInfo 图片信息
type CQImageInfo struct {
	CQResponse
	Data Image `json:"data"`
}

// Image 图片路径
type Image struct {
	FilePath	string	`json:"file"`
}

// CQHTTPStatus CQHTTP插件状态
type CQHTTPStatus struct {
	CQResponse
	Data struct{
		Initialized 	bool					`json:"app_initialized"`
		Enable			bool					`json:"app_enabled"`
		AppGood			bool					`json:"app_good"`
		Online			bool					`json:"online"`
		Good			bool					`json:"good"`
		PluginsGood		map[string]interface{}	`json:"plugins_good"`
	}	`json:"data"`
}

// CQHTTPVersion CQHTTP版本信息
type CQHTTPVersion struct {
	CQResponse
	Data struct {
		Directory	string	`json:"coolq_directory"`
		Edition		string	`json:"coolq_edition"`
		Version		string	`json:"plugin_version"`
		BuildNumber	int		`json:"plugin_build_number"`
		BuildConfig	string	`json:"plugin_build_configuration"`
	}	`json:"data"`
}






// CQGroupList 群组列表信息
type CQGroupList struct {
	CQResponse
	Data	[]GroupItem	`json:"data"`
}

// GroupItem 群组列表元素
type GroupItem struct {
	GroupID		int		`json:"group_id"`
	GroupName	string	`json:"group_name"`
}

// CQGroupMember 群组成员信息
type CQGroupMember struct {
	CQResponse
	Data GroupMemberItem	`json:"data"`
}

// CQGroupMemberList 群组成员列表信息
type CQGroupMemberList struct {
	CQResponse
	Data []GroupMemberItem	`json:"data"`
}

// GroupMemberItem 成员信息
type GroupMemberItem struct {
	GroupID				int		`json:"group_id"`					
	UserID				int		`json:"user_id"`
	NickName			string	`json:"nickname"`
	Card				string	`json:"card"`
	Sex					string	`json:"sex"`
	Age					uint	`json:"age"`
	Area				string	`json:"area"`
	JoinTime			int		`json:"join_time"`
	LastSentTime		int		`json:"last_sent_time"`
	Level				string	`json:"level"`
	Role				string	`json:"role"`
	UnFriendly			bool	`json:"unfriendly"`
	Title				string	`json:"title"`
	TitleExpireTime 	int		`json:"title_expire_time"`
	CardChangeable		bool	`json:"card_changeable"`
}

// CQStrangerInfo 他人信息
type CQStrangerInfo struct {
	CQResponse
	Data	StrangerInfoData	`json:"data"`
}

// StrangerInfoData 他人数据
type StrangerInfoData struct {
	CQLoginInfoData
	Sex	string	`json:"sex"`
	Age	int32	`json:"age"`
}