package common

import "time"

type ReplyProto struct {
	Status        int           `json:"status"`        //状态 0正常，小于0
	BackSinceTime time.Duration `json:"backSinceTime"` //后端时间
}
