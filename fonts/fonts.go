package fonts

import (
	_ "embed"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

//go:embed ttf/OpenSans-Regular.ttf
var RegularTTF []byte

//go:embed ttf/OpenSans-SemiBold.ttf
var SemiBoldTTF []byte

//go:embed ttf/OpenSans-Bold.ttf
var BoldTTF []byte

//go:embed ttf/OpenSans-Italic.ttf
var ItalicTTF []byte

//go:embed ttf/Phosphor.ttf
var IconsTTF []byte

const (
	IconXOffset = 38
	IconYOffset = 2
)

const (
	Regular  = "regular"
	SemiBold = "semibold"
	Bold     = "bold"
	Italic   = "italic"
	Icons    = "icons"
)

func NewCache() (Cache, error) {
	cache := Cache{}

	fonts := []struct {
		ttf  []byte
		name string
	}{
		{RegularTTF, Regular},
		{SemiBoldTTF, SemiBold},
		{BoldTTF, Bold},
		{ItalicTTF, Italic},
		{IconsTTF, Icons},
	}

	for _, font := range fonts {
		err := loadFont(cache, font.ttf, font.name)
		if err != nil {
			return cache, fmt.Errorf("loading font %q: %w", font.name, err)
		}
	}

	return cache, nil
}

func loadFont(fontCache Cache, ttf []byte, name string) error {
	font, err := truetype.Parse(ttf)
	if err != nil {
		return fmt.Errorf("parsing font %v: %w", name, err)
	}

	fontCache.Store(draw2d.FontData{Name: name}, font)

	return nil
}
