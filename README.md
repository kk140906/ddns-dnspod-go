# DDNS DNSPod Go

&emsp;&emsp; 基于 Go 实现的 DNSPod DDNS 脚本，采用腾讯云 API3.0 版本的SDK，功能单一，更新只针对二级域名的 A 记录。

## 配置文件

```
{
  "login_info": {
    "id": "", // API3.0 Secret Id
    "key": "" // API3.0 Secret Key
  },
  "subdomains": [
    "demo.kk9009.com" // 只支持二级域名的更新
  ],
  "flush_after_minutes": 10 // 刷新的间隔时间， 单位是分钟
}
```

DNSPOD API 3.0 的秘钥创建参考腾讯官方首页 [https://cloud.tencent.com/product/api](https://cloud.tencent.com/product/api)