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
	s := make([]TopWord, 0)
	return &WordContainer{
		frequencyMap: m,
		topWords:     s,
	}
}

// GetTopTenWords returns the top ten words as response.
// text request --> validate -> process words in chunks -> updates freq map -> map to top word slice -> sort slice -> return response dto
func (c *WordContainer) GetTopTenWords(text dto.TextRequestDto) ([]dto.TopWordsResponseDto, lib.RestErr) {
	Init() // init or reset word container

	words, err := text.ValidateRequest()
	if err != nil {
		return nil, err
	}

	workers := setGoMaxProcs()
	wordChunks := len(words) / workers

	for i := 0; i < workers; i++ {
		c.wg.Add(1)
		go c.processWords(words[i*wordChunks : (i+1)*wordChunks]) // process words in chunks, calls wg done when finished
	}

	c.wg.Wait()
	c.mapToTopWordsSlice() // frequency map -> top words slice
	c.sortTopTen()         // sort the slice from highest to lowest frequency

	resp := c.makeTopTenResponseDTO() // make response dto, ready to return

	return resp, nil
}
