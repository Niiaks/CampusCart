package file

import (
	"context"
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/rs/zerolog"
)

type Client struct {
	cld    *cloudinary.Cloudinary
	env    string
	logger *zerolog.Logger
}

// DirectUploadPayload contains the data the frontend needs to POST directly to Cloudinary.
type DirectUploadPayload struct {
	UploadURL string            `json:"upload_url"`
	Params    map[string]string `json:"params"`
}

func NewClient(cld *cloudinary.Cloudinary, logger *zerolog.Logger, env string) *Client {
	return &Client{
		cld:    cld,
		logger: logger,
		env:    env,
	}
}

// UploadImage uploads an image to Cloudinary under the given folder, returning secure URL and public ID.
func (c *Client) UploadImage(ctx context.Context, file interface{}, folder string) (string, string, error) {
	uploadFolder := "cc" + c.env + "/" + folder
	c.logger.Info().Str("folder", uploadFolder).Msg("image upload started...")
	resp, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		ResourceType:   "image",
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(true),
		Overwrite:      api.Bool(false),
		Folder:         uploadFolder,
	})

	if err != nil {
		c.logger.Err(err).Msg("an error occurred uploading the image")
		return "", "", err
	}

	return resp.SecureURL, resp.PublicID, nil
}

// UploadVideo uploads a video to Cloudinary and returns secure URL and public ID.
func (c *Client) UploadVideo(ctx context.Context, file interface{}) (string, string, error) {
	c.logger.Info().Msg("video upload started...")
	resp, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		ResourceType:   "video",
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(true),
		Overwrite:      api.Bool(false),
	})
	if err != nil {
		c.logger.Err(err).Msg("an error occurred uploading the video")
		return "", "", err
	}
	return resp.SecureURL, resp.PublicID, nil
}

// DeleteFile removes a Cloudinary asset by public ID and resource type (e.g., "image" or "video").
func (c *Client) DeleteFile(ctx context.Context, publicID string, resourceType string) error {
	c.logger.Info().Str("public_id", publicID).Msg("file deletion started...")
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})
	if err != nil {
		c.logger.Err(err).Str("public_id", publicID).Msg("an error occurred deleting the file")
		return err
	}
	return nil
}

// GenerateDirectUpload builds signed parameters so the client can upload directly to Cloudinary.
// resourceType can be "image", "video", or "auto" (default). folder is appended under the env prefix.
func (c *Client) GenerateDirectUpload(ctx context.Context, folder string, resourceType string) (*DirectUploadPayload, error) {
	if resourceType == "" {
		resourceType = "auto"
	}

	uploadURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s/upload", c.cld.Config.Cloud.CloudName, resourceType)
	folderPath := "cc" + c.env + "/" + folder
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	params := map[string]string{
		"timestamp": timestamp,
		"folder":    folderPath,
	}

	// Build string_to_sign: sorted key=value joined by '&'
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	stringToSign := strings.Join(parts, "&")

	signature := sha1.Sum([]byte(stringToSign + c.cld.Config.Cloud.APISecret))
	params["signature"] = fmt.Sprintf("%x", signature[:])
	params["api_key"] = c.cld.Config.Cloud.APIKey

	return &DirectUploadPayload{
		UploadURL: uploadURL,
		Params:    params,
	}, nil
}
