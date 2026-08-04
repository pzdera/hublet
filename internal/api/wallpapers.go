package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	maxWallpaperBytes = 20 * 1024 * 1024

	wallpaperTimeout = 20 * time.Second
)

type wallpaperFile struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

type wallpaperDownloadRequest struct {
	URL string `json:"url"`
}

func (
	s *Server,
) listWallpapers(
	w http.ResponseWriter,
	_ *http.Request,
) {
	entries, err :=
		os.ReadDir(
			s.wallpaperDir,
		)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to read wallpaper directory",
		)

		return
	}

	wallpapers :=
		make(
			[]wallpaperFile,
			0,
			len(entries),
		)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		extension :=
			strings.ToLower(
				filepath.Ext(
					entry.Name(),
				),
			)

		if !allowedWallpaperExtension(
			extension,
		) {
			continue
		}

		info, err :=
			entry.Info()

		if err != nil {
			continue
		}

		wallpapers = append(
			wallpapers,
			wallpaperFile{
				Filename: entry.Name(),

				URL: "/wallpapers/" +
					url.PathEscape(
						entry.Name(),
					),

				Size: info.Size(),

				Modified: info.ModTime().
					UTC().
					Format(
						time.RFC3339,
					),
			},
		)
	}

	sort.Slice(
		wallpapers,
		func(
			left int,
			right int,
		) bool {
			return wallpapers[left].
				Modified >
				wallpapers[right].
					Modified
		},
	)

	writeJSON(
		w,
		http.StatusOK,
		wallpapers,
	)
}

func (
	s *Server,
) uploadWallpaper(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body =
		http.MaxBytesReader(
			w,
			r.Body,
			maxWallpaperBytes+
				1024*1024,
		)

	if err := r.ParseMultipartForm(
		maxWallpaperBytes,
	); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid upload or file is too large",
		)

		return
	}

	file, header, err :=
		r.FormFile("file")

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"wallpaper file is required",
		)

		return
	}

	defer file.Close()

	wallpaper, err :=
		s.storeWallpaper(
			file,
			header,
		)

	if err != nil {
		writeError(
			w,
			http.StatusUnprocessableEntity,
			err.Error(),
		)

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		wallpaper,
	)
}

func (
	s *Server,
) downloadWallpaper(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body =
		http.MaxBytesReader(
			w,
			r.Body,
			64*1024,
		)

	var request wallpaperDownloadRequest

	decoder :=
		json.NewDecoder(
			r.Body,
		)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(
		&request,
	); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request",
		)

		return
	}

	parsedURL, err :=
		validateRemoteWallpaperURL(
			request.URL,
		)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	ctx, cancel :=
		context.WithTimeout(
			r.Context(),
			wallpaperTimeout,
		)

	defer cancel()

	httpRequest, err :=
		http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			parsedURL.String(),
			nil,
		)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"unable to create download request",
		)

		return
	}

	httpRequest.Header.Set(
		"User-Agent",
		"Hublet-Wallpaper-Downloader/2",
	)

	httpRequest.Header.Set(
		"Accept",
		"image/avif,image/webp,image/png,image/jpeg,image/gif;q=0.9,*/*;q=0.1",
	)

	client := &http.Client{
		Timeout: wallpaperTimeout,

		CheckRedirect: func(
			request *http.Request,
			via []*http.Request,
		) error {
			if len(via) >= 5 {
				return errors.New(
					"too many redirects",
				)
			}

			_, err :=
				validateRemoteWallpaperURL(
					request.URL.String(),
				)

			return err
		},
	}

	response, err :=
		client.Do(
			httpRequest,
		)

	if err != nil {
		writeError(
			w,
			http.StatusBadGateway,
			"unable to download wallpaper",
		)

		return
	}

	defer response.Body.Close()

	if response.StatusCode <
		http.StatusOK ||
		response.StatusCode >=
			http.StatusMultipleChoices {
		writeError(
			w,
			http.StatusBadGateway,
			fmt.Sprintf(
				"remote server returned HTTP %d",
				response.StatusCode,
			),
		)

		return
	}

	if response.ContentLength >
		maxWallpaperBytes {
		writeError(
			w,
			http.StatusRequestEntityTooLarge,
			"wallpaper exceeds the 20 MB limit",
		)

		return
	}

	limitedReader :=
		io.LimitReader(
			response.Body,
			maxWallpaperBytes+1,
		)

	content, err :=
		io.ReadAll(
			limitedReader,
		)

	if err != nil {
		writeError(
			w,
			http.StatusBadGateway,
			"unable to read downloaded wallpaper",
		)

		return
	}

	if len(content) >
		maxWallpaperBytes {
		writeError(
			w,
			http.StatusRequestEntityTooLarge,
			"wallpaper exceeds the 20 MB limit",
		)

		return
	}

	filename :=
		filenameFromURL(
			parsedURL,
		)

	header :=
		&multipart.FileHeader{
			Filename: filename,

			Size: int64(
				len(content),
			),

			Header: make(
				map[string][]string,
			),
		}

	header.Header.Set(
		"Content-Type",
		response.Header.Get(
			"Content-Type",
		),
	)

	wallpaper, err :=
		s.storeWallpaper(
			io.NopCloser(
				bytes.NewReader(
					content,
				),
			),
			header,
		)

	if err != nil {
		writeError(
			w,
			http.StatusUnprocessableEntity,
			err.Error(),
		)

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		wallpaper,
	)
}

