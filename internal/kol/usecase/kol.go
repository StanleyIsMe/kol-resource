package usecase

import (
	"context"
	"errors"
	"fmt"

	"kolresource/internal/kol/domain"
	"kolresource/internal/kol/domain/entities"

	"github.com/google/uuid"
)

type KolUseCaseImpl struct {
	repo      domain.Repository
	emailRepo domain.EmailRepository
}

var _ KolUseCase = (*KolUseCaseImpl)(nil)

func NewKolUseCaseImpl(repo domain.Repository, emailRepo domain.EmailRepository) *KolUseCaseImpl {
	return &KolUseCaseImpl{repo: repo, emailRepo: emailRepo}
}

// CreateKol is responsible for creating a new kol.
func (uc *KolUseCaseImpl) CreateKol(ctx context.Context, param CreateKolParam) error {
	existKol, err := uc.repo.GetKolByEmail(ctx, param.Email)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		return fmt.Errorf("repo.GetKolByEmail error: %w", err)
	}

	if existKol != nil {
		return DuplicatedResourceError{resource: "kol", name: param.Email}
	}

	createKolParams := domain.CreateKolParams{
		Name:           param.Name,
		Email:          param.Email,
		Description:    param.Description,
		SocialMedia:    param.SocialMedia,
		Sex:            param.Sex,
		Enable:         true,
		UpdatedAdminID: param.UpdatedAdminID,
		Tags:           param.Tags,
	}

	_, err = uc.repo.CreateKol(ctx, createKolParams)
	if err != nil {
		return fmt.Errorf("repo.CreateKol error: %w", err)
	}

	return nil
}

// DeleteKolByID is responsible for deleting a kol by id.
func (uc *KolUseCaseImpl) DeleteKolByID(ctx context.Context, kolID uuid.UUID) error {
	if err := uc.repo.DeleteKolByID(ctx, kolID); err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return NotFoundError{resource: "kol", id: kolID.String()}
		}

		return fmt.Errorf("repo.DeleteKolByID error: %w", err)
	}

	return nil
}

func (uc *KolUseCaseImpl) GetKolByID(ctx context.Context, kolID uuid.UUID) (*Kol, error) {
	kolAggregate, err := uc.repo.GetKolWithTagsByID(ctx, kolID)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return nil, NotFoundError{resource: "kol", id: kolID.String()}
		}

		return nil, fmt.Errorf("repo.GetKolWithTagsByID error: %w", err)
	}

	kol := &Kol{
		ID:          kolAggregate.GetKol().ID,
		Name:        kolAggregate.GetKol().Name,
		Email:       kolAggregate.GetKol().Email,
		Description: kolAggregate.GetKol().Description,
		SocialMedia: kolAggregate.GetKol().SocialMedia,
		Sex:         kolAggregate.GetKol().Sex,
		Tags:        make([]Tag, 0, len(kolAggregate.GetTags())),
	}

	for _, tag := range kolAggregate.GetTags() {
		kol.Tags = append(kol.Tags, Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return kol, nil
}

// UpdateKol is responsible for updating a kol.
func (uc *KolUseCaseImpl) UpdateKol(ctx context.Context, param UpdateKolParam) error {
	updateKolParams := domain.UpdateKolParams{
		ID:             param.KolID,
		Name:           param.Name,
		Email:          param.Email,
		Description:    param.Description,
		SocialMedia:    param.SocialMedia,
		Sex:            param.Sex,
		Tags:           param.Tags,
		UpdatedAdminID: param.UpdatedAdminID,
	}

	if _, err := uc.repo.UpdateKol(ctx, updateKolParams); err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return NotFoundError{resource: "kol", id: param.KolID.String()}
		}

		return fmt.Errorf("repo.UpdateKol error: %w", err)
	}

	return nil
}

