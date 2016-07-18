package brdgme

// PlayerTemplater can generate a renderable template for a player.
type PlayerTemplater interface {
	PlayerTemplate(player int) (string, error)
}

// SpectatorTemplater can generate a renderable template for a spectator.
type SpectatorTemplater interface {
	SpectatorTemplate() string
}
