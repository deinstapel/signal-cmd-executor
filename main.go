package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/abiosoft/lineprefix"
)

type Config struct {
	Signals map[string]string `json:"signals"`
}

var signalMap = map[string]syscall.Signal{
	"USR1": syscall.SIGUSR1,
	"USR2": syscall.SIGUSR2,
	"INT":  syscall.SIGINT,
	"ABRT": syscall.SIGABRT,
	"ALRM": syscall.SIGALRM,
	"BUS":  syscall.SIGBUS,
	"CHLD": syscall.SIGCHLD,
	"CLD":  syscall.SIGCLD,
	"CONT": syscall.SIGCONT,
	"FPE":  syscall.SIGFPE,
	"HUP":  syscall.SIGHUP,
	"IO":   syscall.SIGIO,
	"IOT":  syscall.SIGIOT,
}

// handles a single signal
func handleSignal(ctx context.Context, signalName string, command string) {
	c := make(chan os.Signal, 3)
	signalNumber, ok := signalMap[signalName]
	if !ok {
		log.Fatalf("failed to find syscall for signal %v\n", signalName)
	}

	signal.Notify(c, signalNumber)

	prefix := lineprefix.Prefix(fmt.Sprintf("[CMD SIG%v]", signalName))
	stdoutWrapper := lineprefix.New(prefix, lineprefix.Writer(os.Stdout))
	stderrWrapper := lineprefix.New(prefix, lineprefix.Writer(os.Stderr))

outer:
	for {
		select {
		case <-c:
			cmd := exec.Command("/bin/bash", "-c", command)
			cmd.Stderr = stderrWrapper
			cmd.Stdout = stdoutWrapper
			if err := cmd.Run(); err != nil {
				log.Printf("WARNING: program terminated with error: %v\n", err)
			}
		case <-ctx.Done():
			break outer
		}
	}

	signal.Stop(c)
	close(c)
}

// main entry point for the program, loads config file etc.
func main() {
	configFile, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		configFile = "/etc/config.json"
	}

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config file: %v\n", err)
	}

	config := Config{}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		log.Fatalf("failed to deserialize config file: %v\n", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	for signal, command := range config.Signals {
		go handleSignal(ctx, signal, command)
	}

	<-ctx.Done()

	log.Printf("exiting")
}
