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

type commandAction func() (string, []string)

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
	beforeCommand := s.game
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return []string{" IS NOT A COMMAND"}
	}
	lines := s.dispatchCommand(fields)
	if s.game.Status() == wumpus.StatusLost {
		s.lostGame = beforeCommand
	}
	return lines
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
	room, ok := parseCommandRoom(fields)
	if !ok {
		return []string{"CAN'T MOVE THERE"}
	}
	return s.executeTurn(shouldDetonate, false, func() (string, []string) {
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
	return s.executeTurn(shouldDetonate, false, func() (string, []string) {
		result := s.game.Shoot(path)
		return result.RejectedMessage, result.Messages
	})
}

func (s *Session) throwCommand(fields []string) []string {
	room, ok := parseCommandRoom(fields)
	if !ok {
		return []string{"CAN'T THROW THERE"}
	}
	return s.executeTurn(false, false, func() (string, []string) {
		result := s.game.ThrowGrenade(room)
		return result.RejectedMessage, result.Messages
	})
}

func (s *Session) executeTurn(shouldDetonate bool, incrementTurn bool, action commandAction) []string {
	prefix := s.commandTurnMessages()
	if incrementTurn {
		s.turns++
	}
	if s.game.Status() != wumpus.StatusInProgress {
		return s.finishCommand(prefix, false)
	}
	rejected, messages := action()
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
	return s.executeTurn(shouldDetonate, true, func() (string, []string) {
		return "", nil
	})
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
		s.game = s.lostGame.ReplaySnapshot()
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
		return instructions()
	}
	return nil
}

func (s *Session) MarkLostForTest() {
	s.lostGame = s.game
}

