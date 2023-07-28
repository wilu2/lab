package encrypt

import (
	"crypto/md5"
	"fmt"
	"io"

	"gitlab.intsig.net/textin-gateway/pkg/stringx"
)

const salt = "&^@*(" // 固定盐

// PassWordEncrypt 使用固定盐和随机盐和账号密码做 md5
func PassWordEncrypt(password, username, saltRand string) (string, string, error) {
	h := md5.New()
	if saltRand == "" {
		saltRand = stringx.RandString(8) // 随机盐
	}

	val := salt + password + saltRand + username
	if _, err := io.WriteString(h, val); err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), saltRand, nil
}