// ListKols is responsible for searching multiple kols by dynamic filters.
func (uc *KolUseCaseImpl) ListKols(ctx context.Context, param ListKolsParam) ([]*Kol, int, error) {
	kolAggregates, total, err := uc.repo.ListKolWithTagsByFilters(ctx, domain.ListKolWithTagsByFiltersParams{
		Email:    param.Email,
		Name:     param.Name,
		Tag:      param.Tag,
		TagIDs:   param.TagIDs,
		Sex:      param.Sex,
		Page:     param.Page,
		PageSize: param.PageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo.ListKolWithTagsByFilters error: %w", err)
	}

	kols := make([]*Kol, 0, len(kolAggregates))
	for _, kol := range kolAggregates {
		kol := kol

		tags := make([]Tag, 0, len(kol.GetTags()))
		for _, tag := range kol.GetTags() {
			tags = append(tags, Tag{
				ID:   tag.ID,
				Name: tag.Name,
			})
		}

		kols = append(kols, &Kol{
			ID:          kol.GetKol().ID,
			Name:        kol.GetKol().Name,
			Email:       kol.GetKol().Email,
			Description: kol.GetKol().Description,
			SocialMedia: kol.GetKol().SocialMedia,
			Sex:         kol.GetKol().Sex,
			Tags:        tags,
		})
	}

	return kols, total, nil
}

// CreateTag is responsible for creating a new tag.
func (uc *KolUseCaseImpl) CreateTag(ctx context.Context, param CreateTagParam) error {
	existTag, err := uc.repo.GetTagByName(ctx, param.Name)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		return fmt.Errorf("repo.GetTagByName error: %w", err)
	}

	if existTag != nil {
		return DuplicatedResourceError{resource: "tag", name: param.Name}
	}

	_, err = uc.repo.CreateTag(ctx, domain.CreateTagParams{
		Name:           param.Name,
		UpdatedAdminID: param.UpdatedAdminID,
	})

	if err != nil {
		return fmt.Errorf("repo.CreateTag error: %w", err)
	}

	return nil
}

// ListTagsByName is responsible for searching multiple tags by name.
func (uc *KolUseCaseImpl) ListTagsByName(ctx context.Context, name string) ([]*Tag, error) {
	tagEntities, err := uc.repo.ListTagsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("repo.ListTagsByName error: %w", err)
	}

	tags := make([]*Tag, 0, len(tagEntities))
	for _, tag := range tagEntities {
		tags = append(tags, &Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return tags, nil
}

// CreateProduct is responsible for creating a new product.
func (uc *KolUseCaseImpl) CreateProduct(ctx context.Context, param CreateProductParam) error {
	existProduct, err := uc.repo.GetProductByName(ctx, param.Name)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		return fmt.Errorf("repo.GetProductByName error: %w", err)
	}

	if existProduct != nil {
		return DuplicatedResourceError{resource: "product", name: param.Name}
	}

	_, err = uc.repo.CreateProduct(ctx, domain.CreateProductParams{
		Name:           param.Name,
		Description:    param.Description,
		UpdatedAdminID: param.UpdatedAdminID,
	})

	if err != nil {
		return fmt.Errorf("repo.CreateProduct error: %w", err)
	}

	return nil
}

// ListProductsByName is responsible for searching multiple products by name.
func (uc *KolUseCaseImpl) ListProductsByName(ctx context.Context, name string) ([]*Product, error) {
	productEntities, err := uc.repo.ListProductsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("repo.ListProductsByName error: %w", err)
	}

	products := make([]*Product, 0, len(productEntities))
	for _, product := range productEntities {
		products = append(products, &Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
		})
	}

	return products, nil
}

// SendEmail is responsible for sending an email to multiple kols.
func (uc *KolUseCaseImpl) SendEmail(ctx context.Context, param SendEmailParam) error {
	product, err := uc.repo.GetProductByID(ctx, param.ProductID)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return NotFoundError{resource: "product", id: param.ProductID.String()}
		}

		return fmt.Errorf("repo.GetProductByID error: %w", err)
	}

	kols, err := uc.repo.ListKolsByIDs(ctx, param.KolIDs)
	if err != nil {
		return fmt.Errorf("repo.ListKolsByIDs error: %w", err)
	}

	if len(kols) == 0 {
		return NotFoundError{resource: "kol", id: param.KolIDs}
	}

	sendEmailImages := make([]domain.SendEmailImage, 0, len(param.Images))
	for _, img := range param.Images {
		sendEmailImages = append(sendEmailImages, domain.SendEmailImage{
			ContentID: img.ContentID,
			Data:      img.Data,
			ImageType: img.ImageType,
		})
	}
	sendEmailParams := domain.SendEmailParams{
		Subject:  param.Subject,
		Body:     param.EmailContent,
		ToEmails: make([]domain.ToEmail, 0, len(kols)),
		Images:   sendEmailImages,
	}

	for _, kol := range kols {
		sendEmailParams.ToEmails = append(sendEmailParams.ToEmails, domain.ToEmail{
			Email: kol.Email,
			Name:  kol.Name,
		})
	}

	if err := uc.emailRepo.SendEmail(ctx, sendEmailParams); err != nil {
		return fmt.Errorf("emailRepo.SendEmail error: %w", err)
	}

	for _, kol := range kols {
		if _, err := uc.repo.CreateSendEmailLog(ctx, &entities.SendEmailLog{
			AdminID:     param.UpdatedAdminID,
			AdminName:   param.UpdatedAdminName,
			KolID:       kol.ID,
			KolName:     kol.Name,
			Email:       kol.Email,
			ProductID:   param.ProductID,
			ProductName: product.Name,
		}); err != nil {
			return fmt.Errorf("repo.CreateSendEmailLog error: %w", err)
		}
	}

	return nil
}
