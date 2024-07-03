package core

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultRequestHost = `https://www.ipaddress.com/ip-lookup`

var invalidContentErr = errors.New("invalid response content, please check website content")

func GetWebsiteHost(domain string, ch chan<- *HostChan) {
	client := http.Client{Timeout: 5 * time.Second}
	newRequest, err := http.NewRequest("POST", defaultRequestHost, http.NoBody)
	if err != nil {
		ch <- &HostChan{Domain: domain, Err: err}
		return
	}

	newRequest.PostForm = url.Values{"host": {domain}}
	newRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")

	resp, err := client.Do(newRequest)
	if err != nil {
		ch <- &HostChan{Domain: domain, Err: err}
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respHtml, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- &HostChan{Domain: domain, Err: err}
		return
	}

	htmlStr := string(respHtml)
	resultStart := "IP Lookup :"
	resultEnd := fmt.Sprintf("(%s)", domain)

	startIdx := strings.Index(htmlStr, resultStart)
	endIdx := strings.Index(htmlStr, resultEnd)
	if startIdx == -1 || endIdx == -1 {
		fmt.Println("invalid response content, please check website content or contact author")
		ch <- &HostChan{Domain: domain, Err: invalidContentErr}
		return
	}

	host := htmlStr[startIdx+len(resultStart) : endIdx]
	host = strings.TrimSpace(host)
	ch <- &HostChan{Domain: domain, Ip: host}
	return
}

func parseHtml(rawHtml string) string {
	// 使用 goquery 加载 HTML 文档
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHtml))
	if err != nil {
		log.Fatal(err)
	}

	parsedText := ""
	// 提取所有标签内容并判断是否为文本节点
	findTextNodes(doc.Find("body"), &parsedText)
	return parsedText
}

func findTextNodes(selection *goquery.Selection, parsedText *string) {
	selection.Contents().Each(func(_ int, s *goquery.Selection) {
		if s.Nodes != nil {
			if s.Children() == nil {
				for i := 0; i < len(s.Nodes); i++ {
					if s.Nodes[i].FirstChild == nil {
						*parsedText += s.Nodes[i].Data
					}
				}
			} else {
				// 递归遍历子节点
				findTextNodes(s, parsedText)
			}
		}
	})
}
