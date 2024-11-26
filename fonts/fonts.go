package fonts

import _ "embed"

//go:embed ttf/OpenSans-Bold.ttf
var Bold []byte

//go:embed ttf/OpenSans-Regular.ttf
var Regular []byte

//go:embed ttf/OpenSans-Italic.ttf
var Italic []byte

//go:embed ttf/Phosphor.ttf
var Icons []byte

const (
	IconXOffset = 38
	IconYOffset = 2
)
