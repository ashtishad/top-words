package dto

type TopWordsResponseDto struct {
	Word      string `json:"word"`
	Frequency int64  `json:"frequency"`
}
