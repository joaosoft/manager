package color

const Escape = "\033"

type Format int
type Foreground int
type Background int

const (
	NormalIntensityForeground = 30
	HighIntensityForeground   = 90
	NormalIntensityBackground = 40
	HighIntensityBackground   = 100
)

const (
	// format
	FormatReset Format = iota
	FormatBold
	FormatFaint
	FormatItalic
	FormatUnderline
	FormatBlinkSlow
	FormatBlinkRapid
	FormatReverseVideo
	FormatConcealed
	FormatCrossedOut
)

const (
	// foreground text colors
	ForegroundBlack Foreground = iota + NormalIntensityForeground
	ForegroundRed
	ForegroundGreen
	ForegroundYellow
	ForegroundBlue
	ForegroundMagenta
	ForegroundCyan
	ForegroundWhite
)

const (
	// foreground hi-intensity text colors
	ForegroundHiBlack Foreground = iota + HighIntensityForeground
	ForegroundHiRed
	ForegroundHiGreen
	ForegroundHiYellow
	ForegroundHiBlue
	ForegroundHiMagenta
	ForegroundHiCyan
	ForegroundHiWhite
)

const (
	// background text colors
	BackgroundBlack Background = iota + NormalIntensityBackground
	BackgroundRed
	BackgroundGreen
	BackgroundYellow
	BackgroundBlue
	BackgroundMagenta
	BackgroundCyan
	BackgroundWhite
)

const (
	// background hi-intensity text colors
	BackgroundHiBlack Background = iota + HighIntensityBackground
	BackgroundHiRed
	BackgroundHiGreen
	BackgroundHiYellow
	BackgroundHiBlue
	BackgroundHiMagenta
	BackgroundHiCyan
	BackgroundHiWhite
)
