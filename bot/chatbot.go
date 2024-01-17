package bot

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ondbyte/chatbot/bot/adapters/input"
	"github.com/ondbyte/chatbot/bot/adapters/logic"
	"github.com/ondbyte/chatbot/bot/adapters/output"
	"github.com/ondbyte/chatbot/bot/adapters/storage"
	"github.com/ondbyte/chatbot/bot/corpus"
)

const mega = 1024 * 1024

type ChatBot struct {
	PrintMemStats  bool
	InputAdapter   input.InputAdapter
	LogicAdapter   logic.LogicAdapter
	OutputAdapter  output.OutputAdapter
	StorageAdapter storage.StorageAdapter
	Trainer        Trainer
}

func (chatbot *ChatBot) Train(data []*corpus.Corpus) error {
	start := time.Now()
	defer func() {
		fmt.Printf("Elapsed: %s\n", time.Since(start))
	}()

	if chatbot.PrintMemStats {
		go func() {
			for {
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				fmt.Printf("Alloc = %vm\nTotalAlloc = %vm\nSys = %vm\nNumGC = %v\n\n",
					m.Alloc/mega, m.TotalAlloc/mega, m.Sys/mega, m.NumGC)
				time.Sleep(5 * time.Second)
			}
		}()
	}

	if err := chatbot.Trainer.Train(data); err != nil {
		return err
	} else {
		return chatbot.StorageAdapter.Sync()
	}
}

func (chatbot *ChatBot) GetResponse(text string) []logic.Answer {
	if chatbot.LogicAdapter.CanProcess(text) {
		return chatbot.LogicAdapter.Process(text)
	}

	return nil
}
