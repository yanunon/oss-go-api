// Package Aliyun OSS API.
//
package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
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

type Bucket struct {
	Name         string
	CreationDate string
}

type CreateFileGroup struct {
	Part []GroupPart
}

type CompleteMultipartUpload struct {
	Part []Multipart
}

type Multipart struct {
	PartNumber int
	ETag       string
}

type CompleteMultipartUploadResult struct {
	Location string
	Bucket   string
	ETag     string
	Key      string
}

type FileGroup struct {
	Bucket     string
	Key        string
	ETag       string
	FileLength int
	FilePart   CreateFileGroup
}

type CompleteFileGroup struct {
	Bucket string
	Key    string
	Size   int
	ETag   string
}

type Client struct {
	AccessID   string
	AccessKey  string
	Host       string
	HttpClient *http.Client
}

type Buckets struct {
	Bucket []Bucket
}

type GroupPart struct {
	PartNumber int
	PartName   string
	PartSize   int
	ETag       string
}

type initMultipartUploadResult struct {
	Bucket   string
	Key      string
	UploadId string
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

//NewClient returns a new Client given a Host, AccessID and AccessKey.
func NewClient(host, accessId, accessKey string) *Client {
	client := Client{
		Host:       host,
		AccessID:   accessId,
		AccessKey:  accessKey,
		HttpClient: http.DefaultClient,
	}
	return &client
}

func (c *Client) signHeader(req *http.Request, canonicalizedResource string) {
	//format x-oss-
	tmpParams := make(map[string]string)

	for k, v := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			tmpParams[strings.ToLower(k)] = v[0]
		}
	}
	//sort
	vs := newValSorter(tmpParams)
	vs.Sort()

	canonicalizedOSSHeaders := ""
	for i := range vs.Keys {
		canonicalizedOSSHeaders += vs.Keys[i] + ":" + vs.Vals[i] + "\n"
	}

	date := req.Header.Get("Date")
	contentType := req.Header.Get("Content-Type")
	contentMd5 := req.Header.Get("Content-Md5")

	signStr := req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(c.AccessKey)) //sha1.New()
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	authorizationStr := "OSS " + c.AccessID + ":" + signedStr
	//fmt.Println(authorizationStr)
	req.Header.Set("Authorization", authorizationStr)
}

func (c *Client) doRequest(method, path, canonicalizedResource string, params map[string]string, data io.Reader) (resp *http.Response, err error) {
	reqUrl := "http://" + c.Host + path
	req, _ := http.NewRequest(method, reqUrl, data)
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	req.Header.Set("Date", date)
	req.Header.Set("Host", c.Host)

	if params != nil {
		for k, v := range params {
			req.Header.Set(k, v)
		}
	}

	if data != nil {
		req.Header.Set("Content-Length", strconv.Itoa(int(req.ContentLength)))
	}
	c.signHeader(req, canonicalizedResource)
	resp, err = c.HttpClient.Do(req)
	return
}

//Get bucket list. Return a ListAllMyBucketsResult object.
func (c *Client) GetService() (lar ListAllMyBucketsResult, err error) {
	resp, err := c.doRequest("GET", "/", "/", nil, nil)
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

	err = xml.Unmarshal(body, &lar)
	return
}

//Create a new bucket with a name.
func (c *Client) PutBucket(bname string) (err error) {
	resp, err := c.doRequest("PUT", "/"+bname, "/"+bname, nil, nil)
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
	resp, err := c.doRequest("PUT", "/"+bname, "", params, nil)
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

	resp, err := c.doRequest("GET", reqStr, "", nil, nil)
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
	err = xml.Unmarshal(body, &lbr)
	return
}

