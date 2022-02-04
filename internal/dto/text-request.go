package dto

import (
	"github.com/ashtishad/top-words/internal/lib"
	"regexp"
)

type TextRequestDto struct {
	Text string `json:"text"`
}

func (dto *TextRequestDto) ValidateRequest() lib.RestErr {
	if dto.Text == "" {
		return lib.NewBadRequestError("text cannot be empty")
	}

	if len(dto.Text) > lib.MaxTextLength {
		return lib.NewBadRequestError("text is too long, more than 10 million characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(dto.Text) {
		return lib.NewBadRequestError("text contains invalid characters")
	}

	return nil
}
