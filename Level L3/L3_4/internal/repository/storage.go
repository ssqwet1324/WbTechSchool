package repository

import (
	"bytes"
	"context"
	"fmt"
	"image_processor/internal/entity"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/wb-go/wbf/zlog"
)

// UploadPhoto - загрузка исходного фото (из LoadPhoto)
func (repo *Repository) UploadPhoto(ctx context.Context, bucketName string, photo entity.LoadPhoto) error {
	objectName := fmt.Sprintf("%s_original.jpg", photo.ID)
	_, err := repo.Client.PutObject(ctx, bucketName, objectName, photo.Reader, photo.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("UploadPhoto: error uploading photo")
		return err
	}

	return nil
}

// UploadPhotoBytes - загрузка обработанных фото из []byte
func (repo *Repository) UploadPhotoBytes(ctx context.Context, bucketName, objectName string, data []byte) error {
	_, err := repo.Client.PutObject(ctx, bucketName, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("UploadPhotoBytes: error uploading photo")
		return err
	}

	return nil
}

// GetPhotoBytesByVersion - получение байт оригинального или обработанного фото
func (repo *Repository) GetPhotoBytesByVersion(ctx context.Context, bucketName, photoID, version string) ([]byte, error) {
	objectName := fmt.Sprintf("%s_%s.jpg", photoID, version)

	obj, err := repo.Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("GetPhotoBytesByVersion: error getting photo")
		return nil, err
	}
	defer obj.Close()

	zlog.Logger.Info().Str("objectName", objectName).Msg("GetPhotoBytesByVersion: found photo")

	data, err := io.ReadAll(obj)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("GetPhotoBytesByVersion: error reading photo bytes")
		return nil, err
	}

	return data, nil
}

// GetPhotoUrlByVersion - получение presigned URL для фронта
func (repo *Repository) GetPhotoUrlByVersion(ctx context.Context, bucketName, photoID, version string) (string, error) {
	var minioImg entity.MinIOObject
	objectName := fmt.Sprintf("%s_%s.jpg", photoID, version)

	url, err := repo.Client.PresignedGetObject(ctx, bucketName, objectName, time.Hour*24, nil)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("GetPhotoUrlByVersion: error generating presigned URL")
		return "", err
	}

	// Сохраняем нужные части URL
	minioImg.Scheme = url.Scheme
	minioImg.Host = "localhost:9000" // для открытия фото в браузере
	minioImg.Path = url.Path

	// Собираем полный URL
	fullURL := fmt.Sprintf("%s://%s%s", minioImg.Scheme, minioImg.Host, minioImg.Path)

	return fullURL, nil
}

func (repo *Repository) DeletePhotoFromMinIo(ctx context.Context, bucketName string, photo entity.PhotoInfo) error {
	objectName := fmt.Sprintf("%s_%s.jpg", photo.PhotoID, photo.Version)

	err := repo.Client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("DeletePhotoFromMinIo: failed to remove photo")
		return err
	}

	zlog.Logger.Info().Str("objectName", objectName).Msg("DeletePhotoFromMinIo: photo successfully removed")

	return nil
}
