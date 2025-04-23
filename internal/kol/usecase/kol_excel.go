package usecase

import (
	"context"
	"errors"
	"fmt"
	"kolresource/internal/kol"
	"kolresource/internal/kol/domain"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

const (
	rowColumnIndexTag = iota
	rowColumnIndexName
	rowColumnIndexSex
	rowColumnIndexSocialMedia
	rowColumnIndexEmail
)

type ExcelRawData struct {
	Tags        []string `validate:"required"`
	Name        string   `validate:"required,lte=50"`
	Sex         kol.Sex  `validate:"oneof=m f"`
	SocialMedia string   `validate:"gte=0,lte=255"`
	Email       string   `validate:"required,email"`
}

// BatchCreateKolsByXlsx is responsible for creating multiple kols from an excel file.
// TODO: need collect error per row data and log it
func (uc *KolUseCaseImpl) BatchCreateKolsByXlsx(ctx context.Context, param BatchCreateKolsByXlsxParam) error {
	uploadFile, err := param.File.Open()
	if err != nil {
		return fmt.Errorf("file.Open error: %w", err)
	}

	defer uploadFile.Close()

	xlsxFile, err := excelize.OpenReader(uploadFile)
	if err != nil {
		return fmt.Errorf("excelize.OpenReader error: %w", err)
	}

	defer xlsxFile.Close()

	if len(xlsxFile.GetSheetList()) != 1 {
		return fmt.Errorf("xlsxFile.GetSheetList error: %w", err)
	}

	sheetName := xlsxFile.GetSheetList()[0]

	rows, err := xlsxFile.Rows(sheetName)
	if err != nil {
		return fmt.Errorf("xlsxFile.Rows error: %w", err)
	}

	defer rows.Close()

	tempTags := make(map[string]uuid.UUID)
	validate := validator.New()

	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			fmt.Println(err) //nolint:forbidigo

			continue
		}

		if len(cols) < rowColumnIndexEmail {
			fmt.Println("cols length is less than rowColumnIndexEmail") //nolint:forbidigo

			continue
		}

		rawData := &ExcelRawData{
			Tags:        strings.Split(strings.TrimSpace(cols[rowColumnIndexTag]), "/"),
			Name:        strings.TrimSpace(cols[rowColumnIndexName]),
			Sex:         kol.Sex(strings.TrimSpace(cols[rowColumnIndexSex])),
			SocialMedia: strings.TrimSpace(cols[rowColumnIndexSocialMedia]),
			Email:       strings.TrimSpace(cols[rowColumnIndexEmail]),
		}

		if err := validate.Struct(rawData); err != nil {
			fmt.Println(err) //nolint:forbidigo

			continue
		}

		tagIDs := make([]uuid.UUID, 0, len(rawData.Tags))
		for _, tag := range rawData.Tags {
			if tagID, ok := tempTags[tag]; ok {
				tagIDs = append(tagIDs, tagID)

				continue
			}

			existTag, err := uc.repo.GetTagByName(ctx, tag)
			if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
				return fmt.Errorf("repo.GetTagByName error: %w", err)
			}

			if existTag != nil {
				tagIDs = append(tagIDs, existTag.ID)
				tempTags[tag] = existTag.ID

				continue
			}

			tagEntity, err := uc.repo.CreateTag(ctx, domain.CreateTagParams{
				Name:           tag,
				UpdatedAdminID: param.UpdatedAdminID,
			})
			if err != nil {
				return fmt.Errorf("uc.repo.CreateTag error: %w", err)
			}

			tagIDs = append(tagIDs, tagEntity.ID)
			tempTags[tag] = tagEntity.ID
		}

		if err := uc.upsertKol(ctx, UpsertKolParam{
			Tags:           tagIDs,
			Name:           rawData.Name,
			Sex:            rawData.Sex,
			SocialMedia:    rawData.SocialMedia,
			Email:          rawData.Email,
			UpdatedAdminID: param.UpdatedAdminID,
		}); err != nil {
			fmt.Println(err) //nolint:forbidigo

			continue
		}
	}

	return nil
}

type UpsertKolParam struct {
	Tags           []uuid.UUID
	Name           string
	Sex            kol.Sex
	SocialMedia    string
	Email          string
	UpdatedAdminID uuid.UUID
}

// upsertKol is responsible for upserting a kol.
func (uc *KolUseCaseImpl) upsertKol(ctx context.Context, param UpsertKolParam) error {
	existKol, err := uc.repo.GetKolByEmail(ctx, param.Email)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		return fmt.Errorf("repo.GetKolByEmail error: %w", err)
	}

	if existKol != nil {
		if _, err := uc.repo.UpdateKol(ctx, domain.UpdateKolParams{
			ID:             existKol.ID,
			Email:          param.Email,
			Name:           param.Name,
			Sex:            param.Sex,
			SocialMedia:    param.SocialMedia,
			Enable:         true,
			UpdatedAdminID: param.UpdatedAdminID,
			Tags:           param.Tags,
		}); err != nil {
			return fmt.Errorf("uc.repo.UpdateKol error: %w", err)
		}

		return nil
	}

	if _, err := uc.repo.CreateKol(ctx, domain.CreateKolParams{
		Email:          param.Email,
		Name:           param.Name,
		Sex:            param.Sex,
		SocialMedia:    param.SocialMedia,
		Enable:         true,
		UpdatedAdminID: param.UpdatedAdminID,
		Tags:           param.Tags,
	}); err != nil {
		return fmt.Errorf("uc.repo.CreateKol error: %w", err)
	}

	return nil
}
