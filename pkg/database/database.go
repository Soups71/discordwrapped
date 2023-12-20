package database

import (
	"context"
	"discordwrapped/pkg/config"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB     *mongo.Database
	client *mongo.Client
)

type Emoji struct {
	ID            string `bson:"id"`
	Name          string `bson:"name"`
	RequireColons bool   `bson:"require_colons"`
	Managed       bool   `bson:"managed"`
	Animated      bool   `bson:"animated"`
	Available     bool   `bson:"available"`
}
type MessageReactions struct {
	UserID    string `bson:"user_id"`
	MessageID string `bson:"message_id"`
	Emoji     Emoji  `bson:"emoji"`
}
type Message struct {
	ID        string             `bson:"id"`
	Content   string             `bson:"content"`
	Reactions []MessageReactions `bson:"reactions"`
	Timestamp time.Time          `bson:"timestamp"`
}

type result_ranker struct {
	results []indi_result
}
type indi_result struct {
	name  string
	score int64
}

func ConnectDB(dbID string) int {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DBConn))
	if err != nil {
		return -1
	}
	DB = client.Database(dbID)
	return 0
}

func CollectionExists(guildID string, collectionName string) (bool, error) {
	ConnectDB(guildID)
	collections, err := DB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return false, err
	}
	for _, col := range collections {
		if col == collectionName {
			return true, nil
		}
	}
	return false, nil
}

func CreateCollection(guildID string, channelID string) error {
	ConnectDB(guildID)
	err := DB.CreateCollection(context.TODO(), channelID)
	return err
}

func AddMessage(message *discordgo.Message) {
	coll := DB.Collection(message.ChannelID)
	result, err := coll.InsertOne(context.TODO(), message)
	if err != nil {
		fmt.Println("Fuck", result)
	}
}
func MostRecent(guildID string, channelID string) (time.Time, error) {
	ConnectDB(guildID)
	collection := DB.Collection(channelID)
	options := options.FindOne().SetSort(bson.D{{"timestamp", -1}})
	var result Message
	err := collection.FindOne(context.TODO(), bson.M{}, options).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return time.Time{}, fmt.Errorf("collection is empty")
		}
		return time.Time{}, err
	}

	return result.Timestamp, nil

}
func formatResp(title string, result []indi_result) string {
	sort.Slice(result, func(i, j int) bool {
		return result[i].score > result[j].score
	})
	fmt_str := title

	rank := 1
	for _, users := range result {
		fmt_str += fmt.Sprint(rank) + ". " + users.name + ": " + fmt.Sprint(users.score) + "\n"

		rank += 1
	}
	fmt_str += "\n"
	return fmt_str
}

func GetUserMsgGuildActivity(s *discordgo.Session, guildID string) string {
	ConnectDB(guildID)
	// fmt.Println("Here")
	result := []indi_result{}
	channels, _ := s.GuildChannels(guildID)
	users, _ := s.GuildMembers(guildID, "", 1000)
	for _, user := range users {
		if user.User.Bot {
			continue
		}
		filter := bson.M{"author.id": user.User.ID}
		var total_count int64
		total_count = 0
		for _, chann := range channels {
			if chann.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			coll := DB.Collection(chann.ID)
			count, err := coll.CountDocuments(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
			}
			total_count += count
		}
		if total_count != 0 {
			person := indi_result{
				name:  user.User.Username,
				score: total_count,
			}
			result = append(result, person)
		}
	}

	title := "**Number of messages sent:**\n"
	return formatResp(title, result)
}
func GetUserGifGuildActivity(s *discordgo.Session, guildID string) string {
	ConnectDB(guildID)
	// fmt.Println("Here")
	result := []indi_result{}
	channels, _ := s.GuildChannels(guildID)
	users, _ := s.GuildMembers(guildID, "", 1000)
	for _, user := range users {
		if user.User.Bot {
			continue
		}
		filter := bson.M{"author.id": user.User.ID,
			"embeds": bson.M{"$elemMatch": bson.M{"type": "gifv"}}}
		var total_count int64
		total_count = 0
		for _, chann := range channels {
			if chann.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			coll := DB.Collection(chann.ID)
			count, err := coll.CountDocuments(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
			}
			total_count += count
		}
		if total_count != 0 {
			person := indi_result{
				name:  user.User.Username,
				score: total_count,
			}
			result = append(result, person)
		}

	}

	title := "**Number of GIFs sent:**\n"
	// fmt.Println("Here")
	return formatResp(title, result)

}
func GetUserImageGuildActivity(s *discordgo.Session, guildID string) string {
	ConnectDB(guildID)
	result := []indi_result{}
	channels, _ := s.GuildChannels(guildID)
	users, _ := s.GuildMembers(guildID, "", 1000)
	for _, user := range users {
		if user.User.Bot {
			continue
		}
		filter := bson.M{"author.id": user.User.ID,
			"attachments": bson.M{"$elemMatch": bson.M{"contenttype": bson.M{"$regex": "^image"}}}}
		var total_count int64
		total_count = 0
		for _, chann := range channels {
			if chann.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			coll := DB.Collection(chann.ID)
			count, err := coll.CountDocuments(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
			}
			total_count += count
		}

		if total_count != 0 {
			person := indi_result{
				name:  user.User.Username,
				score: total_count,
			}
			result = append(result, person)
		}
	}

	title := "**Number of Images sent:**\n"

	return formatResp(title, result)

}

