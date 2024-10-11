package usecase

import (
	"context"
	"errors"
	"fmt"
	"kolresource/internal/kol/domain"

	"github.com/google/uuid"
)

type KolUseCaseImpl struct {
	repo domain.Repository
}

func NewKolUseCase(repo domain.Repository) KolUseCase {
	return &KolUseCaseImpl{repo: repo}
}

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

func (uc *KolUseCaseImpl) GetKolByID(ctx context.Context, kolID uuid.UUID) (*Kol, error) {
	return nil, nil
}

func (uc *KolUseCaseImpl) UpdateKol(ctx context.Context, param UpdateKolParam) error {
	updateKolParams := domain.UpdateKolParams{
		ID:             param.KolID,
		Name:           param.Name,
		Email:          param.Email,
		Description:    param.Description,
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

func (uc *KolUseCaseImpl) ListKols(ctx context.Context, param ListKolsParam) ([]*Kol, int, error) {
	kolAggregates, total, err := uc.repo.ListKolWithTagsByFilters(ctx, domain.ListKolWithTagsByFiltersParams{
		Email:    param.Email,
		Name:     param.Name,
		Tag:      param.Tag,
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
			Sex:         kol.GetKol().Sex,
			Tags:        tags,
		})
	}

	return kols, total, nil
}

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

func (uc *KolUseCaseImpl) CreateProduct(ctx context.Context, param CreateProductParam) error {
	existProduct, err := uc.repo.GetProductByName(ctx, param.Name)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		return fmt.Errorf("repo.GetProductByName error: %w", err)
	}

	if existProduct != nil {
		return DuplicatedResourceError{resource: "product", name: param.Name}
	}

	_, err = uc.repo.CreateProduct(ctx, domain.CreateProductParams{
		Name:        param.Name,
		Description: param.Description,
	})

	if err != nil {
		return fmt.Errorf("repo.CreateProduct error: %w", err)
	}

	return nil
}

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

func (uc *KolUseCaseImpl) SendEmail(ctx context.Context, param SendEmailParam) error {
	return nil
}
