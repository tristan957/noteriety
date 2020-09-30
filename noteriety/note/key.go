package note

import (
	"strings"

	"github.com/spf13/viper"
)

type Key string

func (k Key) Sanitize() Key {
	temp := strings.Replace(string(k), "/", "", 1)
	return Key(strings.ReplaceAll(temp, "/", "__"))
}

func (k Key) Normalize() Key {
	if string(k)[0] != '/' {
		return Key("/" + string(k))
	}

	return k
}

func (k Key) ToFilePath() string {
	repo := viper.GetString("repository.path")
	return repo + string(k) + ".md"
}

func (k Key) String() string {
	return string(k)
}
