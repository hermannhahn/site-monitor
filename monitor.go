package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	showMenu()
}

// show menu
func showMenu() {
	cls()
	fmt.Println("==========================================================")
	fmt.Println("                     SITE MONITOR")
	fmt.Println("==========================================================")
	fmt.Println()
	fmt.Println("Enter an option number:")
	fmt.Println()
	fmt.Println("(1) Start monitor")
	fmt.Println("(0) Back")
	fmt.Println()
	fmt.Printf("> ")
	menu := getMenuOption()

	switch menu {
	case "0":
		fmt.Println("Exiting...")
		cls()
		os.Exit(0)
	case "1":
		cls()
		checkSites()
	default:
		showMenu()
	}

}

func getMenuOption() string {
	var menuOption string
	fmt.Scan(&menuOption)
	return menuOption
}

// callClear is used to clear terminal screen
func cls() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// checkSites to open file, read and converto to slice of strings.
func checkSites() {
	fmt.Printf("Monitoring sites...")
	// open file
	file, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	// read file
	reader := bufio.NewReader(file)
	// convert to slice of strings
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			break
		}
		lines = append(lines, line)
	}
	// close file
	file.Close()

	// loop through slice of strings
	for {
		for _, line := range lines {
			// check if site is up
			if isUp(line) == false {
				registerLog(line)
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func isUp(url string) bool {
	_, err := http.Get("http://" + url)
	if err != nil {
		_, serr := http.Get("https://" + url)
		if serr != nil {
			return false
		}
	}
	return true
}

func registerLog(url string) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	currentTime := time.Now()
	file.WriteString("[DOWN] " + currentTime.Format("02/01/2006 15:04:05") + " - " + url + "\n")
	file.Close()
}

func showLogs() {
	cls()
	file, err := os.Open("log.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Println(line)
	}
	file.Close()
	fmt.Println()
	fmt.Printf("Press enter to return to menu")
	var answer string
	answer = getMenuOption()
	if answer != "1" {
		showMenu()
	}
}
