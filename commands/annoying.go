package commands

import (
	"math/rand"
	"time"

	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func Annoying(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }
    if user_one, ok := utils.Secondary["user_one"]; ok && m.Author.Username == user_one && randomInt(0, 100) < 5 {
        s.MessageReactionAdd(m.ChannelID, m.ID, "👀")
    }
}

func ReactionsHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
    if m.UserID == s.State.User.ID {
        return
    }

    if m.Emoji.Name == "👀" {
        s.MessageReactionAdd(m.ChannelID, m.MessageID, "👀")
    }
}
