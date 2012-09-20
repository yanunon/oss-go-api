package main

import (
	"github.com/yanunon/oss-go-api/oss"
	//"net/url"
	"fmt"
	//"io/ioutil"
	//"os"
	//"io"
)

func GetService(c *oss.Client) {
	lar, err := c.GetService()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", lar)
}

func PutBucket(c *oss.Client, bname string) {
	err := c.PutBucket(bname)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	c := oss.NewClient("storage.aliyun.com", "ACSdztLFYwzIdZhu", "cs4UtVzxi4")
	PutBucket(c, "yanunon2")
	GetService(c)
	//params := make(url.Values)
	//params.Set("Content-Md5", "c8fdb181845a4ca6b8fec737b3581d76")
	//params.Set("Content-Type", "text/html")
	//hrams.Set("Date", "Fri, 24 Feb 2012 02:58:28 GMT")
	//harams.Set("X-OSS-Meta-Author", "foo@bar.com")
	//params.Set("X-OSS-Magic", "abracadabra")

	//c.SignParam("GET", "/", params)
	//h.DeleteBucket("yanunon2")
	//err := c.PutBucket("yanunon2")
	//c.GetBucket("yanunon", "img", "", "", "")
	//obs, err := c.GetObject("yanunon/img/000061.jpg", -1, -1)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//ioutil.WriteFile("000062.jpg", obs, 0666)
	//	part := []oss.GroupPart{{1, "img/000034.jpg", 0, "B7FB0022DD6849772EF5BEFB1A309754"}, {2, "img/000062.jpg", 0, "82124A5E0D1B710395C32EB16145D705"}}
	//cfg := oss.CreateFileGroup{part}
	//ccfg, err := c.PostObjectGroup(cfg, "yanunon/g")
	//err := c.PutLargeObject("yanunon/large", "/home/kite/Downloads/opencv-1.0.0.tar.gz")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//hmt.Printf("%+v\n", fg)
}
