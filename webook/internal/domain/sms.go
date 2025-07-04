package domain

type AsyncSms struct {
	Id       int64
	Biz      string
	Args     []string
	Numbers  []string
	RetryMax int8
}
