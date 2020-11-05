package storageHandler

type StorageInterface interface {
	Init(messagesReaderChannel chan Message)
	GetNextMessage() Message
}

type Message struct {
	Text string
	Unix int64
}