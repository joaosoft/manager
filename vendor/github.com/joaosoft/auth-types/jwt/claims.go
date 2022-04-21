package jwt

import (
	"fmt"
	"strconv"
	"time"
)

type Claims map[string]interface{}

func (c Claims) Validate() bool {
	now := time.Now().Unix()

	if c.checkExpiredAt(now) && c.checkIssuedAt(now) && c.checkNotBefore(now) {
		return true
	}

	return false
}

func (c Claims) checkExpiredAt(now int64) bool {
	if value, ok := c[ClaimsExpireAtKey]; ok {
		intValue, err := strconv.ParseInt(fmt.Sprintf("%+v", value), 10, 64)
		if err != nil {
			return false
		}

		if intValue >= now {
			return false
		}
	}

	return true
}

func (c Claims) checkIssuedAt(now int64) bool {
	if value, ok := c[ClaimsIssuedAtKey]; ok {
		intValue, err := strconv.ParseInt(fmt.Sprintf("%+v", value), 10, 64)
		if err != nil {
			return false
		}

		if intValue >= now {
			return false
		}
	}
	return true
}

func (c Claims) checkNotBefore(now int64) bool {
	if value, ok := c[ClaimsNotBeforeKey]; ok {
		intValue, err := strconv.ParseInt(fmt.Sprintf("%+v", value), 10, 64)
		if err != nil {
			return false
		}

		if intValue < now {
			return false
		}
	}
	return true
}
