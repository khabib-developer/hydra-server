package file

type FileDto struct {
	ID		 string  `json:"id"`
	Receiver string	 `json:"receiver"`
	Sender   string  `json:"sender"`
	Filename string  `json:"filename"`
	Total    int64	 `json:"total"`
	Size     int64   `json:"size"`
}


type FileChunkDto struct {
	ID    string `json:"id"`
	Body  []byte `json:"body"`
	Index int64  `json:"index"`
}