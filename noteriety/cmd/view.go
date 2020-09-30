package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"git.sr.ht/~tristan957/noteriety/noteriety/note"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/cobra"
)

var inBrowser bool
var asHTML bool

func init() {
	rootCmd.AddCommand(viewCommand)
	viewCommand.Flags().BoolVarP(&inBrowser, "browser", "b", false, "Display note as HTML in a browser")
	viewCommand.Flags().BoolVarP(&asHTML, "html", "m", false, "Force the output to HTML if not displaying the note in the browser")
}

var viewCommand = &cobra.Command{
	Use:   "display",
	Short: "Display content of a note",
	Long:  "Display content of a note. Note can be displayed in $BROWSER or stdout.",
	Args:  cobra.ExactArgs(1),
	Run:   view,
}

func createHTMLFile(note *note.Note) (string, error) {
	data, err := ioutil.ReadFile(note.Key.ToFilePath())
	if err != nil {
		return "", err
	}

	p := parser.NewWithExtensions(parser.Tables | parser.FencedCode | parser.Autolink |
		parser.Strikethrough | parser.Strikethrough | parser.Footnotes |
		parser.HeadingIDs | parser.AutoHeadingIDs | parser.MathJax |
		parser.OrderedListStart | parser.Attributes | parser.SuperSubscript |
		parser.EmptyLinesBreakList | parser.Includes)
	htmlFlags := html.NofollowLinks | html.NoopenerLinks | html.NoreferrerLinks | html.HrefTargetBlank | html.TOC
	htmlOpts := html.RendererOptions{Flags: htmlFlags}
	r := html.NewRenderer(htmlOpts)
	htmlData := markdown.ToHTML(data, p, r)

	temp, err := ioutil.TempFile(os.TempDir(), string(note.Key.Sanitize())+"-*.html")
	if err != nil {
		return "", err
	}
	defer temp.Close()
	_, err = temp.Write(htmlData)
	if err != nil {
		return "", err
	}
	if err = temp.Close(); err != nil {
		return "", err
	}

	return temp.Name(), nil
}

func viewInBrowser(note *note.Note) {
	browser := os.Getenv("BROWSER")
	if browser == "" {
		fmt.Fprintln(os.Stderr, "No $BROWSER set")
		os.Exit(1)
	}
	executable, err := exec.LookPath(browser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find %s in you $PATH\n", browser)
		os.Exit(1)
	}

	htmlFileName, err := createHTMLFile(note)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	var browserCommand *exec.Cmd
	if note.Encrypted {
		browserCommand = exec.Command(executable)
	} else {
		browserCommand = exec.Command(executable, "file://"+htmlFileName)
	}
	browserCommand.Stdin = os.Stdin
	browserCommand.Stdout = os.Stdout
	browserCommand.Stderr = os.Stderr
	browserCommand.Run()
}

func viewInStdout(note *note.Note) {
	var f *os.File
	if asHTML {
		htmlFileName, err := createHTMLFile(note)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		f, err = os.Open(htmlFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	} else {
		// Go compiler pls
		var err error
		f, err = os.Open(note.Key.ToFilePath())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}

	if _, err := io.Copy(os.Stdout, f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func view(cmd *cobra.Command, args []string) {
	key := note.Key(args[0]).Normalize()
	note, err := note.NoteFromKey(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if note == nil {
		fmt.Printf("Note with key of %s does not exist\n", key)
		os.Exit(1)
	}

	if inBrowser {
		viewInBrowser(note)
		return
	}

	viewInStdout(note)
}
