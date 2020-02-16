package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
)

var (
	config       Config
	client       *alidns.Client
	interfaceMap map[string](*IP)
)

func main() {
	for _, subDomain := range config.SubDomains {
		var ip = interfaceMap[subDomain.Interface]
		if ip == nil {
			continue
		}
		if subDomain.Type == "A" && ip.Ipv4 != "" { // ipv4
			setDomainRecord(subDomain, ip.Ipv4)
		} else if subDomain.Type == "AAAA" && ip.Ipv6 != "" { // ipv6
			setDomainRecord(subDomain, ip.Ipv6)
		} else {
			continue
		}
	}
}

func setDomainRecord(subDomain SubDomain, ipAddr string) {
	var sub = fmt.Sprintf("%s.%s", subDomain.RR, config.DomainName)
	var records = getSubDomainRecords(sub)
	if len(records) > 0 { //存在解析记录
		for _, record := range records {
			request := alidns.CreateUpdateDomainRecordRequest()
			request.Scheme = "https"
			request.RecordId = record.RecordId
			request.RR = subDomain.RR
			request.Type = subDomain.Type
			request.Value = ipAddr
			_, err := client.UpdateDomainRecord(request)
			if err != nil {
				logInfo("更新域名解析失败：【%v】-> %v, %v, %v", sub, record.Value, record.RecordId, err)
				return
			}
			if record.Status != "ENABLE" {
				request := alidns.CreateSetDomainRecordStatusRequest()
				request.Scheme = "https"
				request.RecordId = record.RecordId
				request.Status = "ENABLE"
				client.SetDomainRecordStatus(request)
			}
		}
	} else {
		request := alidns.CreateAddDomainRecordRequest()
		request.Scheme = "https"
		request.DomainName = config.DomainName
		request.RR = subDomain.RR
		request.Type = subDomain.Type
		request.Value = ipAddr
		_, err := client.AddDomainRecord(request)
		if err != nil {
			logInfo("添加域名解析失败：【%v】-> %v", sub, err)
			return
		}
	}
	logInfo("域名解析：【%v】-> %v", sub, ipAddr)
}

// 初始化
func init() {
	initConfig()

	// 初始化客户端
	client, _ = alidns.NewClientWithAccessKey("cn-hangzhou", config.AccessKeyId, config.AccessKeySecret)

	// 获取网卡
	interfaceMap = make(map[string](*IP))
	getNetInterfaceInfo()
}

func getNetInterfaceInfo() {
	interfaces, err := net.Interfaces()
	if err != nil {
		logErr("获取网卡信息失败", err)
	}
	for _, inter := range interfaces {
		flags := inter.Flags.String()
		if strings.Contains(flags, "up") && strings.Contains(flags, "broadcast") {
			var ip IP
			adds, _ := inter.Addrs()
			for _, add := range adds {
				if ipnet, ok := add.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.IsGlobalUnicast() {
					var ipAddr = ipnet.IP.String()
					if ipAddr == "" {
						continue
					}
					if ipnet.IP.To4() != nil {
						ip.Ipv4 = ipAddr
					} else if ipnet.IP.To16() != nil && ip.Ipv6 == "" {
						ip.Ipv6 = ipAddr
					}
				}
			}
			interfaceMap[inter.Name] = &ip
		}
	}
}

func getSubDomainRecords(subDomain string) []alidns.Record {
	request := alidns.CreateDescribeSubDomainRecordsRequest()
	request.Scheme = "https"
	request.SubDomain = subDomain
	response, err := client.DescribeSubDomainRecords(request)
	if err != nil {
		logErr("获取子域名解析记录失败: %v", subDomain, err)
	}
	return response.DomainRecords.Record
}

func initConfig() {
	dir, _ := os.Getwd()
	f, err := os.Open(path.Join(dir, "settings.json"))
	if err != nil {
		logErr("无法打开配置文件", err)
	}
	defer f.Close()
	data, _ := ioutil.ReadAll(f)

	if err := json.Unmarshal(data, &config); err != nil {
		logErr("配置文件解析错误", err)
	}
	if config.AccessKeyId == "" || config.AccessKeySecret == "" || config.DomainName == "" || len(config.SubDomains) == 0 {
		logErr("请检查配置文件是否正确")
	}
}

func logInfo(format string, v ...interface{}) {
	log.Println(fmt.Sprintf(format, v...))
}
func logErr(format string, v ...interface{}) {
	log.Println(fmt.Sprintf(format, v...))
	os.Exit(-1)
}
