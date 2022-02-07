package service

import (
	"github.com/ashtishad/top-words/internal/dto"
	"github.com/ashtishad/top-words/internal/lib"
	Map "github.com/ashtishad/top-words/pkg/ordered_map"
	"log"
	"sync"
)

type TopWordsService interface {
	GetTopTenWords(text dto.TextRequestDto) ([]dto.TopWordsResponseDto, lib.RestErr)
}

type TopWord struct {
	word string
	freq int
}

// WordContainer has a slice of words and their frequencies,
// a mutex so that it can be safely accessed by multiple goroutines.
// wg is used to wait for all goroutines to finish before returning.
type WordContainer struct {
	mu           sync.Mutex
	wg           sync.WaitGroup
	frequencyMap *Map.OrderedMap
	topWords     []TopWord
}

// Init initializes/resets the TopWordsService for fresh use.
func Init() TopWordsService {
	s := make([]TopWord, 0)
	m := Map.New()
	//m := make(map[string]int64)
	return &WordContainer{
		frequencyMap: m,
		topWords:     s,
	}
}

// GetTopTenWords returns the top ten words as response.
// text request --> validate -> process words in chunks -> updates freq map -> map to top word slice -> sort slice -> return response dto
func (c *WordContainer) GetTopTenWords(text dto.TextRequestDto) ([]dto.TopWordsResponseDto, lib.RestErr) {
	words, err := text.ValidateRequest()
	if err != nil {
		return nil, err
	}

	workers := getMaxNumCPUs()           // set max number of go routines to use process concurrently
	chunks := calcChunks(words, workers) // calculate word chunks to process by each worker

	// single unit of work for each worker
	processWord := func(words []string) {
		for _, word := range words {
			if word == "a" || len(word) >= 2 {
				c.pushToFrequencyMap(word)
			}
		}
		c.wg.Done()
	}

	// workers will process words concurrently
	for i := 0; i < workers; i++ {
		start := i * chunks
		end := start + chunks
		if end > len(words) {
			end = len(words)
		}
		c.wg.Add(1)
		go processWord(words[start:end])
	}

	c.wg.Wait()

	c.toTopWordsSlice() // frequency map -> top words slice
	log.Printf("Frequency: %v", c.frequencyMap)
	log.Println("----------------------------------------------------")
	log.Printf("Top words: %v", c.topWords)
	c.sortWordsStable() // sort the slice from highest to lowest frequency

	resp := c.makeTopTenResponseDTO() // make top ten response dto, ready to return

	return resp, nil
}