func GetUserMsgChannelActivity(s *discordgo.Session, guildID string, channel_id string) string {
	ConnectDB(guildID)
	result := []indi_result{}
	users, _ := s.GuildMembers(guildID, "", 1000)
	// fmt.Println(users, channeles, errs)
	for _, user := range users {
		if user.User.Bot {
			continue
		}
		filter := bson.M{"author.id": user.User.ID}
		var total_count int64
		total_count = 0
		coll := DB.Collection(channel_id)
		count, err := coll.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		total_count += count
		if total_count != 0 {
			person := indi_result{
				name:  user.User.Username,
				score: total_count,
			}
			result = append(result, person)
		}
	}

	title := "**Number of messages sent:**\n"

	return formatResp(title, result)
}

func GetUserGifChannelActivity(s *discordgo.Session, guildID string, channel_id string) string {
	ConnectDB(guildID)
	result := []indi_result{}
	users, _ := s.GuildMembers(guildID, "", 1000)
	// fmt.Println(users, channeles, errs)
	for _, user := range users {
		if user.User.Bot {
			continue
		}
		filter := bson.M{"author.id": user.User.ID,
			"embeds": bson.M{"$elemMatch": bson.M{"type": "gifv"}}}
		var total_count int64
		total_count = 0

		coll := DB.Collection(channel_id)
		// Count the number of messages matching the filter
		count, err := coll.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		total_count += count
		// fmt.Println(total_count)

		if total_count != 0 {
			person := indi_result{
				name:  user.User.Username,
				score: total_count,
			}
			result = append(result, person)
		}
	}

	title := "**Number of GIFs sent:**\n"

	return formatResp(title, result)

}

func GetUserImageChannelActivity(s *discordgo.Session, guildID string, channel_id string) string {
	ConnectDB(guildID)
	result := []indi_result{}
	users, _ := s.GuildMembers(guildID, "", 1000)
	for _, user := range users {
		if user.User.Bot {
			continue
		}
		filter := bson.M{"author.id": user.User.ID,
			"attachments": bson.M{"$elemMatch": bson.M{"contenttype": bson.M{"$regex": "^image"}}}}
		var total_count int64
		total_count = 0
		coll := DB.Collection(channel_id)
		count, err := coll.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		total_count += count

		if total_count != 0 {
			person := indi_result{
				name:  user.User.Username,
				score: total_count,
			}
			result = append(result, person)
		}
	}
	title := "**Number of Images sent:**\n"

	return formatResp(title, result)

}

// func GetUserReactionsGuildActivity(s *discordgo.Session, guildID string) {

// 	guilds := s.State.Guilds
// 	for _, gui := range guilds {
// 		channels, _ := s.GuildChannels(gui.ID)
// 		users, _ := s.GuildMembers(gui.ID, "", 1000)
// 		// fmt.Println(users, channeles, errs)
// 		for _, user := range users {
// 			if user.User.Bot {
// 				continue
// 			}
// 			filter := bson.M{"author.id": user.User.ID, "reactions": bson.M{"$ne": nil}}
// 			// var total_count int64
// 			// total_count = 0
// 			var emojis []string
// 			for _, chann := range channels {
// 				if chann.Type != discordgo.ChannelTypeGuildText {
// 					continue
// 				}

// 				coll := DB.Collection(chann.ID)
// 				// Count the number of messages matching the filter
// 				// var result Message

// 				// err := coll.FindOne(context.TODO(), filter).Decode(&result)
// 				// if err != nil {
// 				// 	log.Fatal(err)
// 				// }
// 				// fmt.Println(result.Reactions[0].Emoji.Name)
// 				cursor, err := coll.Find(context.TODO(), filter)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				defer cursor.Close(context.TODO())
// 				var messages []Message
// 				if err := cursor.All(context.TODO(), &messages); err != nil {
// 					log.Fatal(err)
// 				}
// 				// fmt.Println(messages)

// 				// // Print the matching documents
// 				// fmt.Println("Matching documents:")

// 				for _, msg := range messages {
// 					// fmt.Println(msg.Reactions[0].Emoji.Name)
// 					emojis = append(emojis, msg.Reactions[0].Emoji.Name)
// 					// fmt.Printf("Message ID: %s\n", msg.ID)
// 					// Access other fields as needed
// 				}
// 				// total_count += count
// 				// fmt.Println(total_count)
// 			}
// 			emojiList := emojis

// 			emojiCounts := make(map[string]int)
// 			for _, emoji := range emojiList {
// 				emojiCounts[emoji]++
// 			}

// 			// Create a slice of counts to sort
// 			var counts []int
// 			for _, count := range emojiCounts {
// 				counts = append(counts, count)
// 			}

// 			// Sort the counts in descending order
// 			sort.Sort(sort.Reverse(sort.IntSlice(counts)))

// 			// Print the sorted list of emojis by count
// 			fmt.Println("Most used emojis to react to messages by", user.User.Username, ": ")
// 			for _, count := range counts {
// 				for emoji, emojiCount := range emojiCounts {
// 					if emojiCount == count {
// 						fmt.Printf("%s: %d\n", emoji, count)
// 						// Mark the emoji as processed to avoid duplicate printing
// 						delete(emojiCounts, emoji)
// 						break
// 					}
// 				}
// 			}

// 			// fmt.Printf("Number of Images sent by user %s: %d\n", user.User.Username, total_count)
// 		}
// 	}
// }
