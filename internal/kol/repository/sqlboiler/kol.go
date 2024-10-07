package sqlboiler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	model "kolresource/internal/db/sqlboiler"
	"kolresource/internal/kol/domain"
	"kolresource/internal/kol/domain/kol"
	"kolresource/internal/kol/entities"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/google/uuid"
)

type KolRepository struct {
	db *sql.DB
}

func NewKolRepository(db *sql.DB) *KolRepository {
	return &KolRepository{db: db}
}

func (repo *KolRepository) GetKolByID(ctx context.Context, id uuid.UUID) (*entities.Kol, error) {
	kolModel, err := model.Kols(qm.Where("id = ?", id)).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, kol.ErrDataNotFound
		}

		return nil, kol.QueryRecordError{Err: err}
	}

	return repo.newKolFromModel(kolModel)
}

func (repo *KolRepository) CreateKol(ctx context.Context, param kol.CreateKolParams) (*entities.Kol, error) {
	kolUUID, err := uuid.NewV7()
	if err != nil {
		return nil, kol.GenerateUUIDError{Err: err}
	}

	kolModel := &model.Kol{
		ID:             kolUUID.String(),
		Name:           param.Name,
		Email:          param.Email,
		Description:    param.Description,
		Sex:            model.Sex(param.Sex),
		Enable:         param.Enable,
		UpdatedAdminID: param.UpdatedAdminID.String(),
	}

	err = kolModel.Insert(ctx, repo.db, boil.Infer())
	if err != nil {
		return nil, kol.InsertRecordError{Err: err}
	}

	return repo.newKolFromModel(kolModel)
}

func (repo *KolRepository) UpdateKol(ctx context.Context, param kol.UpdateKolParams) (*entities.Kol, error) {
	kolModel, err := model.FindKol(ctx, repo.db, param.ID.String())
	if err != nil {
		return nil, kol.QueryRecordError{Err: err}
	}

	kolModel.Name = param.Name
	kolModel.Email = param.Email
	kolModel.Description = param.Description
	kolModel.Sex = model.Sex(param.Sex)
	kolModel.Enable = param.Enable
	kolModel.UpdatedAdminID = param.UpdatedAdminID.String()

	_, err = kolModel.Update(ctx, repo.db, boil.Infer())
	if err != nil {
		return nil, kol.UpdateRecordError{Err: err}
	}

	return repo.newKolFromModel(kolModel)
}

func (repo *KolRepository) DeleteKolByID(ctx context.Context, id uuid.UUID) error {
	kolModel, err := model.FindKol(ctx, repo.db, id.String())
	if err != nil {
		return kol.QueryRecordError{Err: err}
	}

	rows, err := kolModel.Delete(ctx, repo.db)
	if err != nil {
		return kol.DeleteRecordError{Err: err}
	}

	if rows == 0 {
		return kol.ErrDataNotFound
	}

	return nil
}

type KolWithTags struct {
	model.Kol `boil:",bind"`
	Tag       string `boil:"tag"`
	TagID     string `boil:"tag_id"`
}

func (repo *KolRepository) GetKolWithTagsByID(ctx context.Context, id uuid.UUID) (*kol.Kol, error) {
	var kolWithTags []KolWithTags
	err := model.NewQuery(
		qm.Select("kol.*", "tag.name as tag", "tag.id as tag_id"),
		qm.From("kol"),
		qm.InnerJoin("kol_tag ON kol_tag.kol_id = kol.id"),
		qm.InnerJoin("tag ON tag.id = kol_tag.tag_id"),
		qm.Where("kol.id = ?", id.String()),
	).Bind(ctx, repo.db, &kolWithTags)
	if err != nil {
		return nil, kol.QueryRecordError{Err: err}
	}

	if len(kolWithTags) == 0 {
		return nil, nil
	}

	kolEntity, err := repo.newKolFromModel(&kolWithTags[0].Kol)
	if err != nil {
		return nil, fmt.Errorf("failed to create kol from model: %w", err)
	}

	kolAggregate := kol.NewKol(kolEntity)

	for _, tag := range kolWithTags {
		tagUUID, err := uuid.Parse(tag.TagID)
		if err != nil {
			return nil, kol.UUIDInvalidError{Field: "tag_id", UUID: tag.TagID}
		}

		kolAggregate.AppendTag(&entities.Tag{
			ID:   tagUUID,
			Name: tag.Tag,
		})
	}

	return kolAggregate, nil
}

