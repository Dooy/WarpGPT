package test

import (
	"WarpGPT/pkg/env"
	"WarpGPT/pkg/logger"
	"WarpGPT/pkg/tools"
	"net/url"
	"os"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type RegClient struct {
	EmailAddress       string
	Password           string
	Proxy              string
	Client             tls_client.HttpClient
	UserAgent          string
	State              string
	CheckUrl           string
	PUID               string
	Verifier_code      string
	Verifier_challenge string
	//AuthResult         AuthResult
}

type Error tools.Error

func NewError(location string, statusCode int, details string, err error) *Error {
	return (*Error)(tools.NewError(location, statusCode, details, err))
}
func NewRegClient(checkUrl string, puid string) *RegClient {
	auth := &RegClient{
		// EmailAddress: emailAddress,
		// Password:     password,
		CheckUrl:  checkUrl,
		Proxy:     os.Getenv("proxy"),
		PUID:      puid,
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	}
	jar := tls_client.NewCookieJar()
	cookie := &http.Cookie{
		Name:   "_puid",
		Value:  puid,
		Path:   "/",
		Domain: ".openai.com",
	}
	urls, _ := url.Parse("https://openai.com")
	jar.SetCookies(urls, []*http.Cookie{cookie})
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(20),
		tls_client.WithClientProfile(profiles.Chrome_109),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
		tls_client.WithProxyUrl(env.E.Proxy),
	}
	auth.Client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	return auth
}
func (auth *RegClient) Start() *Error {
	logger.Log.Debug("begin")
	target := auth.CheckUrl //"https://" + env.E.OpenaiHost + "/api/auth/csrf"
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return NewError("begin", 0, "", err)
	}
	//req.Header.Set("Host", ""+env.E.OpenaiHost+"")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", auth.UserAgent)
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://"+env.E.OpenaiHost)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	resp, err := auth.Client.Do(req)
	if err != nil {
		return NewError("begin", 0, "", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 302 {
		redirectURL := resp.Header.Get("Location")
		println(redirectURL)
		auth.pOne(redirectURL)
	} else {
		return NewError("begin", 0, "Location error", nil)
	}
	// readBytes, err := io.ReadAll(resp.Body)
	// logger.Log.Debug("ggd", resp.StatusCode, string(readBytes))

	return nil
}
func (auth *RegClient) pOne(tarUrl string) *Error {
	logger.Log.Debug("pOne")

	req, err := http.NewRequest("GET", tarUrl, nil)
	if err != nil {
		return NewError("pOne", 0, "", err)
	}
	req.Header.Set("Host", "auth0.openai.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", auth.UserAgent)
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://"+env.E.OpenaiHost)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	resp, err := auth.Client.Do(req)
	if err != nil {
		return NewError("pOne", 0, "", err)
	}
	defer resp.Body.Close()
	// readBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return NewError("pOne", resp.StatusCode, "", err)
	// }
	// logger.Log.Debug("ggd", resp.StatusCode, string(readBytes))
	logger.Log.Debug("pOne", resp.StatusCode)

	return nil
}
