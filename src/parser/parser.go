package parser

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

type Parser struct {
	logchan chan string
	wg      sync.WaitGroup
}

func New(logchan chan string) *Parser {
	return &Parser{
		logchan: logchan,
		wg:      sync.WaitGroup{},
	}
}

func (p *Parser) Parse(url string) chan Result {

	resultChan := make(chan Result, 1)
	p.wg.Add(1)

	go func() {

		res := newResult()

		nodes, err := p.parseContent(url, "")

		res.SetError(err)
		res.AddNodes(nodes...)

		resultChan <- res
		p.wg.Done()
	}()

	return resultChan
}

func (p *Parser) WaitUntilDone() {
	p.wg.Wait()
}

func (p *Parser) logMsg(msg string) {
	p.logchan <- fmt.Sprint("PARSER: ", msg)
}

func (p *Parser) findTableBody(doc *html.Node) *html.Node {

	var tableBodyNode *html.Node

	var crawler func(*html.Node)

	crawler = func(node *html.Node) {

		if node.Type == html.ElementNode && node.Data == "tbody" {
			tableBodyNode = node
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(doc)
	return tableBodyNode
}

func (p *Parser) parseContent(url string, namePrefix string) ([]*Node, error) {

	var result []*Node
	body, err := p.getContent(url)
	defer body.Close()

	if err != nil {
		return nil, err
	}

	doc, err := p.tryValidHTML(body)

	if err != nil {
		return nil, err
	}

	tableBodyHtmlNode := p.findTableBody(doc)

	result = p.parseTableBody(tableBodyHtmlNode)

	for _, node := range result {
		node.Name = fmt.Sprint(namePrefix, node.Name)
		if node.IsFolder() {
			node.Name += "/"
			newNodes, err := p.parseContent(node.Url, node.Name)
			if err == nil {
				node.AddNodes(newNodes)
			}
		}
	}

	return result, nil

}

func (p *Parser) parseTableBody(table *html.Node) []*Node {
	var result []*Node
	var crawler func(*html.Node)

	trs := 0

	crawler = func(node *html.Node) {

		if node.Type == html.ElementNode && node.Data == "tbody" {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				crawler(child)
			}
		}

		if node.Type == html.ElementNode && node.Data == "tr" {
			if trs == 0 {
				trs++
				return
			}

			crawler(node.FirstChild.NextSibling)
		}

		if node.Type == html.ElementNode && node.Data == "td" && node.FirstChild.Data == "a" && len(node.FirstChild.Attr) >= 1 {
			result = append(result, newNode(node.FirstChild.FirstChild.Data, node.FirstChild.Attr[0].Val))
		}

	}

	crawler(table)

	return result
}

func (p *Parser) tryValidHTML(r io.Reader) (*html.Node, error) {

	n, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (p *Parser) getContent(url string) (io.ReadCloser, error) {
	if url == "" {
		return nil, fmt.Errorf("empty url")
	}

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
