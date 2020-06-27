package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var options [3]string

const MONITOR_TIMES = 3
const DELAY = 1

func main() {
	options[0] = "Sair do programa"
	options[1] = "Iniciar monitoramento"
	options[2] = "Exibir logs"

	fmt.Println("Seja bem vindo ao sistema de monitoramento de sites")

	for {
		fmt.Println("Escolha uma das opções para prosseguir")
		for index, value := range options {
			fmt.Println(index, "-", value)
		}

		choosedOption := 0

		fmt.Scan(&choosedOption)
		fmt.Println("A opção escolhida foi:", sanitizeOptions(choosedOption))

		run(choosedOption)
	}
}

func sanitizeOptions(option int) string {
	return fmt.Sprintf("%d - %s", option, options[option])
}

func run(option int) {
	switch option {
	case 1:
		monitor()
		break
	case 2:
		fmt.Println("Exibindo logs...")
		printLogs()
	case 0:
		fmt.Println("Saindo do programa")
		os.Exit(0)
	default:
		fmt.Println("Comando não identificado")
		os.Exit(-1)
	}
}

func monitor() {
	fmt.Println("Monitorando...")
	urls := loadUrlsFromFile("urls.txt")

	for i := 0; i < MONITOR_TIMES; i++ {
		for _, url := range urls {
			monitorUrl(url)
		}
		time.Sleep(DELAY * time.Second)
	}
}

func monitorUrl(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == 200 {
		fmt.Println(fmt.Sprintf("[SUCCESS]Site %s carregado com sucesso (%d) \n", url, res.StatusCode))
		saveLog(url, true)
	} else {
		fmt.Println(fmt.Sprintf("[ERROR]Site %s apresentou erro (%d) \n", url, res.StatusCode))
		saveLog(url, false)
	}
}

func loadUrlsFromFile(filepath string) []string {
	var urls []string

	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err == io.EOF {
			break
		}

		urls = append(urls, line)
	}

	file.Close()
	return urls
}

func saveLog(url string, status bool) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(fmt.Sprintf("[%s][%s] %s \n", sanitizeStatus(status), time.Now().Format("01/02/2006 15:04:05"), url))
	file.Close()
}

func sanitizeStatus(status bool) string {
	if status {
		return "ON"
	}
	return "OFF"
}

func printLogs() {
	file, err := ioutil.ReadFile("logs.txt")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(file))
}
