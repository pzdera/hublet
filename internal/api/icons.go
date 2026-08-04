package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	maxIconBytes = 2 * 1024 * 1024

	iconDownloadTimeout = 15 * time.Second

	dashboardIconsPNGBase = "https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png"
)

var dashboardIconSlugPattern = regexp.MustCompile(
	`^[a-z0-9]+(?:-[a-z0-9]+)*$`,
)

type iconFile struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

type dashboardIconRequest struct {
	Value string `json:"value"`
}

func (s *Server) listIcons(
	w http.ResponseWriter,
	_ *http.Request,
) {
	entries, err :=
		os.ReadDir(s.iconDir)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to read icon directory",
		)

		return
	}

	icons := make(
		[]iconFile,
		0,
		len(entries),
	)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.ToLower(
			filepath.Ext(entry.Name()),
		) != ".png" {
			continue
		}

		info, err := entry.Info()

		if err != nil {
			continue
		}

		icons = append(
			icons,
			iconFile{
				Filename: entry.Name(),

				URL: "/icons/" +
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
		icons,
		func(
			left int,
			right int,
		) bool {
			return icons[left].Modified >
				icons[right].Modified
		},
	)

	writeJSON(
		w,
		http.StatusOK,
		icons,
	)
}

func (s *Server) uploadIcon(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxIconBytes+
			512*1024,
	)

	if err := r.ParseMultipartForm(
		maxIconBytes,
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
			"PNG icon file is required",
		)

		return
	}

	defer file.Close()

	icon, err :=
		s.storeUploadedIcon(
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
		icon,
	)
}

func (s *Server) downloadDashboardIcon(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		32*1024,
	)

	var request dashboardIconRequest

	decoder :=
		json.NewDecoder(r.Body)

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

	slug, err :=
		dashboardIconSlug(
			request.Value,
		)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	iconURL :=
		fmt.Sprintf(
			"%s/%s.png",
			dashboardIconsPNGBase,
			slug,
		)

	ctx, cancel :=
		context.WithTimeout(
			r.Context(),
			iconDownloadTimeout,
		)

	defer cancel()

	httpRequest, err :=
		http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			iconURL,
			nil,
		)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to create icon request",
		)

		return
	}

	httpRequest.Header.Set(
		"User-Agent",
		"Hublet-v2-Icon-Downloader/2",
	)

	httpRequest.Header.Set(
		"Accept",
		"image/png",
	)

	client := &http.Client{
		Timeout: iconDownloadTimeout,
	}

	response, err :=
		client.Do(httpRequest)

	if err != nil {
		writeError(
			w,
			http.StatusBadGateway,
			"unable to download Dashboard Icon",
		)

		return
	}

	defer response.Body.Close()

	if response.StatusCode ==
		http.StatusNotFound {
		writeError(
			w,
			http.StatusNotFound,
			"Dashboard Icon was not found",
		)

		return
	}

	if response.StatusCode <
		http.StatusOK ||
		response.StatusCode >=
			http.StatusMultipleChoices {
		writeError(
			w,
			http.StatusBadGateway,
			fmt.Sprintf(
				"Dashboard Icons returned HTTP %d",
				response.StatusCode,
			),
		)

		return
	}

	content, err :=
		readLimitedIcon(
			response.Body,
		)

	if err != nil {
		writeError(
			w,
			http.StatusUnprocessableEntity,
			err.Error(),
		)

		return
	}

	if err := validatePNG(content); err != nil {
		writeError(
			w,
			http.StatusUnprocessableEntity,
			err.Error(),
		)

		return
	}

	filename :=
		"dashboard-" +
			slug +
			".png"

	path :=
		filepath.Join(
			s.iconDir,
			filename,
		)

	if err := os.WriteFile(
		path,
		content,
		0o644,
	); err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to save icon",
		)

		return
	}

	icon, err :=
		iconFileFromPath(
			path,
			filename,
		)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to inspect saved icon",
		)

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		icon,
	)
}

