package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.Bot || m.Content != "!help" {
        return
    }
    commands := "!upstatus - Golpea a nuestros servidores para ver si reaccionan.\n\n" +
    "LW-12 - Consigue el ticket 12 del proyecto \"LW\". Se puede cambiar \"LW\" por cualquier otro proyecto y 12 por cualquier otro ticket. Por ejemplo: PWA-4 consigue el ticket 4 del proyecto \"PWA\".\n\n" +
    "!tickets-LW-Error - Busca en el proyecto \"LW\" los tickets que estén en estado \"Error\". Se puede cambiar \"LW\" por cualquier otro proyecto y \"Error\" por cualquier otro estado. Por ejemplo: !tickets-PWA-Proceso busca en el proyecto \"PWA\" los tickets que estén en estado \"Proceso\".\n\n" +
    "!config - Muestra los datos del servidor de discord y del canal actual.\n\n" +
    "!help - Muestra este mensaje."
    s.ChannelMessageSend(m.ChannelID, commands)
}

func GetGuildData(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.Bot || m.Content != "!config" {return}
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Guild ID: %v\nChannel ID: %v", m.GuildID, m.ChannelID))
}
