package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const defaultRequestHost = `https://www.ipaddress.com/ip-lookup`

func GetWebsiteHost(domain string, ch chan<- *HostChan) {
	// 创建表单数据
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("host", domain)
	_ = writer.Close()

	// 创建HTTP请求
	req, err := http.NewRequest(http.MethodPost, defaultRequestHost, &body)
	if err != nil {
		fmt.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送HTTP请求
	client := http.Client{}
	resp, err := client.Do(req)
	defer func(response *http.Response) {
		if response != nil {
			_ = response.Body.Close()
		}
	}(resp)
	if err != nil {
		fmt.Printf("Failed to send request: %v", err)
		return
	}

	// 读取响应内容
	respHtml, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v", err)
		ch <- &HostChan{Domain: domain, Err: err}
		return
	}

	htmlStr := string(respHtml)
	ipAddr, singleErr := handleSingleResult(htmlStr, domain)
	if singleErr == nil {
		ch <- &HostChan{Domain: domain, Ip: ipAddr}
		return
	}

	ipAddrList, err := handleMultiResult(htmlStr)
	if err != nil {
		ch <- &HostChan{Domain: domain, Err: err}
		return
	}

	if len(ipAddrList) == 0 {
		return
	}

	ch <- &HostChan{Domain: domain, Ip: ipAddrList[0]}
	return
}

func handleSingleResult(htmlStr string, domain string) (ipAddr string, err error) {
	resultStart := "IP Address Lookup :"
	resultEnd := fmt.Sprintf("(%s)", domain)

	startIdx := strings.Index(htmlStr, resultStart)
	endIdx := strings.Index(htmlStr, resultEnd)
	if startIdx == -1 || endIdx == -1 {
		err = errors.New("invalid single result")
		return
	}

	ipAddr = htmlStr[startIdx+len(resultStart) : endIdx]
	ipAddr = strings.TrimSpace(ipAddr)
	return
}

func handleMultiResult(htmlStr string) (ipAddrList []string, err error) {
	defaultTipSentence := "Please select one from the list below for detailed lookup results."
	if !strings.Contains(htmlStr, defaultTipSentence) {
		err = errors.New("invalid multi result")
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		err = errors.New("invalid multi result: failed to parse html")
		return
	}

	doc.Find("form").Each(func(i int, selection *goquery.Selection) {
		val, exists := selection.Attr("action")
		if !exists || val != defaultRequestHost {
			return
		}

		val, exists = selection.Attr("class")
		if !exists || val != "form" {
			return
		}

		val, exists = selection.Attr("method")
		if !exists || strings.ToLower(val) != strings.ToLower(http.MethodPost) {
			return
		}

		selection.Find("a").Each(func(i int, selection *goquery.Selection) {
			ipAddrList = append(ipAddrList, selection.Text())
		})
	})

	return
}
