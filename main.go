package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	helpMessage = "以下のコマンドを入力することでマイクラのサーバ内の状態がわかるぺこ！\n" +
		"* `login-list`\n\t現在ログインしているユーザのIDが分かるぺこ！\n" +
		"* `user-list` \n\t過去にログインしたユーザのIDが分かるぺこ！\n" +
		"* `last-died` \n\t直近で死亡したユーザのIDが分かるぺこ！\n"
)

var (
	Token             = "Bot "
	BotName           = ""
	BotAppName        = ""
	stopBot           = make(chan bool)
	vcsession         *discordgo.VoiceConnection
	HelloWorld        = "!helloworld"
	ChannelVoiceJoin  = "!vcjoin"
	ChannelVoiceLeave = "!vcleave"
)

func initCredentials() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Not Found .env")
	}
	BotName = os.Getenv("CLIENT_ID")
	Token += os.Getenv("TOKEN")
	BotAppName = os.Getenv("BOTNAME")
}

func main() {
	initCredentials()
	discord, err := discordgo.New(Token)
	discord.Token = Token
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}
	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}
	defer discord.Close()

	fmt.Println("Listening...")
	<-stopBot //プログラムが終了しないようロック
	discord.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Message.Content)
	if m.Message.Author.Username == BotAppName {
		return
	}
	if m.Message.Content == "help" {
		sendMessage(s, m.ChannelID, helpMessage)
		return
	}
	if m.Message.Content == "login-list" {
		sendMessage(s, m.ChannelID, joinPeko("ログインしている人のID", 1))
		return
	}
	if m.Message.Content == "user-list" {
		sendMessage(s, m.ChannelID, joinPeko("過去にログインしたユーザ", 1))
		return
	}
	sendMessage(s, m.ChannelID, joinPeko("そのコマンドは無い", 2))
}

func joinPeko(msg string, typ int) string {
	if typ == 1 {
		return msg + "ぺこ！"
	} else if typ == 2 {
		return msg + "ぺこぉ.."
	}
	return ""
}

func sendMessage(s *discordgo.Session, channelID, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}
