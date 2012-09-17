package main

import (
	"github.com/yanunon/oss-go-api/oss"
	//"net/url"
	"fmt"
)

func main() {
	c := oss.NewClient("storage.aliyun.com", "ACSdztLFYwzIdZhu", "cs4UtVzxi4")
	//params := make(url.Values)
	//params.Set("Content-Md5", "c8fdb181845a4ca6b8fec737b3581d76")
	//params.Set("Content-Type", "text/html")
	//params.Set("Date", "Fri, 24 Feb 2012 02:58:28 GMT")
	//params.Set("X-OSS-Meta-Author", "foo@bar.com")
	//params.Set("X-OSS-Magic", "abracadabra")

	//c.SignParam("GET", "/", params)
	//c.DeleteBucket("yanunon2")
	err := c.PutBucket("yanunon2")
	if err != nil {
		fmt.Printf("PutBucket error:%s\n", err)
	}
	err = c.PutBucketACL("yanunon2", oss.ACL_PUBLIC_RW)
	if err != nil {
		fmt.Printf("PutBucketACL error:%s\n", err)
	}
	err = c.DeleteBucket("yanunon2")
	if err != nil {
		fmt.Printf("DeleteBucket error:%s\n", err)
	}
	c.GetService()
}
