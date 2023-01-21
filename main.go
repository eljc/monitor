package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoring = 3
const delay = 5

func main() {
	for {
		showMenu()

		option := readCommand()

		switch option {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("Show Log")
			printLogs()
		case 0:
			fmt.Println("Exit")
			os.Exit(0)
		default:
			fmt.Println("Unidentified option.")
			os.Exit(-1)
		}

	}
}

func showMenu() {
	messageWelcome := "Welcome to Monitoring website"
	fmt.Println(messageWelcome + "\n")
	fmt.Println("Choose an option:")
	fmt.Println("1 - Start Monitoring")
	fmt.Println("2 - Show Logs")
	fmt.Println("0 - Exit")
}

func readCommand() int {
	var optionSelected int
	fmt.Scan(&optionSelected)
	fmt.Println("Option selected:", optionSelected)

	return optionSelected
}

func startMonitoring() {
	fmt.Println("Scan sites...")

	sites := readSitesFromFile()

	for i := 0; i < monitoring; i++ {
		for i, site := range sites {
			fmt.Println("Try access: ", i, "->", site)
			testSite(site)
		}
	}
	time.Sleep((delay * time.Second))
	fmt.Println("")

}

func readSitesFromFile() []string {
	var sites []string
	contentFile, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	leitor := bufio.NewReader(contentFile)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			break
		}

	}

	contentFile.Close()
	return sites
}

func testSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Print("Something went wrong.")
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "success!")
		saveLog(site, true)
	} else {
		fmt.Println("Site:", site, "error. Status Code:", resp.StatusCode)
		saveLog(site, false)
	}
}

func saveLog(site string, status bool) {

	fileContent, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	fileContent.WriteString(time.Now().Format("01-02-2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	fileContent.Close()
}

func printLogs() {

	fileContent, err := os.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(fileContent))

}
