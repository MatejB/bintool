package bintool

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var remoteAddsRegex = regexp.MustCompile("^(.+@[^:]+):.*")

// TransferRemote will make remote transfer from to.
// Remote address is expected if form user@address:location.
// If customPort is 0 default one is used.
func TransferRemote(from, to string, customPort int) {
	// which one is remote?
	remote := to
	if ix := strings.IndexRune(from, '@'); ix != -1 {
		remote = from
	}
	remoteMatches := remoteAddsRegex.FindStringSubmatch(remote)
	if len(remoteMatches) < 2 {
		log.Fatalf("could not find remote address in %q", remote)
	}

	// assume rsync
	port := ""
	if customPort != 0 {
		port = fmt.Sprintf("-e 'ssh -p %d' ", customPort)
	}
	execCmd := fmt.Sprintf("rsync %s-zhr %s %s", port, from, to)

	// check rsync
	checkArgs := make([]string, 0)
	if customPort != 0 {
		checkArgs = append(checkArgs, fmt.Sprintf("-p%d", customPort))
	}
	checkArgs = append(checkArgs, remoteMatches[1])
	checkArgs = append(checkArgs, "which rsync")

	err := exec.Command("ssh", checkArgs...).Run()
	if err != nil {
		// fallback to scp
		port := ""
		if customPort != 0 {
			port = fmt.Sprintf("-P %d ", customPort)
		}
		execCmd = fmt.Sprintf("scp %s-r %s %s", port, from, to)
	}

	cp := strings.FieldsFunc(execCmd, inQouteSplit())

	cleanQoutedParams(cp[1:])
	c := exec.Command(cp[0], cp[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	fmt.Printf("%s %s\n", "Executing", execCmd)

	err = c.Run()
	if err != nil {
		log.Fatalf("failed running cmds %q: %s", execCmd, err)
	}
}
