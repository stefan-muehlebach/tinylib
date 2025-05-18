// Bietet einen einfachen Zugriff auf die Go-Fonts aber auch auf eine Reihe
// von OpenSource-Schriten.
package fonts

import (
    "embed"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed FontFiles/*.ttf
var fontFiles embed.FS

// Erstellt einen neuen Fontface, der bspw. bei der Methode [SetFontFace]
// verwendet werden kann. textFont ist ein Pointer auf einen OpenType-Font
// Siehe auch Array [Names] für eine Liste aller Fonts, die in diesem Package
// angeboten werden.
func NewFace(fontName string, size float64) font.Face {
    ttf, err := fontFiles.ReadFile(Map[fontName])
    if err != nil {
        println("ReadFile: ", err.Error())
    }
    textFont, err := opentype.Parse(ttf)
    if err != nil {
        println("opentype.Parse: ", err.Error())
    }
	face, _ := opentype.NewFace(textFont,
		&opentype.FaceOptions{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingFull,
		})
	return face
}
