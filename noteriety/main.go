package main

import (
	"fmt"
	"os"

	"git.sr.ht/~tristan957/noteriety/noteriety/cmd"
	"github.com/spf13/viper"
)

//go:generate go-bindata -nometadata -nocompress -nomemcopy -o note/bindata/bindata.go -pkg bindata sql/migrations/

func main() {
	viper.SetConfigName("noteriety")
	viper.SetConfigType("yaml")

	viper.SetDefault("encrypt", false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	cmd.Execute()
}
