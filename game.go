package brdgme

import (
	"io"
)

// Gamer is a playable game.
//
// At a minimum a game only has a log, which works for incredibly simple games.
// More complex games will want to implement one or more rendering interfaces
// such as `PlayerTemplater` or `PlayerRenderer`.
type Gamer interface {
	Name() string
	Identifier() string
	Start(players int) ([]Log, error)
	AvailableCommands(player int) []CommandDescription
	Command(
		player int,
		input io.Reader,
		playerNames []string,
	) (logs []Log, remaining io.Reader, err error)
	IsFinished() (finished bool, winners []int)
	WhoseTurn() []int
}

// A CommandDescription is a description of a command
type CommandDescription struct {
	Name        string
	Description string
	Example     string
}

// Eliminator can have eliminated players.
type Eliminator interface {
	Eliminated() []int
}
