package main

import (
	"bytes"
	"net/http"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/gin-gonic/gin"
)

type Greet struct{
	Data	string	`json:"data"`
}

var matchtext = [4]string{"lighton", "lightoff", "aircon_on", "aircon_off"}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("PORT must be set")
	}

	bot, err := linebot.New(
			os.Getenv("CHANNEL_SECRET"),
			os.Getenv("CHANNEL_TOKEN"),
		)
	if err != nil {
		fmt.Println(err)
		return
	}

	beebottetoken := os.Getenv("BEETOKEN")

	route := gin.Default()

	route.POST("/post", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage{
				message := event.Message.(*linebot.TextMessage)
				if message != nil {
					text_d := message.Text
					goodtext := inarray(matchtext, text_d)

					if goodtext {
						postbeebotte(text_d, beebottetoken)
					}
					replytext := fmt.Sprintf("Hello")
					switch text_d {
						case "lighton":
							replytext := fmt.Sprintf("Light ON")
						case "lightoff":
							replytext := fmt.Sprintf("Light OFF")
						case "aircon_on":
							replytext := fmt.Sprintf("Aircon ON")
						case "aircon_off":
							replytext := fmt.Sprintf("Aircon OFF")
						default:
							replytext := fmt.Sprintf("That command is not found")
						}

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replytext)).Do(); err != nil {
						fmt.Println(err)
					}
				}
			}
		}


	})
	route.Run(":" + port)




}

func inarray(lis [4]string, text string) bool{
	for _, value := range lis {
		if value == text {
			return true
		}
	}
	return false
}

func postbeebotte(text, token string) {
	greet := Greet{
		Data:	text,
	}

	greetjson, _ := json.Marshal(greet)
	rep, err := http.NewRequest("POST", "http://api.beebotte.com/v1/data/publish/Home_NET/Line", bytes.NewBuffer([]byte(greetjson)))

	if err != nil {
		fmt.Println(err)
	}
	rep.Header.Set("Content-Type", "application/json")
	rep.Header.Add("X-Auth-Token", token)

	fmt.Println(rep.Body)

	client := new(http.Client)
	resp, err := client.Do(rep)

	if err != nil {
		fmt.Println(err)
	}


	io.Copy(os.Stdout, resp.Body)

	defer resp.Body.Close()
}
