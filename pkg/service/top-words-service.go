package service

import (
	"sync"
)

type TopWordsService interface {
	GetWordFrequency() int64
	GetTopTenWords(text string) []byte
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
	WordFrequency map[string]int64
	TopWords      []TopWord
}

// Init initializes the TopWordsService for use.
func Init() TopWordsService {
	m := make(map[string]int64)
	s := make([]TopWord, 0)
	return &WordContainer{
		WordFrequency: m,
		TopWords:      s,
	}
}

// GetWordFrequency returns frequency of words
func (c *WordContainer) GetWordFrequency() int64 {
	var f int64

	for _, v := range c.TopWords {
		f = f + v.freq
	}

	return f
}

func (c *WordContainer) GetTopTenWords(text string) []byte {
	//TODO implement me
	panic("implement me")
}
