# aliyun_ddns
#### 获取绑定网卡的IP（包括IPV6）

## 使用方法
1. 下载相应平台指向文件
2. 编辑 **settings.json**
```json
{
  "AccessKeyId": "xxxxxx",
  "AccessKeySecret": "xxxxxx",
  "DomainName": "xxx.com",
  "SubDomains": [
    {
      "Type": "A",
      "RR": "sub1",
      "Interface": "eth0"
    },
    {
      "Type": "AAAA",
      "RR": "sub2",
      "Interface": "eth1"
    }
  ]
}
```
3. 执行 **./aliyun_ddns_xxx**

| 类别  | 说明  |
| ------------ | ------------ |
| AccessKeyId | 阿里云平台  |
| AccessKeySecret |  阿里云平台 |
| DomainName | 域名 |
| Type | 记录类型（A：ipv4，AAAA：ipv6） |
| RR | 主机记录  |
| Interface | 网卡名称（如 eth0，以太网） |
