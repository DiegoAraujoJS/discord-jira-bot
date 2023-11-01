package commands

import (

	"github.com/bwmarrin/discordgo"
)

func ReactionsHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
    if m.Member.User.Bot {
        return
    }

    if m.Emoji.Name == "👀" {
        s.MessageReactionAdd(m.ChannelID, m.MessageID, "👀")
    }
}
