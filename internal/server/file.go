package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/khabib-developer/hydra-server/pkg/file"
	"github.com/khabib-developer/hydra-server/pkg/websocket"
)

type FileMetadata struct {
	ID       string `json:"id"`
	Receiver string `json:"receiver"`
	Filename string `json:"filename"`
	Total    int64  `json:"total"`
	Size     int64  `json:"size"`
}

type PermanentFileData struct {
	ID       string
	Receiver *User
	Sender   *User
	Index    int64
	Total    int64
	Filename string
	Size     int64
	File     *os.File
}

func (server *Server) TransferFileMetadata(payload json.RawMessage, sender *User) {
	var fileMetadata FileMetadata

	if err := json.Unmarshal(payload, &fileMetadata); err != nil {
		// Handle error
		println(err)
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Unsupported type of file metadata")
		return
	}

	var receiver *User

	for _, user := range server.Users {
		if fileMetadata.Receiver == user.Username {
			receiver = user
		}
	}

	if receiver == nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Username not found")
		return
	}

	if sender.PermanentFile != nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "You can send only one file in the same time")
		return
	}

	if receiver.PermanentFile != nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Receiver can accept only one file in the same time.")
		return
	}

	permanentFileData := &PermanentFileData{
		ID:       fileMetadata.ID,
		Receiver: receiver,
		Sender:   sender,
		Index:    0,
		Total:    fileMetadata.Total,
		Filename: fileMetadata.Filename,
		Size:     fileMetadata.Size,
	}

	sender.PermanentFile = permanentFileData
	receiver.PermanentFile = permanentFileData

	fileDto := file.FileDto{
		ID:       fileMetadata.ID,
		Receiver: receiver.Username,
		Sender:   sender.Username,
		Filename: fileMetadata.Filename,
		Total:    fileMetadata.Total,
		Size:     fileMetadata.Size,
	}

	fileDtoJson, err := json.Marshal(fileDto)
	if err != nil {
		fmt.Println(err)
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Unsupported type of file metadata")
		return
	}

	server.sendMessage(receiver.SafeConn, websocket.MessageTypeFile, fileDtoJson)
}

func (server *Server) TransferFileChunk(payload json.RawMessage, sender *User) {

	fileChunkDto := file.FileChunkDto{}

	err := json.Unmarshal(payload, &fileChunkDto)

	if err != nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Unsupported type of file")
		return
	}

	if sender.PermanentFile == nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "User is not expecting file")
		return
	}

	if sender.PermanentFile.Receiver == nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "User is not expecting file")
		return
	}

	if sender.PermanentFile.ID != fileChunkDto.ID {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "You can send only one file in the same time")
		return
	}

	if sender.PermanentFile.Index >= sender.PermanentFile.Total {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "File has been sent completely")
		return
	}

	sender.PermanentFile.Index++

	server.sendMessage(sender.PermanentFile.Receiver.SafeConn, websocket.MessageTypeFileChunk, payload)

	if sender.PermanentFile.Total == sender.PermanentFile.Index {
		sender.PermanentFile.Receiver.PermanentFile = nil
		sender.PermanentFile = nil
	}
}

func (server *Server) CancelFileTransfer(_ json.RawMessage, sender *User) {
	if sender.PermanentFile == nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "You are not transfering file")
		return
	}

	sender.PermanentFile = nil
	sender.PermanentFile.Receiver.PermanentFile = nil

	if sender == sender.PermanentFile.Sender {
		server.sendRawMessage(sender.PermanentFile.Receiver.SafeConn, websocket.MessageTypeCancel, sender.PermanentFile.ID)
	} else {
		server.sendRawMessage(sender.PermanentFile.Sender.SafeConn, websocket.MessageTypeCancel, sender.PermanentFile.ID)
	}

	server.sendRawMessage(sender.PermanentFile.Sender.SafeConn, websocket.MessageTypeInfo, "File transfer has been canceled")
	server.sendRawMessage(sender.PermanentFile.Receiver.SafeConn, websocket.MessageTypeInfo, "File transfer has been canceled")
}