func (repo *KolRepository) ListKolWithTagsByFilters(ctx context.Context, param kol.ListKolWithTagsByFiltersParams) ([]*kol.Kol, int, error) {
	var kolWithTags []KolWithTags

	query := []qm.QueryMod{
		qm.Select("kol.*", "tag.name as tag", "tag.id as tag_id"),
		qm.From("kol"),
		qm.InnerJoin("kol_tag ON kol_tag.kol_id = kol.id"),
		qm.InnerJoin("tag ON tag.id = kol_tag.tag_id"),
		qm.Limit(param.PageSize),
		qm.Offset((param.Page - 1) * param.PageSize),
	}

	if param.Tag != nil {
		query = append(query, qm.Where("tag.name LIKE %?%", *param.Tag))
	}

	if param.Sex != nil {
		query = append(query, qm.Where("kol.sex = ?", *param.Sex))
	}

	if param.Email != nil {
		query = append(query, qm.Where("kol.email LIKE %?%", *param.Email))
	}

	if param.Name != nil {
		query = append(query, qm.Where("kol.name LIKE %?%", *param.Name))
	}

	err := model.NewQuery(query...).Bind(ctx, repo.db, &kolWithTags)
	if err != nil {
		return nil, 0, kol.QueryRecordError{Err: err}
	}

	count, err := repo.countKolWithTagsByFilters(ctx, param)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count kol with tags by filters: %w", err)
	}

	kols, err := repo.newKolWithTagsFromModel(kolWithTags)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create kol with tags from model: %w", err)
	}

	return kols, count, nil
}

type Count struct {
	Count int `boil:"count"`
}

func (repo *KolRepository) countKolWithTagsByFilters(ctx context.Context, param kol.ListKolWithTagsByFiltersParams) (int, error) {
	var count Count

	query := []qm.QueryMod{
		qm.Select("count(*) as count"),
		qm.From("kol"),
		qm.InnerJoin("kol_tag ON kol_tag.kol_id = kol.id"),
		qm.InnerJoin("tag ON tag.id = kol_tag.tag_id"),
		qm.GroupBy("kol.id"),
	}

	if param.Tag != nil {
		query = append(query, qm.Where("tag.name LIKE %?%", *param.Tag))
	}

	if param.Sex != nil {
		query = append(query, qm.Where("kol.sex = ?", *param.Sex))
	}

	if param.Email != nil {
		query = append(query, qm.Where("kol.email LIKE %?%", *param.Email))
	}

	if param.Name != nil {
		query = append(query, qm.Where("kol.name LIKE %?%", *param.Name))
	}

	err := model.NewQuery(query...).Bind(ctx, repo.db, &count)
	if err != nil {
		return 0, kol.QueryRecordError{Err: err}
	}

	return count.Count, nil
}

func (repo *KolRepository) CreateTag(ctx context.Context, param kol.CreateTagParams) (*entities.Tag, error) {
	tagUUID, err := uuid.NewV7()
	if err != nil {
		return nil, kol.GenerateUUIDError{Err: err}
	}

	tagModel := &model.Tag{
		ID:             tagUUID.String(),
		Name:           param.Name,
		UpdatedAdminID: param.UpdatedAdminID.String(),
	}

	err = tagModel.Insert(ctx, repo.db, boil.Infer())
	if err != nil {
		return nil, kol.InsertRecordError{Err: err}
	}

	return repo.newTagFromModel(tagModel)
}

func (repo *KolRepository) GetTagByID(ctx context.Context, id uuid.UUID) (*entities.Tag, error) {
	tagModel, err := model.Tags(qm.Where("id = ?", id)).One(ctx, repo.db)
	if err != nil {
		return nil, kol.QueryRecordError{Err: err}
	}

	return repo.newTagFromModel(tagModel)
}

func (repo *KolRepository) DeleteTagByID(ctx context.Context, id uuid.UUID) error {
	tagModel, err := model.FindTag(ctx, repo.db, id.String())
	if err != nil {
		return kol.QueryRecordError{Err: err}
	}

	rows, err := tagModel.Delete(ctx, repo.db)
	if err != nil {
		return kol.DeleteRecordError{Err: err}
	}

	if rows == 0 {
		return kol.ErrDataNotFound
	}

	return nil
}

