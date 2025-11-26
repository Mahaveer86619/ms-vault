package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/png"
	"log"

	"github.com/MuhammadSaim/goavatar"
)

type AvatarService struct{}

func NewAvatarService() *AvatarService {
	return &AvatarService{}
}

func (s *AvatarService) GenerateHash(uniqueInput string) string {
	hasher := md5.New()
	hasher.Write([]byte(uniqueInput))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *AvatarService) GenerateAvatarImage(hash string) (image.Image, error) {
	avatar := goavatar.Make(hash, goavatar.WithSize(128))
	return avatar, nil
}

func (s *AvatarService) ImageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		log.Printf("Failed to encode image to PNG: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *AvatarService) GetAvatarURL(hash string) string {
	return "/avatar/" + hash
}
