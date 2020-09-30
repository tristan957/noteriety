package note

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	// NoterietyDBName is the name of the Noteriety DB file
	NoterietyDBName = "noteriety.db"
)

// DB is a handle to the SQLite DB
var DB *sql.DB

func init() {
	var dbLocation string
	if repository := viper.GetString("repository.path"); len(repository) == 0 {
		dbLocation = NoterietyDBName
	} else {
		dbLocation = path.Join(repository, NoterietyDBName)
	}
	_, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