func (repo *KolRepository) CreateProduct(ctx context.Context, param kol.CreateProductParams) (*entities.Product, error) {
	productUUID, err := uuid.NewV7()
	if err != nil {
		return nil, kol.GenerateUUIDError{Err: err}
	}

	productModel := &model.Product{
		ID:             productUUID.String(),
		Name:           param.Name,
		Description:    param.Description,
		UpdatedAdminID: param.UpdatedAdminID.String(),
	}

	err = productModel.Insert(ctx, repo.db, boil.Infer())
	if err != nil {
		return nil, kol.InsertRecordError{Err: err}
	}

	return repo.newProductFromModel(productModel)
}

func (repo *KolRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	productModel, err := model.Products(qm.Where("id = ?", id)).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, kol.ErrDataNotFound
		}

		return nil, kol.QueryRecordError{Err: err}
	}

	return repo.newProductFromModel(productModel)
}

func (repo *KolRepository) DeleteProductByID(ctx context.Context, id uuid.UUID) error {
	productModel, err := model.FindProduct(ctx, repo.db, id.String())
	if err != nil {
		return kol.QueryRecordError{Err: err}
	}

	rows, err := productModel.Delete(ctx, repo.db)
	if err != nil {
		return kol.DeleteRecordError{Err: err}
	}

	if rows == 0 {
		return kol.ErrDataNotFound
	}

	return nil
}

func (repo *KolRepository) CreateSendEmailLog(ctx context.Context, sendEmailLog *entities.SendEmailLog) (*entities.SendEmailLog, error) {
	sendEmailLogUUID, err := uuid.NewV7()
	if err != nil {
		return nil, kol.GenerateUUIDError{Err: err}
	}

	sendEmailLogModel := &model.SendEmailLog{
		ID:          sendEmailLogUUID.String(),
		KolID:       sendEmailLog.KolID.String(),
		Email:       sendEmailLog.Email,
		AdminID:     sendEmailLog.AdminID.String(),
		AdminName:   sendEmailLog.AdminName,
		ProductID:   sendEmailLog.ProductID.String(),
		ProductName: sendEmailLog.ProductName,
	}

	if err := sendEmailLogModel.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return nil, kol.InsertRecordError{Err: err}
	}

	return sendEmailLog, nil
}

func (repo *KolRepository) ListSendEmailLogsByFilter(ctx context.Context, param kol.ListSendEmailLogsByFilterParams) ([]*entities.SendEmailLog, int, error) {
	count, err := repo.countSendEmailLogsByFilter(ctx, param)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count send email logs by filter: %w", err)
	}

	if count == 0 {
		return nil, 0, nil
	}

	query := []qm.QueryMod{
		qm.Limit(param.PageSize),
		qm.Offset((param.Page - 1) * param.PageSize),
		qm.OrderBy("created_at DESC"),
	}

	if param.Email != nil {
		query = append(query, qm.Where("email LIKE %?%", *param.Email))
	}

	if param.ProductName != nil {
		query = append(query, qm.Where("product_name LIKE %?%", *param.ProductName))
	}

	if param.AdminName != nil {
		query = append(query, qm.Where("admin_id Like %?%", *param.AdminName))
	}

	if param.KolName != nil {
		query = append(query, qm.Where("kol_name LIKE %?%", *param.KolName))
	}

	modelLogs, err := model.SendEmailLogs(query...).All(ctx, repo.db)
	if err != nil {
		return nil, 0, kol.QueryRecordError{Err: err}
	}

	logs := make([]*entities.SendEmailLog, len(modelLogs))
	for index, modelLog := range modelLogs {
		modelLog := modelLog

		sendEmailLogUUID, err := uuid.Parse(modelLog.ID)
		if err != nil {
			return nil, 0, kol.UUIDInvalidError{Field: "id", UUID: modelLog.ID}
		}

		kolID, err := uuid.Parse(modelLog.KolID)
		if err != nil {
			return nil, 0, kol.UUIDInvalidError{Field: "kol_id", UUID: modelLog.KolID}
		}

		adminID, err := uuid.Parse(modelLog.AdminID)
		if err != nil {
			return nil, 0, kol.UUIDInvalidError{Field: "admin_id", UUID: modelLog.AdminID}
		}

		productID, err := uuid.Parse(modelLog.ProductID)
		if err != nil {
			return nil, 0, kol.UUIDInvalidError{Field: "product_id", UUID: modelLog.ProductID}
		}

		logs[index] = &entities.SendEmailLog{
			ID:          sendEmailLogUUID,
			KolID:       kolID,
			Email:       modelLog.Email,
			AdminID:     adminID,
			AdminName:   modelLog.AdminName,
			ProductID:   productID,
			ProductName: modelLog.ProductName,
		}
	}

	return logs, count, nil
}

