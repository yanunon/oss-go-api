package main

import (
	"github.com/yanunon/oss-go-api/oss"
	"net/url"
)

func main(){
	c := oss.NewClient("storage.aliyun.com", "44CF9590006BF252F707", "OtxrzxIsfpFjA7SwPzILwy8Bw21TLhquhboDYROV")
	params := make(url.Values)
	params.Set("Content-Md5", "c8fdb181845a4ca6b8fec737b3581d76")
	params.Set("Content-Type", "text/html")
	params.Set("Date", "Thu, 17 Nov 2005 18:49:58 GMT")
	params.Set("X-OSS-Meta-Author", "foo@bar.com")
	params.Set("X-OSS-Magic", "abracadabra")

	c.SignParam("PUT", "/quotes/nelson", params)
}
