package httpClient

import (
	"cloudMonitor/package/utils"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func ScoreUploader(dataServer string, waitGroup *sync.WaitGroup, scoreCh chan float32) {
	defer waitGroup.Done()

	for {
		score := <-scoreCh
		s := strconv.FormatFloat(float64(score), 'f', 2, 32)
		uri := utils.ConcatStrings(dataServer, s)
		err := sendScore(uri)
		if nil != err {
			fmt.Println("Failed to upload the score to the server. Reason:", err)
		}
	}
}

func sendScore(uri string) error {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	response, err := client.Get(uri)
	if nil != err {
		return err
	}

	defer response.Body.Close()
	if 200 != response.StatusCode {
		body, err := io.ReadAll(response.Body)
		if nil != err {
			return err
		}
		return errors.New(string(body))
	}

	return nil
}
