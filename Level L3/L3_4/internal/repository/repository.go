package repository

import (
	"context"
	"database/sql"
	"image_processor/internal/config"
	"image_processor/internal/entity"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

// Repository - структура для работы с бд
type Repository struct {
	DB     *dbpg.DB
	Client *minio.Client
	master *dbpg.Options
	cfg    *config.ServiceConfig
}

// New - конструктор для репы
func New(masterDSN string, options *dbpg.Options, cfg *config.ServiceConfig) *Repository {
	masterDB, err := sql.Open("postgres", masterDSN)
	if err != nil {
		log.Fatalf("failed to open master DB: %v", err)
	}

	masterDB.SetMaxOpenConns(options.MaxOpenConns)
	masterDB.SetMaxIdleConns(options.MaxIdleConns)
	masterDB.SetConnMaxLifetime(options.ConnMaxLifetime)

	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSl,
	})
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to connect to minio")
		return nil
	}

	db := &dbpg.DB{Master: masterDB}

	return &Repository{
		DB:     db,
		Client: minioClient,
		master: options,
		cfg:    cfg,
	}
}

// AddPhotoInfo - добавить информацию о фото в бд
func (repo *Repository) AddPhotoInfo(ctx context.Context, photo entity.LoadPhoto) error {
	query := `Insert into photos (id, name, status, created_at) values ($1, $2, $3, $4)`

	_, err := repo.DB.ExecContext(ctx, query, photo.ID, photo.Name, photo.Status, photo.CreatedAt)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: AddPhotoInfo: failed to add photo info in DB")
		return err
	}

	return nil
}

func (repo *Repository) ChangeStatus(ctx context.Context, photoID, status string) error {
	query := `Update photos set status = $1 where id = $2`

	_, err := repo.DB.ExecContext(ctx, query, status, photoID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: ChangeStatus: failed to update status in DB")
		return err
	}

	return nil
}

func (repo *Repository) DeletePhoto(ctx context.Context, photoID string) error {
	query := `DELETE FROM photos WHERE id = $1`
	_, err := repo.DB.ExecContext(ctx, query, photoID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: DeletePhoto: failed to delete photo in DB")
		return err
	}

	return nil
}
