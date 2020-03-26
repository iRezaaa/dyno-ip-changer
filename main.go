package main

import (
	"net/http"
	"fmt"
	"time"
	"flag"
)

type Domain struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Token             string `json:"token"`
	State             string `json:"state"`
	IpV4Address       string `json:"ipv4Address"`
	LastUpdate        string `json:"updatedOn"`
	Group             string `json:"group"`
	TTL               int    `json:"ttl"`
	Ipv4              bool   `json:"Ipv4"`
	Ipv6              bool   `json:"ipv6"`
	Ipv4WildcardAlias bool   `json:"ipv4WildcardAlias"`
}

type UpdateDomainIPRequest struct {
	Name              string `json:"name"`
	Group             string `json:"group"`
	Ipv4Address       string `json:"ipv4Address"`
	Ipv6Address       string `json:"ipv6Address"`
	TTL               int    `json:"ttl"`
	Ipv4              bool   `json:"ipv4"`
	Ipv6              bool   `json:"ipv6"`
	Ipv4WildcardAlias bool   `json:"ipv4WildcardAlias"`
	Ipv6WildcardAlias bool   `json:"ipv6WildcardAlias"`
	AllowZoneTransfer bool   `json:"allowZoneTransfer"`
	DnsSec            bool   `json:"dnssec"`
}

type DnsRequestResponse struct {
	StatusCode int      `json:"statusCode"`
	Domains    []Domain `json:"domains"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Type       string `json:"type"`
	Message    string `json:"message"`
}

func main() {
	var checkDomainName string
	var apiKey string
	var port int
	var checkInterval int

	flag.StringVar(&checkDomainName, "domain", "", "check domain name to change ip address if blocked!")
	flag.StringVar(&apiKey, "apikey", "", "dynu api key, check : https://www.dynu.com/ControlPanel/APICredentials")
	flag.IntVar(&port, "port", 0, "your port to check")
	flag.IntVar(&checkInterval, "interval", 30, "check domain ip is blocked time in seconds")
	flag.Parse()

	if checkDomainName == "" {
		fmt.Println(fmt.Sprintf("you should enter -domain , a domain name to check blocked or not , if blocked change ip address"))
		return
	}

	if apiKey == "" {
		fmt.Println(fmt.Sprintf("you should enter -apikey , if not have a api key check : https://www.dynu.com/ControlPanel/APICredentials"))
		return
	}

	if port == 0 {
		fmt.Println(fmt.Sprintf("you should enter -port , your bot port"))
		return
	}

	client := &http.Client{}

	for {
		fmt.Println("Checking internet connection...")
		connected := checkInternetConnection()
		if !connected {
			fmt.Println("Check your internet connection, this application can't work without internet connection.")
			fmt.Println("------------------------------------------------")
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		}

		fmt.Println("Getting registered domain list...")
		domainList, err := getDomainList(client, apiKey)

		if err != nil {
			handleError(GetDomainList, err)
			fmt.Println("------------------------------------------------")
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		}

		if len(domainList.Domains) == 0 {
			fmt.Println("You are not registered a domain in dynu! , cannot continue with empty domain list in dynu panel!")
			return
		}

		var domain *Domain

		for _, o := range domainList.Domains {
			if o.Name == checkDomainName {
				domain = &o
				break
			}
		}

		if domain == nil {
			fmt.Println(fmt.Sprintf("%s domain name not exist in your dynu panel! exmaple domain name : %s", checkDomainName, "workplace.mywire.org"))
			return
		}

		if domain.IpV4Address == "" {
			fmt.Println(fmt.Sprintf("%s have no ipV4Address! please set a ip to domain in dynu panel", checkDomainName))
			return
		}

		fmt.Println(fmt.Sprintf("Checking %s:%d is blocked or not!", domain.IpV4Address, port))
		blocked := checkIPAddressIsBlocked(domain.IpV4Address, port)

		if blocked {
			if checkInternetConnection() {
				fmt.Println(fmt.Sprintf("this ip %s:%d is blocked!", domain.IpV4Address, port))
				newIP, err := getNewIPAddress()
				if err != nil {
					handleError(GetNewIPFromFile, err)
					return
				}

				err = updateDomainIP(client, apiKey, domain, newIP)

				if err != nil {
					handleError(UpdateDomainIP, err)
					return
				}

				fmt.Println(fmt.Sprintf("Update domain ip is successufly , %s domain's new ip is %s", checkDomainName, newIP))
			} else {
				fmt.Println("Check your internet connection, this application can't work without internet connection.")
				fmt.Println("------------------------------------------------")
				time.Sleep(time.Duration(checkInterval) * time.Second)
			}
		} else {
			fmt.Println(fmt.Sprintf("Everything is ok! your ip : %s not blocked", domain.IpV4Address))
		}

		fmt.Println("------------------------------------------------")
		time.Sleep(time.Duration(checkInterval) * time.Second)
	}
}
