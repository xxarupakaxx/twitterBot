package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/xxarupakaxx/twitterBot/model"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	HELPMESSAGE                              = `コマンド一覧
help でコマンド一覧を表示できます
これから追加予定
makePlayList 指定したチャンネルのプレイリストを作るor 指定した動画たちをプレイリストにする 
upload 動画をアップロード
delete 動画を削除
mychannel 自分のアカウント情報
googleDriveに動画を保存する
二つの動画をつなげて返す
できたら画像加工もできたらいいね
weather code で天気予報取得
位置情報を取得して近くのお店を表示`
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("oioio")
	}

	e:=echo.New()
	e.Use(middleware.Logger())
	e.POST("/callback", lineHandler)
	e.Start(":"+port)

}

func lineHandler(c echo.Context) error {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Response().WriteHeader(http.StatusBadRequest)
		}else {
			c.Response().WriteHeader(http.StatusInternalServerError)
		}
		return err
	}

	events,err:=bot.ParseRequest(c.Request())
	if err != nil {
		log.Fatal(err)
	}


	for _,event:=range events{
		/*if event.Type==linebot.EventTypeFollow {
			userId:=event.Source.UserID
			user,err:=bot.GetProfile(userId).Do()
			if err!=nil{
				log.Fatalf("Failed in getting user :%v",err)
			}



			userData:=domain.User{
				Id:         event.Source.UserID,
				DisplayName: user.DisplayName,
				IdType:     string(event.Source.Type),
				Timestamp:  event.Timestamp,
				ReplyToken: event.ReplyToken,
			}
			err=db.QueryRow("SELECT * from user where user.id=$1",userData.Id).Scan(&userData.Id, &userData.DisplayName,&userData.IdType,&userData.Timestamp,&userData.ReplyToken)
			if err != nil {
				_,err=db.Exec("INSERT INTO user VALUES (?,?,?,?,?)",userData.Id,userData.DisplayName,userData.IdType,userData.Timestamp,userData.ReplyToken)
				if err != nil {
					log.Fatalf("Couldnot add user:%v",err)
				}
			}

			text:=user.DisplayName+"さん\n"+HELPMESSAGE
			if _,err:=bot.PushMessage(userId,linebot.NewTextMessage(text),linebot.NewStickerMessage("8522","16581267")).Do();err!=nil{
				log.Fatalf("Failed in Pushing message:%v",err)
			}

		}*/
		switch message := event.Message.(type) {
		case *linebot.LocationMessage:
			model.SendRestoInfo(bot,event)
		case *linebot.TextMessage:
			user,_:=bot.GetProfile(event.Source.UserID).Do()
			if message.Text == "help" {
				text:=user.DisplayName+"さん\n"+HELPMESSAGE
				if _,err:=bot.ReplyMessage(event.ReplyToken,linebot.NewTextMessage(text)).Do();err!=nil{
					log.Fatalf("Failed in Replying message:%v",err)
				}

			}
			if message.Text == "quick" {
				quickRep(bot,event)
			}

			if strings.Contains(message.Text, "weather") {
				msg:=message.Text
				cityName:=msg[len("weather "):]
				data:=model.PrefCode(cityName)
				if data == nil {
					log.Fatalf("data is nil %v",err)
				}
				model.SendWeather(bot,event,data.ID)
			}

		case *linebot.VideoMessage:
			if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("未実装")).Do(); err != nil {
				log.Fatalf("Failed in getting url:%v",err)
			}
		}
	}
	query:=c.QueryParam("item")
	return c.String(http.StatusOK,query)
	//return c.String(http.StatusOK,"OK")

}

func quickRep(bot *linebot.Client, event *linebot.Event) {

	resp:=linebot.NewQuickReplyItems(
		&linebot.QuickReplyButton{
			ImageURL: "https://www.jma.go.jp/jma/kishou/img/logo2.jpg",
			Action:   linebot.NewMessageAction("天気", "weather 140010"),
		},
		&linebot.QuickReplyButton{
			ImageURL: "https://cdn.icon-icons.com/icons2/729/PNG/512/twitter_icon-icons.com_62751.png",
			Action:   linebot.NewDatetimePickerAction("time","&item=1111","time","00:00","23:59","00:00"),
		},
		&linebot.QuickReplyButton{
			ImageURL: "https://upload.wikimedia.org/wikipedia/commons/3/39/%E3%81%BD-red.png",
			Action:   linebot.NewPostbackAction("label","&item=1112","","displayText"),
		},
		&linebot.QuickReplyButton{
			ImageURL: "https://img.icons8.com/material/452/camera--v1.png",
			Action:   linebot.NewCameraAction("camera"),
		},
		&linebot.QuickReplyButton{
			ImageURL: "https://icooon-mono.com/i/icon_16245/icon_162450_256.png",
			Action:   linebot.NewCameraRollAction("CameraRoll"),
		},
		&linebot.QuickReplyButton{
			ImageURL: "https://icooon-mono.com/i/icon_15694/icon_156940_256.png",
			Action:   linebot.NewLocationAction("Location"),
		},
	)
	re:=linebot.NewTextMessage("helperer").WithQuickReplies(resp)
	if _, err := bot.ReplyMessage(event.ReplyToken,re).Do(); err != nil {
		log.Fatalf("QuickError:%v",err)
	}


}