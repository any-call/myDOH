package myDOH

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestDoh(t *testing.T) {
	domain := "_acme-challenge.proxy.stcoin.uk" // 改成你自己的 TXT 域名

	client := &http.Client{
		Timeout: 8 * time.Second,
	}
	providers := []DoHProvider{
		AliDoH,
		TencentDoH,
		CloudflareDoH,
	}
	ctx := context.Background()
	t.Log("=== GET /dns-query?dns=... ===")
	for _, p := range providers {
		txts, err := queryTXTByDoHGET(ctx, client, p, domain)
		if err != nil {
			t.Logf("[%s] err: %v\n", p.Name, err)
			continue
		}
		t.Logf("[%s] TXT: %#v\n", p.Name, txts)
	}
	t.Log("\n=== POST application/dns-message ===")
	for _, p := range providers {
		txts, err := queryTXTByDoHPOST(ctx, client, p, domain)
		if err != nil {
			t.Logf("[%s] err: %v\n", p.Name, err)
			continue
		}
		t.Logf("[%s] TXT: %#v\n", p.Name, txts)
	}

	t.Log("run ok ")
}

func TestQueryTxt(t *testing.T) {
	ret, err := QueryTxt(nil, time.Second*5, "sslink.api.prefix.stcoin.uk")
	if err != nil {
		t.Error(err)
		return
	}

	ret2, err := QueryTxtByProvider(nil, time.Second*5, "sslink.api.prefix.stcoin.uk", CloudflareDoH)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ret is :", ret, "ret 2 :", ret2)
}
