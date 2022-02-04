package service

import (
	"github.com/ashtishad/top-words/internal/dto"
	"runtime"
	"sort"
)

// processWords processes the words in the given slice.
func (c *WordContainer) processWords(words []string) {
	defer c.wg.Done()

	for _, word := range words {
		c.mu.Lock()
		c.frequencyMap[word]++
		c.mu.Unlock()
	}
}

// mapToTopWordsSlice maps the word frequency to the top ten words.
func (c *WordContainer) mapToTopWordsSlice() {
	for k, v := range c.frequencyMap {
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

// makeTopTenResponseDTO creates the top ten response DTO.
func (c *WordContainer) makeTopTenResponseDTO() []dto.TopWordsResponseDto {
	resp := make([]dto.TopWordsResponseDto, 0, 10)
	for _, w := range c.topWords {
		resp = append(resp, dto.TopWordsResponseDto{
			Word:      w.word,
			Frequency: w.freq,
		})
	}
	return resp
}

// setGoMaxProcs sets the maximum number of CPUs that can be executing concurrently.
// returns the number, default is 1.
func setGoMaxProcs() int {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
	return n
}
