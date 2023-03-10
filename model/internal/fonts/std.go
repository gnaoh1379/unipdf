/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package fonts

import (
	"github.com/gnaoh1379/unipdf/core"
	"github.com/gnaoh1379/unipdf/internal/textencoding"
)

// StdFontName is a name of a standard font.
type StdFontName string

// FontWeight specified font weight.
type FontWeight int

// Font weights
const (
	FontWeightMedium FontWeight = iota // Medium
	FontWeightBold                     // Bold
	FontWeightRoman                    // Roman
)

// Descriptor describes geometric properties of a font.
type Descriptor struct {
	Name        StdFontName
	Family      string
	Weight      FontWeight
	Flags       uint
	BBox        [4]float64
	ItalicAngle float64
	Ascent      float64
	Descent     float64
	CapHeight   float64
	XHeight     float64
	StemV       float64
	StemH       float64
}

var stdFonts = make(map[StdFontName]func() StdFont)

// IsStdFont check if a name is registered for a standard font.
func IsStdFont(name StdFontName) bool {
	_, ok := stdFonts[name]
	return ok
}

// NewStdFontByName creates a new StdFont by registered name. See RegisterStdFont.
func NewStdFontByName(name StdFontName) (StdFont, bool) {
	fnc, ok := stdFonts[name]
	if !ok {
		return StdFont{}, false
	}
	return fnc(), true
}

// RegisterStdFont registers a given StdFont constructor by font name. Font can then be created with NewStdFontByName.
func RegisterStdFont(name StdFontName, fnc func() StdFont, aliases ...StdFontName) {
	if _, ok := stdFonts[name]; ok {
		panic("font already registered: " + string(name))
	}
	stdFonts[name] = fnc
	for _, alias := range aliases {
		RegisterStdFont(alias, fnc)
	}
}

var _ Font = StdFont{}

// StdFont represents one of the built-in fonts and it is assumed that every reader has access to it.
type StdFont struct {
	desc    Descriptor
	metrics map[rune]CharMetrics
	encoder textencoding.TextEncoder
}

// NewStdFont returns a new instance of the font with a default encoder set (StandardEncoding).
func NewStdFont(desc Descriptor, metrics map[rune]CharMetrics) StdFont {
	return NewStdFontWithEncoding(desc, metrics, textencoding.NewStandardEncoder())
}

// NewStdFontWithEncoding returns a new instance of the font with a specified encoder.
func NewStdFontWithEncoding(desc Descriptor, metrics map[rune]CharMetrics, encoder textencoding.TextEncoder) StdFont {
	var nbsp rune = 0xA0
	if _, ok := metrics[nbsp]; !ok {
		// Use same metrics for 0xA0 (no-break space) and 0x20 (space).
		metrics[nbsp] = metrics[0x20]
	}

	return StdFont{
		desc:    desc,
		metrics: metrics,
		encoder: encoder,
	}
}

// Name returns a PDF name of the font.
func (font StdFont) Name() string {
	return string(font.desc.Name)
}

// Encoder returns the font's text encoder.
func (font StdFont) Encoder() textencoding.TextEncoder {
	return font.encoder
}

// GetRuneMetrics returns character metrics for a given rune.
func (font StdFont) GetRuneMetrics(r rune) (CharMetrics, bool) {
	metrics, has := font.metrics[r]
	return metrics, has
}

// GetMetricsTable is a method specific to standard fonts. It returns the metrics table of all glyphs.
// Caller should not modify the table.
func (font StdFont) GetMetricsTable() map[rune]CharMetrics {
	return font.metrics
}

// Descriptor returns a font descriptor.
func (font StdFont) Descriptor() Descriptor {
	return font.desc
}

// ToPdfObject returns a primitive PDF object representation of the font.
func (font StdFont) ToPdfObject() core.PdfObject {
	fontDict := core.MakeDict()
	fontDict.Set("Type", core.MakeName("Font"))
	fontDict.Set("Subtype", core.MakeName("Type1"))
	fontDict.Set("BaseFont", core.MakeName(font.Name()))
	fontDict.Set("Encoding", font.encoder.ToPdfObject())

	return core.MakeIndirectObject(fontDict)
}

// type1CommonRunes is list of runes common for some Type1 fonts. Used to unpack character metrics.
var type1CommonRunes = []rune{
	'A', '??', '??', '??', '??', '??', '??', '??', '??', '??',
	'??', 'B', 'C', '??', '??', '??', 'D', '??', '??', '???',
	'E', '??', '??', '??', '??', '??', '??', '??', '??', '??',
	'???', 'F', 'G', '??', '??', 'H', 'I', '??', '??', '??',
	'??', '??', '??', '??', 'J', 'K', '??', 'L', '??', '??',
	'??', '??', 'M', 'N', '??', '??', '??', '??', 'O', '??',
	'??', '??', '??', '??', '??', '??', '??', '??', 'P', 'Q',
	'R', '??', '??', '??', 'S', '??', '??', '??', '??', 'T',
	'??', '??', '??', 'U', '??', '??', '??', '??', '??', '??',
	'??', '??', 'V', 'W', 'X', 'Y', '??', '??', 'Z', '??',
	'??', '??', 'a', '??', '??', '??', '??', '??', '??', '??',
	'??', '&', '??', '??', '^', '~', '*', '@', '??', 'b',
	'\\', '|', '{', '}', '[', ']', '??', '??', '???', 'c',
	'??', '??', '??', '??', '??', '??', '??', ':', ',', '\uf6c3',
	'??', '??', 'd', '???', '???', '??', '??', '??', '??', '??',
	'$', '??', '??', 'e', '??', '??', '??', '??', '??', '??',
	'8', '???', '??', '???', '???', '??', '=', '??', '!', '??',
	'f', '???', '5', '???', '??', '4', '???', 'g', '??', '??',
	'??', '`', '>', '???', '??', '??', '???', '???', 'h', '??',
	'-', 'i', '??', '??', '??', '??', '??', '??', 'j', 'k',
	'??', 'l', '??', '??', '??', '<', '???', '??', '???', '??',
	'm', '??', '???', '??', '??', 'n', '??', '??', '??', '9',
	'???', '??', '#', 'o', '??', '??', '??', '??', '??', '??',
	'??', '??', '1', '??', '??', '??', '??', '??', '??', '??',
	'p', '??', '(', ')', '???', '%', '.', '??', '???', '+',
	'??', 'q', '?', '??', '"', '???', '???', '???', '???', '???',
	'???', '\'', 'r', '??', '???', '??', '??', '??', '??', 's',
	'??', '??', '??', '??', '??', ';', '7', '6', '/', ' ',
	'??', '???', 't', '??', '??', '??', '3', '??', '??', '??',
	'???', '2', '??', 'u', '??', '??', '??', '??', '??', '??',
	'_', '??', '??', 'v', 'w', 'x', 'y', '??', '??', '??',
	'z', '??', '??', '??', '0',
}
