package bintool

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Execute(cmd string) {
	fmt.Printf("%s %s\n", "Executing", cmd)
	cp := strings.FieldsFunc(cmd, inQouteSplit())

	args := []string{}
	if len(cp) > 1 {
		args = cp[1:]
		cleanQoutedParams(args)
	}

	c := exec.Command(cp[0], args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err := c.Run()
	if err != nil {
		log.Fatalf("failed running cmds %q: %s", cp, err)
	}
}

func inQouteSplit() func(r rune) bool {
	// TODO bug, make qoute split type aware, aka.
	// from ' to ' not brake on " if first encounter
	inQoute := false
	return func(r rune) bool {
		if r == '\'' || r == '"' {
			inQoute = !inQoute
		}
		return inQoute == false && r == ' '
	}
}

func cleanQoutedParams(in []string) {
	for i := range in {
		in[i] = strings.Trim(in[i], "\"'")
	}
}
