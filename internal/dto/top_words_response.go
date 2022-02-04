package dto

import "github.com/ashtishad/top-words/pkg/service"

type TopWordsResponseDto struct {
	Words     []service.TopWord `json:"words"`
	Frequency []int64           `json:"frequency"`
}
