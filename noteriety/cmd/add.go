package cmd

import (
	"fmt"
	"os"

	"git.sr.ht/~tristan957/noteriety/noteriety/note"
	"github.com/spf13/cobra"
)

var encrypted bool

func init() {
	rootCmd.AddCommand(addCommand)
	addCommand.Flags().BoolVarP(&encrypted, "encrypted", "e", false, "Whether the note is encrypted")
}

var addCommand = &cobra.Command{
	Use:  "add",
	Args: cobra.ExactValidArgs(1),
	Run:  add,
}

func add(cmd *cobra.Command, args []string) {
	key := note.Key(args[0]).Sanitize()
	n, err := note.NoteFromKey(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if n != nil {
		fmt.Println("Note already exists")
		os.Exit(0)
	}

	if err = note.CreateNote(key, encrypted, nil, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
