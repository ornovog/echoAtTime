package storageHandler

import(
	"github.com/go-redis/redis/v7"
	"time"
)

var (
	redisServerAddress = "localhost:6379"
	messagesSetName = "messagesSet"
	readingDelay = time.Duration(1)
)

type storageRedisHandler struct {
	messagesReader chan Message
	redisClient *redis.Client
}

func NewStorageHandler() storageRedisHandler {
	return storageRedisHandler{}
}

func (handler *storageRedisHandler) Init(messagesReaderChannel chan Message){
	handler.messagesReader = messagesReaderChannel

	options := redis.Options{Addr: redisServerAddress}
	handler.redisClient = redis.NewClient(&options)

	go handler.readAndStoreMessages()
}

func (handler *storageRedisHandler) readAndStoreMessages(){
	for message := range handler.messagesReader{
		handler.storeMessage(message)
	}
}
func (handler *storageRedisHandler) storeMessage(m Message){
	redisSortedSetMember := redis.Z{Member: m.Text, Score: float64(m.Unix)}
	handler.redisClient.ZAdd(messagesSetName, &redisSortedSetMember)
}


func (handler *storageRedisHandler) GetNextMessage() Message {
	redisSortedSetMember, err := handler.redisClient.BZPopMin(time.Second, messagesSetName).Result()
	for ;redisSortedSetMember == nil || err != nil ;{
		<-time.After(readingDelay * time.Second)
		if err!=nil{
			println(err.Error())
		}
		redisSortedSetMember, err = handler.redisClient.BZPopMin(time.Second, messagesSetName).Result()
	}

	unix := int64(redisSortedSetMember.Score)
	text := redisSortedSetMember.Member.(string)

	return  Message{Unix: unix, Text: text}
}
