package util

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/api"
)

func Get[Option any](cfg api.Config, key string) (option Option, found bool) {
	option_ := cfg.Get(key)
	if option_ == nil {
		return
	}
	option, found = option_.(Option)
	return option, found
}

func ExposeConfig(cfg api.Config) string {
	all := cfg.All()
	if len(all) == 0 {
		return "Config{}"
	}
	b := &strings.Builder{}
	for key, value := range all {
		b.WriteString(key + ": " + fmt.Sprint(value) + ", ")
	}
	res := b.String()
	// remove trailing ", "
	return "Config{" + res[:len(res)-2] + "}"
}
