package uploader

import (
	"accompany-sdk/pkg/str"
	"accompany-sdk/pkg/ternary"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var supportFilters = []string{"-1024_square", "-512_square", "-avatar", "-maxsize700", "-maxsize800", "-square_500", "-thumb", "-thumb1000", "-thumb_500", "-fix_square_1024"}
var supportImages = []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}

// BuildImageURLWithFilter build image url with filter
func BuildImageURLWithFilter(remoteURL string, filter, storageDomain string) string {
	if !str.HasPrefixes(remoteURL, []string{
		"https://ssl.aicode.cc/",
		ternary.If(storageDomain == "", "https://ssl.aicode.cc/", storageDomain),
	}) {
		return remoteURL
	}

	if str.HasSuffixes(strings.ToLower(remoteURL), supportImages) {
		remoteURL = remoteURL + "-" + filter
	} else if str.HasSuffixes(strings.ToLower(remoteURL), supportFilters) {
		segs := strings.Split(remoteURL, "-")
		segs[len(segs)-1] = filter

		remoteURL = strings.Join(segs, "-")
	}

	return remoteURL
}

// RemoveImageFilter 移除图片的 filter
func RemoveImageFilter(imageURL string) string {
	if str.HasSuffixes(strings.ToLower(imageURL), supportFilters) {
		segs := strings.Split(imageURL, "-")
		segs = segs[:len(segs)-1]

		return strings.Join(segs, "-")
	}

	return imageURL
}

var (
	ErrFileForbidden = fmt.Errorf("文件违规已被禁用")
)

// DownloadRemoteFile download remote file to local
func DownloadRemoteFile(ctx context.Context, remoteURL string) (string, error) {
	if str.HasSuffixes(strings.ToLower(remoteURL), supportImages) {
		remoteURL = remoteURL + "-thumb"
	}

	resp, err := http.Get(remoteURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusForbidden {
			return "", ErrFileForbidden
		}
		return "", fmt.Errorf("download remote file failed: [%d] %s", resp.StatusCode, resp.Status)
	}

	prefix, _ := uuid.GenerateUUID()
	savePath := filepath.Join(os.TempDir(), prefix+"-"+filepath.Base(remoteURL))
	f, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}

	return savePath, nil
}

func DownloadRemoteFileAsBase64(ctx context.Context, remoteURL string) (string, error) {
	if str.HasSuffixes(strings.ToLower(remoteURL), supportImages) {
		remoteURL = remoteURL + "-thumb"
	}

	resp, err := http.Get(remoteURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusForbidden {
			return "", ErrFileForbidden
		}
		return "", fmt.Errorf("download remote file failed: [%d] %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(data)
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

func DownloadRemoteFileAsBase64Raw(ctx context.Context, remoteURL string, fromQiniu bool) (imageData string, mimeType string, err error) {
	if fromQiniu && str.HasSuffixes(strings.ToLower(remoteURL), supportImages) {
		remoteURL = remoteURL + "-thumb"
	}

	resp, err := http.Get(remoteURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusForbidden {
			return "", "", ErrFileForbidden
		}
		return "", "", fmt.Errorf("download remote file failed: [%d] %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(data), http.DetectContentType(data), nil
}

type ImageInfo struct {
	// Size 文件大小，单位：Bytes
	Size int64 `json:"size"`
	// Format 图片类型，如png、jpeg、gif、bmp等
	Format string `json:"format"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

func QueryImageInfo(imageURL string) (*ImageInfo, error) {
	resp, err := http.Get(RemoveImageFilter(imageURL) + "?imageInfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusForbidden {
			return nil, ErrFileForbidden
		}
		return nil, fmt.Errorf("query remote file info failed: [%d] %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info ImageInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
