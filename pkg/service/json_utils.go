package service

import (
	"encoding/json"
	"github.com/ashtishad/top-words/internal/dto"
	"github.com/ashtishad/top-words/internal/lib"
	"io"
)

// ToJSON takes a dto.Response and writes it to the writer in JSON format
func ToJSON(w io.Writer, resp []dto.TopWordsResponseDto) lib.RestErr {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return lib.NewInternalServerError("Error encoding the response", err)
	}
	return nil
}

// FromJSON converts json to TextRequestDto
func FromJSON(data io.ReadCloser) (dto.TextRequestDto, lib.RestErr) {
	var text dto.TextRequestDto
	if data == nil {
		return text, lib.NewBadRequestError("request body cannot be empty")
	}

	if err := json.NewDecoder(data).Decode(&text); err != nil {
		return text, lib.NewInternalServerError("error while decoding request body", err)
	}

	return text, nil
}
