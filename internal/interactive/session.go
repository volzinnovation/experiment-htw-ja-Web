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
	turns    int
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

func (s *Session) SetTurnCount(turns int) {
	s.turns = turns
}

func (s *Session) TurnCount() int {
	return s.turns
}

func (s *Session) DisplayTurn() []string {
	lines := append([]string{}, s.game.TurnWarnings()...)
	setup := s.game.Setup()
	exits, _ := wumpus.NewCave().Exits(setup.Player)
	lines = append(lines,
		fmt.Sprintf("YOU ARE IN ROOM %d", setup.Player),
		fmt.Sprintf("TUNNELS LEAD TO %d %d %d", exits[0], exits[1], exits[2]),
		fmt.Sprintf("ARROWS LEFT: %d", s.game.Arrows()),
		s.prompt(),
	)
	return lines
}

func (s *Session) BeginTurn() []string {
	lines := append([]string{}, s.game.ResolveJumpingWumpusTurn().Messages...)
	return append(lines, s.DisplayTurn()...)
}

func (s *Session) prompt() string {
	if s.game.CarriesGrenade() {
		return "SHOOT, MOVE OR THROW (S-M-T)?"
	}
	return "SHOOT OR MOVE (S-M)?"
}

func (s *Session) EnterCommand(command string) []string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return []string{" IS NOT A COMMAND"}
	}
	return s.dispatchCommand(fields)
}

func (s *Session) dispatchCommand(fields []string) []string {
	action := strings.ToLower(fields[0])
	switch action {
	case "m":
		return s.moveCommand(fields)
	case "s":
		return s.shootCommand(fields)
	case "t":
		return s.throwCommand(fields)
	case "r", "rest":
		return s.restCommand(fields)
	default:
		return []string{strings.ToUpper(fields[0]) + " IS NOT A COMMAND"}
	}
}

func (s *Session) moveCommand(fields []string) []string {
	shouldDetonate := s.hasPendingGrenade()
	return s.roomCommand(fields, "CAN'T MOVE THERE", shouldDetonate, func(room int) (string, []string) {
		result := s.game.Move(room)
		return result.RejectedMessage, result.Messages
	})
}

func (s *Session) shootCommand(fields []string) []string {
	shouldDetonate := s.hasPendingGrenade()
	path, ok := parseShotPath(fields[1:])
	if !ok {
		return []string{"CAN'T SHOOT THERE"}
	}
	prefix := s.commandTurnMessages()
	if s.game.Status() != wumpus.StatusInProgress {
		return s.finishCommand(prefix, false)
	}
	result := s.game.Shoot(path)
	return s.finishCommand(append(prefix, result.Messages...), shouldDetonate)
}

func (s *Session) throwCommand(fields []string) []string {
	return s.roomCommand(fields, "CAN'T THROW THERE", false, func(room int) (string, []string) {
		result := s.game.ThrowGrenade(room)
		return result.RejectedMessage, result.Messages
	})
}

func (s *Session) roomCommand(fields []string, invalidMessage string, shouldDetonate bool, action func(int) (string, []string)) []string {
	room, ok := parseCommandRoom(fields)
	if !ok {
		return []string{invalidMessage}
	}
	prefix := s.commandTurnMessages()
	if s.game.Status() != wumpus.StatusInProgress {
		return s.finishCommand(prefix, false)
	}
	rejected, messages := action(room)
	if rejected != "" {
		return []string{rejected}
	}
	return s.finishCommand(append(prefix, messages...), shouldDetonate)
}

func (s *Session) restCommand(fields []string) []string {
	if len(fields) != 1 {
		return []string{strings.ToUpper(fields[0]) + " IS NOT A COMMAND"}
	}
	shouldDetonate := s.hasPendingGrenade()
	messages := s.commandTurnMessages()
	s.turns++
	return s.finishCommand(messages, shouldDetonate)
}

func parseCommandRoom(fields []string) (int, bool) {
	if len(fields) != 2 {
		return 0, false
	}
	room, err := strconv.Atoi(fields[1])
	return room, err == nil
}

func (s *Session) commandTurnMessages() []string {
	return s.game.ResolveJumpingWumpusTurn().Messages
}

