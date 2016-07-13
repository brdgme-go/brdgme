package brdgme

import (
	"image"

	"github.com/llgcode/draw2d"
)

// PlayerTemplater can generate a renderable template for a player.
type PlayerTemplater interface {
	PlayerTemplate(player int) (string, error)
}

// SpectatorTemplater can generate a renderable template for a spectator.
type SpectatorTemplater interface {
	SpectatorTemplate() string
}

// PlayerRenderer can render graphically for a player.
type PlayerRenderer interface {
	PlayerRender(
		player int,
		gc draw2d.GraphicContext,
		bounds image.Rectangle,
	) error
}

// SpectatorRenderer can render graphically for a spectator.
type SpectatorRenderer interface {
	SpectatorRender(
		gc draw2d.GraphicContext,
		bounds image.Rectangle,
	)
}
