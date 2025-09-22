package captcha

import (
	"context"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"

	"time"
)

const captchaKeyPrefix = "captcha:"

type captchaStore struct {
	redis *redis.Client
}

func New(redis *redis.Client) *base64Captcha.Captcha {
	return base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, &captchaStore{redis})
}

func (c *captchaStore) Set(id string, value string) error {
	return c.redis.
		SetEx(context.Background(), captchaKeyPrefix+id, value, time.Minute*1).
		Err()
}

func (c *captchaStore) Get(id string, clear bool) string {
	var value string
	var err error
	if clear {
		value, err = c.redis.GetDel(context.Background(), captchaKeyPrefix+id).Result()
	} else {
		value, err = c.redis.Get(context.Background(), captchaKeyPrefix+id).Result()
	}
	if err != nil {
		return ""
	}
	return value
}

func (c *captchaStore) Verify(id, answer string, clear bool) bool {
	value := c.Get(id, clear)
	return answer == value
}
