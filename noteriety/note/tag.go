package note

type Tag struct{}

func CreateTag(tag string) error {
	if _, err := DB.Query("INSERT INTO tag (name) VALUES (?)", tag); err != nil {
		return err
	}

	return nil
}
