package brdgme

import "image/color"

// Default colors for use in brdgme games.
var (
	ColorRed        = color.RGBA{244, 67, 54, 255}
	ColorPink       = color.RGBA{233, 30, 99, 255}
	ColorPurple     = color.RGBA{156, 39, 176, 255}
	ColorDeepPurple = color.RGBA{103, 58, 183, 255}
	ColorIndigo     = color.RGBA{63, 81, 181, 255}
	ColorBlue       = color.RGBA{33, 150, 243, 255}
	ColorLightBlue  = color.RGBA{3, 169, 244, 255}
	ColorCyan       = color.RGBA{0, 188, 212, 255}
	ColorTeal       = color.RGBA{0, 150, 136, 255}
	ColorGreen      = color.RGBA{76, 175, 80, 255}
	ColorLightGreen = color.RGBA{139, 195, 74, 255}
	ColorLime       = color.RGBA{205, 220, 57, 255}
	ColorYellow     = color.RGBA{255, 235, 59, 255}
	ColorAmber      = color.RGBA{255, 193, 7, 255}
	ColorOrange     = color.RGBA{255, 152, 0, 255}
	ColorDeepOrange = color.RGBA{255, 87, 34, 255}
	ColorBrown      = color.RGBA{121, 85, 72, 255}
	ColorGrey       = color.RGBA{158, 158, 158, 255}
	ColorBlueGrey   = color.RGBA{96, 125, 139, 255}
	ColorWhite      = color.RGBA{255, 255, 255, 255}
	ColorBlack      = color.RGBA{0, 0, 0, 255}
)

// PlayerColors are a subset of the default colours suitable for player
// coloring.  These should be used in correct order to match log rendering.
var PlayerColors = []color.Color{
	ColorGreen,
	ColorRed,
	ColorBlue,
	ColorOrange,
	ColorPurple,
	ColorBrown,
	ColorBlueGrey,
}

// PlayerColor gets the player colour for a given player number.
func PlayerColor(p int) color.Color {
	return PlayerColors[p%len(PlayerColors)]
}
