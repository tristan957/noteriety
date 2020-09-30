package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"git.sr.ht/~tristan957/noteriety/noteriety/note"
	"git.sr.ht/~tristan957/noteriety/noteriety/note/bindata"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCommand)
}

var initCommand = &cobra.Command{
	Use: "init",
	Run: _init,
}

func _init(cmd *cobra.Command, args []string) {
	executable, err := exec.LookPath("sqlite3")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to find sqlite3 in your $PATH")
		os.Exit(1)
	}

	var dbLocation string
	if len(noterietyRepo) == 0 {
		dbLocation = note.NoterietyDBName
	} else {
		dbLocation = path.Join(noterietyRepo, note.NoterietyDBName)
	}

	createDBCommand := exec.Command(executable, dbLocation, ".database")
	if err = createDBCommand.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create noteriety.db")
		os.Exit(2)
	}

	sql, err := bindata.Asset("sql/migrations/0-noteriety.sql")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	createSchemaCommand := exec.Command(executable, dbLocation, string(sql))
	if err = createSchemaCommand.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create database schema")
		os.Exit(2)
	}
}
