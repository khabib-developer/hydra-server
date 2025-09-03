package channel

type ChannelDto struct {
	Name    string `json:"name"`
	Members int    `json:"members"`
	Mine    bool   `json:"mine"`
}

type ChannelMessageDto struct {
	Channel string `json:"channel"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}