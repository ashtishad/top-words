package dto

import (
	"github.com/ashtishad/top-words/internal/lib"
	"regexp"
	"strings"
)

type TextRequestDto struct {
	Text string `json:"text"`
}

// ValidateRequest validates the TextRequestDto, returns list of words split by space
// and returns error if input is empty, or contains non-alphanumeric characters or more than 10 million characters
func (dto *TextRequestDto) ValidateRequest() ([]string, lib.RestErr) {
	if dto.Text == "" {
		return nil, lib.NewBadRequestError("text cannot be empty")
	}

	if len(dto.Text) > lib.MaxTextLength {
		return nil, lib.NewBadRequestError("text is too long, more than 10 million characters")
	}

	r := regexp.MustCompile(`[^a-zA-Z\-'’]`)

	dto.Text = strings.ToLower(dto.Text)
	dto.Text = r.ReplaceAllString(dto.Text, " ")
	words := strings.Split(dto.Text, " ")

	return words, nil
}
