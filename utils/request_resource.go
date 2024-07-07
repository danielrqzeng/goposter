package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"
	//project
)

var (
	httpPool *http.Client
)

// createHTTPClient for connection re-use
func createHTTPClientPool() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 5,
		},
		Timeout: time.Duration(2) * time.Second,
	}
	return client
}

func getHttpRequestPool() *http.Client {
	return httpPool
}

func DeferWhenOnRequestDone() func(retCode *int, msg, url, rBody, wBody *string) {
	begin := time.Now()
	return func(retCode *int, msg, url, rBody, wBody *string) {
		elapsed := time.Since(begin).String()
		_ = elapsed
		//logger.ApiLogger.Info(*url + "|" + elapsed + "|" + utils.Num2Str(*retCode) + "|" + *msg + "|" + *rBody + "|" + *wBody)
	}
}

//RequestResource 请求资源（一般来说是指获取图片资源）
func RequestResource(resUrl string) (body io.ReadCloser, contentType string, err error) {
	statusCode := -1
	errMsg := ""
	urlPath := resUrl
	rBody := ""
	wBody := ""
	defer DeferWhenOnRequestDone()(&statusCode, &errMsg, &urlPath, &rBody, &wBody)

	if resUrl == "" {
		err = fmt.Errorf("resUrl is none")
		errMsg = err.Error()
		return
	}

	//下载资源回来
	pool := getHttpRequestPool()
	objRsp, err := pool.Get(resUrl)
	if err != nil {
		errMsg = err.Error()
		return
	}
	//defer objRsp.Body.Close()
	statusCode = objRsp.StatusCode
	if statusCode != http.StatusOK {
		err = fmt.Errorf("not right to get,statusCode=%d", statusCode)
		return
	}

	contentType = objRsp.Header.Get("content-type") //Content-Type
	//成功
	body = objRsp.Body
	statusCode = 200
	return
}

func init() {
	httpPool = createHTTPClientPool()

	time.AfterFunc(time.Duration(1)*time.Second, func() {
	})
}
