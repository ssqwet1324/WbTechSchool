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

const minioOpen = "localhost:9000"

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
	// кладем фото в minio
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

	// проверка на существования такого объекта в minio
	err := repo.checkPhoto(ctx, bucketName, objectName)
	if err != nil {
		return nil, err
	}

	// получаем объект из minio
	obj, err := repo.Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("GetPhotoBytesByVersion: error getting photo")
		return nil, err
	}
	defer func(obj *minio.Object) {
		err := obj.Close()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("GetPhotoBytesByVersion: error closing object")
			return
		}
	}(obj)

	zlog.Logger.Info().Str("objectName", objectName).Msg("GetPhotoBytesByVersion: found photo")

	// читаем байты
	data, err := io.ReadAll(obj)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("GetPhotoBytesByVersion: error reading photo bytes")
		return nil, err
	}

	return data, nil
}

// GetPhotoURLByVersion - получение presigned URL для фронта
func (repo *Repository) GetPhotoURLByVersion(ctx context.Context, bucketName, photoID, version string) (string, error) {
	objectName := fmt.Sprintf("%s_%s.jpg", photoID, version)

	// Проверяем наличие объекта
	err := repo.checkPhoto(ctx, bucketName, objectName)
	if err != nil {
		return "", err
	}

	// создаем объект-ссылку для формирования url
	url, err := repo.Client.PresignedGetObject(ctx, bucketName, objectName, time.Hour*24, nil)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("GetPhotoURLByVersion: error generating presigned URL")
		return "", err
	}

	// берем нужные поля для формирования url
	var minioImg entity.MinIOObject
	minioImg.Scheme = url.Scheme
	minioImg.Host = minioOpen // для открытия фото в браузере
	minioImg.Path = url.Path

	// формируем url
	fullURL := fmt.Sprintf("%s://%s%s", minioImg.Scheme, minioImg.Host, minioImg.Path)

	return fullURL, nil
}

// DeletePhotoFromMinIo - удалить фото из minio
func (repo *Repository) DeletePhotoFromMinIo(ctx context.Context, bucketName string, photo entity.PhotoInfo) error {
	// получаем все объекты с таким id
	opts := minio.ListObjectsOptions{
		Recursive:    true,
		Prefix:       photo.PhotoID,
		WithVersions: true,
	}

	// пробегаемся по всем объектам и удаляем их
	for object := range repo.Client.ListObjects(ctx, bucketName, opts) {
		if object.Err != nil {
			zlog.Logger.Error().Err(object.Err).Msg("DeletePhotoFromMinIo: error listing objects")
			continue
		}
		err := repo.Client.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{
			VersionID:   object.VersionID,
			ForceDelete: true,
		})
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("DeletePhotoFromMinIo: error deleting object")
			return err
		}
	}

	return nil
}

// checkPhoto - проверяем наличие объекта в minio
func (repo *Repository) checkPhoto(ctx context.Context, bucketName, objectName string) error {
	_, err := repo.Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" || minio.ToErrorResponse(err).Code == "NoSuchObject" {
			zlog.Logger.Error().Str("objectName", objectName).Msg("GetPhotoURLByVersion: photo not found")
			return err
		}
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("GetPhotoURLByVersion: error checking object existence")
		return err
	}

	return nil
}