func (
	s *Server,
) deleteWallpaper(
	w http.ResponseWriter,
	r *http.Request,
) {
	filename :=
		r.PathValue(
			"filename",
		)

	safeFilename, err :=
		sanitizeExistingFilename(
			filename,
		)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	path :=
		filepath.Join(
			s.wallpaperDir,
			safeFilename,
		)

	if err := os.Remove(
		path,
	); err != nil {
		if os.IsNotExist(err) {
			writeError(
				w,
				http.StatusNotFound,
				"wallpaper not found",
			)

			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"unable to delete wallpaper",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]bool{
			"success": true,
		},
	)
}

func (
	s *Server,
) storeWallpaper(
	file io.Reader,
	header *multipart.FileHeader,
) (
	wallpaperFile,
	error,
) {
	limitedReader :=
		io.LimitReader(
			file,
			maxWallpaperBytes+1,
		)

	content, err :=
		io.ReadAll(
			limitedReader,
		)

	if err != nil {
		return wallpaperFile{},
			errors.New(
				"unable to read wallpaper",
			)
	}

	if len(content) == 0 {
		return wallpaperFile{},
			errors.New(
				"wallpaper is empty",
			)
	}

	if len(content) >
		maxWallpaperBytes {
		return wallpaperFile{},
			errors.New(
				"wallpaper exceeds the 20 MB limit",
			)
	}

	contentType :=
		http.DetectContentType(
			content,
		)

	extension, err :=
		extensionForContentType(
			contentType,
		)

	if err != nil {
		return wallpaperFile{},
			err
	}

	baseName :=
		sanitizeBaseName(
			header.Filename,
		)

	randomSuffix, err :=
		randomHex(6)

	if err != nil {
		return wallpaperFile{},
			errors.New(
				"unable to generate wallpaper filename",
			)
	}

	filename :=
		fmt.Sprintf(
			"%s-%s%s",
			baseName,
			randomSuffix,
			extension,
		)

	path :=
		filepath.Join(
			s.wallpaperDir,
			filename,
		)

	if err := os.WriteFile(
		path,
		content,
		0o644,
	); err != nil {
		return wallpaperFile{},
			errors.New(
				"unable to save wallpaper",
			)
	}

	info, err :=
		os.Stat(path)

	if err != nil {
		_ = os.Remove(path)

		return wallpaperFile{},
			errors.New(
				"unable to inspect saved wallpaper",
			)
	}

	return wallpaperFile{
		Filename: filename,

		URL: "/wallpapers/" +
			url.PathEscape(
				filename,
			),

		Size: info.Size(),

		Modified: info.ModTime().
			UTC().
			Format(
				time.RFC3339,
			),
	}, nil
}

