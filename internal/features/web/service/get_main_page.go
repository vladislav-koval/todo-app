package web_service

import (
	"fmt"
)

func (s *WebService) GetMainPage() ([]byte, error) {
	html, err := s.webRepository.GetFile("index.html")

	if err != nil {
		return nil, fmt.Errorf("get file from repo: %w", err)
	}

	return html, nil
}
