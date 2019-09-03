package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	username  = "admin"
	token     = "54eefa7dbd5bcf852c52fecd816f2a315c61832c"
	targetURL = "http://localhost:28080/file/fastupload"
	filehash  = "no_such_file_hash"
	filename  = "just_for_test"
)

func test_upload() {

	resp, err := http.PostForm(targetURL, url.Values{
		"username": {username},
		"token":    {token},
		"filehash": {filehash},
		"filename": {filename},
	})
	log.Printf("error: %+v\n", err)
	log.Printf("resp: %+v\n", resp)
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		log.Printf("parseBodyErr: %+v\n", err)
		if err == nil {
			log.Printf("parseBody: %+v\n", string(body))
		}
	}
}

func main() {
	test_upload()
}
