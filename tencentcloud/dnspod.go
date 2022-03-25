package tencentcloud

import (
	"errors"
	"github.com/kk140906/ddns-dnspod-go/utils"
	"github.com/kk140906/ddns-dnspod-go/zap_wrapper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentclouderrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"go.uber.org/zap"
	"strings"
)

type Client struct {
	*dnspod.Client
	Domain map[string][]string
}

type RecordResponse struct {
	RecordId uint64
	TTL      uint64
	Value    string
}

const dnspodEndPoint = "dnspod.tencentcloudapi.com"

func NewClient() *Client {
	config := utils.DDNSConfig{}
	_ = utils.ReadConfig(&config)
	domains := make(map[string][]string, 1)
	for _, subdomain := range config.Subdomains {
		info := strings.Split(subdomain, ".")
		if len(info) == 3 {
			domain := strings.Join(info[1:], ".")
			domains[domain] = append(domains[domain], info[0])
		}
	}
	credential := common.NewCredential(config.LoginInfo.Id, config.LoginInfo.Key)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = dnspodEndPoint
	client, _ := dnspod.NewClient(credential, "", cpf)
	return &Client{client, domains}
}

func (client *Client) GetDomainRecordList(domain string) (map[string]RecordResponse, error) {
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(domain)
	response, err := client.DescribeRecordList(request)
	if _, ok := err.(*tencentclouderrors.TencentCloudSDKError); ok {
		zap_wrapper.DefaultLogger.Info("An API error has returned", zap.String("error", err.Error()))
		return nil, errors.New("api call error")
	}
	if err != nil {
		return nil, err
	}
	recordList := make(map[string]RecordResponse)

	for _, record := range response.Response.RecordList {
		name := *record.Name
		id := *record.RecordId
		value := *record.Value
		ttl := *record.TTL

		recordList[name] = RecordResponse{
			RecordId: id,
			TTL:      ttl,
			Value:    value,
		}
	}
	return recordList, nil
}

func (client *Client) UpdateDomainRecord(domain, subdomain, value string,
	record RecordResponse) (*dnspod.ModifyRecordResponse, error) {
	request := dnspod.NewModifyRecordRequest()
	request.Domain = common.StringPtr(domain)
	request.RecordType = common.StringPtr("A")
	request.SubDomain = common.StringPtr(subdomain)
	request.RecordId = common.Uint64Ptr(record.RecordId)
	request.TTL = common.Uint64Ptr(record.TTL)
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(value)
	response, err := client.ModifyRecord(request)
	if _, ok := err.(*tencentclouderrors.TencentCloudSDKError); ok {
		zap_wrapper.DefaultLogger.Error("An API error has returned", zap.String("error", err.Error()))
		return nil, errors.New("api call error")
	}
	if err != nil {
		zap_wrapper.DefaultLogger.Error("modify record request failed", zap.String("error", err.Error()))
		return nil, err
	}
	return response, nil
}
