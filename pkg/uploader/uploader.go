package uploader

import (
	"accompany-sdk/pkg/misc"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"accompany-sdk/pkg/must"
	"accompany-sdk/pkg/ternary"
	"github.com/hashicorp/go-uuid"
	"github.com/openimsdk/tools/log"
	qiniuAuth "github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/qiniu/go-sdk/v7/storage"
)

// DefaultUploadExpireAfterDays 默认上传文件过期时间，0 表示永不过期
const DefaultUploadExpireAfterDays = 0

type Uploader struct {
	conf       *Config
	baseURL    string
	httpClient *http.Client
}

type Config struct {
	// qiniu
	StorageAppKey       string `json:"storage_appkey" yaml:"storage_appkey"`
	StorageAppSecret    string `json:"-" yaml:"storage_secret"`
	StorageBucket       string `json:"storage_bucket" yaml:"storage_bucket"`
	StorageCallback     string `json:"storage_callback" yaml:"storage_callback"`
	StorageCallbackHost string `json:"storage_callback_host" yaml:"storage_callback_host"`
	StorageDomain       string `json:"storage_domain" yaml:"storage_domain"`
	StorageRegion       string `json:"storage_region" yaml:"storage_region"`
}

func NewUploader(conf *Config) *Uploader {
	client := &http.Client{Timeout: 120 * time.Second}
	return &Uploader{conf: conf, baseURL: conf.StorageDomain, httpClient: client}
}

func New(conf *Config) *Uploader {
	return &Uploader{conf: conf, baseURL: conf.StorageDomain, httpClient: &http.Client{Timeout: 120 * time.Second}}
}

type UploadInit struct {
	Filename string `json:"filename"`
	Token    string `json:"token"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	URL      string `json:"url"`
	// Uploaded To inform the client whether the file has been uploaded before.
	Uploaded bool `json:"uploaded,omitempty"`
}

const (
	UploadUsageAvatar    = "avatar"
	UploadUsageImageChat = "chat"
	UploadUsageDocument  = "document"
)

type UploadCallback struct {
	Key     string `json:"key"`
	Hash    string `json:"hash"`
	Fsize   int64  `json:"fsize"`
	Bucket  string `json:"bucket"`
	Name    string `json:"name"`
	UID     int64  `json:"uid"`
	Channel string `json:"channel"`
}

func (cb UploadCallback) ToJSON() string {
	data, _ := json.Marshal(cb)
	return string(data)
}

// Init 文件上传初始化，生成上传凭证
func (u *Uploader) Init(filename string, uid int, usage string, maxSizeInMB int64, expireAfterDays int, enableCallback bool, channel string) UploadInit {
	putPolicy := storage.PutPolicy{
		Scope:           u.conf.StorageBucket,
		FsizeLimit:      1024 * 1024 * maxSizeInMB,
		DeleteAfterDays: expireAfterDays,
	}

	if enableCallback {
		putPolicy.CallbackHost = u.conf.StorageCallbackHost
		putPolicy.CallbackURL = u.conf.StorageCallback
		putPolicy.CallbackBodyType = "application/json"
		putPolicy.CallbackBody = fmt.Sprintf(
			`{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":%s,"uid":%d,"usage":"%s","channel":"%s"}`,
			strconv.Quote(filename),
			uid,
			usage,
			channel,
		)
	}

	mac := qiniuAuth.New(u.conf.StorageAppKey, u.conf.StorageAppSecret)

	var publicUrl, key string
	switch usage {
	case UploadUsageAvatar:
		key = fmt.Sprintf("ai-server/%d/avatar/ugc%s%s", uid, must.Must(uuid.GenerateUUID()), misc.FileExt(filename))
		publicUrl = fmt.Sprintf("%s/%s-avatar", u.baseURL, key)
	default:
		key = fmt.Sprintf("ai-server/%d/%s/ugc%s%s", uid, time.Now().Format("20060102"), must.Must(uuid.GenerateUUID()), misc.FileExt(filename))
		publicUrl = fmt.Sprintf("%s/%s", u.baseURL, key)
	}

	return UploadInit{
		Filename: filename,
		Token:    putPolicy.UploadToken(mac),
		Bucket:   u.conf.StorageBucket,
		Key:      key,
		URL:      publicUrl,
	}
}

// Upload 上传文件
func (u *Uploader) Upload(ctx context.Context, init UploadInit) (string, error) {
	cfg := storage.Config{}

	region, ok := storage.GetRegionByID(storage.RegionID(u.conf.StorageRegion))
	if !ok {
		return "", fmt.Errorf("invalid storage region: %s", u.conf.StorageRegion)
	}

	cfg.Region = &region
	cfg.UseHTTPS = true
	cfg.UseCdnDomains = true

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := formUploader.PutFile(ctx, &ret, init.Token, init.Key, init.Filename, nil)
	if err != nil {
		return "", err
	}

	return init.URL, nil
}

