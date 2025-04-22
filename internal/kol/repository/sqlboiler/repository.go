package sqlboiler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	model "kolresource/internal/db/sqlboiler"
	"kolresource/internal/kol"
	"kolresource/internal/kol/domain"
	"kolresource/internal/kol/domain/entities"

	"github.com/volatiletech/null/v9"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/google/uuid"
)

type KolRepository struct {
	db *sql.DB
}

var _ domain.Repository = (*KolRepository)(nil)

func NewKolRepository(db *sql.DB) *KolRepository {
	return &KolRepository{db: db}
}

func (repo *KolRepository) GetKolByID(ctx context.Context, id uuid.UUID) (*entities.Kol, error) {
	kolModel, err := model.Kols(qm.Where("id = ?", id.String())).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, domain.QueryRecordError{Err: err}
	}

	return repo.newKolFromModel(kolModel)
}

func (repo *KolRepository) GetKolByEmail(ctx context.Context, email string) (*entities.Kol, error) {
	kolModel, err := model.Kols(qm.Where("email = ?", email)).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, domain.QueryRecordError{Err: err}
	}

	return repo.newKolFromModel(kolModel)
}

func (repo *KolRepository) CreateKol(ctx context.Context, param domain.CreateKolParams) (*entities.Kol, error) {
	kolUUID, err := uuid.NewV7()
	if err != nil {
		return nil, domain.GenerateUUIDError{Err: err}
	}

	kolModel := &model.Kol{
		ID:             kolUUID.String(),
		Name:           param.Name,
		Email:          param.Email,
		Description:    param.Description,
		SocialMedia:    param.SocialMedia,
		Sex:            model.Sex(param.Sex),
		Enable:         param.Enable,
		UpdatedAdminID: param.UpdatedAdminID.String(),
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var txErr error

	defer func() {
		if txErr != nil {
			tx.Rollback() //nolint:errcheck
		}
	}()

	if txErr = kolModel.Insert(ctx, tx, boil.Infer()); txErr != nil {
		return nil, domain.InsertRecordError{Err: err}
	}

	for _, tagID := range param.Tags {
		kolTagModel := &model.KolTag{
			KolID:          kolModel.ID,
			TagID:          tagID.String(),
			UpdatedAdminID: param.UpdatedAdminID.String(),
		}

		if txErr = kolTagModel.Insert(ctx, tx, boil.Infer()); txErr != nil {
			return nil, domain.InsertRecordError{Err: err}
		}
	}

	if txErr = tx.Commit(); txErr != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return repo.newKolFromModel(kolModel)
}

func (repo *KolRepository) UpdateKol(ctx context.Context, param domain.UpdateKolParams) (*entities.Kol, error) {
	kolModel, err := model.Kols(qm.Where("id = ?", param.ID.String())).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, domain.QueryRecordError{Err: err}
	}

	kolModel.Name = param.Name
	kolModel.Email = param.Email
	kolModel.Description = param.Description
	kolModel.SocialMedia = param.SocialMedia
	kolModel.Sex = model.Sex(param.Sex)
	kolModel.Enable = param.Enable
	kolModel.UpdatedAdminID = param.UpdatedAdminID.String()

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var txErr error

	defer func() {
		if txErr != nil {
			tx.Rollback() //nolint:errcheck
		}
	}()

	_, txErr = kolModel.Update(ctx, tx, boil.Infer())
	if txErr != nil {
		return nil, domain.UpdateRecordError{Err: err}
	}

	if len(param.Tags) > 0 {
		_, txErr = model.KolTags(qm.Where("kol_id = ?", kolModel.ID)).DeleteAll(ctx, tx)
		if txErr != nil {
			return nil, domain.DeleteRecordError{Err: err}
		}
	}

	for _, tagID := range param.Tags {
		kolTagModel := &model.KolTag{
			KolID:          kolModel.ID,
			TagID:          tagID.String(),
			UpdatedAdminID: param.UpdatedAdminID.String(),
		}

		txErr = kolTagModel.Insert(ctx, tx, boil.Infer())
		if txErr != nil {
			return nil, domain.InsertRecordError{Err: err}
		}
	}

	if txErr = tx.Commit(); txErr != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return repo.newKolFromModel(kolModel)
}

