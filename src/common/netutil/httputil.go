package netutil

import (
	"errors"
	"net/http"
	"net"
	"time"
	"strings"
	"fmt"
	"io/ioutil"
)


var ActionHttpHandleErr error = errors.New("http server handle msg err")
var ActionMsgIsDiscardErr error = errors.New("msg is filter and discard err")

type RespMsg struct {
	ErrNo  int    `json:"errno"`
	ErrMsg string `json:"errmsg"`
}

type HttpErrInfo struct {
	HttpErr string
	ErrCode int
	CostMs  int64
}

const (
	RESP_CODE_OK             = 0
	HTTP_CONNECT_TIMEOUT_MS  = 200
	HTTP_RESPONSE_TIMEOUT_MS = 3000
	HTTP_REQU_RETRY_TIMES    = 3
)

func DoHttpPost(url string, data string, connTimeoutMs int, serveTimeoutMs int, httpErrInfo *HttpErrInfo) ([]byte, error) {
	beginTime := time.Now().UnixNano() / int64(time.Millisecond)
	defer func() {
		endTime := time.Now().UnixNano() / int64(time.Millisecond)
		httpErrInfo.CostMs = endTime - beginTime
	}()
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				return c, nil
			},
		},
	}

	body := strings.NewReader(data)
	reqest, _ := http.NewRequest("POST", url, body)
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(reqest)
	if err != nil {
		httpErrInfo.HttpErr = "post_err"
		httpErrInfo.ErrCode = 400
		err = errors.New(fmt.Sprintf("http failed, POST url:%s, reason:%s", url, err.Error()))
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		httpErrInfo.HttpErr = "response_err"
		httpErrInfo.ErrCode = response.StatusCode
		err = errors.New(fmt.Sprintf("http status code errorcode, POST url:%s, code:%d", url, response.StatusCode))
		return nil, err
	}

	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		httpErrInfo.HttpErr = "response_body_err"
		httpErrInfo.ErrCode = 200
		err = errors.New(fmt.Sprintf("cannot read http response, POST url:%s, reason:%s", url, err.Error()))
		return nil, err
	}
	httpErrInfo.HttpErr = "ok"
	httpErrInfo.ErrCode = 200
	return res_body, nil
}

func DoHttpGet(url string, connTimeoutMs int, serveTimeoutMs int) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				return c, nil
			},
		},
	}

	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Accept-Encoding", "text/plain; charset=utf-8")
	reqest.Header.Add("User-Agent", "")
	reqest.Header.Add("Connection", "Close")
	response, err := client.Do(reqest)
	if err != nil {
		err = errors.New(fmt.Sprintf("http failed, GET url:%s, reason:%s", url, err.Error()))
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("http status code errorcode, GET url:%s, code:%d", url, response.StatusCode))
		return nil, err
	}

	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("cannot read http response, GET url:%s, reason:%s", url, err.Error()))
		return nil, err
	}
	return res_body, nil
}
