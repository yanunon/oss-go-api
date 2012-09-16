package oss

import (
//	"encoding/xml"
	"encoding/base64"
	"fmt"
//	"net/url"
	"net/http"
	"strings"
	"sort"
	"bytes"
	"crypto/sha1"
	"crypto/hmac"
	"io"
	"io/ioutil"
	"hash"
	"time"
	"log"
)

type Client struct {
	AccessID string
	AccessKey string
	Host string
}

type Bucket struct {
	Name string
	CreationDate string
}

type ValSorter struct {
	Keys []string
	Vals []string
}


func NewClient(host, accessId, accessKey string) (*Client) {
	client := Client{Host:host, AccessID:accessId, AccessKey:accessKey}
	return &client
}

func (c *Client) SignParam(method, resource string, params http.Header) {
	//format x-oss-
	tmpParams := make(map[string]string)

	for k, v := range params {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			tmpParams[strings.ToLower(k)] = v[0]
		}
	}
	//sort
	valSorter := NewValSorter(tmpParams)
	valSorter.Sort()

	canonicalizedOSSHeaders := ""
	for i := range(valSorter.Keys) {
		canonicalizedOSSHeaders += valSorter.Keys[i] + ":" + valSorter.Vals[i] + "\n"
	}

	date := params.Get("Date")
	contentType := params.Get("Content-Type")
	contentMd5 := params.Get("Content-Md5")

	signStr := method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders  + resource
	h := hmac.New(func() hash.Hash {return sha1.New()}, []byte(c.AccessKey)) //sha1.New()
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	authorizationStr := "OSS " + c.AccessID + ":" + signedStr
	fmt.Println(authorizationStr)
	params.Set("Authorization", authorizationStr)
}

func (c *Client) GetService() {
	httpClient := http.DefaultClient
	reqUrl := "http://" + c.Host + "/"
	req, _ := http.NewRequest("GET", reqUrl, nil)
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	//date = "Sun, 16 Sep 2012 12:18:36 GMT"
	req.Header.Set("Date", date)
	req.Header.Set("Host", c.Host)
	//req.Header.Set("Authorization", c.AccessID)
	c.SignParam("GET", "/", req.Header)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	respbytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respbytes))
	//fmt.Println(date)

}

func NewValSorter(m map[string]string) *ValSorter {
	vs := &ValSorter {
		Keys: make([]string, 0, len(m)),
		Vals: make([]string, 0, len(m)),
	}

	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Vals = append(vs.Vals, v)
	}
	return vs
}

func (vs *ValSorter) Sort() {
	sort.Sort(vs)
}

func (vs *ValSorter) Len() int {
	return len(vs.Vals)
}

func (vs *ValSorter) Less(i, j int) bool {
	return bytes.Compare([]byte(vs.Keys[i]), []byte(vs.Keys[j])) < 0
}

func (vs *ValSorter) Swap(i, j int) {
	vs.Vals[i], vs.Vals[j] = vs.Vals[j], vs.Vals[i]
	vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}
