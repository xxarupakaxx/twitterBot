package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo"
)

func CreateCRCToken(crcToken string) string {
	mac := hmac.New(sha256.New, []byte(cs)) //csはConsumer Secret
	mac.Write([]byte(crcToken))
	return "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func main() {
	e:=echo.New()
	e.POST("/callback", func(c echo.Context) error {
		return nil
	})
}