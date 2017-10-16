// Package bintool is a collection of helper for building cli tools.
//
// TODO:
// - introduce colors
// - refactor remember file interaction
// - add tests
// - refactor failing logic
// - make ability to remove/edit previous saved values
// - ability to pass all inputs as argument
package bintool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func Ask(what, defaultVal, rememberFile string) string {
	rdata := make(map[string][]string, 0)
	var answer string

	if rememberFile != "" {
		if _, err := os.Stat(rememberFile); err == nil {
			fd, err := ioutil.ReadFile(rememberFile)
			if err != nil {
				log.Fatalf("failed reading %s: %s", rememberFile, err)
			}
			err = json.Unmarshal(fd, &rdata)
			if err != nil {
				log.Fatalf("failed to unmarshal %s: %s", rememberFile, err)
			}
		}

		defer func() {
			if _, ok := rdata[what]; ok {
				for _, v := range rdata[what] {
					if v == answer {
						goto save
					}
				}
			} else {
				rdata[what] = make([]string, 0)
			}
			rdata[what] = append(rdata[what], answer)

		save:
			fd, err := json.Marshal(rdata)
			if err != nil {
				log.Fatalf("failed marshaling remembers data: %s", err)
			}

			err = ioutil.WriteFile(rememberFile, fd, 0664)
			if err != nil {
				log.Fatalf("failed writing to %s: %s", rememberFile, err)
			}
		}()
	}

askQuestion:
	switch {
	case len(rdata[what]) > 0:
		fmt.Printf("%s\n", what)
		for i, v := range rdata[what] {
			fmt.Printf("\u001b[33m%d)\u001b[0m %s\n", i+1, v)
		}
		switch {
		case defaultVal != "":
			// TODO bug, what if user whant to enter number as
			// value not as selection of previous value.
			fmt.Printf("%s [%s]: ", "Chose previous value by number or enter new", defaultVal)
			fmt.Scanln(&answer)
			if answer == "" && defaultVal != "" {
				answer = defaultVal
			}
		default:
			fmt.Printf("%s ", "Chose previous value by number or enter new:")
			fmt.Scanln(&answer)
		}

		if numAnswer, err := strconv.Atoi(answer); err == nil {
			if len(rdata[what]) >= numAnswer {
				answer = rdata[what][numAnswer-1]
			}
		}
	case defaultVal != "":
		fmt.Printf("%s [%s]: ", what, defaultVal)
		fmt.Scanln(&answer)
		if answer == "" && defaultVal != "" {
			answer = defaultVal
		}
	default:
		fmt.Printf("%s ", what+":")
		fmt.Scanln(&answer)
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		goto askQuestion
	}

	return answer
}

func Confirm(what string, yesDef bool) bool {
	def := "n"
	if yesDef {
		def = "y"
	}

ask:
	fmt.Printf("%s (%s/%s)[%s]: ", what, "y", "n", def)

	var answer string
	fmt.Scanln(&answer)
	if answer == "" {
		return yesDef
	}
	answer = strings.ToLower(answer)
	if answer != "y" && answer != "n" {
		goto ask
	}

	return answer == "y"
}