func parseShotPath(values []string) ([]int, bool) {
	if !validShotPathLength(len(values)) {
		return nil, false
	}
	var path []int
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
	if err != nil {
		return 0, false
	}
	if room < 1 || room > 20 {
		return 0, false
	}
	return room, true
}

func (s *Session) finishCommand(messages []string, shouldDetonate bool) []string {
	lines := append([]string{}, messages...)
	if shouldDetonate {
		lines = append(lines, s.game.DetonateGrenade()...)
	}
	if s.game.Status() == wumpus.StatusLost {
		s.lostGame = s.game
		return append(lines, "SAME SET UP (Y-N)?")
	}
	return lines
}

func (s *Session) hasPendingGrenade() bool {
	_, ok := s.game.PendingGrenadeRoom()
	return ok
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

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T08:32:47-05:00","module_hash":"c7bdc4549a83854cc1824b2a9032e559b19cee6e2e1f65c5e8f3b190792b791e","functions":[{"id":"func/NewSession","name":"NewSession","line":19,"end_line":21,"hash":"c8c79e4c5e1027ada492328b99c919d0495af2f241f5cc80aa0349fdac5da99f"},{"id":"func/NewSessionWithGame","name":"NewSessionWithGame","line":23,"end_line":25,"hash":"cf0c3db4bd5cd94cc98b1f5570d38e0f01e52d4c958290e3b2bbfca5b987c28b"},{"id":"func/NewSessionWithSeed","name":"NewSessionWithSeed","line":27,"end_line":30,"hash":"d164ac2ef226451bcc11f8cc9d15f1423a829c9c4c8997ad74a6b4ef788a10e2"},{"id":"func/Session.Game","name":"Session.Game","line":32,"end_line":34,"hash":"795de1a673f230936a72ee82210fb1fb11679cf9a033493878509b7ef5bd9347"},{"id":"func/Session.DisplayTurn","name":"Session.DisplayTurn","line":36,"end_line":47,"hash":"ae27de9fa038ac2c2780972615457e1b27b352d1ac2ce90802517981037d25ec"},{"id":"func/Session.EnterCommand","name":"Session.EnterCommand","line":49,"end_line":63,"hash":"883ebc98ef1d68d9cf63e397b284e78b6a216b7588bb45497c5f5742713587c5"},{"id":"func/Session.moveCommand","name":"Session.moveCommand","line":65,"end_line":78,"hash":"0a8b22ae3a42c127495b2de14e41ee55653b148d0a3d061c393db7a530d39000"},{"id":"func/Session.shootCommand","name":"Session.shootCommand","line":80,"end_line":87,"hash":"b4ae4b44b4790ab7ff4732fae67c2591b3e3a16298bc9e36009f64c8857c4a6f"},{"id":"func/parseShotPath","name":"parseShotPath","line":89,"end_line":102,"hash":"ce1353165231d053fb636808073dd34c6a0e0e7030b600945a3222cd51950c7c"},{"id":"func/validShotPathLength","name":"validShotPathLength","line":104,"end_line":106,"hash":"7d64f126f87f9b256d21430251f8a11872229496cb463956612463c4bf9712a8"},{"id":"func/parseRoom","name":"parseRoom","line":108,"end_line":117,"hash":"cafb5a76f1606e6a442309f0d1b9f50a95dd2ad8f42d87f6a5f2b6b6e884b464"},{"id":"func/Session.finishCommand","name":"Session.finishCommand","line":119,"end_line":125,"hash":"ccbdb3162f25a29496eba6c390cdee46ec4fe194b709a6d57a34581d43f3c7ae"},{"id":"func/Session.AnswerSameSetup","name":"Session.AnswerSameSetup","line":127,"end_line":139,"hash":"6243c961e6e7ebc0f4e54fb70cd37e87a6a740786b2809989b7fdde312620144"},{"id":"func/Session.AnswerInstructions","name":"Session.AnswerInstructions","line":141,"end_line":146,"hash":"82835edd4b56b27547be487cc944fb6b43755b1793a0f611a03ca7a78aafe217"},{"id":"func/Session.MarkLostForTest","name":"Session.MarkLostForTest","line":148,"end_line":150,"hash":"847078d1f1668277f6ab17a5524a09865f6e0bdce1a3498e911779f1427573c0"}]}
// mutate4go-manifest-end
