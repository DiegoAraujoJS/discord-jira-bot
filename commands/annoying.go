package commands

import (

	"github.com/bwmarrin/discordgo"
)

func ReactionsHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
    if m.UserID == s.State.User.ID {
        return
    }

    if m.Emoji.Name == "👀" {
        s.MessageReactionAdd(m.ChannelID, m.MessageID, "👀")
    }
}
