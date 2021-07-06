package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestGetBeeKey(t *testing.T) {

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://192.168.2.12:8010/getAddressName")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	fmt.Println(result.String())
	fmt.Println(result.String()[0:8])
}
