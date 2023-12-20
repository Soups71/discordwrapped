package bot

import (
	"discordwrapped/pkg/config"
	"discordwrapped/pkg/database"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var BotId string

func Start() {

	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		log.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		log.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)
	err = goBot.Open()

	if err != nil {
		log.Println(err.Error())
		return
	}
	goBot.UpdateGameStatus(0, config.BotPrefix)
	log.Println("Bot Started")
}

func initDataset(s *discordgo.Session, m *discordgo.MessageCreate) {
	currTime := time.Now()
	yearBefore := currTime.AddDate(-1, 0, 0)

	channels, _ := s.GuildChannels(m.GuildID)
	if len(channels) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Initializing Database... (This might take a LONG time...)")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Updating Database... (This might take a LONG time...)")
	}
	for _, chann := range channels {
		if chann.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		exists, err := database.CollectionExists(m.GuildID, chann.ID)
		if err != nil {
			log.Fatalln(err)
		}
		if !exists {
			log.Printf("Loading channel data for %s...\n", chann.Name)
			database.CreateCollection(m.GuildID, chann.ID)
			getChannelMessages(s, m.GuildID, chann.ID, yearBefore)
			log.Printf("Channel data for %s loaded!\n", chann.Name)
		} else {
			lastTime, err := database.MostRecent(m.GuildID, m.ChannelID)
			// fmt.Println(m.ChannelID)
			if err != nil {
				log.Printf("Loading channel data from past year for %s...\n", chann.Name)
				getChannelMessages(s, m.GuildID, chann.ID, yearBefore)
			} else {
				log.Printf("Loading new channel data for %s...\n", chann.Name)
				getChannelMessages(s, m.GuildID, chann.ID, lastTime)
			}

		}
	}

	s.ChannelMessageSend(m.ChannelID, "Database Ready!")
}

func generateServerReport(s *discordgo.Session, guildID string, chanID string) {
	test1 := database.GetUserMsgGuildActivity(s, guildID)
	test2 := database.GetUserGifGuildActivity(s, guildID)
	test3 := database.GetUserImageGuildActivity(s, guildID)
	s.ChannelMessageSend(chanID, "**DISCORD WRAPPED**: A look back on the server's activity over the past year\n\n")
	s.ChannelMessageSend(chanID, test1+test2+test3)
}
func generateChannelReport(s *discordgo.Session, guildID string, chanID string) {
	test1 := database.GetUserMsgChannelActivity(s, guildID, chanID)
	test2 := database.GetUserGifChannelActivity(s, guildID, chanID)
	test3 := database.GetUserImageChannelActivity(s, guildID, chanID)
	s.ChannelMessageSend(chanID, "**DISCORD WRAPPED**: A look back on the server's activity over the past year\n\n")
	s.ChannelMessageSend(chanID, test1+test2+test3)
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotId {
		return
	}
	if strings.Contains(m.Content, config.BotPrefix) {
		log.Printf("Message recieved from %s, \"%s\"\n", m.Author, m.Content)
		split := strings.Split(m.Content, " ")
		if len(split) != 2 {
			cwd, err := os.Getwd()
			fullPath := filepath.Join(string(cwd), "help.txt")
			file, err := os.Open(fullPath)
			if err != nil {
				log.Fatalln("Can't open the help file")
			}
			defer file.Close()
			content, err := ioutil.ReadAll(file)
			if err != nil {
				log.Fatalln("Can't read the help file")
			}
			s.ChannelMessageSend(m.ChannelID, string(content))
			return
		}
		if strings.Compare(split[0], config.BotPrefix) == 0 {

			initDataset(s, m)

			if strings.Compare(split[1], "channel") == 0 {
				generateChannelReport(s, m.GuildID, m.ChannelID)
				log.Printf("Providing channel report for %s\n", m.ChannelID)
			} else if strings.Compare(split[1], "server") == 0 {
				generateServerReport(s, m.GuildID, m.ChannelID)
				log.Printf("Providing server report for %s\n", m.GuildID)
			} else {

				log.Println("Command send could not be parsed")
				s.ChannelMessageSend(m.ChannelID, "It appears that you have not provided the appropriate parameters.")
			}

		}
	}

}

func getChannelMessages(s *discordgo.Session, guildID string, channelID string, yearBefore time.Time) error {
	database.ConnectDB(guildID)
	count := 0
	messages, err := s.ChannelMessages(channelID, 100, "", "", "")

	if err != nil {
		return err
	}
	// fmt.Println(len(messages), yearBefore)
	if len(messages) <= 0 {
		return nil
	}
	for messages[0].Timestamp.After(yearBefore) {
		for _, msg := range messages {
			if msg.Timestamp.After(yearBefore) {
				database.AddMessage(msg)
			}
		}
		count += 100
		messages, err = s.ChannelMessages(channelID, 100, messages[len(messages)-1].ID, "", "")
		if err != nil {
			log.Fatalln(err)
		}
		if len(messages) == 0 {
			break
		}
	}

	return nil
}
