package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"os"
)

func CreateCRCToken(crcToken string) string {
	mac := hmac.New(sha256.New, []byte(os.Getenv("")))
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

	e.POST("/callback", func(c echo.Context) error {
		return nil
	})
	e.Start(":"+port)
}