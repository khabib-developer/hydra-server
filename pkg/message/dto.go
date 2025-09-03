package message

type SendMessageDto struct {
	Receiver string `json:"receiver"`
	Message  string	`json:"message"`
}

type ReceiveMessageDto struct {
	Sender string `json:"sender"`
	Message  string	`json:"message"`
}
