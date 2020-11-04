package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

const prefix = "고"

func init() {
	flag.StringVar(&Token, "t", "NzczMzQ1NTA3MjIwMjU4ODI2.X6H4IA.JgQ9DsxpMjoieR8JukG0L1nDWPo", "")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.UpdateStatus(0, "1개 서버에 참여")

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func commandCheck(userInp string, needCmd string) bool {
	if userInp == prefix + needCmd {
		return true
	}

	return false
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if commandCheck(m.Content, "차트") {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed {
			Title: "Title",
			Description: "Description",
			Author: &discordgo.MessageEmbedAuthor {
				Name: "Author Name",
				URL: "https://github.com/JominJun",
				IconURL: "https://imgur.com/roQhXpQ.png",
			},
			Image: &discordgo.MessageEmbedImage {
				URL: "https://imgur.com/roQhXpQ.png",
			},
			Fields: []*discordgo.MessageEmbedField {
				&discordgo.MessageEmbedField {
					Name: "Field Name1",
					Value: "Field Value1",
				},
				&discordgo.MessageEmbedField {
					Name: "Field Name2",
					Value: "Field Value2",
				},
			},
			Footer: &discordgo.MessageEmbedFooter {
				Text: "Footer Text",
				IconURL: "https://imgur.com/roQhXpQ.png",
			},
			Color: 123456,
			Type: discordgo.EmbedTypeRich,
		})

		text := fmt.Sprintf("현재 채널: %s", m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, text)
	}
}