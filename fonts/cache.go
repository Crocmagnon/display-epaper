package fonts

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

type Cache map[string]*truetype.Font

func (fc Cache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

func (fc Cache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("font %s is not stored in font cache", fd.Name)
	}
	return font, nil
}
