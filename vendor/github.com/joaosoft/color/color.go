package color

import (
	"fmt"
)

func WithColor(text string, format Format, foreground Foreground, background Background, params ...interface{}) string {
	return fmt.Sprintf(fmt.Sprintf("\r%s[%d;%d;%dm%s%s[%dm", Escape, format, foreground, background, text, Escape, FormatReset), params...)
}
