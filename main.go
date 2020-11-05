package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/dustin/go-humanize"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

const prefix = "고"

// CompanyResponse for Company API
type CompanyResponse struct {
	Result []struct {
		Seq					int
		Name				string
		Description	string
		StockValue	int
	}
	Status	int
}

func init() {
	flag.StringVar(&Token, "t", "NzczMzQ1NTA3MjIwMjU4ODI2.X6H4IA.gn0aRgX641pCVJxWQHC_B_NsZHY", "")
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
		req, err := http.NewRequest("GET", "http://api.localhost:8081/v1/company", nil)
		if err != nil {
			panic(err)
		}

		req.Header.Add("Authorization", "키")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		bytes, _ := ioutil.ReadAll(resp.Body)

		var res CompanyResponse
		json.Unmarshal(bytes, &res)
		
		var companyFields []*discordgo.MessageEmbedField
		for _, company := range(res.Result) {
			companyField := discordgo.MessageEmbedField {
				Name: company.Name,
				Value: fmt.Sprintf("%s원", humanize.Comma(int64(company.StockValue))),
				Inline: true,
			}

			companyFields = append(companyFields, &companyField)
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed {
			Author: &discordgo.MessageEmbedAuthor {
				Name: "GO! CHART",
				URL: "https://github.com/JominJun",
				IconURL: "https://imgur.com/roQhXpQ.png",
			},
			Fields: companyFields,
			Color: 3463297,
			Type: discordgo.EmbedTypeRich,
		})
	}
}