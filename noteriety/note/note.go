package note

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Note struct {
	Key          Key
	Encrypted    bool
	CreatedAtUTC time.Time
	UpdatedAtUTC time.Time
	Tags         []Tag
	Collections  []Collection
}

func NoteFromKey(key Key) (*Note, error) {
	var encrypted sql.NullBool
	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	if err := DB.QueryRow("SELECT encrypted, created_at, updated_at FROM note WHERE key == ?", key).Scan(&encrypted, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if !encrypted.Valid {
		return nil, fmt.Errorf("Unknown value for note.encrypted where key equals %s", key)
	}

	if !createdAt.Valid {
		return nil, fmt.Errorf("Unknown value for note.created_at where key equals %s", key)
	}

	if !updatedAt.Valid {
		return nil, fmt.Errorf("Unknown value for note.updated_at where key equals %s", key)
	}

	return &Note{key, encrypted.Bool, createdAt.Time, updatedAt.Time, nil, nil}, nil
}

func CreateNote(key Key, encrypted bool, tags []string, collections []string) error {
	path := key.ToFilePath()
	if info, err := os.Stat(path); info != nil || (err != nil && err != os.ErrNotExist) {
		if info != nil {
			return os.ErrExist
		}

		if err != os.ErrNotExist {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	// TODO: if encrypted write out the armor signature
	if err = file.Close(); err != nil {
		return err
	}

	return AddNote(key, encrypted, tags, collections)
}

func AddNote(key Key, encrypted bool, tags []string, collections []string) error {
	if _, err := DB.Query("INSERT INTO note (key, encrypted) VALUES (?, ?)", key, encrypted); err != nil {
		return err
	}

	return nil
}
