package usecase

import (
	"context"
	"fmt"
	"image_processor/internal/config"
	"image_processor/internal/entity"
	"strconv"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/wb-go/wbf/zlog"
)

// RepositoryProvider - интерфейс репозитория
type RepositoryProvider interface {
	AddPhotoInfo(ctx context.Context, photo entity.LoadPhoto) error
	UploadPhoto(ctx context.Context, bucketName string, photo entity.LoadPhoto) error
	UploadPhotoBytes(ctx context.Context, bucketName, objectName string, data []byte) error
	ChangeStatus(ctx context.Context, photoID, status string) error
	GetPhotoBytesByVersion(ctx context.Context, bucketName, photoID, version string) ([]byte, error)
	GetPhotoURLByVersion(ctx context.Context, bucketName, photoID, version string) (string, error)
	DeletePhoto(ctx context.Context, photoID string) error
	DeletePhotoFromMinIo(ctx context.Context, bucketName string, photo entity.PhotoInfo) error
}

// UseCase - структура бизнес-логики
type UseCase struct {
	repo RepositoryProvider
	cfg  *config.Config
}

// New - конструктор
func New(repo RepositoryProvider, cfg *config.Config) *UseCase {
	return &UseCase{repo: repo, cfg: cfg}
}

// generateIDPhoto - генерация id для фото
func generateIDPhoto() string {
	return uuid.New().String()
}

// AddPhoto - добавление исходного фото
func (uc *UseCase) AddPhoto(ctx context.Context, photo entity.LoadPhoto) (string, error) {
	photoID := generateIDPhoto()
	photo.ID = photoID

	// сохраняем метаданные в БД
	if err := uc.repo.AddPhotoInfo(ctx, photo); err != nil {
		zlog.Logger.Error().Err(err).Msg("AddPhoto: failed adding photo info")
		return "", err
	}

	// загружаем оригинальное фото в MinIO
	if err := uc.repo.UploadPhoto(ctx, uc.cfg.BucketName, photo); err != nil {
		zlog.Logger.Error().Err(err).Msg("AddPhoto: failed uploading photo")
		if err := uc.repo.DeletePhoto(ctx, photoID); err != nil {
			zlog.Logger.Error().Err(err).Msg("AddPhoto: failed deleting photo from db")
			return "", err
		}
		return "", err
	}

	photo.Status = "loaded"
	// обновляем статус
	if err := uc.repo.ChangeStatus(ctx, photo.ID, photo.Status); err != nil {
		zlog.Logger.Error().Err(err).Msg("AddPhoto: failed updating status photo")
	}

	zlog.Logger.Info().Str("photo", photoID).Msg("AddPhoto: photo added successfully")

	return photo.ID, nil
}

// PhotoProcessing - обработка фото (resize, miniature, watermark) и возврат presigned URL
func (uc *UseCase) PhotoProcessing(ctx context.Context, photo entity.PhotoInfo) (string, error) {
	photo.BucketName = uc.cfg.BucketName

	var width, height int
	var err error

	// проверяем что ширина не пустая
	if photo.Width != "" {
		width, err = strconv.Atoi(photo.Width)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("PhotoProcessing: failed to convert width")
			return "", err
		}
	}

	// проверяем что высота не пустая
	if photo.Height != "" {
		height, err = strconv.Atoi(photo.Height)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("PhotoProcessing: failed to convert height")
			return "", err
		}
	}

	// Берём оригинальное фото из MinIO
	data, err := uc.repo.GetPhotoBytesByVersion(ctx, photo.BucketName, photo.PhotoID, "original")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("PhotoProcessing: failed to get original photo bytes")
		return "", err
	}

	// создаем объект для bimg
	var processedImage []byte
	img := bimg.NewImage(data)

	// Обрабатываем по версии
	switch photo.Version {
	case "resize":
		// изменяем размер фото
		processedImage, err = img.Resize(width, height)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("PhotoProcessing: failed to resize image")
			return "", err
		}
	case "miniature":
		// создаем миниатюру
		processedImage, err = img.Resize(width, height)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("PhotoProcessing: failed to resize image")
			return "", err
		}
	case "watermark":
		// Сначала конвертируем изображение в формат с альфа-каналом (RGBA)
		options := bimg.Options{
			Type:       bimg.JPEG,
			Embed:      true,                               // добавляет альфа-канал
			Background: bimg.Color{R: 255, G: 255, B: 255}, // белый фон
			Watermark: bimg.Watermark{
				Text: "Watermark",
			},
		}

		// итоговое изображение
		processedImage, err = img.Process(options)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("PhotoProcessing: failed to watermark")
			return "", err
		}

	default:
		// если ввели неподдерживаемый формат редактирования
		zlog.Logger.Error().Msg("PhotoProcessing: unsupported photo version")
		return "", err
	}

	// Загружаем обработанное фото обратно в MinIO
	objectName := fmt.Sprintf("%s_%s.jpg", photo.PhotoID, photo.Version)
	if err := uc.repo.UploadPhotoBytes(ctx, photo.BucketName, objectName, processedImage); err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("PhotoProcessing: failed to upload processed photo")
		return "", err
	}

	// Генерируем presigned URL для фронта
	url, err := uc.repo.GetPhotoURLByVersion(ctx, photo.BucketName, photo.PhotoID, photo.Version)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("objectName", objectName).Msg("PhotoProcessing: failed to generate presigned URL")
		return "", err
	}

	return url, nil
}

// GetProcessedImg - получить url по версии фото
func (uc *UseCase) GetProcessedImg(ctx context.Context, info entity.PhotoInfo) (string, error) {
	info.BucketName = uc.cfg.BucketName
	imgProcessedURL, err := uc.repo.GetPhotoURLByVersion(ctx, info.BucketName, info.PhotoID, info.Version)
	if err != nil || imgProcessedURL == "" {
		zlog.Logger.Error().Err(err).Msg("GetProcessedImg: failed to get processed url")
		return "", err
	}

	zlog.Logger.Info().Str("finally url", imgProcessedURL).Msg("GetProcessedImg: processed url")

	return imgProcessedURL, nil
}

// DeletePhoto - удалить фото
func (uc *UseCase) DeletePhoto(ctx context.Context, photo entity.PhotoInfo) error {
	// удаляем из бд
	err := uc.repo.DeletePhoto(ctx, photo.PhotoID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("DeletePhoto: failed deleting photo from database")
		return err
	}

	// удаляем фото из MinIo
	err = uc.repo.DeletePhotoFromMinIo(ctx, uc.cfg.BucketName, photo)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("DeletePhoto: failed deleting photo from minIo")
		return err
	}

	return nil
}