func (repo *KolRepository) countSendEmailLogsByFilter(ctx context.Context, param kol.ListSendEmailLogsByFilterParams) (int, error) {
	var count Count

	query := []qm.QueryMod{
		qm.Select("count(*) as count"),
		qm.From("send_email_log"),
	}

	if param.Email != nil {
		query = append(query, qm.Where("email LIKE %?%", *param.Email))
	}

	if param.ProductName != nil {
		query = append(query, qm.Where("product_name LIKE %?%", *param.ProductName))
	}

	if param.AdminName != nil {
		query = append(query, qm.Where("admin_id Like %?%", *param.AdminName))
	}

	if param.KolName != nil {
		query = append(query, qm.Where("kol_name LIKE %?%", *param.KolName))
	}

	err := model.NewQuery(query...).Bind(ctx, repo.db, &count)
	if err != nil {
		return 0, kol.QueryRecordError{Err: err}
	}

	return count.Count, nil
}

func (repo *KolRepository) newKolFromModel(kolModel *model.Kol) (*entities.Kol, error) {
	kolUUID, err := uuid.Parse(kolModel.ID)
	if err != nil {
		return nil, kol.UUIDInvalidError{Field: "id", UUID: kolModel.ID}
	}

	sex := domain.Sex(kolModel.Sex)
	if !sex.IsValid() {
		return nil, kol.SexInvalidError{Sex: string(kolModel.Sex)}
	}

	updateAdminUUID, err := uuid.Parse(kolModel.UpdatedAdminID)
	if err != nil {
		return nil, kol.UUIDInvalidError{Field: "update_admin_id", UUID: kolModel.UpdatedAdminID}
	}

	return &entities.Kol{
		ID:             kolUUID,
		Name:           kolModel.Name,
		Email:          kolModel.Email,
		Description:    kolModel.Description,
		Sex:            sex,
		Enable:         kolModel.Enable,
		UpdatedAdminID: updateAdminUUID,
	}, nil
}

func (repo *KolRepository) newTagFromModel(tagModel *model.Tag) (*entities.Tag, error) {
	tagUUID, err := uuid.Parse(tagModel.ID)
	if err != nil {
		return nil, kol.UUIDInvalidError{Field: "id", UUID: tagModel.ID}
	}

	return &entities.Tag{
		ID:   tagUUID,
		Name: tagModel.Name,
	}, nil
}

func (repo *KolRepository) newProductFromModel(productModel *model.Product) (*entities.Product, error) {
	productUUID, err := uuid.Parse(productModel.ID)
	if err != nil {
		return nil, kol.UUIDInvalidError{Field: "id", UUID: productModel.ID}
	}

	return &entities.Product{
		ID:          productUUID,
		Name:        productModel.Name,
		Description: productModel.Description,
	}, nil
}

func (repo *KolRepository) newKolWithTagsFromModel(kolWithTags []KolWithTags) ([]*kol.Kol, error) {
	kolMap := make(map[string]*kol.Kol)
	kols := make([]*kol.Kol, 0)

	for _, kolWithTag := range kolWithTags {
		if _, ok := kolMap[kolWithTag.ID]; !ok {
			kolEntity, err := repo.newKolFromModel(&kolWithTag.Kol)
			if err != nil {
				return nil, err
			}

			kolMap[kolWithTag.ID] = kol.NewKol(kolEntity)
			kols = append(kols, kolMap[kolWithTag.ID])
		}

		tagUUID, err := uuid.Parse(kolWithTag.TagID)
		if err != nil {
			return nil, kol.UUIDInvalidError{Field: "tag_id", UUID: kolWithTag.TagID}
		}

		kolMap[kolWithTag.ID].AppendTag(&entities.Tag{
			ID:   tagUUID,
			Name: kolWithTag.Tag,
		})
	}

	return kols, nil
}
