package main

import (
	"file_tree_downloader/downloader"
	"file_tree_downloader/parser"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	urlf := flag.String("url", "", "url to parse")

	flag.Parse()

	if *urlf == "" {
		fmt.Println("url can't be left empty")
		return
	}

	sigkillchan := make(chan os.Signal, 1)
	signal.Notify(sigkillchan, syscall.SIGINT, syscall.SIGTERM)

	logchan := make(chan string)
	donechan := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go monitor(logchan, donechan, wg)

	p := parser.New(logchan)

	res1 := p.Parse(*urlf)

	go readResult(res1, logchan, donechan)

	p.WaitUntilDone()

	wg.Wait()
}

func monitor(logchan chan string, donechan chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case msg := <-logchan:
			fmt.Println(msg)

		case <-donechan:
			wg.Done()
			return
		}
	}
}

func readResult(resChan chan parser.Result, logchan chan string, donechan chan struct{}) {
	res := <-resChan

	if res.Error != nil {
		logchan <- res.Error.Error()
		return
	}

	logchan <- "FOUND FILE TREE"
	logNodes(res.Nodes, logchan, "")

	cwd, err := os.Getwd()
	if err != nil {
		logchan <- res.Error.Error()
		return
	}

	if err = downloader.DownloadFiles(res.Nodes, cwd, logchan); err != nil {
		logchan <- err.Error()
	}

	close(donechan)
}

func logNodes(nodes []*parser.Node, logchan chan string, prefix string) {
	for _, n := range nodes {
		logchan <- fmt.Sprint(prefix, "name:", n.Name, "\t", "url:", n.Url)

		logNodes(n.Nodes, logchan, fmt.Sprint(prefix+" "))
	}
}
