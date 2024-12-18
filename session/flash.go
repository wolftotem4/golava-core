package session

func GetFlashErrors(s *Store) map[string]string {
	errors, ok := s.Get("errors")
	if !ok {
		return nil
	}

	value, _ := errors.(map[string]string)
	return value
}

func GetFlashError(s *Store, key string) string {
	errors := GetFlashErrors(s)
	if errors == nil {
		return ""
	}

	return errors[key]
}
