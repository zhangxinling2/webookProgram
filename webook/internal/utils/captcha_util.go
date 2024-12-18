package utils

import "github.com/mojocn/base64Captcha"

// 验证码工具类
type StringCaptcha struct {
	Captcha *base64Captcha.Captcha
}

// 创建验证码
func NewCaptcha() *StringCaptcha {
	// store
	store := base64Captcha.DefaultMemStore

	// 包含数字和字母的字符集
	source := "123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

	// driver
	driver := base64Captcha.NewDriverString(
		80,     // height int
		240,    // width int
		6,      // noiseCount int
		1,      // showLineOptions int
		4,      // length int
		source, // source string
		nil,    // bgColor *color.RGBA
		nil,    // fontsStorage FontsStorage
		nil,    // fonts []string
	)

	captcha := base64Captcha.NewCaptcha(driver, store)
	return &StringCaptcha{
		Captcha: captcha,
	}
}

// 生成验证码
func (stringCaptcha *StringCaptcha) Generate() (string, string, string) {
	id, b64s, answer, _ := stringCaptcha.Captcha.Generate()
	return id, b64s, answer
}

// 验证验证码
func (stringCaptcha *StringCaptcha) Verify(id string, answer string) bool {
	return stringCaptcha.Captcha.Verify(id, answer, true)
}
