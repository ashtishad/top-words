package dto

type TopWordsResponseDto struct {
	Word      string `json:"word"`
	Frequency int64  `json:"frequency"`
}

type ResponseDto struct {
	TopWords []TopWordsResponseDto `json:"top_words"`
}
