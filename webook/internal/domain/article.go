package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Status  ArticleStatus
	Author  Author
}
type Author struct {
	Id   int64
	Name string
}
type ArticleStatus uint8

const (
	ArticleUnknown ArticleStatus = iota
	ArticleUnPublished
	ArticlePublished
	ArticlePrivate
)

func (a ArticleStatus) IsValid() bool {
	return a != ArticleUnknown
}
func (a ArticleStatus) ToUint8() uint8 {
	return uint8(a)
}
