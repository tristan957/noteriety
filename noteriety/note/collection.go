package note

type Collection struct{}

func CreateCollection(collection string) error {
	if _, err := DB.Query("INSERT INTO collection (name) VALUES (?)", collection); err != nil {
		return err
	}

	return nil
}
