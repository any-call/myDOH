package myDOH

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/miekg/dns"
)

type DoHProvider struct {
	Name     string
	Endpoint string
}

var (
	AliDoH = DoHProvider{
		Name:     "aliyun",
		Endpoint: "https://dns.alidns.com/dns-query",
	}
	TencentDoH = DoHProvider{
		Name:     "tencent",
		Endpoint: "https://doh.pub/dns-query",
	}
	CloudflareDoH = DoHProvider{
		Name:     "cloudflare",
		Endpoint: "https://cloudflare-dns.com/dns-query",
	}
)

func buildDNSQueryWire(name string, qtype uint16) ([]byte, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), qtype)
	msg.RecursionDesired = true
	return msg.Pack()
}

// GET + dns=base64url(wireformat)

func queryTXTByDoHGET(ctx context.Context, client *http.Client, provider DoHProvider, domain string) ([]string, error) {
	wire, err := buildDNSQueryWire(domain, dns.TypeTXT)
	if err != nil {
		return nil, fmt.Errorf("[%s] pack dns query: %w", provider.Name, err)
	}

	enc := base64.RawURLEncoding.EncodeToString(wire)
	u, err := url.Parse(provider.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("[%s] parse endpoint: %w", provider.Name, err)
	}
	q := u.Query()
	q.Set("dns", enc)
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("[%s] build request: %w", provider.Name, err)
	}
	req.Header.Set("Accept", "application/dns-message")
	req.Header.Set("User-Agent", "go-doh-txt-test/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[%s] do request: %w", provider.Name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bs, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("[%s] http status=%d body=%s", provider.Name, resp.StatusCode, string(bs))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] read body: %w", provider.Name, err)
	}
	var reply dns.Msg
	if err := reply.Unpack(body); err != nil {
		return nil, fmt.Errorf("[%s] unpack dns response: %w", provider.Name, err)
	}
	if reply.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("[%s] dns rcode=%s", provider.Name, dns.RcodeToString[reply.Rcode])
	}
	var txts []string
	for _, ans := range reply.Answer {
		if rr, ok := ans.(*dns.TXT); ok {
			txts = append(txts, strings.Join(rr.Txt, ""))
		}
	}
	return txts, nil
}

// POST + wireformat body
func queryTXTByDoHPOST(ctx context.Context, client *http.Client, provider DoHProvider, domain string) ([]string, error) {
	wire, err := buildDNSQueryWire(domain, dns.TypeTXT)
	if err != nil {
		return nil, fmt.Errorf("[%s] pack dns query: %w", provider.Name, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, provider.Endpoint, bytes.NewReader(wire))
	if err != nil {
		return nil, fmt.Errorf("[%s] build request: %w", provider.Name, err)
	}
	req.Header.Set("Accept", "application/dns-message")
	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("User-Agent", "go-doh-txt-test/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[%s] do request: %w", provider.Name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bs, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("[%s] http status=%d body=%s", provider.Name, resp.StatusCode, string(bs))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] read body: %w", provider.Name, err)
	}
	var reply dns.Msg
	if err := reply.Unpack(body); err != nil {
		return nil, fmt.Errorf("[%s] unpack dns response: %w", provider.Name, err)
	}
	if reply.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("[%s] dns rcode=%s", provider.Name, dns.RcodeToString[reply.Rcode])
	}
	var txts []string
	for _, ans := range reply.Answer {
		if rr, ok := ans.(*dns.TXT); ok {
			txts = append(txts, strings.Join(rr.Txt, ""))
		}
	}
	return txts, nil

}
