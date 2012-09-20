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

func HeadObject(c *oss.Client, opath string) {
	header, err := c.HeadObject(opath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", header)
}

func DeleteMultipleObject(c *oss.Client, bname string, keys []string) {
	err := c.DeleteMultipleObject(bname, keys)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Delete multiple object ok")
}

func PostObjectGroup(c *oss.Client, cfg oss.CreateFileGroup, gpath string) {
	cofg, err := c.PostObjectGroup(cfg, gpath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", cofg)
}

func GetObjectGroupIndex(c *oss.Client, gpath string) {
	fg, err := c.GetObjectGroupIndex(gpath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", fg)
}

func AbortMultipartUpload(c *oss.Client, opath, uploadId string) {
	err := c.AbortMultipartUpload(opath, uploadId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("abort multipart upload ok")
}

func PutLargeObject(c *oss.Client, opath, fpath string) {
	err := c.PutLargeObject(opath, fpath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("put large object ok")
}

func ListMultipartUpload(c *oss.Client, bname string) {
	lmur, err := c.ListMultipartUpload(bname, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", lmur)
}

func ListParts(c *oss.Client, opath, uploadId string) {
	lpr, err := c.ListParts(opath, uploadId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", lpr)
}

func main() {
	c := oss.NewClient("storage.aliyun.com", "ACSdztLFYwzIdZhu", "cs4UtVzxi4", 10)
	PutLargeObject(c, "yanunon/android-ndk-r7b-linux-x86.2.tar.bz2", "/home/kite/Downloads/android-ndk-r7b-linux-x86.tar.bz2")
	//ListMultipartUpload(c, "yanunon")
}
