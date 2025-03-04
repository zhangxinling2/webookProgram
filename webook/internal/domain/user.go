package domain

import "time"

type User struct {
	Id           int64
	Email        string
	Password     string
	Phone        string
	NickName     string
	Birth        time.Time
	Introduction string
	CTime        time.Time
}
type UserInfo struct {
	NickName     string
	Birth        string
	Introduction string
}
