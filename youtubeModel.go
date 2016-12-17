package youtuber

import "fmt"

type Youtube struct {
	String string `json:"string"`
}

func ReturnYoutuberErrorMsg(errYoutuber map[string]interface{}) error {
	s, ok := errYoutuber["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error from Youtuber server, type coercion error: %t", ok)
	}
	m, ok := s["message"]
	if !ok {
		return fmt.Errorf("error from Youtuber server, map retrieval error: %t", ok)
	}
	message, ok := m.(string)
	if !ok {
		return fmt.Errorf("error from Youtuber server, unable to convert to string: %tgit st", ok)
	}
	return fmt.Errorf("error from Youtuber server: %s", message)
}