func (repo *KolRepository) DeleteKolByID(ctx context.Context, id uuid.UUID) error {
	kolModel, err := model.Kols(qm.Where("id = ?", id.String())).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrDataNotFound
		}

		return domain.QueryRecordError{Err: err}
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var txErr error

	defer func() {
		if txErr != nil {
			tx.Rollback() //nolint:errcheck
		}
	}()

	rows, txErr := kolModel.Delete(ctx, tx)
	if txErr != nil {
		return domain.DeleteRecordError{Err: err}
	}

	if rows == 0 {
		return domain.ErrDataNotFound
	}

	if _, txErr = model.KolTags(qm.Where("kol_id = ?", kolModel.ID)).DeleteAll(ctx, tx); txErr != nil {
		return domain.DeleteRecordError{Err: err}
	}

	if txErr = tx.Commit(); txErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type KolWithTags struct {
	model.Kol `boil:",bind"`
	Tag       null.String `boil:"tag"`
	TagID     null.String `boil:"tag_id"`
}

func (repo *KolRepository) GetKolWithTagsByID(ctx context.Context, id uuid.UUID) (*domain.Kol, error) {
	var kolWithTags []KolWithTags
	err := model.NewQuery(
		qm.Select("kol.*", "tag.name as tag", "tag.id as tag_id"),
		qm.From("kol"),
		qm.InnerJoin("kol_tag ON kol_tag.kol_id = kol.id"),
		qm.InnerJoin("tag ON tag.id = kol_tag.tag_id"),
		qm.Where("kol.id = ?", id.String()),
	).Bind(ctx, repo.db, &kolWithTags)
	if err != nil {
		return nil, domain.QueryRecordError{Err: err}
	}

	if len(kolWithTags) == 0 {
		return nil, nil
	}

	kolEntity, err := repo.newKolFromModel(&kolWithTags[0].Kol)
	if err != nil {
		return nil, fmt.Errorf("failed to create kol from model: %w", err)
	}

	kolAggregate := domain.NewKol(kolEntity)

	for _, tag := range kolWithTags {
		tagUUID, err := uuid.Parse(tag.TagID.String)
		if err != nil {
			return nil, domain.UUIDInvalidError{Field: "tag_id", UUID: tag.TagID.String}
		}

		kolAggregate.AppendTag(&entities.Tag{
			ID:   tagUUID,
			Name: tag.Tag.String,
		})
	}

	return kolAggregate, nil
}

func (repo *KolRepository) ListKolsByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.Kol, error) {
	idArr := make([]interface{}, len(ids))
	for index, id := range ids {
		idArr[index] = id.String()
	}

	query := qm.WhereIn("id IN ?", idArr...)
	models, err := model.Kols(query).All(ctx, repo.db)
	if err != nil {
		return nil, domain.QueryRecordError{Err: err}
	}

	entityKols := make([]*entities.Kol, len(models))
	for index, model := range models {
		model := model
		entityKols[index], err = repo.newKolFromModel(model)
		if err != nil {
			return nil, fmt.Errorf("failed to create kol from model: %w", err)
		}
	}

	return entityKols, nil
}

