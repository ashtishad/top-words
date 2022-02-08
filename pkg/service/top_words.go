package service

import (
	"github.com/ashtishad/top-words/internal/dto"
	"github.com/ashtishad/top-words/internal/lib"
	"sync"
)

type TopWordsService interface {
	GetTopTenWords(text dto.TextRequestDto) ([]dto.TopWordsResponseDto, lib.RestErr)
}

type TopWord struct {
	word string
	freq int64
}

// WordContainer has a slice of words and their frequencies,
// a mutex so that it can be safely accessed by multiple goroutines.
// wg is used to wait for all goroutines to finish before returning.
type WordContainer struct {
	mu           sync.Mutex
	wg           sync.WaitGroup
	frequencyMap map[string]int64
	topWords     []TopWord
}

// Init initializes/resets the TopWordsService for fresh use.
func Init() TopWordsService {
	m := make(map[string]int64)
	return &WordContainer{
		frequencyMap: m,
	}
}

// GetTopTenWords returns the top ten words as response.
// text request --> validate -> process words in chunks -> updates freq map -> map to top word slice -> sort slice -> return response dto
func (c *WordContainer) GetTopTenWords(text dto.TextRequestDto) ([]dto.TopWordsResponseDto, lib.RestErr) {
	words, err := text.ValidateRequest()
	if err != nil {
		return nil, err
	}

	workers := getMaxNumCPUs()                // set max number of go routines to use process concurrently
	chunkSize := getChunkSize(words, workers) // calculate total words processed by each go routine(workers)

	c.wg.Add(workers)

	for i := 0; i < workers; i++ {
		// set which word chunk to be processed by individual worker
		start, end := getChunkRange(chunkSize, i, len(words))

		// launch a go routine to process the chunk
		go func(words []string) {
			defer c.wg.Done()
			for _, word := range words {
				if word == "a" || len(word) >= 2 {
					c.pushToFrequencyMap(word)
				}
			}
		}(words[start:end])
	}

	c.wg.Wait()

	c.toTopWordsSlice() // frequency map -> top words slice
	c.sortWords()       // sort the slice from highest to lowest frequency,if frequency is same, sort alphabetically

	resp := c.makeTopTenResponseDTO() // make top ten response dto, ready to return

	return resp, nil
}
