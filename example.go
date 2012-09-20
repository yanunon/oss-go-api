/*
 * Copyright (c) 2012, Yang Junyong <yanunon@gmail.com>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the Google Inc. nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

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
	c := oss.NewClient("storage.aliyun.com", "ACCESS_ID", "ACCESS_KEY", 10)
}
