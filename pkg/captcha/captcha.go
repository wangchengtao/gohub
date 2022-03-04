package captcha

import (
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"sync"

	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

var once sync.Once

var internalCaptcha *Captcha

func NewCaptcha() *Captcha {
	once.Do(func() {
		// 初始化对象
		internalCaptcha = &Captcha{}

		// 使用全局 redis 对象, 并配置 key
		store := RedisStore{
			RedisClient: redis.Redis,
			KeyPrefix:   config.GetString("app.name" + ":captcha"),
		}

		// 配置 base64Captcha 驱动信息
		driver := base64Captcha.NewDriverDigit(
			config.GetInt("captcha.height"),
			config.GetInt("captcha.width"),
			config.GetInt("captcha.length"),
			config.GetFloat("captcha.maxskew"), //数字最大倾斜度
			config.GetInt("captcha.dotcount"),
		)

		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, &store)
	})

	return internalCaptcha
}

func (c *Captcha) GenerateCaptcha() (id string, b64s string, err error) {
	return c.Base64Captcha.Generate()
}

func (c *Captcha) VerifyCaptcha(id, answer string) (match bool) {
	if !app.IsProduction() && id == "uphicoo" {
		return true
	}

	return c.Base64Captcha.Verify(id, answer, false)
}
