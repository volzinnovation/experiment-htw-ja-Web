package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunStartsPlayableConsole(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"--qa-seed=1973"}, strings.NewReader("n\n"), &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{
		"QA MODE ENABLED: SEEDED SETUP",
		"INSTRUCTIONS (Y-N)?",
		"YOU ARE IN ROOM",
		"TUNNELS LEAD TO",
		"ARROWS LEFT: 5",
		"SHOOT OR MOVE (S-M)?",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestRunPrintsFullInstructions(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"--qa-seed=1973"}, strings.NewReader("Y\n"), &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{
		"WELCOME TO 'HUNT THE WUMPUS'",
		"THE WUMPUS LIVES IN A CAVE OF 20 ROOMS",
		"HAZARDS:",
		"SUPER BATS",
		"WHEN YOU ARE ONE ROOM AWAY FROM WUMPUS OR HAZARD",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestRunRevealsStateAndAppliesQASetupCommands(t *testing.T) {
	var stdout, stderr bytes.Buffer
	input := strings.NewReader(strings.Join([]string{
		"n",
		"qa set player 1",
		"qa set wumpus 2",
		"qa set pits 11 12",
		"qa set bats 13 14",
		"",
	}, "\n"))

	code := run([]string{"--qa-reveal-state", "--qa-seed=1973"}, input, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{
		"QA MODE ENABLED: STATE REVEALED, SEEDED SETUP",
		"QA SET: PLAYER=1",
		"QA SET: WUMPUS=2",
		"QA SET: PITS=11,12",
		"QA SET: BATS=13,14",
		"QA STATE: PLAYER=1 WUMPUS=2",
		"I SMELL A WUMPUS",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestRunInertHazardsSuppressesPitLoss(t *testing.T) {
	var stdout, stderr bytes.Buffer
	input := strings.NewReader(strings.Join([]string{
		"n",
		"qa set player 1",
		"qa set pits 2 12",
		"m 2",
		"",
	}, "\n"))

	code := run([]string{"--qa-inert-hazards", "--qa-seed=1973"}, input, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "QA INERT: PIT IGNORED") {
		t.Fatalf("output missing inert pit message:\n%s", output)
	}
	if strings.Contains(output, "HA HA HA - YOU LOSE!") {
		t.Fatalf("inert hazard should not lose:\n%s", output)
	}
}
