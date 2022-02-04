package service

import (
	"github.com/ashtishad/top-words/internal/dto"
	"github.com/ashtishad/top-words/internal/lib"
	"sort"
	"sync"
)

type TopWordsService interface {
	GetTopTenWords(text dto.TextRequestDto) (dto.TopWordsResponseDto, lib.RestErr)
}

type TopWord struct {
	word string
	freq int64
}

// WordContainer has a slice of words and their frequencies,
// a mutex so that it can be safely accessed by multiple goroutines.
// wg is used to wait for all goroutines to finish before returning.
type WordContainer struct {
	mu            sync.Mutex
	wg            sync.WaitGroup
	wordFrequency map[string]int64
	topWords      []TopWord
}

// Init initializes the TopWordsService for use.
func Init() TopWordsService {
	m := make(map[string]int64)
	s := make([]TopWord, 0)
	return &WordContainer{
		wordFrequency: m,
		topWords:      s,
	}
}

func (c *WordContainer) GetTopTenWords(text dto.TextRequestDto) (dto.TopWordsResponseDto, lib.RestErr) {
	words, err := text.ValidateRequest()
	if err != nil {
		return dto.TopWordsResponseDto{}, err
	}

	workers := lib.SetGoMaxProcs()
	wordChunks := len(words) / workers

	for i := 0; i < workers; i++ {
		c.wg.Add(1)
		go c.processWords(words[i*wordChunks : (i+1)*wordChunks])
	}

	c.mapToTopTen()
	c.sortTopTen()

	resp := dto.TopWordsResponseDto{
		Words:     c.topWords,
		Frequency: c.getWordFrequency(),
	}
	return resp, nil
}

// processWords processes the words in the given slice.
func (c *WordContainer) processWords(words []string) {
	defer c.wg.Done()

	for _, word := range words {
		c.mu.Lock()
		c.wordFrequency[word]++
		c.mu.Unlock()
	}
}

// mapToTopTen maps the word frequency to the top ten words.
func (c *WordContainer) mapToTopTen() {
	for k, v := range c.wordFrequency {
		c.topWords = append(c.topWords, TopWord{
			word: k,
			freq: v,
		})
	}
}

// sortTopTen sorts the top ten words in descending order.
func (c *WordContainer) sortTopTen() {
	sort.Slice(c.topWords, func(i, j int) bool {
		return c.topWords[i].freq > c.topWords[j].freq
	})
}

// getWordFrequency returns the word frequency.
func (c *WordContainer) getWordFrequency() []int64 {
	freq := make([]int64, 0, 10)
	for _, v := range c.topWords {
		freq = append(freq, c.wordFrequency[v.word])
	}
	return freq
}
