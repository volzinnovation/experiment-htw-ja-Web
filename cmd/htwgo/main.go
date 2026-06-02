package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"htwgo/internal/interactive"
	"htwgo/internal/wumpus"
)

type config struct {
	seed         int64
	hasSeed      bool
	revealState  bool
	inertHazards bool
}

type app struct {
	session *interactive.Session
	config  config
	out     io.Writer
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, err := parseConfig(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	gameSeed := cfg.seed
	if !cfg.hasSeed {
		gameSeed = time.Now().UnixNano()
	}
	session := interactive.NewSessionWithSeed(gameSeed)
	session.Game().SetInertHazards(cfg.inertHazards)
	a := app{session: session, config: cfg, out: stdout}
	return a.play(stdin)
}

func parseConfig(args []string) (config, error) {
	flags := flag.NewFlagSet("htwgo", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	seed := flags.Int64("qa-seed", 0, "seed deterministic setup")
	reveal := flags.Bool("qa-reveal-state", false, "print hidden state each turn")
	inert := flags.Bool("qa-inert-hazards", false, "keep hazards visible but non-fatal")
	if err := flags.Parse(args); err != nil {
		return config{}, err
	}
	cfg := config{seed: *seed, revealState: *reveal, inertHazards: *inert}
	for _, arg := range args {
		if arg == "--qa-seed" || strings.HasPrefix(arg, "--qa-seed=") {
			cfg.hasSeed = true
		}
	}
	return cfg, nil
}

func (a *app) play(stdin io.Reader) int {
	scanner := bufio.NewScanner(stdin)
	if a.qaEnabled() {
		fmt.Fprintln(a.out, a.qaBanner())
	}
	fmt.Fprintln(a.out, "INSTRUCTIONS (Y-N)?")
	if !scanner.Scan() {
		return 0
	}
	printLines(a.out, a.session.AnswerInstructions(scanner.Text()))
	a.printTurn(a.session.BeginTurn())
	for scanner.Scan() {
		command := scanner.Text()
		if a.handleQACommand(command) {
			a.printTurn(a.session.DisplayTurn())
			continue
		}
		lines := a.session.EnterCommand(command)
		printLines(a.out, lines)
		switch a.session.Game().Status() {
		case wumpus.StatusLost:
			if !scanner.Scan() {
				return 0
			}
			a.session.AnswerSameSetup(scanner.Text())
			a.session.Game().SetInertHazards(a.config.inertHazards)
			a.printTurn(a.session.BeginTurn())
		case wumpus.StatusWon:
			return 0
		default:
			a.printTurn(a.session.DisplayTurn())
		}
	}
	return 0
}

func (a app) qaBanner() string {
	var modes []string
	if a.config.inertHazards {
		modes = append(modes, "HAZARDS INERT")
	}
	if a.config.revealState {
		modes = append(modes, "STATE REVEALED")
	}
	if a.config.hasSeed {
		modes = append(modes, "SEEDED SETUP")
	}
	return "QA MODE ENABLED: " + strings.Join(modes, ", ")
}

func (a *app) printTurn(lines []string) {
	if a.config.revealState {
		fmt.Fprintln(a.out, qaState(a.session.Game()))
	}
	printLines(a.out, lines)
}

func printLines(out io.Writer, lines []string) {
	for _, line := range lines {
		fmt.Fprintln(out, line)
	}
}

func qaState(game *wumpus.Game) string {
	setup := game.Setup()
	grenade := "none"
	if room, ok := game.GrenadeRoom(); ok {
		grenade = strconv.Itoa(room)
	}
	pending := "none"
	fuse := "none"
	if room, ok := game.PendingGrenadeRoom(); ok {
		pending = strconv.Itoa(room)
		fuse = "1"
	}
	wumpusState := "awake"
	if game.WumpusAsleep() {
		wumpusState = "asleep"
	}
	return fmt.Sprintf(
		"QA STATE: PLAYER=%d WUMPUS=%d WUMPUS_STATE=%s PITS=%s BATS=%s HHG=%s CARRYING_HHG=%t PENDING_HHG=%s FUSE=%s ARROWS=%d",
		setup.Player,
		setup.Wumpus,
		wumpusState,
		joinRooms(setup.Pits),
		joinRooms(setup.Bats),
		grenade,
		game.CarriesGrenade(),
		pending,
		fuse,
		game.Arrows(),
	)
}

func joinRooms(rooms []int) string {
	values := make([]string, 0, len(rooms))
	for _, room := range rooms {
		values = append(values, strconv.Itoa(room))
	}
	return strings.Join(values, ",")
}

func (a *app) handleQACommand(command string) bool {
	fields := strings.Fields(command)
	if len(fields) == 0 || strings.ToLower(fields[0]) != "qa" {
		return false
	}
	if !a.qaEnabled() {
		printLines(a.out, a.session.EnterCommand(command))
		return true
	}
	message, ok := a.applyQACommand(fields)
	if !ok {
		fmt.Fprintln(a.out, strings.ToUpper(command)+" IS NOT A COMMAND")
		return true
	}
	fmt.Fprintln(a.out, message)
	return true
}

func (a app) qaEnabled() bool {
	return a.config.revealState || a.config.inertHazards || a.config.hasSeed
}

func (a *app) applyQACommand(fields []string) (string, bool) {
	if len(fields) < 4 || strings.ToLower(fields[1]) != "set" {
		return "", false
	}
	target := strings.ToLower(fields[2])
	game := a.session.Game()
	switch target {
	case "player":
		room, ok := parseOneRoom(fields[3:])
		if !ok {
			return "", false
		}
		if err := game.SetPlayerRoomForQA(room); err != nil {
			return "", false
		}
		return fmt.Sprintf("QA SET: PLAYER=%d", room), true
	case "wumpus":
		room, ok := parseOneRoom(fields[3:])
		if !ok {
			return "", false
		}
		if err := game.SetWumpusRoomForQA(room); err != nil {
			return "", false
		}
		return fmt.Sprintf("QA SET: WUMPUS=%d", room), true
	case "pits":
		first, second, ok := parseTwoRooms(fields[3:])
		if !ok {
			return "", false
		}
		if err := game.SetPitRoomsForQA(first, second); err != nil {
			return "", false
		}
		return fmt.Sprintf("QA SET: PITS=%d,%d", first, second), true
	case "bats":
		first, second, ok := parseTwoRooms(fields[3:])
		if !ok {
			return "", false
		}
		if err := game.SetBatRoomsForQA(first, second); err != nil {
			return "", false
		}
		return fmt.Sprintf("QA SET: BATS=%d,%d", first, second), true
	case "hhg":
		if len(fields) == 4 && strings.EqualFold(fields[3], "none") {
			game.ClearGrenadeRoom()
			return "QA SET: HHG=none", true
		}
		room, ok := parseOneRoom(fields[3:])
		if !ok {
			return "", false
		}
		if err := game.SetGrenadeRoomForQA(room); err != nil {
			return "", false
		}
		return fmt.Sprintf("QA SET: HHG=%d", room), true
	case "arrows":
		arrows, ok := parseOneInt(fields[3:])
		if !ok {
			return "", false
		}
		if err := game.SetArrowsForQA(arrows); err != nil {
			return "", false
		}
		return fmt.Sprintf("QA SET: ARROWS=%d", arrows), true
	default:
		return "", false
	}
}

func parseOneRoom(fields []string) (int, bool) {
	room, ok := parseOneInt(fields)
	return room, ok && room >= 1 && room <= 20
}

func parseTwoRooms(fields []string) (int, int, bool) {
	if len(fields) != 2 {
		return 0, 0, false
	}
	first, err := strconv.Atoi(fields[0])
	if err != nil || first < 1 || first > 20 {
		return 0, 0, false
	}
	second, err := strconv.Atoi(fields[1])
	if err != nil || second < 1 || second > 20 {
		return 0, 0, false
	}
	return first, second, true
}

func parseOneInt(fields []string) (int, bool) {
	if len(fields) != 1 {
		return 0, false
	}
	value, err := strconv.Atoi(fields[0])
	return value, err == nil
}
