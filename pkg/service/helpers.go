package service

import (
	"github.com/ashtishad/top-words/internal/dto"
	"runtime"
	"sort"
)

// pushToFrequencyMap updates the frequency map with the word and its frequency.
func (c *WordContainer) pushToFrequencyMap(word string) {
	// avoid concurrent map read and write
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.frequencyMap[word] > 0 {
		c.frequencyMap[word]++
	} else {
		c.frequencyMap[word] = 1
	}
}

// toTopWordsSlice maps the word frequency to the top ten words.
func (c *WordContainer) toTopWordsSlice() {
	for k, v := range c.frequencyMap {
		c.topWords = append(c.topWords, TopWord{
			word: k,
			freq: v,
		})
	}
}

// sortWords sorts the top ten words by frequency highest to lowest.
// if frequency is same, sort by alphabetical order.
func (c *WordContainer) sortWords() {
	sort.Slice(c.topWords, func(i, j int) bool {
		if c.topWords[i].freq == c.topWords[j].freq {
			return c.topWords[i].word < c.topWords[j].word
		}
		return c.topWords[i].freq > c.topWords[j].freq
	})
}

// makeTopTenResponseDTO creates the top ten response DTO.
func (c *WordContainer) makeTopTenResponseDTO() []dto.TopWordsResponseDto {
	n := len(c.topWords)
	if n > 10 {
		n = 10
	}
	resp := make([]dto.TopWordsResponseDto, n)
	for i := 0; i < n; i++ {
		resp[i] = dto.TopWordsResponseDto{
			Word:      c.topWords[i].word,
			Frequency: c.topWords[i].freq,
		}
	}
	return resp
}

// getMaxNumCPUs sets the maximum number of CPUs that can be executing concurrently.
// returns the number, default is 1.
func getMaxNumCPUs() int {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
	return n
}

// getChunkSize calculates the number of chunks to process by each worker.
// total word chunks = total number of words / number of workers
// if there is a remainder, add one more chunk to the last worker
// Why? - to ensure that the last worker processes the remaining words.
func getChunkSize(words []string, workers int) int {
	chunks := len(words) / workers
	if len(words)%workers != 0 {
		chunks++
	}
	return chunks
}

// getChunkSize calculates the range of words to be processed by each go routine(worker)
func getChunkRange(size int, i int, len int) (int, int) {
	start := i * size
	end := start + size
	if end > len {
		end = len
	}
	return start, end
}