func instructions() []string {
	return []string{
		"WELCOME TO 'HUNT THE WUMPUS'",
		"",
		"THE WUMPUS LIVES IN A CAVE OF 20 ROOMS: EACH ROOM HAS 3 TUNNELS LEADING TO OTHER",
		"ROOMS. (LOOK AT A DODECAHEDRON TO SEE HOW THIS WORKS. IF YOU DON'T KNOW WHAT A",
		"DODECAHEDRON IS, ASK SOMEONE)",
		"",
		"***",
		"HAZARDS:",
		"",
		"BOTTOMLESS PITS - TWO ROOMS HAVE BOTTOMLESS PITS IN THEM",
		"IF YOU GO THERE: YOU FALL INTO THE PIT (& LOSE!)",
		"",
		"SUPER BATS  - TWO OTHER ROOMS HAVE SUPER BATS. IF YOU GO THERE, A BAT GRABS YOU",
		"AND TAKES YOU TO SOME OTHER ROOM AT RANDOM. (WHICH MIGHT BE TROUBLESOME)",
		"",
		"WUMPUS:",
		"",
		"THE WUMPUS IS NOT BOTHERED BY THE HAZARDS (HE HAS SUCKER FEET AND IS TOO BIG FOR",
		"A BAT TO LIFT). USUALLY HE IS ASLEEP. TWO THINGS WAKE HIM UP: YOUR ENTERING HIS",
		"ROOM OR YOUR SHOOTING AN ARROW.",
		"",
		"IF THE WUMPUS WAKES, HE MOVES (P=0.75) ONE ROOM OR STAYS STILL (P=0.25). AFTER",
		"THAT, IF HE IS WHERE YOU ARE, HE EATS YOU UP (& YOU LOSE!)",
		"",
		"YOU:",
		"",
		"EACH TURN YOU MAY MOVE OR SHOOT A CROOKED ARROW",
		"MOVING: YOU CAN GO ONE ROOM (THRU ONE TUNNEL)",
		"ARROWS: YOU HAVE 5 ARROWS. YOU LOSE WHEN YOU RUN OUT.",
		"",
		"EACH ARROW CAN GO FROM 1 TO 5 ROOMS: YOU AIM BY TELLING THE COMPUTER THE ROOMS",
		"YOU WANT THE ARROW TO GO TO. IF THE ARROW CAN'T GO THAT WAY (IE NO TUNNEL) IT",
		"MOVES AT RANDOM TO THE NEXT ROOM.",
		"",
		"IF THE ARROW HITS THE WUMPUS: YOU WIN.",
		"",
		"IF THE ARROW HITS YOU: YOU LOSE.",
		"",
		"WARNINGS:",
		"",
		"WHEN YOU ARE ONE ROOM AWAY FROM WUMPUS OR HAZARD, THE COMPUTER SAYS:",
		"",
		"WUMPUS - 'I SMELL A WUMPUS'",
		"",
		"BAT - 'BATS NEARBY'",
		"",
		"PIT - 'I FEEL A DRAFT'",
	}
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T09:55:58-05:00","module_hash":"4667ef3a3a5b693e591c1703a6dad614ddc94284039650fa5a7bcbf097e081d3","functions":[{"id":"func/NewSession","name":"NewSession","line":20,"end_line":22,"hash":"c8c79e4c5e1027ada492328b99c919d0495af2f241f5cc80aa0349fdac5da99f"},{"id":"func/NewSessionWithGame","name":"NewSessionWithGame","line":24,"end_line":26,"hash":"cf0c3db4bd5cd94cc98b1f5570d38e0f01e52d4c958290e3b2bbfca5b987c28b"},{"id":"func/NewSessionWithSeed","name":"NewSessionWithSeed","line":28,"end_line":31,"hash":"d164ac2ef226451bcc11f8cc9d15f1423a829c9c4c8997ad74a6b4ef788a10e2"},{"id":"func/Session.Game","name":"Session.Game","line":33,"end_line":35,"hash":"795de1a673f230936a72ee82210fb1fb11679cf9a033493878509b7ef5bd9347"},{"id":"func/Session.SetTurnCount","name":"Session.SetTurnCount","line":37,"end_line":39,"hash":"53cda1db3effcc7e5aa5b1c0ccf681f4d852254503b89615270a33f14a0094e7"},{"id":"func/Session.TurnCount","name":"Session.TurnCount","line":41,"end_line":43,"hash":"4f429976284318ed05b04415d9af0630db57d0e74981476aa7840cc2ed561429"},{"id":"func/Session.DisplayTurn","name":"Session.DisplayTurn","line":45,"end_line":56,"hash":"0c0375b5997ba7676e62a9205dee3729b7ce8ea99c847f18143a266dd6dcfc02"},{"id":"func/Session.BeginTurn","name":"Session.BeginTurn","line":58,"end_line":61,"hash":"62bad99e61e52d18be126401d0be8e49e001228e963324bbb1c8170cedaddb0c"},{"id":"func/Session.prompt","name":"Session.prompt","line":63,"end_line":68,"hash":"8cb0be9ad3a4c496ead356aa8e69ed9b69f76013221acb579a4a83a9c2bfc2db"},{"id":"func/Session.EnterCommand","name":"Session.EnterCommand","line":70,"end_line":76,"hash":"8de304a8fcd3758f9551672fc99d26ff7b2c2eb6d83d9632a642f5687773e9d8"},{"id":"func/Session.dispatchCommand","name":"Session.dispatchCommand","line":78,"end_line":92,"hash":"cb34368414e1ed4e561fcefc7e8982e264c9fad0bb4221d853ed406a214342c3"},{"id":"func/Session.moveCommand","name":"Session.moveCommand","line":94,"end_line":100,"hash":"eb67a7d36d92b838639e59f189881cb53b801ed8269806cda014076838124f0b"},{"id":"func/Session.shootCommand","name":"Session.shootCommand","line":102,"end_line":114,"hash":"2be3bbba1a1c02f7c13f6293ea8aea1ecce4a7e586a6942cae62f699f7ee600b"},{"id":"func/Session.throwCommand","name":"Session.throwCommand","line":116,"end_line":121,"hash":"95e1328fa7a9ed4e404e5ec5364a9a437d993d337a3ecafb6d04232a2c425c43"},{"id":"func/Session.roomCommand","name":"Session.roomCommand","line":123,"end_line":137,"hash":"57f7d05b2b34ec5fe0c80291755d4c9f122ec75828670f2b6eabfa307912d17d"},{"id":"func/Session.restCommand","name":"Session.restCommand","line":139,"end_line":147,"hash":"37e022516b7b4cff973e44e76327290cdb2a054856e4453477bd90d4b89eb198"},{"id":"func/parseCommandRoom","name":"parseCommandRoom","line":149,"end_line":155,"hash":"4891e0d6e0148b7083bba8d01cda627db21647b9076954e7fa746d33c5df5a7b"},{"id":"func/Session.commandTurnMessages","name":"Session.commandTurnMessages","line":157,"end_line":159,"hash":"8811f0acd0f7f9e2c353b668b46f06e64abd8b5148dfa528bbc0b3cbf1951ba0"},{"id":"func/parseShotPath","name":"parseShotPath","line":161,"end_line":174,"hash":"ce1353165231d053fb636808073dd34c6a0e0e7030b600945a3222cd51950c7c"},{"id":"func/validShotPathLength","name":"validShotPathLength","line":176,"end_line":178,"hash":"7d64f126f87f9b256d21430251f8a11872229496cb463956612463c4bf9712a8"},{"id":"func/parseRoom","name":"parseRoom","line":180,"end_line":189,"hash":"cafb5a76f1606e6a442309f0d1b9f50a95dd2ad8f42d87f6a5f2b6b6e884b464"},{"id":"func/Session.finishCommand","name":"Session.finishCommand","line":191,"end_line":201,"hash":"cbc3b54cbb97c49d333d1e66002be69dee84323ffd93209ed4181588b6302d24"},{"id":"func/Session.hasPendingGrenade","name":"Session.hasPendingGrenade","line":203,"end_line":206,"hash":"a2b77d228d49370abdee71d24c7bc623c142491e4da2d755a298294d761f657c"},{"id":"func/Session.AnswerSameSetup","name":"Session.AnswerSameSetup","line":208,"end_line":220,"hash":"6243c961e6e7ebc0f4e54fb70cd37e87a6a740786b2809989b7fdde312620144"},{"id":"func/Session.AnswerInstructions","name":"Session.AnswerInstructions","line":222,"end_line":227,"hash":"82835edd4b56b27547be487cc944fb6b43755b1793a0f611a03ca7a78aafe217"},{"id":"func/Session.MarkLostForTest","name":"Session.MarkLostForTest","line":229,"end_line":231,"hash":"847078d1f1668277f6ab17a5524a09865f6e0bdce1a3498e911779f1427573c0"}]}
// mutate4go-manifest-end
