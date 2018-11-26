package log

import (
	"github.com/yhhaiua/log4go"
	"github.com/yhhaiua/log4go/examples/log/vo"
	"time"
)

func LogRegister()  {
	var vo vo.RegisterLogVO
	vo.SetChannelid("1")
	vo.SetPlatName("1")
	vo.SetTime(time.Now().Unix())
	vo.SetUsername("hhy")
	log4go.InfoLog(REGISTER,&vo)
}
