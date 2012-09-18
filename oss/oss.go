package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ACL_PUBLIC_RW = "public-read-write"
	ACL_PUBLIC_R  = "public-read"
	ACL_PRIVATE   = "private"
)

type AccessControlList struct {
	Grant string
}
type AccessControlPolicy struct {
	Owner             Owner
	AccessControlList AccessControlList
}

type Client struct {
	AccessID   string
	AccessKey  string
	Host       string
	HttpClient *http.Client
}

type Bucket struct {
	Name         string
	CreationDate string
}

type Buckets struct {
	Bucket []Bucket
}

type ListAllMyBucketsResult struct {
	Owner   Owner
	Buckets Buckets
}

type ListBucketResult struct {
	Name        string
	Prefix      string
	Marker      string
	MaxKeys     int
	Delimiter   string
	IsTruncated bool
	Contents    []Object
}

type Object struct {
	Key          string
	LastModified string
	ETag         string
	Type         string
	Size         int
	StorageClass string
	Owner        Owner
}

type Owner struct {
	ID          string
	DisplayName string
}

type valSorter struct {
	Keys []string
	Vals []string
}

func NewClient(host, accessId, accessKey string) *Client {
	client := Client{
		Host:       host,
		AccessID:   accessId,
		AccessKey:  accessKey,
		HttpClient: http.DefaultClient,
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
	vs := NewvalSorter(tmpParams)
	vs.Sort()

	canonicalizedOSSHeaders := ""
	for i := range vs.Keys {
		canonicalizedOSSHeaders += vs.Keys[i] + ":" + vs.Vals[i] + "\n"
	}

	date := req.Header.Get("Date")
	contentType := req.Header.Get("Content-Type")
	contentMd5 := req.Header.Get("Content-Md5")

	canonicalizedResource := req.URL.Path
	query := req.URL.Query()
	if _, ok := query["acl"]; ok {
		canonicalizedResource += "?" + "acl"
	}
	signStr := req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(c.AccessKey)) //sha1.New()
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	authorizationStr := "OSS " + c.AccessID + ":" + signedStr
	//fmt.Println(authorizationStr)
	req.Header.Set("Authorization", authorizationStr)
}

func (c *Client) doRequest(method, path string, params map[string]string) (resp *http.Response, err error) {
	reqUrl := "http://" + c.Host + path
	req, _ := http.NewRequest(method, reqUrl, nil)
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
	resp, err = c.HttpClient.Do(req)
	return
}

//Get bucket list
func (c *Client) GetService() (lar ListAllMyBucketsResult, err error) {
	resp, err := c.doRequest("GET", "/", nil)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		fmt.Println(string(body))
		return
	}

	xml.Unmarshal(body, &lar)
	return
}

func (c *Client) PutBucket(bname string) (err error) {
	resp, err := c.doRequest("PUT", "/"+bname, nil)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Println(string(body))
	}
	return
}

func (c *Client) PutBucketACL(bname, acl string) (err error) {
	params := map[string]string{"x-oss-acl": acl}
	resp, err := c.doRequest("PUT", "/"+bname, params)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Println(string(body))
	}
	return
}

func (c *Client) GetBucket(bname, prefix, marker, delimiter, maxkeys string) (lbr ListBucketResult, err error) {
	reqStr := "/" + bname
	query := map[string]string{}
	if prefix != "" {
		query["prefix"] = prefix
	}

	if marker != "" {
		query["marker"] = marker
	}

	if delimiter != "" {
		query["delimiter"] = delimiter
	}

	if maxkeys != "" {
		query["max-keys"] = maxkeys
	}

	if len(query) > 0 {
		reqStr += "?"
		for k, v := range query {
			reqStr += k + "=" + v + "&"
		}
	}

	resp, err := c.doRequest("GET", reqStr, nil)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		fmt.Println(string(body))
		return
	}
	xml.Unmarshal(body, &lbr)
	return
}

func (c *Client) GetBucketACL(bname string) (acl AccessControlPolicy, err error) {
	resp, err := c.doRequest("GET", "/"+bname+"?acl", nil)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		fmt.Println(string(body))
		return
	}

	xml.Unmarshal(body, &acl)
	return
}

func (c *Client) DeleteBucket(bname string) (err error) {
	return c.DeleteObject(bname)
}

func (c *Client) CopyObject(src, dst string) (err error) {
	if strings.HasPrefix(src, "/") == false {
		src = "/" + src
	}
	if strings.HasPrefix(dst, "/") == false {
		dst = "/" + dst
	}
	params := map[string]string{"x-oss-copy-source": src}
	resp, err := c.doRequest("PUT", dst, params)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Println(string(body))
	}
	return
}

func (c *Client) DeleteObject(opath string) (err error) {
	if strings.HasPrefix(opath, "/") == false {
		opath = "/" + opath
	}
	resp, err := c.doRequest("DELETE", opath, nil)
	if err != nil {
		return
	}

	if resp.StatusCode != 204 {
		err = errors.New(resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Println(string(body))
	}
	return
}

//Download object in opath.
//If rangeStart > -1 and rangeEnd > -1, download the object partially.
func (c *Client) GetObject(opath string, rangeStart, rangeEnd int) (obytes []byte, err error) {
	if strings.HasPrefix(opath, "/") == false {
		opath = "/" + opath
	}

	params := map[string]string{}
	if rangeStart > -1 && rangeEnd > -1 {
		params["range"] = "bytes=" + strconv.Itoa(rangeStart) + "-" + strconv.Itoa(rangeEnd)
	}

	resp, err := c.doRequest("GET", opath, params)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		err = errors.New(resp.Status)
		fmt.Println(string(body))
		return
	}
	//fmt.Println(string(body))
	obytes = body
	return
}

func (c *Client) PutObject(opath string, filepath string) (err error) {
	if strings.HasPrefix(opath, "/") == false {
		opath = "/" + opath
	}

	reqUrl := "http://" + c.Host + opath
	buffer := new(bytes.Buffer)

	fh, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer fh.Close()
	io.Copy(buffer, fh)

	contentType := http.DetectContentType(buffer.Bytes())

	req, err := http.NewRequest("PUT", reqUrl, buffer)
	if err != nil {
		return
	}

	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	req.Header.Set("Date", date)
	req.Header.Set("Host", c.Host)
	req.Header.Set("Content-Length", strconv.Itoa(int(req.ContentLength)))
	req.Header.Set("Content-Type", contentType)

	//req.Header.Set("Authorization", c.AccessID)
	//c.SignParam("GET", "/", req.Header)
	c.signHeader(req)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		fmt.Println(string(body))
		return
	}
	fmt.Println(string(body))
	return

}

func NewvalSorter(m map[string]string) *valSorter {
	vs := &valSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]string, 0, len(m)),
	}

	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Vals = append(vs.Vals, v)
	}
	return vs
}

func (vs *valSorter) Sort() {
	sort.Sort(vs)
}

func (vs *valSorter) Len() int {
	return len(vs.Vals)
}

func (vs *valSorter) Less(i, j int) bool {
	return bytes.Compare([]byte(vs.Keys[i]), []byte(vs.Keys[j])) < 0
}

func (vs *valSorter) Swap(i, j int) {
	vs.Vals[i], vs.Vals[j] = vs.Vals[j], vs.Vals[i]
	vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}
