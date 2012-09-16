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
	HttpClient *http.Client
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
	client := Client{
		Host : host,
		AccessID : accessId,
		AccessKey : accessKey,
		HttpClient : http.DefaultClient,
	}
	return &client
}

func (c *Client) signHeader(req *http.Request) {
	//format x-oss-
	tmpParams := make(map[string]string)

	for k, v := range req.Header {
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

	date := req.Header.Get("Date")
	contentType := req.Header.Get("Content-Type")
	contentMd5 := req.Header.Get("Content-Md5")

	signStr := req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders  + req.URL.Path
	h := hmac.New(func() hash.Hash {return sha1.New()}, []byte(c.AccessKey)) //sha1.New()
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	authorizationStr := "OSS " + c.AccessID + ":" + signedStr
	fmt.Println(authorizationStr)
	req.Header.Set("Authorization", authorizationStr)
}

func (c *Client) getRequest(method, path string, params map[string]string) (req *http.Request) {
	reqUrl := "http://" + c.Host + path
	req, _ = http.NewRequest(method, reqUrl, nil)
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	req.Header.Set("Date", date)
	req.Header.Set("Host", c.Host)
	if params != nil {
		for k, v := range params {
			req.Header.Set(k, v)
		}
	}
	//req.Header.Set("Authorization", c.AccessID)
	//c.SignParam("GET", "/", req.Header)
	c.signHeader(req)
	return
}

//Get bucket list
func (c *Client) GetService() {
	req := c.getRequest("GET", "/", nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	respbytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respbytes))
	//fmt.Println(date)

}


func (c *Client) PutBucket(bname string) {
	req := c.getRequest("PUT", "/"+bname, nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	respbytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respbytes))
}

func (c *Client) PutBucketACL(bname, acl string) {
}

func (c *Client) DeleteBucket(bname string) {
	req := c.getRequest("DELETE", "/"+bname, nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	respbytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respbytes))
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
