package file

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/rs/zerolog"
)

type Client struct {
	cld    *cloudinary.Cloudinary
	logger *zerolog.Logger
}

func NewClient(cld *cloudinary.Cloudinary, logger *zerolog.Logger) *Client {
	return &Client{
		cld:    cld,
		logger: logger,
	}
}

func (c *Client) UploadImage(ctx context.Context, file interface{}) (string, string, error) {
	c.logger.Info().Msg("image upload started...")
	resp, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		ResourceType:    "image",
		UseFilename:     api.Bool(true),
		UniqueFilename:  api.Bool(true),
		Moderation:      "webpurify",
		QualityAnalysis: api.Bool(true),
		Overwrite:       api.Bool(false),
	})
	if err != nil {
		c.logger.Err(err).Msg("an error occurred uploading the image")
		return "", "", err
	}
	return resp.SecureURL, resp.PublicID, nil
}

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
