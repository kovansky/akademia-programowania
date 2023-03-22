package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/kovansky/wp-reddit/fetcher"
	"io"
	"log"
	"os"
	"strings"
)

const (
	VERSION = "0.0.1"
	AUTHOR  = "Stanisław Kowański <s25477@pjwstk.edu.pl>"
)

func main() {
	defer details()

	logger := log.New(os.Stdout, "", log.Ltime)
	ctx := context.WithValue(context.Background(), fetcher.CtxLoggerKey, logger)

	print(color.New(color.Bold).Sprintf("- Reddit Fetcher v%s -\n", VERSION))

	var (
		subs string
		out  int

		output io.Writer
	)

	print("Wypisz subreddity, z których chcesz pobrać informacje - oddzielane przecinkami (np. golang,webdev):\n-> ")

	reader := bufio.NewReader(os.Stdin)
	subs, err := reader.ReadString('\n')
	if err != nil {
		logger.Printf("%s\n%#v\n", color.HiRedString("Error occurred while parsing user input:"), err)
	}

readOutput:
	print("Wybierz output (1 - plik, 2 - konsola):\n-> ")
	_, err = fmt.Scan(&out)

	switch out {
	case 1:
		output, err = os.Create("output")
		if err != nil {
			logger.Printf("%s\n%#v\n", color.HiRedString("Error occurred while creating file output:"), err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(output.(*os.File))
	case 2:
		output = os.Stdout
	default:
		print(color.RedString("%d nie jest prawidłowym wyborem.\n", out))
		goto readOutput
	}

	reddit := fetcher.RedditFetcherImpl{Subreddits: strings.Split(strings.TrimSpace(subs), ",")}
	err = reddit.Fetch(ctx)
	if err != nil {
		logger.Printf("%s\n%s\n", color.HiRedString("Error(s) occurred while fetching from reddit:"), err.Error())
	}
	err = reddit.Save(output)
	if err != nil {
		logger.Printf("%s\n%#v\n", color.HiRedString("Error occurred while writing output:"), err)
	}
}

func details() {
	print(color.HiBlackString("--> Przygotowano na Akademię Programowania WP przez %s", AUTHOR))
}