// UploadRemoteFile 上传远程文件（先下载，后上传）
func (u *Uploader) UploadRemoteFile(ctx context.Context, url string, uid int, expiredAfterDays int, ext string, breakWall bool) (string, error) {
	res, err := u.uploadRemoteFile(ctx, url, uid, expiredAfterDays, ext, breakWall)
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		return u.uploadRemoteFile(ctx, url, uid, expiredAfterDays, ext, breakWall)
	}

	return res, nil
}

func (u *Uploader) uploadRemoteFile(ctx context.Context, url string, uid int, expiredAfterDays int, ext string, breakWall bool) (string, error) {
	client := ternary.If(breakWall, u.httpClient, &http.Client{Timeout: 120 * time.Second})
	resp, err := client.Get(url)
	if err != nil {
		time.Sleep(500 * time.Millisecond)

		resp, err = client.Get(url)
		if err != nil {
			return "", fmt.Errorf("download remote file failed: %w", err)
		}
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read remote file failed: %w", err)
	}

	return u.UploadStream(ctx, uid, expiredAfterDays, data, ext)
}

// UploadFile 上传文件
func (u *Uploader) UploadFile(ctx context.Context, uid int, expiredAfterDays int, filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file failed: %w", err)
	}

	return u.UploadStream(ctx, uid, expiredAfterDays, data, misc.FileExt(filePath))
}

// UploadStream 上传文件流
func (u *Uploader) UploadStream(ctx context.Context, uid int, expireAfterDays int, data []byte, ext string) (string, error) {
	res, err := u.uploadStream(ctx, uid, expireAfterDays, data, ext)
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		return u.uploadStream(ctx, uid, expireAfterDays, data, ext)
	}

	return res, nil
}

func (u *Uploader) uploadStream(ctx context.Context, uid int, expireAfterDays int, data []byte, ext string) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope:           u.conf.StorageBucket,
		FsizeLimit:      1024 * 1024 * 20,
		DeleteAfterDays: expireAfterDays,
	}

	if u.conf.StorageCallback != "" {
		putPolicy.CallbackHost = u.conf.StorageCallbackHost
		putPolicy.CallbackURL = u.conf.StorageCallback
		putPolicy.CallbackBodyType = "application/json"
		putPolicy.CallbackBody = fmt.Sprintf(`{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)","uid":%d,"usage":"%s","channel":"server"}`, uid, "")
	}

	mac := qiniuAuth.New(u.conf.StorageAppKey, u.conf.StorageAppSecret)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	region, ok := storage.GetRegionByID(storage.RegionID(u.conf.StorageRegion))
	if !ok {
		return "", fmt.Errorf("invalid storage region: %s", u.conf.StorageRegion)
	}

	cfg.Region = &region
	cfg.UseHTTPS = true
	cfg.UseCdnDomains = true

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	key := fmt.Sprintf("ai-server/%d/%s/aigc%s.%s", uid, time.Now().Format("20060102"), must.Must(uuid.GenerateUUID()), ext)

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	err := formUploader.Put(ctx, &ret, upToken, key, bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		return "", fmt.Errorf("upload file failed: %w", err)
	}

	return fmt.Sprintf("%s/%s", u.baseURL, key), nil
}

// RemoveFile 删除文件
func (u *Uploader) RemoveFile(ctx context.Context, pathWithoutURLPrefix string) error {
	log.ZInfo(ctx, "删除文件", "path", pathWithoutURLPrefix)

	mac := qiniuAuth.New(u.conf.StorageAppKey, u.conf.StorageAppSecret)
	cfg := storage.Config{
		UseHTTPS: true,
	}

	bucketManager := storage.NewBucketManager(mac, &cfg)
	return bucketManager.Delete(u.conf.StorageBucket, pathWithoutURLPrefix)
}

// ForbidFile 禁用文件
func (u *Uploader) ForbidFile(ctx context.Context, pathWithoutURLPrefix string) error {
	log.ZInfo(ctx, "禁用文件", "path", pathWithoutURLPrefix)

	mac := qiniuAuth.New(u.conf.StorageAppKey, u.conf.StorageAppSecret)
	cfg := storage.Config{
		UseHTTPS: true,
	}

	bucketManager := storage.NewBucketManager(mac, &cfg)
	return bucketManager.UpdateObjectStatus(u.conf.StorageBucket, pathWithoutURLPrefix, false)
}

// RefreshCDN 刷新 CDN 缓存
func (u *Uploader) RefreshCDN(ctx context.Context, urls []string) (cdn.RefreshResp, error) {
	mac := qiniuAuth.New(u.conf.StorageAppKey, u.conf.StorageAppSecret)

	cdnManager := cdn.NewCdnManager(mac)
	return cdnManager.RefreshUrls(urls)
}

// MakePrivateURL 生成私有文件访问 URL
func (u *Uploader) MakePrivateURL(key string, ttl time.Duration) string {
	mac := qiniuAuth.New(u.conf.StorageAppKey, u.conf.StorageAppSecret)
	deadline := time.Now().Add(ttl).Unix() // 24 小时有效期

	return storage.MakePrivateURL(mac, u.baseURL, key, deadline)
}
