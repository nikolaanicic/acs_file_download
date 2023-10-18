package downloader

import (
	"file_tree_downloader/parser"
	"file_tree_downloader/writer"
	"fmt"
	"io"
	"net/http"
	"path"
)

func logMsgFormatted(logchan chan string, msg string) {
	logchan <- fmt.Sprint("DOWNLOADER: ", msg)
}

func logMsg(logchan chan string, msg string) {
	logchan <- msg
}

func DownloadFiles(nodes []*parser.Node, basePath string, logchan chan string) error {

	logMsgFormatted(logchan, "starting to download")
	for _, n := range nodes {
		if n.IsFolder() {
			writer.CreateDirectory(path.Join(basePath, n.Name))
			DownloadFiles(n.Nodes, basePath, logchan)
		} else {

			body, err := getFile(n.Url)
			if err != nil {
				logMsg(logchan, fmt.Sprint(n.Name, "...Fail"))
				return err
			}

			content, err := io.ReadAll(body)
			if err != nil {
				logMsg(logchan, fmt.Sprint(n.Name, "...Fail"))
				return err
			}

			writer.WriteFile(n.Name, content)
			logMsg(logchan, fmt.Sprint(n.Name, "...Done"))
		}

	}

	return nil
}

func getFile(url string) (io.ReadCloser, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