func (repo *KolRepository) ListKolWithTagsByFilters(ctx context.Context, param domain.ListKolWithTagsByFiltersParams) ([]*domain.Kol, int, error) {
	var kolWithTags []KolWithTags

	// due to the sqlboiler was not friendly for sub query, we need to use sub query to get kol ids first
	subQuery := []qm.QueryMod{
		qm.Select("kol.id as id"),
		qm.LeftOuterJoin("kol_tag ON kol_tag.kol_id = kol.id"),
		qm.LeftOuterJoin("tag ON tag.id = kol_tag.tag_id"),
	}

	if len(param.TagIDs) > 0 {
		tagIDs := make([]interface{}, len(param.TagIDs))
		for index, tagID := range param.TagIDs {
			tagIDs[index] = tagID.String()
		}

		subQuery = append(subQuery, qm.WhereIn("tag.id IN ?", tagIDs...))
	}

	if param.Tag != nil {
		subQuery = append(subQuery, qm.Where("tag.name LIKE ?", "%"+*param.Tag+"%"))
	}

	if param.Sex != nil {
		subQuery = append(subQuery, qm.Where("kol.sex = ?", *param.Sex))
	}

	if param.Email != nil {
		subQuery = append(subQuery, qm.Where("kol.email LIKE ?", "%"+*param.Email+"%"))
	}

	if param.Name != nil {
		subQuery = append(subQuery, qm.Where("kol.name LIKE ?", "%"+*param.Name+"%"))
	}

	subQuery = append(subQuery,
		qm.GroupBy("kol.id"),
		qm.Limit(param.PageSize),
		qm.Offset((param.Page-1)*param.PageSize),
	)

	subKols, err := model.Kols(subQuery...).All(ctx, repo.db)
	if err != nil {
		return nil, 0, domain.QueryRecordError{Err: err}
	}

	subQueryKolIDs := make([]interface{}, len(subKols))
	for index, kol := range subKols {
		subQueryKolIDs[index] = kol.ID
	}

	mainQuery := []qm.QueryMod{
		qm.Select("kol.*", "tag.name as tag", "tag.id as tag_id"),
		qm.From("kol"),
		qm.LeftOuterJoin("kol_tag ON kol_tag.kol_id = kol.id"),
		qm.LeftOuterJoin("tag ON tag.id = kol_tag.tag_id"),
		qm.WhereIn("kol.id IN ?", subQueryKolIDs...),
	}

	if err := model.NewQuery(mainQuery...).Bind(ctx, repo.db, &kolWithTags); err != nil {
		return nil, 0, domain.QueryRecordError{Err: err}
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

func (repo *KolRepository) countKolWithTagsByFilters(ctx context.Context, param domain.ListKolWithTagsByFiltersParams) (int, error) {
	var count Count

	query := []qm.QueryMod{
		qm.Select("count(distinct kol.id) as count"),
		qm.From("kol"),
		qm.LeftOuterJoin("kol_tag ON kol_tag.kol_id = kol.id"),
		qm.LeftOuterJoin("tag ON tag.id = kol_tag.tag_id"),
	}

	if len(param.TagIDs) > 0 {
		tagIDs := make([]interface{}, len(param.TagIDs))
		for index, tagID := range param.TagIDs {
			tagIDs[index] = tagID.String()
		}

		query = append(query, qm.WhereIn("tag.id IN ?", tagIDs...))
	}

	if param.Tag != nil {
		query = append(query, qm.Where("tag.name LIKE ?", "%"+*param.Tag+"%"))
	}

	if param.Sex != nil {
		query = append(query, qm.Where("kol.sex = ?", *param.Sex))
	}

	if param.Email != nil {
		query = append(query, qm.Where("kol.email LIKE ?", "%"+*param.Email+"%"))
	}

	if param.Name != nil {
		query = append(query, qm.Where("kol.name LIKE ?", "%"+*param.Name+"%"))
	}

	err := model.NewQuery(query...).Bind(ctx, repo.db, &count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, domain.QueryRecordError{Err: err}
	}

	return count.Count, nil
}

func (repo *KolRepository) CreateTag(ctx context.Context, param domain.CreateTagParams) (*entities.Tag, error) {
	tagUUID, err := uuid.NewV7()
	if err != nil {
		return nil, domain.GenerateUUIDError{Err: err}
	}

	tagModel := &model.Tag{
		ID:             tagUUID.String(),
		Name:           param.Name,
		UpdatedAdminID: param.UpdatedAdminID.String(),
	}

	err = tagModel.Insert(ctx, repo.db, boil.Infer())
	if err != nil {
		return nil, domain.InsertRecordError{Err: err}
	}

	return repo.newTagFromModel(tagModel)
}

func (repo *KolRepository) GetTagByID(ctx context.Context, id uuid.UUID) (*entities.Tag, error) {
	tagModel, err := model.Tags(qm.Where("id = ?", id)).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, domain.QueryRecordError{Err: err}
	}

	return repo.newTagFromModel(tagModel)
}

func (repo *KolRepository) GetTagByName(ctx context.Context, name string) (*entities.Tag, error) {
	tagModel, err := model.Tags(qm.Where("name = ?", name)).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, domain.QueryRecordError{Err: err}
	}

	return repo.newTagFromModel(tagModel)
}

func (repo *KolRepository) DeleteTagByID(ctx context.Context, id uuid.UUID) error {
	tagModel, err := model.FindTag(ctx, repo.db, id.String())
	if err != nil {
		return domain.QueryRecordError{Err: err}
	}

	rows, err := tagModel.Delete(ctx, repo.db)
	if err != nil {
		return domain.DeleteRecordError{Err: err}
	}

	if rows == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}

func (repo *KolRepository) ListTagsByName(ctx context.Context, name string) ([]*entities.Tag, error) {
	tagModels, err := model.Tags(qm.Where("name LIKE ?", "%"+name+"%")).All(ctx, repo.db)
	if err != nil {
		return nil, domain.QueryRecordError{Err: err}
	}

	tags := make([]*entities.Tag, len(tagModels))
	for index, tagModel := range tagModels {
		tags[index], err = repo.newTagFromModel(tagModel)
		if err != nil {
			return nil, fmt.Errorf("failed to create tag from model: %w", err)
		}
	}

	return tags, nil
}

