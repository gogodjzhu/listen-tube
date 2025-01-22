package ioutil

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ChanWriter chan string

func (cw ChanWriter) Write(p []byte) (n int, err error) {
    str := string(p)   // 将字节切片转换为字符串
    cw <- str          // 将字符串写入通道
    return len(p), nil // 返回写入的字节数
}

func DownloadFile(url string, output string) error {
    // create a custom HTTP client with TLS verification disabled
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}

    // download the file
    resp, err := client.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // create output path if not exists
    if err := os.MkdirAll(filepath.Dir(output), os.ModePerm); err != nil {
        return err
    }

    // create the file
    out, err := os.Create(output)
    if err != nil {
        return err
    }
    defer out.Close()

    // download the file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}