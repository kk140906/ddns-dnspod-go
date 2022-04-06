package utils

import (
	"crypto/tls"
	"encoding/json"
	"github.com/kk140906/ddns-dnspod-go/zap_wrapper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"regexp"
)

type DDNSConfig struct {
	LoginInfo struct {
		Id  string `json:"id"`
		Key string `json:"key"`
	} `json:"login_info"`
	Subdomains        []string `json:"subdomains"`
	FlushAfterMinutes int      `json:"flush_after_minutes"`
}

func ReadConfig(config *DDNSConfig) error {
	data, err := ioutil.ReadFile("configure.json")
	if err != nil {
		zap_wrapper.DefaultLogger.Error("read configure.json file failed")
		return err
	}
	return json.Unmarshal(data, config)
}

func GetIp() string {
	// https://ip.cn/api/index?ip=%22%22&type=0
	response, err := httpGetRequest("https://api.ip.sb/ip")
	// zap_wrapper.DefaultLogger.Debug(fmt.Sprintln(response))
	if err != nil {
		zap_wrapper.DefaultLogger.Error("get ip address failed", zap.String("error", err.Error()))
		return ""
	}

	data := make([]byte, 256)
	pos, err := response.Body.Read(data)
	// zap_wrapper.DefaultLogger.Debug(string(data))
	if err != nil {
		zap_wrapper.DefaultLogger.Error("get ip address failed", zap.String("error", err.Error()))
		return ""
	}
	return string(data[:pos-1])
}

func httpGetRequest(url string) (*http.Response, error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	}}
	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,"+
		"*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("sec-ch-ua-platform", "Windows")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36")
	return client.Do(req)
}

func ValidIp(ip string) bool {
	regex := regexp.MustCompile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
	return regex.MatchString(ip)
}
