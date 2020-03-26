package main

import (
	"fmt"
	"net/http"
	"time"
	"net"
	"os"
	"bufio"
	"strings"
	"io/ioutil"
)

func handleError(requestType RequestType, err error) {
	switch requestType {
	case GetDomainList:
		fmt.Println(fmt.Sprintf("Error on getting domain list, this is details {%s}", err.Error()))
		return
	case CheckDomainIsBlocked:
		fmt.Println(fmt.Sprintf("Error while checking domain is blocked or not! this is details {%s}", err.Error()))
		return
	case UpdateDomainIP:
		fmt.Println(fmt.Sprintf("Error while update domain IP, this is details {%s}", err.Error()))
		return
	case GetNewIPFromFile:
		fmt.Println(fmt.Sprintf("Error while get new ip from iplist.txt, this is details {%s}", err.Error()))
		return
	}
}

func checkInternetConnection() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

func checkIPAddressIsBlocked(hostname string, port int) bool {
	seconds := 30
	timeOut := time.Duration(seconds) * time.Second

	connection, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", hostname, port), timeOut)
	defer func() {
		if connection != nil {
			connection.Close()
		}
	}()

	if err != nil {
		return true
	}

	return false
}

func getNewIPAddress() (string, error) {
	file, err := os.Open("iplist.txt")
	if err != nil {
		return "", &CustomError{
			ReadFile,
			"Failed to open iplist.txt file",
		}
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	if len(fileTextLines) == 0 {
		return "", &CustomError{
			ReadFile,
			"You must to enter at least one ip address! in iplist.txt",
		}
	}

	var chooseIP = ""
	var chooseIPIndex = -1
	for index, checkIP := range fileTextLines {
		if !strings.Contains(checkIP, "-used") {
			chooseIP = strings.Trim(checkIP, " ")
			chooseIPIndex = index
			break
		}
	}

	if chooseIP == "" {
		return "", &CustomError{
			ReadFile,
			"All ips are used , update iplist.txt file with new ips",
		}
	}

	fileTextLines[chooseIPIndex] = fmt.Sprintf("%s -used", chooseIP)
	output := strings.Join(fileTextLines, "\n")
	err = ioutil.WriteFile("iplist.txt", []byte(output), 0644)
	if err != nil {
		return "", &CustomError{
			WriteFile,
			"Error while update iplist.txt file!",
		}
	}

	return chooseIP, nil
}