func (c *Client) GetBucketACL(bname string) (acl AccessControlPolicy, err error) {
	reqStr := "/" + bname + "?acl"
	resp, err := c.doRequest("GET", reqStr, reqStr, nil, nil)
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

	err = xml.Unmarshal(body, &acl)
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
	resp, err := c.doRequest("PUT", dst, "", params, nil)
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
	resp, err := c.doRequest("DELETE", opath, "", nil, nil)
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

	resp, err := c.doRequest("GET", opath, "", params, nil)
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

	//reqUrl := "http://" + c.Host + opath
	buffer := new(bytes.Buffer)

	fh, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer fh.Close()
	io.Copy(buffer, fh)

	contentType := http.DetectContentType(buffer.Bytes())
	params := map[string]string {}
	params["Content-Type"] = contentType

	resp, err := c.doRequest("PUT", opath, "", params, buffer)
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


func (c *Client) initMultipartUpload(opath string) (imur initMultipartUploadResult, err error) {
	resp, err := c.doRequest("POST", opath+"?uploads", opath+"?uploads", nil, nil)
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

	err = xml.Unmarshal(body, &imur)
	return
}

func (c *Client) uploadWorker(file *os.File, start, length, idx int, opath, uploadId string) (part Multipart, err error) {
	buffer := new(bytes.Buffer)
	file.Seek(int64(start), 0)
	io.CopyN(buffer, file, int64(length))
	h := md5.New()
	h.Write(buffer.Bytes())
	md5sum := fmt.Sprintf("%x", h.Sum(nil))
	md5sum = "\"" + strings.ToUpper(md5sum) + "\""

	reqStr := opath + "?partNumber=" + strconv.Itoa(idx) + "&uploadId=" + uploadId

	resp, err := c.doRequest("PUT", reqStr, reqStr, nil, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Printf("resp status:%s\n", err)
		fmt.Println(string(body))
		return
	}

	ETag := resp.Header.Get("ETag")
	if ETag != md5sum {
		fmt.Printf("ETag:%s != md5sum %s\n", ETag, md5sum)
	}
	part.ETag = ETag
	part.PartNumber = idx
	return
}


func (c *Client) uploadPart(imur initMultipartUploadResult, opath, filepath string) (cmu CompleteMultipartUpload, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}

	buffer_len := 5 << 20
	fi, err := file.Stat()
	if err != nil {
		return
	}
	file_len := int(fi.Size())
	thread_num := (file_len + buffer_len - 1) / buffer_len
	fmt.Printf("thread_num:%d\n", thread_num)
	for i := 0; i < thread_num; i++ {
		var part Multipart
		if i == thread_num-1 {
			last_len := file_len - buffer_len*i
			part, err = c.uploadWorker(file, i*buffer_len, last_len, i+1, opath, imur.UploadId)
		} else {
			part, err = c.uploadWorker(file, i*buffer_len, buffer_len, i+1, opath, imur.UploadId)
		}
		cmu.Part = append(cmu.Part, part)
	}
	return

}

func (c *Client) completeMultipartUpload(cmu CompleteMultipartUpload, opath, uploadId string) (cmur CompleteMultipartUploadResult, err error) {
	bs, err := xml.Marshal(cmu)
	if err != nil {
		return
	}

	reqStr := opath + "?uploadId=" + uploadId

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	resp, err := c.doRequest("POST", reqStr, reqStr, nil, buffer)
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
	err = xml.Unmarshal(body, &cmur)
	return
}

func (c *Client) PutLargeObject(opath string, filepath string) (err error) {
	if strings.HasPrefix(opath, "/") == false {
		opath = "/" + opath
	}

	imur, err := c.initMultipartUpload(opath)
	fmt.Printf("%+v\n", imur)
	imu, err := c.uploadPart(imur, opath, filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = c.completeMultipartUpload(imu, opath, imur.UploadId)
	return

}

func (c *Client) HeadObject(opath string) (header http.Header, err error) {
	if strings.HasPrefix(opath, "/") == false {
		opath = "/" + opath
	}
	resp, err := c.doRequest("HEAD", opath, "", nil, nil)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return
	}
	header = resp.Header
	return
}

func (c *Client) DeleteMultipleObject(bname string, onames []string) (err error) {
	return
}

func (c *Client) PostObjectGroup(cfg CreateFileGroup, opath string) (completefg CompleteFileGroup, err error) {
	//part := []GroupPart{{1, "11", "111"}, {2, "22", "222"}, {3, "33", "333"}}
	//fg := CreateFileGroup{Part:part}
	bs, err := xml.Marshal(cfg)
	if err != nil {
		return
	}

	if strings.HasPrefix(opath, "/") == false {
		opath = "/" + opath
	}

	reqStr := opath + "?group"
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	resp, err := c.doRequest("POST", reqStr, reqStr, nil, buffer)
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
	err = xml.Unmarshal(body, &completefg)
	return
}

func (c *Client) GetObjectGroupIndex(opath string) (fg FileGroup, err error) {
	params := map[string]string{"x-oss-file-group": ""}
	if strings.HasPrefix(opath, "/") == false {
		opath = "/" + opath
	}
	resp, err := c.doRequest("GET", opath, "", params, nil)
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

	err = xml.Unmarshal(body, &fg)
	return
}

func (c *Client) GetObjectGroup(opath string, rangeStart, rangeEnd int) (obytes []byte, err error) {
	return c.GetObject(opath, rangeStart, rangeEnd)
}

func (c *Client) HeadObjectGroup(opath string) (header http.Header, err error) {
	return c.HeadObject(opath)
}

func (c *Client) DeleteObjectGroup(bname string) (err error) {
	return c.DeleteObject(bname)
}

func newValSorter(m map[string]string) *valSorter {
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
