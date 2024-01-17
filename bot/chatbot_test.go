package bot_test

import (
	"testing"

	"github.com/ondbyte/chatbot/bot"
	"github.com/ondbyte/chatbot/bot/adapters/logic"
	"github.com/ondbyte/chatbot/bot/adapters/storage"
	"github.com/ondbyte/chatbot/bot/corpus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestBot(t *testing.T) {
	assert := assert.New(t)
	// training
	storeFile := t.TempDir() + "/test.gob"
	store, err := storage.NewSeparatedMemoryStorage(storeFile)
	if !assert.NoError(err) {
		return
	}
	// create training files
	data := `
categories:
- API
- NLP
conversations:
- - carrots 5 kg?
  - '{"carrots":{"qty":5,"unit":"kgs"}}'
- - i need carrot 1 kg?
  - '{"carrot":{"qty":1,"unit":"kgs"}}'
- - beetroot 5 kg?
  - '{"beetroot":{"qty":5,"unit":"kgs"}}'
- - cuccumber?
  - '{"cuccumber":{"qty":0,"unit":""}}'`
	c := &corpus.Corpus{}
	err = yaml.Unmarshal([]byte(data), c)
	if !assert.NoError(err) {
		return
	}
	b := &bot.ChatBot{
		PrintMemStats:  true,
		Trainer:        bot.NewCorpusTrainer(store),
		StorageAdapter: store,
		LogicAdapter:   logic.NewClosestMatch(store, 1),
	}
	err = b.Train([]*corpus.Corpus{c})
	if !assert.NoError(err) {
		return
	}
	answers := b.GetResponse("cuccumber")
	if !assert.NotEmpty(answers) {
		return
	}
	assert.Equal(answers[0].Content, `{"cuccumber":{"qty":0,"unit":""}}`)
}