func validateRemoteWallpaperURL(
	value string,
) (*url.URL, error) {
	parsed, err :=
		url.ParseRequestURI(
			strings.TrimSpace(
				value,
			),
		)

	if err != nil {
		return nil,
			errors.New(
				"invalid wallpaper URL",
			)
	}

	if parsed.Scheme != "http" &&
		parsed.Scheme != "https" {
		return nil,
			errors.New(
				"wallpaper URL must use http or https",
			)
	}

	if parsed.Hostname() == "" {
		return nil,
			errors.New(
				"wallpaper URL requires a hostname",
			)
	}

	addresses, err :=
		net.DefaultResolver.
			LookupNetIP(
				context.Background(),
				"ip",
				parsed.Hostname(),
			)

	if err != nil {
		return nil,
			errors.New(
				"unable to resolve wallpaper hostname",
			)
	}

	for _, address := range addresses {
		if unsafeRemoteAddress(
			address,
		) {
			return nil,
				errors.New(
					"wallpaper URL may not target a private or local address",
				)
		}
	}

	return parsed, nil
}

func unsafeRemoteAddress(
	address netip.Addr,
) bool {
	return address.IsLoopback() ||
		address.IsPrivate() ||
		address.IsLinkLocalUnicast() ||
		address.IsLinkLocalMulticast() ||
		address.IsMulticast() ||
		address.IsUnspecified()
}

func sanitizeExistingFilename(
	value string,
) (string, error) {
	filename :=
		filepath.Base(
			strings.TrimSpace(
				value,
			),
		)

	if filename == "" ||
		filename == "." ||
		filename != value {
		return "",
			errors.New(
				"invalid wallpaper filename",
			)
	}

	if !allowedWallpaperExtension(
		strings.ToLower(
			filepath.Ext(
				filename,
			),
		),
	) {
		return "",
			errors.New(
				"unsupported wallpaper extension",
			)
	}

	return filename, nil
}

func sanitizeBaseName(
	value string,
) string {
	base :=
		strings.TrimSuffix(
			filepath.Base(value),
			filepath.Ext(value),
		)

	base =
		strings.ToLower(
			strings.TrimSpace(
				base,
			),
		)

	var builder strings.Builder

	lastDash := false

	for _, character := range base {
		isLetter :=
			character >= 'a' &&
				character <= 'z'

		isNumber :=
			character >= '0' &&
				character <= '9'

		if isLetter ||
			isNumber {
			builder.WriteRune(
				character,
			)

			lastDash = false

			continue
		}

		if !lastDash {
			builder.WriteByte('-')
			lastDash = true
		}
	}

	result :=
		strings.Trim(
			builder.String(),
			"-",
		)

	if result == "" {
		return "wallpaper"
	}

	if len(result) > 48 {
		result = result[:48]
		result =
			strings.TrimRight(
				result,
				"-",
			)
	}

	return result
}

func filenameFromURL(
	parsedURL *url.URL,
) string {
	filename :=
		filepath.Base(
			parsedURL.Path,
		)

	if filename == "" ||
		filename == "." ||
		filename == "/" {
		return "downloaded-wallpaper"
	}

	return filename
}

func extensionForContentType(
	contentType string,
) (string, error) {
	normalized :=
		strings.ToLower(
			strings.TrimSpace(
				strings.Split(
					contentType,
					";",
				)[0],
			),
		)

	switch normalized {
	case "image/jpeg":
		return ".jpg", nil

	case "image/png":
		return ".png", nil

	case "image/webp":
		return ".webp", nil

	case "image/gif":
		return ".gif", nil

	default:
		return "",
			fmt.Errorf(
				"unsupported wallpaper type %q",
				normalized,
			)
	}
}

func allowedWallpaperExtension(
	extension string,
) bool {
	switch extension {
	case ".jpg",
		".jpeg",
		".png",
		".webp",
		".gif":
		return true

	default:
		return false
	}
}

func randomHex(
	byteCount int,
) (string, error) {
	buffer :=
		make(
			[]byte,
			byteCount,
		)

	if _, err := rand.Read(
		buffer,
	); err != nil {
		return "", err
	}

	return hex.EncodeToString(
		buffer,
	), nil
}
