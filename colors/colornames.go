package colors

import (
	"image/color"
	"math/rand/v2"
)

func RandColor() color.RGBA {
	return Map[Names[rand.N(len(Names))]]
}

var (
	// AliceBlue            = color.RGBA{0xF0, 0xF8, 0xFF, 0xFF}
	// AntiqueWhite         = color.RGBA{0xFA, 0xEB, 0xD7, 0xFF}
	// Aqua                 = color.RGBA{0x00, 0xFF, 0xFF, 0xFF}
	// Aquamarine           = color.RGBA{0x7F, 0xFF, 0xD4, 0xFF}
	// Azure                = color.RGBA{0xF0, 0xFF, 0xFF, 0xFF}
	// Beige                = color.RGBA{0xF5, 0xF5, 0xDC, 0xFF}
	// Bisque               = color.RGBA{0xFF, 0xE4, 0xC4, 0xFF}
	Black                = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	// BlanchedAlmond       = color.RGBA{0xFF, 0xEB, 0xCD, 0xFF}
	Blue                 = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
	// BlueViolet           = color.RGBA{0x8A, 0x2B, 0xE2, 0xFF}
	// Brown                = color.RGBA{0xA5, 0x2A, 0x2A, 0xFF}
	// BurlyWood            = color.RGBA{0xDE, 0xB8, 0x87, 0xFF}
	// CadetBlue            = color.RGBA{0x5F, 0x9E, 0xA0, 0xFF}
	// Chartreuse           = color.RGBA{0x7F, 0xFF, 0x00, 0xFF}
	// Chocolate            = color.RGBA{0xD2, 0x69, 0x1E, 0xFF}
	// Coral                = color.RGBA{0xFF, 0x7F, 0x50, 0xFF}
	// CornflowerBlue       = color.RGBA{0x64, 0x95, 0xED, 0xFF}
	// Cornsilk             = color.RGBA{0xFF, 0xF8, 0xDC, 0xFF}
	// Crimson              = color.RGBA{0xDC, 0x14, 0x3C, 0xFF}
	Cyan                 = color.RGBA{0x00, 0xFF, 0xFF, 0xFF}
	// DarkBlue             = color.RGBA{0x00, 0x00, 0x8B, 0xFF}
	// DarkCyan             = color.RGBA{0x00, 0x8B, 0x8B, 0xFF}
	// DarkGoldenrod        = color.RGBA{0xB8, 0x86, 0x0B, 0xFF}
	// DarkGray             = color.RGBA{0xA9, 0xA9, 0xA9, 0xFF}
	// DarkGreen            = color.RGBA{0x00, 0x64, 0x00, 0xFF}
	// DarkGrey             = color.RGBA{0xA9, 0xA9, 0xA9, 0xFF}
	// DarkKhaki            = color.RGBA{0xBD, 0xB7, 0x6B, 0xFF}
	// DarkMagenta          = color.RGBA{0x8B, 0x00, 0x8B, 0xFF}
	// DarkOliveGreen       = color.RGBA{0x55, 0x6B, 0x2F, 0xFF}
	// DarkOrange           = color.RGBA{0xFF, 0x8C, 0x00, 0xFF}
	// DarkOrchid           = color.RGBA{0x99, 0x32, 0xCC, 0xFF}
	// DarkRed              = color.RGBA{0x8B, 0x00, 0x00, 0xFF}
	// DarkSalmon           = color.RGBA{0xE9, 0x96, 0x7A, 0xFF}
	// DarkSeaGreen         = color.RGBA{0x8F, 0xBC, 0x8F, 0xFF}
	// DarkSlateBlue        = color.RGBA{0x48, 0x3D, 0x8B, 0xFF}
	// DarkSlateGray        = color.RGBA{0x2F, 0x4F, 0x4F, 0xFF}
	// DarkSlateGrey        = color.RGBA{0x2F, 0x4F, 0x4F, 0xFF}
	// DarkTurquoise        = color.RGBA{0x00, 0xCE, 0xD1, 0xFF}
	// DarkViolet           = color.RGBA{0x94, 0x00, 0xD3, 0xFF}
	// DeepPink             = color.RGBA{0xFF, 0x14, 0x93, 0xFF}
	// DeepSkyBlue          = color.RGBA{0x00, 0xBF, 0xFF, 0xFF}
	// DimGray              = color.RGBA{0x69, 0x69, 0x69, 0xFF}
	// DimGrey              = color.RGBA{0x69, 0x69, 0x69, 0xFF}
	// DodgerBlue           = color.RGBA{0x1E, 0x90, 0xFF, 0xFF}
	// FireBrick            = color.RGBA{0xB2, 0x22, 0x22, 0xFF}
	// FloralWhite          = color.RGBA{0xFF, 0xFA, 0xF0, 0xFF}
	// ForestGreen          = color.RGBA{0x22, 0x8B, 0x22, 0xFF}
	// Fuchsia              = color.RGBA{0xFF, 0x00, 0xFF, 0xFF}
	// Gainsboro            = color.RGBA{0xDC, 0xDC, 0xDC, 0xFF}
	// GhostWhite           = color.RGBA{0xF8, 0xF8, 0xFF, 0xFF}
	// Gold                 = color.RGBA{0xFF, 0xD7, 0x00, 0xFF}
	// Goldenrod            = color.RGBA{0xDA, 0xA5, 0x20, 0xFF}
	// Gray                 = color.RGBA{0x80, 0x80, 0x80, 0xFF}
	Green                = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	// GreenYellow          = color.RGBA{0xAD, 0xFF, 0x2F, 0xFF}
	// Grey                 = color.RGBA{0x80, 0x80, 0x80, 0xFF}
	// Honeydew             = color.RGBA{0xF0, 0xFF, 0xF0, 0xFF}
	// HotPink              = color.RGBA{0xFF, 0x69, 0xB4, 0xFF}
	// IndianRed            = color.RGBA{0xCD, 0x5C, 0x5C, 0xFF}
	// Indigo               = color.RGBA{0x4B, 0x00, 0x82, 0xFF}
	// Ivory                = color.RGBA{0xFF, 0xFF, 0xF0, 0xFF}
	// Khaki                = color.RGBA{0xF0, 0xE6, 0x8C, 0xFF}
	// Lavender             = color.RGBA{0xE6, 0xE6, 0xFA, 0xFF}
	// LavenderBlush        = color.RGBA{0xFF, 0xF0, 0xF5, 0xFF}
	// LawnGreen            = color.RGBA{0x7C, 0xFC, 0x00, 0xFF}
	// LemonChiffon         = color.RGBA{0xFF, 0xFA, 0xCD, 0xFF}
	// LightBlue            = color.RGBA{0xAD, 0xD8, 0xE6, 0xFF}
	// LightCoral           = color.RGBA{0xF0, 0x80, 0x80, 0xFF}
	// LightCyan            = color.RGBA{0xE0, 0xFF, 0xFF, 0xFF}
	// LightGoldenrodYellow = color.RGBA{0xFA, 0xFA, 0xD2, 0xFF}
	// LightGray            = color.RGBA{0xD3, 0xD3, 0xD3, 0xFF}
	// LightGreen           = color.RGBA{0x90, 0xEE, 0x90, 0xFF}
	// LightGrey            = color.RGBA{0xD3, 0xD3, 0xD3, 0xFF}
	// LightPink            = color.RGBA{0xFF, 0xB6, 0xC1, 0xFF}
	// LightSalmon          = color.RGBA{0xFF, 0xA0, 0x7A, 0xFF}
	// LightSeaGreen        = color.RGBA{0x20, 0xB2, 0xAA, 0xFF}
	// LightSkyBlue         = color.RGBA{0x87, 0xCE, 0xFA, 0xFF}
	// LightSlateGray       = color.RGBA{0x77, 0x88, 0x99, 0xFF}
	// LightSlateGrey       = color.RGBA{0x77, 0x88, 0x99, 0xFF}
	// LightSteelBlue       = color.RGBA{0xB0, 0xC4, 0xDE, 0xFF}
	// LightYellow          = color.RGBA{0xFF, 0xFF, 0xE0, 0xFF}
	// Lime                 = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	// LimeGreen            = color.RGBA{0x32, 0xCD, 0x32, 0xFF}
	// Linen                = color.RGBA{0xFA, 0xF0, 0xE6, 0xFF}
	Magenta              = color.RGBA{0xFF, 0x00, 0xFF, 0xFF}
	// Maroon               = color.RGBA{0x80, 0x00, 0x00, 0xFF}
	// MediumAquamarine     = color.RGBA{0x66, 0xCD, 0xAA, 0xFF}
	// MediumBlue           = color.RGBA{0x00, 0x00, 0xCD, 0xFF}
	// MediumOrchid         = color.RGBA{0xBA, 0x55, 0xD3, 0xFF}
	// MediumPurple         = color.RGBA{0x93, 0x70, 0xDB, 0xFF}
	// MediumSeaGreen       = color.RGBA{0x3C, 0xB3, 0x71, 0xFF}
	// MediumSlateBlue      = color.RGBA{0x7B, 0x68, 0xEE, 0xFF}
	// MediumSpringGreen    = color.RGBA{0x00, 0xFA, 0x9A, 0xFF}
	// MediumTurquoise      = color.RGBA{0x48, 0xD1, 0xCC, 0xFF}
	// MediumVioletRed      = color.RGBA{0xC7, 0x15, 0x85, 0xFF}
	// MidnightBlue         = color.RGBA{0x19, 0x19, 0x70, 0xFF}
	// MintCream            = color.RGBA{0xF5, 0xFF, 0xFA, 0xFF}
	// MistyRose            = color.RGBA{0xFF, 0xE4, 0xE1, 0xFF}
	// Moccasin             = color.RGBA{0xFF, 0xE4, 0xB5, 0xFF}
	// NavajoWhite          = color.RGBA{0xFF, 0xDE, 0xAD, 0xFF}
	// Navy                 = color.RGBA{0x00, 0x00, 0x80, 0xFF}
	// OldLace              = color.RGBA{0xFD, 0xF5, 0xE6, 0xFF}
	// Olive                = color.RGBA{0x80, 0x80, 0x00, 0xFF}
	// OliveDrab            = color.RGBA{0x6B, 0x8E, 0x23, 0xFF}
	// Orange               = color.RGBA{0xFF, 0xA5, 0x00, 0xFF}
	// OrangeRed            = color.RGBA{0xFF, 0x45, 0x00, 0xFF}
	// Orchid               = color.RGBA{0xDA, 0x70, 0xD6, 0xFF}
	// PaleGoldenrod        = color.RGBA{0xEE, 0xE8, 0xAA, 0xFF}
	// PaleGreen            = color.RGBA{0x98, 0xFB, 0x98, 0xFF}
	// PaleTurquoise        = color.RGBA{0xAF, 0xEE, 0xEE, 0xFF}
	// PaleVioletRed        = color.RGBA{0xDB, 0x70, 0x93, 0xFF}
	// PapayaWhip           = color.RGBA{0xFF, 0xEF, 0xD5, 0xFF}
	// PeachPuff            = color.RGBA{0xFF, 0xDA, 0xB9, 0xFF}
	// Peru                 = color.RGBA{0xCD, 0x85, 0x3F, 0xFF}
	// Pink                 = color.RGBA{0xFF, 0xC0, 0xCB, 0xFF}
	// Plum                 = color.RGBA{0xDD, 0xA0, 0xDD, 0xFF}
	// PowderBlue           = color.RGBA{0xB0, 0xE0, 0xE6, 0xFF}
	// Purple               = color.RGBA{0x80, 0x00, 0x80, 0xFF}
	Red                  = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	// RosyBrown            = color.RGBA{0xBC, 0x8F, 0x8F, 0xFF}
	// RoyalBlue            = color.RGBA{0x41, 0x69, 0xE1, 0xFF}
	// SaddleBrown          = color.RGBA{0x8B, 0x45, 0x13, 0xFF}
	// Salmon               = color.RGBA{0xFA, 0x80, 0x72, 0xFF}
	// SandyBrown           = color.RGBA{0xF4, 0xA4, 0x60, 0xFF}
	// SeaGreen             = color.RGBA{0x2E, 0x8B, 0x57, 0xFF}
	// Seashell             = color.RGBA{0xFF, 0xF5, 0xEE, 0xFF}
	// Sienna               = color.RGBA{0xA0, 0x52, 0x2D, 0xFF}
	// Silver               = color.RGBA{0xC0, 0xC0, 0xC0, 0xFF}
	// SkyBlue              = color.RGBA{0x87, 0xCE, 0xEB, 0xFF}
	// SlateBlue            = color.RGBA{0x6A, 0x5A, 0xCD, 0xFF}
	// SlateGray            = color.RGBA{0x70, 0x80, 0x90, 0xFF}
	// SlateGrey            = color.RGBA{0x70, 0x80, 0x90, 0xFF}
	// Snow                 = color.RGBA{0xFF, 0xFA, 0xFA, 0xFF}
	// SpringGreen          = color.RGBA{0x00, 0xFF, 0x7F, 0xFF}
	// SteelBlue            = color.RGBA{0x46, 0x82, 0xB4, 0xFF}
	// Tan                  = color.RGBA{0xD2, 0xB4, 0x8C, 0xFF}
	// Teal                 = color.RGBA{0x00, 0x80, 0x80, 0xFF}
	// Thistle              = color.RGBA{0xD8, 0xBF, 0xD8, 0xFF}
	// Tomato               = color.RGBA{0xFF, 0x63, 0x47, 0xFF}
	// Turquoise            = color.RGBA{0x40, 0xE0, 0xD0, 0xFF}
	// Violet               = color.RGBA{0xEE, 0x82, 0xEE, 0xFF}
	// Wheat                = color.RGBA{0xF5, 0xDE, 0xB3, 0xFF}
	White                = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	// WhiteSmoke           = color.RGBA{0xF5, 0xF5, 0xF5, 0xFF}
	Yellow               = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
	// YellowGreen          = color.RGBA{0x9A, 0xCD, 0x32, 0xFF}

	GoGopherBlue         = color.RGBA{0x00, 0xAD, 0xD8, 0xFF}
	GoLightBlue          = color.RGBA{0x5D, 0xC9, 0xE2, 0xFF}
	GoAqua               = color.RGBA{0x00, 0xA2, 0x9C, 0xFF}
	GoFuchsia            = color.RGBA{0xCE, 0x32, 0x62, 0xFF}
	GoYellow             = color.RGBA{0xFD, 0xDD, 0x00, 0xFF}
	GoTeal               = color.RGBA{0x00, 0x75, 0x8D, 0xFF}
	GoDimGray            = color.RGBA{0x55, 0x57, 0x59, 0xFF}
	GoIndigo             = color.RGBA{0x40, 0x2B, 0x56, 0xFF}
	GoLightGray          = color.RGBA{0xDB, 0xD9, 0xD6, 0xFF}

	Map = map[string]color.RGBA{
		// "AliceBlue":            AliceBlue,
		// "AntiqueWhite":         AntiqueWhite,
		// "Aqua":                 Aqua,
		// "Aquamarine":           Aquamarine,
		// "Azure":                Azure,
		// "Beige":                Beige,
		// "Bisque":               Bisque,
		"Black":                Black,
		// "BlanchedAlmond":       BlanchedAlmond,
		"Blue":                 Blue,
		// "BlueViolet":           BlueViolet,
		// "Brown":                Brown,
		// "BurlyWood":            BurlyWood,
		// "CadetBlue":            CadetBlue,
		// "Chartreuse":           Chartreuse,
		// "Chocolate":            Chocolate,
		// "Coral":                Coral,
		// "CornflowerBlue":       CornflowerBlue,
		// "Cornsilk":             Cornsilk,
		// "Crimson":              Crimson,
		"Cyan":                 Cyan,
		// "DarkBlue":             DarkBlue,
		// "DarkCyan":             DarkCyan,
		// "DarkGoldenrod":        DarkGoldenrod,
		// "DarkGray":             DarkGray,
		// "DarkGreen":            DarkGreen,
		// "DarkGrey":             DarkGrey,
		// "DarkKhaki":            DarkKhaki,
		// "DarkMagenta":          DarkMagenta,
		// "DarkOliveGreen":       DarkOliveGreen,
		// "DarkOrange":           DarkOrange,
		// "DarkOrchid":           DarkOrchid,
		// "DarkRed":              DarkRed,
		// "DarkSalmon":           DarkSalmon,
		// "DarkSeaGreen":         DarkSeaGreen,
		// "DarkSlateBlue":        DarkSlateBlue,
		// "DarkSlateGray":        DarkSlateGray,
		// "DarkSlateGrey":        DarkSlateGrey,
		// "DarkTurquoise":        DarkTurquoise,
		// "DarkViolet":           DarkViolet,
		// "DeepPink":             DeepPink,
		// "DeepSkyBlue":          DeepSkyBlue,
		// "DimGray":              DimGray,
		// "DimGrey":              DimGrey,
		// "DodgerBlue":           DodgerBlue,
		// "FireBrick":            FireBrick,
		// "FloralWhite":          FloralWhite,
		// "ForestGreen":          ForestGreen,
		// "Fuchsia":              Fuchsia,
		// "Gainsboro":            Gainsboro,
		// "GhostWhite":           GhostWhite,
		// "Gold":                 Gold,
		// "Goldenrod":            Goldenrod,
		// "Gray":                 Gray,
		"Green":                Green,
		// "GreenYellow":          GreenYellow,
		// "Grey":                 Grey,
		// "Honeydew":             Honeydew,
		// "HotPink":              HotPink,
		// "IndianRed":            IndianRed,
		// "Indigo":               Indigo,
		// "Ivory":                Ivory,
		// "Khaki":                Khaki,
		// "Lavender":             Lavender,
		// "LavenderBlush":        LavenderBlush,
		// "LawnGreen":            LawnGreen,
		// "LemonChiffon":         LemonChiffon,
		// "LightBlue":            LightBlue,
		// "LightCoral":           LightCoral,
		// "LightCyan":            LightCyan,
		// "LightGoldenrodYellow": LightGoldenrodYellow,
		// "LightGray":            LightGray,
		// "LightGreen":           LightGreen,
		// "LightGrey":            LightGrey,
		// "LightPink":            LightPink,
		// "LightSalmon":          LightSalmon,
		// "LightSeaGreen":        LightSeaGreen,
		// "LightSkyBlue":         LightSkyBlue,
		// "LightSlateGray":       LightSlateGray,
		// "LightSlateGrey":       LightSlateGrey,
		// "LightSteelBlue":       LightSteelBlue,
		// "LightYellow":          LightYellow,
		// "Lime":                 Lime,
		// "LimeGreen":            LimeGreen,
		// "Linen":                Linen,
		"Magenta":              Magenta,
		// "Maroon":               Maroon,
		// "MediumAquamarine":     MediumAquamarine,
		// "MediumBlue":           MediumBlue,
		// "MediumOrchid":         MediumOrchid,
		// "MediumPurple":         MediumPurple,
		// "MediumSeaGreen":       MediumSeaGreen,
		// "MediumSlateBlue":      MediumSlateBlue,
		// "MediumSpringGreen":    MediumSpringGreen,
		// "MediumTurquoise":      MediumTurquoise,
		// "MediumVioletRed":      MediumVioletRed,
		// "MidnightBlue":         MidnightBlue,
		// "MintCream":            MintCream,
		// "MistyRose":            MistyRose,
		// "Moccasin":             Moccasin,
		// "NavajoWhite":          NavajoWhite,
		// "Navy":                 Navy,
		// "OldLace":              OldLace,
		// "Olive":                Olive,
		// "OliveDrab":            OliveDrab,
		// "Orange":               Orange,
		// "OrangeRed":            OrangeRed,
		// "Orchid":               Orchid,
		// "PaleGoldenrod":        PaleGoldenrod,
		// "PaleGreen":            PaleGreen,
		// "PaleTurquoise":        PaleTurquoise,
		// "PaleVioletRed":        PaleVioletRed,
		// "PapayaWhip":           PapayaWhip,
		// "PeachPuff":            PeachPuff,
		// "Peru":                 Peru,
		// "Pink":                 Pink,
		// "Plum":                 Plum,
		// "PowderBlue":           PowderBlue,
		// "Purple":               Purple,
		"Red":                  Red,
		// "RosyBrown":            RosyBrown,
		// "RoyalBlue":            RoyalBlue,
		// "SaddleBrown":          SaddleBrown,
		// "Salmon":               Salmon,
		// "SandyBrown":           SandyBrown,
		// "SeaGreen":             SeaGreen,
		// "Seashell":             Seashell,
		// "Sienna":               Sienna,
		// "Silver":               Silver,
		// "SkyBlue":              SkyBlue,
		// "SlateBlue":            SlateBlue,
		// "SlateGray":            SlateGray,
		// "SlateGrey":            SlateGrey,
		// "Snow":                 Snow,
		// "SpringGreen":          SpringGreen,
		// "SteelBlue":            SteelBlue,
		// "Tan":                  Tan,
		// "Teal":                 Teal,
		// "Thistle":              Thistle,
		// "Tomato":               Tomato,
		// "Turquoise":            Turquoise,
		// "Violet":               Violet,
		// "Wheat":                Wheat,
		"White":                White,
		// "WhiteSmoke":           WhiteSmoke,
		"Yellow":               Yellow,
		// "YellowGreen":          YellowGreen,
	}

	Names = []string{
		// "AliceBlue",
		// "AntiqueWhite",
		// "Aqua",
		// "Aquamarine",
		// "Azure",
		// "Beige",
		// "Bisque",
		"Black",
		// "BlanchedAlmond",
		"Blue",
		// "BlueViolet",
		// "Brown",
		// "BurlyWood",
		// "CadetBlue",
		// "Chartreuse",
		// "Chocolate",
		// "Coral",
		// "CornflowerBlue",
		// "Cornsilk",
		// "Crimson",
		"Cyan",
		// "DarkBlue",
		// "DarkCyan",
		// "DarkGoldenrod",
		// "DarkGray",
		// "DarkGreen",
		// "DarkGrey",
		// "DarkKhaki",
		// "DarkMagenta",
		// "DarkOliveGreen",
		// "DarkOrange",
		// "DarkOrchid",
		// "DarkRed",
		// "DarkSalmon",
		// "DarkSeaGreen",
		// "DarkSlateBlue",
		// "DarkSlateGray",
		// "DarkSlateGrey",
		// "DarkTurquoise",
		// "DarkViolet",
		// "DeepPink",
		// "DeepSkyBlue",
		// "DimGray",
		// "DimGrey",
		// "DodgerBlue",
		// "FireBrick",
		// "FloralWhite",
		// "ForestGreen",
		// "Fuchsia",
		// "Gainsboro",
		// "GhostWhite",
		// "Gold",
		// "Goldenrod",
		// "Gray",
		"Green",
		// "GreenYellow",
		// "Grey",
		// "Honeydew",
		// "HotPink",
		// "IndianRed",
		// "Indigo",
		// "Ivory",
		// "Khaki",
		// "Lavender",
		// "LavenderBlush",
		// "LawnGreen",
		// "LemonChiffon",
		// "LightBlue",
		// "LightCoral",
		// "LightCyan",
		// "LightGoldenrodYellow",
		// "LightGray",
		// "LightGreen",
		// "LightGrey",
		// "LightPink",
		// "LightSalmon",
		// "LightSeaGreen",
		// "LightSkyBlue",
		// "LightSlateGray",
		// "LightSlateGrey",
		// "LightSteelBlue",
		// "LightYellow",
		// "Lime",
		// "LimeGreen",
		// "Linen",
		"Magenta",
		// "Maroon",
		// "MediumAquamarine",
		// "MediumBlue",
		// "MediumOrchid",
		// "MediumPurple",
		// "MediumSeaGreen",
		// "MediumSlateBlue",
		// "MediumSpringGreen",
		// "MediumTurquoise",
		// "MediumVioletRed",
		// "MidnightBlue",
		// "MintCream",
		// "MistyRose",
		// "Moccasin",
		// "NavajoWhite",
		// "Navy",
		// "OldLace",
		// "Olive",
		// "OliveDrab",
		// "Orange",
		// "OrangeRed",
		// "Orchid",
		// "PaleGoldenrod",
		// "PaleGreen",
		// "PaleTurquoise",
		// "PaleVioletRed",
		// "PapayaWhip",
		// "PeachPuff",
		// "Peru",
		// "Pink",
		// "Plum",
		// "PowderBlue",
		// "Purple",
		"Red",
		// "RosyBrown",
		// "RoyalBlue",
		// "SaddleBrown",
		// "Salmon",
		// "SandyBrown",
		// "SeaGreen",
		// "Seashell",
		// "Sienna",
		// "Silver",
		// "SkyBlue",
		// "SlateBlue",
		// "SlateGray",
		// "SlateGrey",
		// "Snow",
		// "SpringGreen",
		// "SteelBlue",
		// "Tan",
		// "Teal",
		// "Thistle",
		// "Tomato",
		// "Turquoise",
		// "Violet",
		// "Wheat",
		"White",
		// "WhiteSmoke",
		"Yellow",
		// "YellowGreen",
	}
)
