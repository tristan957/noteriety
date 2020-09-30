package cmd

import (
	"fmt"
	"os"

	"git.sr.ht/~tristan957/noteriety/noteriety/note"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newCommand)
	addCommand.Flags().BoolP("encrypted", "e", false, "Whether the note is encrypted")
	viper.BindPFlag("notes.encrypt", addCommand.Flags().Lookup("encrypted"))
}

var newCommand = &cobra.Command{
	Use:  "new",
	Args: cobra.ExactValidArgs(1),
	Run:  add,
}

func new(cmd *cobra.Command, args []string) {
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