func (s *Server) deleteIcon(
	w http.ResponseWriter,
	r *http.Request,
) {
	filename, err :=
		sanitizeIconFilename(
			r.PathValue(
				"filename",
			),
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
			s.iconDir,
			filename,
		)

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			writeError(
				w,
				http.StatusNotFound,
				"icon not found",
			)

			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"unable to delete icon",
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

func (s *Server) storeUploadedIcon(
	file multipart.File,
	header *multipart.FileHeader,
) (
	iconFile,
	error,
) {
	content, err :=
		readLimitedIcon(file)

	if err != nil {
		return iconFile{}, err
	}

	if err := validatePNG(content); err != nil {
		return iconFile{}, err
	}

	baseName :=
		sanitizeIconBaseName(
			header.Filename,
		)

	randomSuffix, err :=
		randomIconHex(5)

	if err != nil {
		return iconFile{},
			errors.New(
				"unable to generate icon filename",
			)
	}

	filename :=
		fmt.Sprintf(
			"%s-%s.png",
			baseName,
			randomSuffix,
		)

	path :=
		filepath.Join(
			s.iconDir,
			filename,
		)

	if err := os.WriteFile(
		path,
		content,
		0o644,
	); err != nil {
		return iconFile{},
			errors.New(
				"unable to save icon",
			)
	}

	return iconFileFromPath(
		path,
		filename,
	)
}

func readLimitedIcon(
	reader io.Reader,
) ([]byte, error) {
	limited :=
		io.LimitReader(
			reader,
			maxIconBytes+1,
		)

	content, err :=
		io.ReadAll(limited)

	if err != nil {
		return nil,
			errors.New(
				"unable to read icon",
			)
	}

	if len(content) == 0 {
		return nil,
			errors.New(
				"icon is empty",
			)
	}

	if len(content) >
		maxIconBytes {
		return nil,
			errors.New(
				"icon exceeds the 2 MB limit",
			)
	}

	return content, nil
}

func validatePNG(
	content []byte,
) error {
	contentType :=
		http.DetectContentType(
			content,
		)

	if contentType != "image/png" {
		return fmt.Errorf(
			"only PNG icons are supported, received %s",
			contentType,
		)
	}

	return nil
}

func dashboardIconSlug(
	value string,
) (string, error) {
	trimmed :=
		strings.TrimSpace(
			value,
		)

	if trimmed == "" {
		return "",
			errors.New(
				"Dashboard Icons name or URL is required",
			)
	}

	slug := trimmed

	if parsed, err :=
		url.Parse(trimmed); err == nil &&
		parsed.Hostname() != "" {
		path :=
			strings.Trim(
				parsed.Path,
				"/",
			)

		parts :=
			strings.Split(
				path,
				"/",
			)

		switch {
		case len(parts) >= 2 &&
			parts[0] == "icons":
			slug = parts[1]

		case len(parts) >= 2 &&
			parts[len(parts)-2] == "png":
			slug =
				parts[len(parts)-1]

		default:
			slug =
				parts[len(parts)-1]
		}
	}

	slug =
		strings.TrimSuffix(
			strings.ToLower(
				strings.TrimSpace(
					slug,
				),
			),
			".png",
		)

	slug =
		strings.ReplaceAll(
			slug,
			"_",
			"-",
		)

	if !dashboardIconSlugPattern.
		MatchString(slug) {
		return "",
			errors.New(
				"invalid Dashboard Icons name",
			)
	}

	return slug, nil
}

func sanitizeIconFilename(
	value string,
) (string, error) {
	trimmed :=
		strings.TrimSpace(value)

	filename :=
		filepath.Base(trimmed)

	if filename == "" ||
		filename == "." ||
		filename != trimmed {
		return "",
			errors.New(
				"invalid icon filename",
			)
	}

	if strings.ToLower(
		filepath.Ext(filename),
	) != ".png" {
		return "",
			errors.New(
				"only PNG icons may be deleted",
			)
	}

	return filename, nil
}

func sanitizeIconBaseName(
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

		if isLetter || isNumber {
			builder.WriteRune(character)
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
		return "custom-icon"
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

func iconFileFromPath(
	path string,
	filename string,
) (
	iconFile,
	error,
) {
	info, err := os.Stat(path)

	if err != nil {
		return iconFile{}, err
	}

	return iconFile{
		Filename: filename,

		URL: "/icons/" +
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

func randomIconHex(
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
