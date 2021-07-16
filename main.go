package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"os"
)

func CreateCRCToken(crcToken string) string {
	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	mac.Write([]byte(crcToken))
	return "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("oioio")
	}

	e:=echo.New()
	e.Use(middleware.Logger())
	e.POST("/twitter_webhook",  HandlerCrcCheck)
	e.POST("/callback", func(c echo.Context) error {
		return nil
	})
	e.Start(":"+port)
}
func HandlerCrcCheck(c echo.Context) error{
	// Requestを受ける
	req := GetCrcCheckRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// CrcTokenを生成し、Responseに詰める
	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	mac.Write([]byte(req.CrcToken))
	res := GetCrcCheckResponse{
		Token: "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil)),
	}
	// Responseを返す
	return c.JSON(http.StatusOK, res)
}

type GetCrcCheckRequest struct {
	CrcToken string `json:"crc_token" form:"crc_token" binding:"required"`
}

type GetCrcCheckResponse struct {
	Token string `json:"response_token"`
}