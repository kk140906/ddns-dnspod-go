package utils

import (
	"encoding/json"
	"github.com/kk140906/ddns-dnspod-go/zap_wrapper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
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
	response, err := http.Get("https://api.ip.sb/ip")
	if err != nil {
		zap_wrapper.DefaultLogger.Error("get ip address failed", zap.String("error", err.Error()))
		return ""
	}

	data := make([]byte, 256)
	pos, err := response.Body.Read(data)
	if err != nil {
		zap_wrapper.DefaultLogger.Error("get ip address failed", zap.String("error", err.Error()))
		return ""
	}
	return string(data[:pos-1])
}
