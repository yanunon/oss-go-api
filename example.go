package main

import (
	"fmt"
	"github.com/yanunon/oss-go-api/oss"
	"log"
	"os"
)

func GetService(c *oss.Client) {
	lar, err := c.GetService()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", lar)
}

func PutBucket(c *oss.Client, bname string) {
	err := c.PutBucket(bname)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Put bucket ok")
}

func GetBucket(c *oss.Client, bname string) {
	lbr, err := c.GetBucket(bname, "", "", "", "")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", lbr)
}

func PutBucketACL(c *oss.Client, bname string, acl string) {
	err := c.PutBucketACL(bname, acl)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Put bucket acl ok")
}

func GetBucketACL(c *oss.Client, bname string) {
	acl, err := c.GetBucketACL(bname)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", acl)
}

func CopyObject(c *oss.Client, dst, src string) {
	err := c.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Copy object ok")
}

func DeleteObject(c *oss.Client, opath string) {
	err := c.DeleteObject(opath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Delete object ok")
}

func GetObject(c *oss.Client, fpath, opath string) {
	bytes, err := c.GetObject(opath, -1, -1)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Create(fpath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	file.Write(bytes)
	fmt.Println("Get object ok")
}

func PutObject(c *oss.Client, opath, fpath string) {
	err := c.PutObject(opath, fpath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Put object ok")
}

func main() {
	c := oss.NewClient("storage.aliyun.com", "ACSdztLFYwzIdZhu", "cs4UtVzxi4")
	PutObject(c, "/yanunon/img/sk3.jpg", "/home/kite/Dropbox/sk3.jpg")
	GetObject(c, "sk3.jpg", "/yanunon/img/sk3.jpg")
	//	part := []oss.GroupPart{{1, "img/000034.jpg", 0, "B7FB0022DD6849772EF5BEFB1A309754"}, {2, "img/000062.jpg", 0, "82124A5E0D1B710395C32EB16145D705"}}
	//cfg := oss.CreateFileGroup{part}
	//ccfg, err := c.PostObjectGroup(cfg, "yanunon/g")
	//err := c.PutLargeObject("yanunon/large", "/home/kite/Downloads/opencv-1.0.0.tar.gz")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//hmt.Printf("%+v\n", fg)
}
