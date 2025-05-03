package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Snap struct {
	Name string
	IP   string
}

func main() {
	snaps := []Snap{
		{"Snap1", "192.168.1.9:22"},
		{"Snap2", "192.168.1.46:22"},
	}

	username := "MOAMN"
	password := "Jpanzer2"

	reader := bufio.NewReader(os.Stdin)

	for len(snaps) > 0 {
		fmt.Println("\nConnected Snap Devices:")
		for i, snap := range snaps {
			fmt.Printf("[%d] %s (%s)\n", i+1, snap.Name, snap.IP)
		}

		fmt.Print("Enter number of Snap to shut down (or 0 to exit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 0 || choice > len(snaps) {
			fmt.Println("❌ Invalid choice. Try again.")
			continue
		}

		if choice == 0 {
			fmt.Println("👋 Exiting.")
			break
		}

		selected := snaps[choice-1]
		fmt.Printf("⚠️ Sending shutdown command to %s (%s)...\n", selected.Name, selected.IP)

		if shutdownSnap(selected.IP, username, password) {
			// Remove from list
			snaps = append(snaps[:choice-1], snaps[choice:]...)
		}
	}

	fmt.Println("✅ All selected Snaps have been shut down.")
}

func shutdownSnap(addr, user, pass string) bool {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("Could not connect to %s: %v\n", addr, err)
		return false
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Printf("Failed to create session for %s: %v\n", addr, err)
		return false
	}
	defer session.Close()

	cmd := `shutdown /s /t 0`
	err = session.Run(cmd)
	if err != nil {
		log.Printf("Command failed on %s: %v\n", addr, err)
		return false
	}

	fmt.Printf("Shutdown command sent to %s\n", addr)
	return true
}