func (repo *KolRepository) CreateProduct(ctx context.Context, param domain.CreateProductParams) (*entities.Product, error) {
	productUUID, err := uuid.NewV7()
	if err != nil {
		return nil, domain.GenerateUUIDError{Err: err}
	}

	productModel := &model.Product{
		ID:             productUUID.String(),
		Name:           param.Name,
		Description:    param.Description,
		UpdatedAdminID: param.UpdatedAdminID.String(),
	}

	err = productModel.Insert(ctx, repo.db, boil.Infer())
	if err != nil {
		return nil, domain.InsertRecordError{Err: err}
	}

	return repo.newProductFromModel(productModel)
}

func (repo *KolRepository) ListProductsByName(ctx context.Context, name string) ([]*entities.Product, error) {
	productModels, err := model.Products(qm.Where("name LIKE ?", "%"+name+"%")).All(ctx, repo.db)
	if err != nil {
		return nil, domain.QueryRecordError{Err: err}
	}

	products := make([]*entities.Product, len(productModels))
	for index, productModel := range productModels {
		products[index], err = repo.newProductFromModel(productModel)
		if err != nil {
			return nil, fmt.Errorf("failed to create product from model: %w", err)
		}
	}

	return products, nil
}

func (repo *KolRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	productModel, err := model.Products(qm.Where("id = ?", id)).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, domain.QueryRecordError{Err: err}
	}

	return repo.newProductFromModel(productModel)
}

func (repo *KolRepository) GetProductByName(ctx context.Context, name string) (*entities.Product, error) {
	productModel, err := model.Products(qm.Where("name = ?", name)).One(ctx, repo.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, domain.QueryRecordError{Err: err}
	}

	return repo.newProductFromModel(productModel)
}

func (repo *KolRepository) DeleteProductByID(ctx context.Context, id uuid.UUID) error {
	productModel, err := model.FindProduct(ctx, repo.db, id.String())
	if err != nil {
		return domain.QueryRecordError{Err: err}
	}

	rows, err := productModel.Delete(ctx, repo.db)
	if err != nil {
		return domain.DeleteRecordError{Err: err}
	}

	if rows == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}

func (repo *KolRepository) CreateSendEmailLog(ctx context.Context, sendEmailLog *entities.SendEmailLog) (*entities.SendEmailLog, error) {
	sendEmailLogUUID, err := uuid.NewV7()
	if err != nil {
		return nil, domain.GenerateUUIDError{Err: err}
	}

	sendEmailLogModel := &model.SendEmailLog{
		ID:          sendEmailLogUUID.String(),
		KolID:       sendEmailLog.KolID.String(),
		KolName:     sendEmailLog.KolName,
		Email:       sendEmailLog.Email,
		AdminID:     sendEmailLog.AdminID.String(),
		AdminName:   sendEmailLog.AdminName,
		ProductID:   sendEmailLog.ProductID.String(),
		ProductName: sendEmailLog.ProductName,
	}

	if err := sendEmailLogModel.Insert(ctx, repo.db, boil.Infer()); err != nil {
		return nil, domain.InsertRecordError{Err: err}
	}

	return sendEmailLog, nil
}

