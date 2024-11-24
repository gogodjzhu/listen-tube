package http

import (
    "io"
    "math/rand"
    "net/http"
    net_url "net/url"
)

// HttpGet performs an HTTP GET request with optional proxies.
func HttpGet(proxies []string, url string) (string, error) {
    client := &http.Client{}
    if len(proxies) > 0 {
        randomIdx := rand.Intn(len(proxies))
        proxyURL, _ := net_url.Parse(proxies[randomIdx])
        client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
    }

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
    req.Header.Set("Accept-Language", "en")
    req.AddCookie(&http.Cookie{Name: "CONSENT", Value: "YES+cb", Domain: ".youtube.com"})

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(body), nil
}