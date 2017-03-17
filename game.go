package brdgme

// Gamer is a playable game.
type Gamer interface {
	Start(players int) ([]Log, error)
	Command(
		player int,
		input string,
		playerNames []string,
	) (logs []Log, remaining string, err error)
	IsFinished() bool
	Winners() []int
	WhoseTurn() []int
	Render(player *int) string
}

// Eliminator is a game where players can be eliminated.
type Eliminator interface {
	Eliminated() []int
}