func (repo *KolRepository) ListSendEmailLogsByFilter(ctx context.Context, param domain.ListSendEmailLogsByFilterParams) ([]*entities.SendEmailLog, int, error) {
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
		query = append(query, qm.Where("email LIKE ?", "%"+*param.Email+"%"))
	}

	if param.ProductName != nil {
		query = append(query, qm.Where("product_name LIKE ?", "%"+*param.ProductName+"%"))
	}

	if param.AdminName != nil {
		query = append(query, qm.Where("admin_id Like ?", "%"+*param.AdminName+"%"))
	}

	if param.KolName != nil {
		query = append(query, qm.Where("kol_name LIKE ?", "%"+*param.KolName+"%"))
	}

	modelLogs, err := model.SendEmailLogs(query...).All(ctx, repo.db)
	if err != nil {
		return nil, 0, domain.QueryRecordError{Err: err}
	}

	logs := make([]*entities.SendEmailLog, len(modelLogs))
	for index, modelLog := range modelLogs {
		modelLog := modelLog

		sendEmailLogUUID, err := uuid.Parse(modelLog.ID)
		if err != nil {
			return nil, 0, domain.UUIDInvalidError{Field: "id", UUID: modelLog.ID}
		}

		kolID, err := uuid.Parse(modelLog.KolID)
		if err != nil {
			return nil, 0, domain.UUIDInvalidError{Field: "kol_id", UUID: modelLog.KolID}
		}

		adminID, err := uuid.Parse(modelLog.AdminID)
		if err != nil {
			return nil, 0, domain.UUIDInvalidError{Field: "admin_id", UUID: modelLog.AdminID}
		}

		productID, err := uuid.Parse(modelLog.ProductID)
		if err != nil {
			return nil, 0, domain.UUIDInvalidError{Field: "product_id", UUID: modelLog.ProductID}
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

func (repo *KolRepository) countSendEmailLogsByFilter(ctx context.Context, param domain.ListSendEmailLogsByFilterParams) (int, error) {
	var count Count

	query := []qm.QueryMod{
		qm.Select("count(*) as count"),
		qm.From("send_email_log"),
	}

	if param.Email != nil {
		query = append(query, qm.Where("email LIKE ?", "%"+*param.Email+"%"))
	}

	if param.ProductName != nil {
		query = append(query, qm.Where("product_name LIKE ?", "%"+*param.ProductName+"%"))
	}

	if param.AdminName != nil {
		query = append(query, qm.Where("admin_id Like ?", "%"+*param.AdminName+"%"))
	}

	if param.KolName != nil {
		query = append(query, qm.Where("kol_name LIKE ?", "%"+*param.KolName+"%"))
	}

	err := model.NewQuery(query...).Bind(ctx, repo.db, &count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, domain.QueryRecordError{Err: err}
	}

	return count.Count, nil
}

func (repo *KolRepository) newKolFromModel(kolModel *model.Kol) (*entities.Kol, error) {
	kolUUID, err := uuid.Parse(kolModel.ID)
	if err != nil {
		return nil, domain.UUIDInvalidError{Field: "id", UUID: kolModel.ID}
	}

	sex := kol.Sex(kolModel.Sex)
	if !sex.IsValid() {
		return nil, domain.SexInvalidError{Sex: string(kolModel.Sex)}
	}

	updateAdminUUID, err := uuid.Parse(kolModel.UpdatedAdminID)
	if err != nil {
		return nil, domain.UUIDInvalidError{Field: "update_admin_id", UUID: kolModel.UpdatedAdminID}
	}

	return &entities.Kol{
		ID:             kolUUID,
		Name:           kolModel.Name,
		Email:          kolModel.Email,
		Description:    kolModel.Description,
		SocialMedia:    kolModel.SocialMedia,
		Sex:            sex,
		Enable:         kolModel.Enable,
		UpdatedAdminID: updateAdminUUID,
	}, nil
}

func (repo *KolRepository) newTagFromModel(tagModel *model.Tag) (*entities.Tag, error) {
	tagUUID, err := uuid.Parse(tagModel.ID)
	if err != nil {
		return nil, domain.UUIDInvalidError{Field: "id", UUID: tagModel.ID}
	}

	return &entities.Tag{
		ID:   tagUUID,
		Name: tagModel.Name,
	}, nil
}

func (repo *KolRepository) newProductFromModel(productModel *model.Product) (*entities.Product, error) {
	productUUID, err := uuid.Parse(productModel.ID)
	if err != nil {
		return nil, domain.UUIDInvalidError{Field: "id", UUID: productModel.ID}
	}

	return &entities.Product{
		ID:          productUUID,
		Name:        productModel.Name,
		Description: productModel.Description,
	}, nil
}

func (repo *KolRepository) newKolWithTagsFromModel(kolWithTags []KolWithTags) ([]*domain.Kol, error) {
	kolMap := make(map[string]*domain.Kol)
	kols := make([]*domain.Kol, 0)

	for _, kolWithTag := range kolWithTags {
		if _, ok := kolMap[kolWithTag.ID]; !ok {
			kolEntity, err := repo.newKolFromModel(&kolWithTag.Kol)
			if err != nil {
				return nil, err
			}

			kolMap[kolWithTag.ID] = domain.NewKol(kolEntity)
			kols = append(kols, kolMap[kolWithTag.ID])
		}

		tagUUID, err := uuid.Parse(kolWithTag.TagID.String)
		if err != nil {
			return nil, domain.UUIDInvalidError{Field: "tag_id", UUID: kolWithTag.TagID.String}
		}

		kolMap[kolWithTag.ID].AppendTag(&entities.Tag{
			ID:   tagUUID,
			Name: kolWithTag.Tag.String,
		})
	}

	return kols, nil
}
