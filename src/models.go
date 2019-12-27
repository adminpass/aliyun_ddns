package main

type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	DomainName      string
	SubDomains      []SubDomain
}

type SubDomain struct {
	Type      string // A:ipv4  AAAA:ipv6
	RR        string
	Interface string
}

type IP struct {
	Ipv4 string
	Ipv6 string
}
