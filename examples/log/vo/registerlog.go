package vo

type RegisterLogVO struct {
	Time int64			`json:"time"`
	UserName string		`json:"username"`
	PlatName string		`json:"platName"`
	Channelid string	`json:"channelid"`
}

func (vo *RegisterLogVO)SetTime(time int64)  {
	vo.Time = time
}

func (vo *RegisterLogVO)SetUsername(username string)  {
	vo.UserName = username
}

func (vo *RegisterLogVO)SetPlatName(platName string)  {
	vo.PlatName = platName
}
func (vo *RegisterLogVO)SetChannelid(channelid string)  {
	vo.Channelid = channelid
}
