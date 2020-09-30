package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"git.sr.ht/~tristan957/noteriety/noteriety/note"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/ssh/terminal"
)

var encrypt bool
var truncate bool
var passphrase string

func init() {
	rootCmd.AddCommand(editCommand)

	editCommand.Flags().BoolVarP(&encrypt, "encrypt", "e", false, "Encrypt the file after editing")
	viper.BindPFlag("notes.encrypt", editCommand.Flags().Lookup("encrypt"))
	editCommand.Flags().BoolVarP(&truncate, "truncate", "t", false, "Truncate the current data in the file")
	editCommand.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase used for encrypting and decrypting the note")
}

var editCommand = &cobra.Command{
	Use:  "edit",
	Args: cobra.ExactArgs(1),
	Run:  edit,
}

func edit(cmd *cobra.Command, args []string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		fmt.Fprintln(os.Stderr, "$EDITOR environment variable not set")
		os.Exit(1)
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find %s in you $PATH\n", editor)
		os.Exit(1)
	}

	key := note.Key(args[0]).Sanitize()
	note, err := note.NoteFromKey(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if note == nil {
		fmt.Printf("Note with key of %s does not exist\n", key)
	}

	flags := os.O_RDWR | os.O_APPEND
	if truncate {
		flags |= os.O_TRUNC
	}
	noteFile, err := os.OpenFile(note.Key.ToFilePath(), flags, os.ModeAppend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer noteFile.Close()

	var editorCommand *exec.Cmd
	var temp *os.File
	var rawPassphrase []byte
	if viper.GetBool("encrypt") || note.Encrypted {
		temp, err = ioutil.TempFile(os.TempDir(), note.Key.String()+"-*.md")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		if passphrase == "" {
			fmt.Print("Enter passphrase for encryption: ")
			rawPassphrase, err = terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
		} else {
			rawPassphrase = []byte(passphrase)
		}

		if note.Encrypted {
			guard, err := armor.Decode(bufio.NewReader(noteFile))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
			md, err := openpgp.ReadMessage(guard.Body, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
				if symmetric {
					return rawPassphrase, err
				}

				return nil, nil
			}, nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
			if _, err = io.Copy(bufio.NewWriter(temp), md.UnverifiedBody); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		}

		editorCommand = exec.Command(executable, temp.Name())
	} else {
		editorCommand = exec.Command(executable, note.Key.ToFilePath())
	}
	editorCommand.Stdin = os.Stdin
	editorCommand.Stdout = os.Stdout
	editorCommand.Stderr = os.Stderr
	if err = editorCommand.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// https://gist.github.com/cryptix/a56c66c4ef8f91e96875

	// When using a passphrase
	if viper.GetBool("encrypt") {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		// Setting up encryption flow
		noteBuf := bufio.NewWriter(noteFile)
		guard, err := armor.Encode(noteBuf, openpgp.SignatureType, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		plaintext, err := openpgp.SymmetricallyEncrypt(guard, rawPassphrase, nil, nil)
		if err != nil {
			guard.Close()
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer guard.Close()
		defer plaintext.Close()
		if _, err = io.Copy(plaintext, bufio.NewReader(temp)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		// Cleanup
		if err = plaintext.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if err = guard.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if err = noteBuf.Flush(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}

	// Cleanup
	if temp != nil {
		if err = temp.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if err = os.Remove(temp.Name()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}

	if err = noteFile.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
