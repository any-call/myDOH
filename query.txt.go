package myDOH

import (
	"context"
	"net/http"
	"time"
)

func QueryTxt(ctx context.Context, timeout time.Duration, domain string) ([]string, error) {
	//优先查阿里，再查腾讯，//最后醒cloud flare
	if ctx == nil {
		ctx = context.Background()
	}
	client := &http.Client{
		Timeout: timeout,
	}

	ret, err := queryTXTByDoHGET(ctx, client, AliDoH, domain)
	if err != nil {
		if ret, err = queryTXTByDoHGET(ctx, client, TencentDoH, domain); err != nil {
			return queryTXTByDoHGET(ctx, client, CloudflareDoH, domain)
		} else {
			return ret, nil
		}
	} else {
		return ret, nil
	}

}

func QueryTxtByProvider(ctx context.Context, timeout time.Duration, domain string, provider DoHProvider) ([]string, error) {
	//优先查阿里，再查腾讯，//最后醒cloud flare
	if ctx == nil {
		ctx = context.Background()
	}
	client := &http.Client{
		Timeout: timeout,
	}

	return queryTXTByDoHGET(ctx, client, provider, domain)
}
