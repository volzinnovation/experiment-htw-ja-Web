package interactive

import (
	"fmt"
	"strconv"
	"strings"

	"htwgo/internal/wumpus"
)

type Session struct {
	game     wumpus.Game
	lostGame wumpus.Game
	seed     int64
	hasSeed  bool
	nextSeed int64
}

func NewSession() *Session {
	return &Session{}
}

func NewSessionWithGame(game wumpus.Game) *Session {
	return &Session{game: game}
}

func NewSessionWithSeed(seed int64) *Session {
	game, _ := wumpus.NewGame(seed)
	return &Session{game: game, seed: seed, hasSeed: true, nextSeed: seed + 1}
}

func (s *Session) Game() *wumpus.Game {
	return &s.game
}

func (s *Session) DisplayTurn() []string {
	lines := append([]string{}, s.game.TurnWarnings()...)
	setup := s.game.Setup()
	exits, _ := wumpus.NewCave().Exits(setup.Player)
	lines = append(lines,
		fmt.Sprintf("YOU ARE IN ROOM %d", setup.Player),
		fmt.Sprintf("TUNNELS LEAD TO %d %d %d", exits[0], exits[1], exits[2]),
		fmt.Sprintf("ARROWS LEFT: %d", s.game.Arrows()),
		"SHOOT OR MOVE (S-M)?",
	)
	return lines
}

func (s *Session) EnterCommand(command string) []string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return []string{" IS NOT A COMMAND"}
	}
	action := strings.ToLower(fields[0])
	switch action {
	case "m":
		return s.moveCommand(fields)
	case "s":
		return s.shootCommand(fields)
	default:
		return []string{strings.ToUpper(fields[0]) + " IS NOT A COMMAND"}
	}
}

func (s *Session) moveCommand(fields []string) []string {
	if len(fields) != 2 {
		return []string{"CAN'T MOVE THERE"}
	}
	room, err := strconv.Atoi(fields[1])
	if err != nil {
		return []string{"CAN'T MOVE THERE"}
	}
	result := s.game.Move(room)
	if result.RejectedMessage != "" {
		return []string{result.RejectedMessage}
	}
	return s.finishCommand(result.Messages)
}

func (s *Session) shootCommand(fields []string) []string {
	path, ok := parseShotPath(fields[1:])
	if !ok {
		return []string{"CAN'T SHOOT THERE"}
	}
	result := s.game.Shoot(path)
	return s.finishCommand(result.Messages)
}

func parseShotPath(values []string) ([]int, bool) {
	if !validShotPathLength(len(values)) {
		return nil, false
	}
	path := make([]int, 0, len(values))
	for _, value := range values {
		room, ok := parseRoom(value)
		if !ok {
			return nil, false
		}
		path = append(path, room)
	}
	return path, true
}

func validShotPathLength(length int) bool {
	return length >= 1 && length <= 5
}

func parseRoom(value string) (int, bool) {
	room, err := strconv.Atoi(value)
	if err != nil || room < 1 || room > 20 {
		return 0, false
	}
	return room, true
}

func (s *Session) finishCommand(messages []string) []string {
	if s.game.Status() == wumpus.StatusLost {
		s.lostGame = s.game
		return append(append([]string{}, messages...), "SAME SET UP (Y-N)?")
	}
	return append([]string{}, messages...)
}

func (s *Session) AnswerSameSetup(answer string) {
	if strings.EqualFold(answer, "y") {
		s.game = s.lostGame.ReplaySameSetup()
		return
	}
	if s.hasSeed {
		game, err := wumpus.NewGame(s.nextSeed)
		if err == nil {
			s.nextSeed++
			s.game = game
		}
	}
}

func (s *Session) AnswerInstructions(answer string) []string {
	if strings.EqualFold(answer, "y") {
		return []string{"WELCOME TO 'HUNT THE WUMPUS'"}
	}
	return nil
}

func (s *Session) MarkLostForTest() {
	s.lostGame = s.game
}
