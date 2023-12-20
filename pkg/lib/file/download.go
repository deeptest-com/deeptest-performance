package _fileUtils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func Download(url string, dst string) (err error) {
	fmt.Printf("DownloadToFile From: %s to %s.\n", url, dst)

	err = MkDirIfNeeded(filepath.Dir(dst))
	if err != nil {
		fmt.Printf("file deny %v", err)
	}
	url = strings.ReplaceAll(url, "\\", "/")
	var data []byte
	data, err = HTTPDownload(url)
	if err == nil {
		fmt.Printf("file downloaded %s", url)

		RmDir(dst)

		err = WriteDownloadFile(dst, data)
		if err == nil {
			fmt.Printf("file %s saved to %s", url, dst)
		}
	}

	return
}

func HTTPDownload(uri string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(uri)
	//res, err := http.Get(uri)
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf(err.Error())
	}
	return d, err
}

func WriteDownloadFile(dst string, d []byte) error {
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		fmt.Printf(err.Error())
	}
	return err
}
