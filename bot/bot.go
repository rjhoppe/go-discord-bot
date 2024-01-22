package bot

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"rjhoppe/go-discord-bot/config"
	"time"

	"github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *discordgo.Session
var buffer = make([][]byte, 0)

// Need a space between "Bot" and end of quotation
func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running")

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	textChanID := "1198295595962605699"
	voiceChanID := "1057016031022940291"
	guildID := "1057016030561570907"

	if m.Author.ID == BotID {
		return
	}

	// Need to decouple these from the messageHandler
	go eventWait(s, 0, textChanID, "Bing Bong Bing Bong, 15 mins until Bing Bong")
	go eventWait(s, 300, textChanID, "Bing Bong Bing Bong 10 minutes left!")
	go eventWait(s, 600, textChanID, "Oh boy, oh boy, only 5 more minutes until Bing Bong!")
	go eventWait(s, 828, textChanID, "HERE I COME!!!!! BINGBONGBINGBONGBINGBONG")
	go loadSound()
	go grandFinale(s, 830, guildID, voiceChanID)

	// Fix this
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func eventWait(s *discordgo.Session, t int, textChanID string, message string) {

	convertT := time.Duration(t) * time.Second
	time.Sleep(convertT)
	s.ChannelMessageSend(textChanID, message)
}

func loadSound() error {
	file, err := os.Open("bingbong.dca")
	if err != nil {
		fmt.Println("Error opening audio file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

func grandFinale(s *discordgo.Session, t int, guildID string, voiceChanID string) (err error) {

	vc, err := s.ChannelVoiceJoin(guildID, voiceChanID, false, false)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)

	// Send the buffer data
	fmt.Println("Playing sound now...")
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	fmt.Println("Sound complete...")
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	fmt.Println("Disconnecting now")
	vc.Disconnect()

	return nil

}
