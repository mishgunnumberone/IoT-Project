package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Temperature and Humidity"),
		tgbotapi.NewKeyboardButton("Water Level"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Led On"),
		tgbotapi.NewKeyboardButton("Led Off"),
		),

)

const ip  = "http://172.16.0.69:5000"

func main() {
	bot, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Additional goroutine for checking the water level and notify user about water level
	go checkWaterLevel(155662803, bot)

	for update := range updates {

		if update.Message == nil { // ignore any non-Message Updates

			continue
		}

			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! I'm an IoT bot, waiting for your orders")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			case "Temperature and Humidity":
				sendTemperatureAndHumidity(&update, bot)

			case "Water Level":
				sendWaterLevel(&update, bot)

			case "Led Off":
				turnLedOff(&update, bot)

			case "Led On":
				turnLedOn(&update, bot)
			}

		}

}

func sendTemperatureAndHumidity(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg string

	res, err := http.Get(ip + "/air")
	if err != nil {
		log.Print(err)
		sendTelegramConnectionErrorMessage(update, bot)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		sendTelegramReadingErrorMessage(update, bot)
		return
	}

	var dat map[string]int

	if err := json.Unmarshal(body, &dat); err != nil {
		log.Print(err)
		sendTelegramDecodingErrorMessage(update, bot)
		return
	}

	temp := dat["temperature"]

	hum := dat["humidity"]

	msg = fmt.Sprintf("🌡 Temperature: %d℃\n💧 Humidity: %d%%", temp, hum)

	sendTelegramMessage(msg, update, bot)
}

func sendWaterLevel(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg string

	res, err := http.Get(ip +"/water_level")

	if err != nil {
		sendTelegramConnectionErrorMessage(update, bot)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sendTelegramReadingErrorMessage(update, bot)
	}

	var dat map[string]int

	if err := json.Unmarshal(body, &dat); err != nil {
		sendTelegramDecodingErrorMessage(update, bot)
	}

	waterLevel := dat["waterLevel"]

	if waterLevel <= 300 {
		msg = fmt.Sprintf("No any leaks detected\nWater Level: %d", waterLevel)
	} else {
		msg = fmt.Sprintf("❗️ATTENTION❗ Water leak detected\nWater Level: %d", waterLevel)
	}

	sendTelegramMessage(msg, update, bot)
}

func turnLedOff(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg string

	res, err := http.Get(ip + "/turn_light_off")
	if err != nil {
		sendTelegramConnectionErrorMessage(update, bot)
	}

	if res.StatusCode == 200 {
		msg = "Led is turned off"
	} else {
		msg = "Request wasn't processed properly. Please try it later"
	}

	sendTelegramMessage(msg, update, bot)
}

func turnLedOn(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg string

	res, err := http.Get(ip + "/turn_light_on")
	if err != nil {
		sendTelegramConnectionErrorMessage(update, bot)
	}

	if res.StatusCode == 200 {
		msg = "💡 Led is turned on"
	} else {
		msg = "Request wasn't processed properly. Please try it later"
	}

	sendTelegramMessage(msg, update, bot)
}

func sendTelegramMessage(msg string, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var botMessage = tgbotapi.NewMessage(update.Message.Chat.ID, msg)
	bot.Send(botMessage)
}

func sendTelegramConnectionErrorMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg = "Can't connect to a sensor. Please try it later"
	sendTelegramMessage(msg, update, bot)
}

func sendTelegramReadingErrorMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg = "Error in reading the sensor data. Please try it later"
	sendTelegramMessage(msg, update, bot)
}

func sendTelegramDecodingErrorMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg = "Error in decoding the sensor data. Please try again later"
	sendTelegramMessage(msg, update, bot)
}


// Not ideal solution, right one is to have a db with a corresponding chatID and IoT Device
func checkWaterLevel(chatID int64, bot *tgbotapi.BotAPI) {
	for {
		time.Sleep(15 * time.Second)
		log.Printf("Enter a checkWaterLevel")
		var msg string

		res, err := http.Get(ip + "/water_level")

		if err != nil{
		return
	}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil{
		return
	}

		var dat map[string]int

		if err := json.Unmarshal(body, &dat); err != nil{
		return
	}

		waterLevel := dat["waterLevel"]

		if waterLevel > 300{
		msg = fmt.Sprintf("❗️ATTENTION❗ Water leak detected\nWater Level: %d", waterLevel)
	} else{
		return
	}

		var botMessage = tgbotapi.NewMessage(chatID, msg)
		bot.Send(botMessage)
	}
}
