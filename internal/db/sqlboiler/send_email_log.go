// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// SendEmailLog is an object representing the database table.
type SendEmailLog struct {
	ID          string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	KolID       string    `boil:"kol_id" json:"kol_id" toml:"kol_id" yaml:"kol_id"`
	Email       string    `boil:"email" json:"email" toml:"email" yaml:"email"`
	AdminID     string    `boil:"admin_id" json:"admin_id" toml:"admin_id" yaml:"admin_id"`
	AdminName   string    `boil:"admin_name" json:"admin_name" toml:"admin_name" yaml:"admin_name"`
	ProductID   string    `boil:"product_id" json:"product_id" toml:"product_id" yaml:"product_id"`
	ProductName string    `boil:"product_name" json:"product_name" toml:"product_name" yaml:"product_name"`
	CreatedAt   time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *sendEmailLogR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sendEmailLogL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SendEmailLogColumns = struct {
	ID          string
	KolID       string
	Email       string
	AdminID     string
	AdminName   string
	ProductID   string
	ProductName string
	CreatedAt   string
}{
	ID:          "id",
	KolID:       "kol_id",
	Email:       "email",
	AdminID:     "admin_id",
	AdminName:   "admin_name",
	ProductID:   "product_id",
	ProductName: "product_name",
	CreatedAt:   "created_at",
}

var SendEmailLogTableColumns = struct {
	ID          string
	KolID       string
	Email       string
	AdminID     string
	AdminName   string
	ProductID   string
	ProductName string
	CreatedAt   string
}{
	ID:          "send_email_log.id",
	KolID:       "send_email_log.kol_id",
	Email:       "send_email_log.email",
	AdminID:     "send_email_log.admin_id",
	AdminName:   "send_email_log.admin_name",
	ProductID:   "send_email_log.product_id",
	ProductName: "send_email_log.product_name",
	CreatedAt:   "send_email_log.created_at",
}

// Generated where

var SendEmailLogWhere = struct {
	ID          whereHelperstring
	KolID       whereHelperstring
	Email       whereHelperstring
	AdminID     whereHelperstring
	AdminName   whereHelperstring
	ProductID   whereHelperstring
	ProductName whereHelperstring
	CreatedAt   whereHelpertime_Time
}{
	ID:          whereHelperstring{field: "\"send_email_log\".\"id\""},
	KolID:       whereHelperstring{field: "\"send_email_log\".\"kol_id\""},
	Email:       whereHelperstring{field: "\"send_email_log\".\"email\""},
	AdminID:     whereHelperstring{field: "\"send_email_log\".\"admin_id\""},
	AdminName:   whereHelperstring{field: "\"send_email_log\".\"admin_name\""},
	ProductID:   whereHelperstring{field: "\"send_email_log\".\"product_id\""},
	ProductName: whereHelperstring{field: "\"send_email_log\".\"product_name\""},
	CreatedAt:   whereHelpertime_Time{field: "\"send_email_log\".\"created_at\""},
}

// SendEmailLogRels is where relationship names are stored.
var SendEmailLogRels = struct {
}{}

// sendEmailLogR is where relationships are stored.
type sendEmailLogR struct {
}

// NewStruct creates a new relationship struct
func (*sendEmailLogR) NewStruct() *sendEmailLogR {
	return &sendEmailLogR{}
}

// sendEmailLogL is where Load methods for each relationship are stored.
type sendEmailLogL struct{}

var (
	sendEmailLogAllColumns            = []string{"id", "kol_id", "email", "admin_id", "admin_name", "product_id", "product_name", "created_at"}
	sendEmailLogColumnsWithoutDefault = []string{"kol_id", "email", "admin_id", "admin_name", "product_id", "product_name"}
	sendEmailLogColumnsWithDefault    = []string{"id", "created_at"}
	sendEmailLogPrimaryKeyColumns     = []string{"id"}
	sendEmailLogGeneratedColumns      = []string{}
)

type (
	// SendEmailLogSlice is an alias for a slice of pointers to SendEmailLog.
	// This should almost always be used instead of []SendEmailLog.
	SendEmailLogSlice []*SendEmailLog
	// SendEmailLogHook is the signature for custom SendEmailLog hook methods
	SendEmailLogHook func(context.Context, boil.ContextExecutor, *SendEmailLog) error

	sendEmailLogQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sendEmailLogType                 = reflect.TypeOf(&SendEmailLog{})
	sendEmailLogMapping              = queries.MakeStructMapping(sendEmailLogType)
	sendEmailLogPrimaryKeyMapping, _ = queries.BindMapping(sendEmailLogType, sendEmailLogMapping, sendEmailLogPrimaryKeyColumns)
	sendEmailLogInsertCacheMut       sync.RWMutex
	sendEmailLogInsertCache          = make(map[string]insertCache)
	sendEmailLogUpdateCacheMut       sync.RWMutex
	sendEmailLogUpdateCache          = make(map[string]updateCache)
	sendEmailLogUpsertCacheMut       sync.RWMutex
	sendEmailLogUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var sendEmailLogAfterSelectMu sync.Mutex
var sendEmailLogAfterSelectHooks []SendEmailLogHook

var sendEmailLogBeforeInsertMu sync.Mutex
var sendEmailLogBeforeInsertHooks []SendEmailLogHook
var sendEmailLogAfterInsertMu sync.Mutex
var sendEmailLogAfterInsertHooks []SendEmailLogHook

var sendEmailLogBeforeUpdateMu sync.Mutex
var sendEmailLogBeforeUpdateHooks []SendEmailLogHook
var sendEmailLogAfterUpdateMu sync.Mutex
var sendEmailLogAfterUpdateHooks []SendEmailLogHook

var sendEmailLogBeforeDeleteMu sync.Mutex
var sendEmailLogBeforeDeleteHooks []SendEmailLogHook
var sendEmailLogAfterDeleteMu sync.Mutex
var sendEmailLogAfterDeleteHooks []SendEmailLogHook

var sendEmailLogBeforeUpsertMu sync.Mutex
var sendEmailLogBeforeUpsertHooks []SendEmailLogHook
var sendEmailLogAfterUpsertMu sync.Mutex
var sendEmailLogAfterUpsertHooks []SendEmailLogHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SendEmailLog) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SendEmailLog) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SendEmailLog) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SendEmailLog) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SendEmailLog) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SendEmailLog) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SendEmailLog) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SendEmailLog) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SendEmailLog) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sendEmailLogAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSendEmailLogHook registers your hook function for all future operations.
func AddSendEmailLogHook(hookPoint boil.HookPoint, sendEmailLogHook SendEmailLogHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		sendEmailLogAfterSelectMu.Lock()
		sendEmailLogAfterSelectHooks = append(sendEmailLogAfterSelectHooks, sendEmailLogHook)
		sendEmailLogAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		sendEmailLogBeforeInsertMu.Lock()
		sendEmailLogBeforeInsertHooks = append(sendEmailLogBeforeInsertHooks, sendEmailLogHook)
		sendEmailLogBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		sendEmailLogAfterInsertMu.Lock()
		sendEmailLogAfterInsertHooks = append(sendEmailLogAfterInsertHooks, sendEmailLogHook)
		sendEmailLogAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		sendEmailLogBeforeUpdateMu.Lock()
		sendEmailLogBeforeUpdateHooks = append(sendEmailLogBeforeUpdateHooks, sendEmailLogHook)
		sendEmailLogBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		sendEmailLogAfterUpdateMu.Lock()
		sendEmailLogAfterUpdateHooks = append(sendEmailLogAfterUpdateHooks, sendEmailLogHook)
		sendEmailLogAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		sendEmailLogBeforeDeleteMu.Lock()
		sendEmailLogBeforeDeleteHooks = append(sendEmailLogBeforeDeleteHooks, sendEmailLogHook)
		sendEmailLogBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		sendEmailLogAfterDeleteMu.Lock()
		sendEmailLogAfterDeleteHooks = append(sendEmailLogAfterDeleteHooks, sendEmailLogHook)
		sendEmailLogAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		sendEmailLogBeforeUpsertMu.Lock()
		sendEmailLogBeforeUpsertHooks = append(sendEmailLogBeforeUpsertHooks, sendEmailLogHook)
		sendEmailLogBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		sendEmailLogAfterUpsertMu.Lock()
		sendEmailLogAfterUpsertHooks = append(sendEmailLogAfterUpsertHooks, sendEmailLogHook)
		sendEmailLogAfterUpsertMu.Unlock()
	}
}

// One returns a single sendEmailLog record from the query.
func (q sendEmailLogQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SendEmailLog, error) {
	o := &SendEmailLog{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to execute a one query for send_email_log")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all SendEmailLog records from the query.
func (q sendEmailLogQuery) All(ctx context.Context, exec boil.ContextExecutor) (SendEmailLogSlice, error) {
	var o []*SendEmailLog

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to SendEmailLog slice")
	}

	if len(sendEmailLogAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all SendEmailLog records in the query.
func (q sendEmailLogQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count send_email_log rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q sendEmailLogQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if send_email_log exists")
	}

	return count > 0, nil
}

// SendEmailLogs retrieves all the records using an executor.
func SendEmailLogs(mods ...qm.QueryMod) sendEmailLogQuery {
	mods = append(mods, qm.From("\"send_email_log\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"send_email_log\".*"})
	}

	return sendEmailLogQuery{q}
}

// FindSendEmailLog retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSendEmailLog(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*SendEmailLog, error) {
	sendEmailLogObj := &SendEmailLog{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"send_email_log\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, sendEmailLogObj)
	if err != nil {
		return nil, errors.Wrap(err, "model: unable to select from send_email_log")
	}

	if err = sendEmailLogObj.doAfterSelectHooks(ctx, exec); err != nil {
		return sendEmailLogObj, err
	}

	return sendEmailLogObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SendEmailLog) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no send_email_log provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(sendEmailLogColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	sendEmailLogInsertCacheMut.RLock()
	cache, cached := sendEmailLogInsertCache[key]
	sendEmailLogInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			sendEmailLogAllColumns,
			sendEmailLogColumnsWithDefault,
			sendEmailLogColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(sendEmailLogType, sendEmailLogMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sendEmailLogType, sendEmailLogMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"send_email_log\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"send_email_log\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into send_email_log")
	}

	if !cached {
		sendEmailLogInsertCacheMut.Lock()
		sendEmailLogInsertCache[key] = cache
		sendEmailLogInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the SendEmailLog.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SendEmailLog) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	sendEmailLogUpdateCacheMut.RLock()
	cache, cached := sendEmailLogUpdateCache[key]
	sendEmailLogUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			sendEmailLogAllColumns,
			sendEmailLogPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update send_email_log, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"send_email_log\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sendEmailLogPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sendEmailLogType, sendEmailLogMapping, append(wl, sendEmailLogPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update send_email_log row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for send_email_log")
	}

	if !cached {
		sendEmailLogUpdateCacheMut.Lock()
		sendEmailLogUpdateCache[key] = cache
		sendEmailLogUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q sendEmailLogQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for send_email_log")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for send_email_log")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SendEmailLogSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("model: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sendEmailLogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"send_email_log\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, sendEmailLogPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in sendEmailLog slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all sendEmailLog")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SendEmailLog) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no send_email_log provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(sendEmailLogColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	sendEmailLogUpsertCacheMut.RLock()
	cache, cached := sendEmailLogUpsertCache[key]
	sendEmailLogUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			sendEmailLogAllColumns,
			sendEmailLogColumnsWithDefault,
			sendEmailLogColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			sendEmailLogAllColumns,
			sendEmailLogPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert send_email_log, could not build update column list")
		}

		ret := strmangle.SetComplement(sendEmailLogAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(sendEmailLogPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert send_email_log, could not build conflict column list")
			}

			conflict = make([]string, len(sendEmailLogPrimaryKeyColumns))
			copy(conflict, sendEmailLogPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"send_email_log\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(sendEmailLogType, sendEmailLogMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sendEmailLogType, sendEmailLogMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "model: unable to upsert send_email_log")
	}

	if !cached {
		sendEmailLogUpsertCacheMut.Lock()
		sendEmailLogUpsertCache[key] = cache
		sendEmailLogUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single SendEmailLog record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SendEmailLog) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no SendEmailLog provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sendEmailLogPrimaryKeyMapping)
	sql := "DELETE FROM \"send_email_log\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from send_email_log")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for send_email_log")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q sendEmailLogQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no sendEmailLogQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from send_email_log")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for send_email_log")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SendEmailLogSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(sendEmailLogBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sendEmailLogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"send_email_log\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sendEmailLogPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from sendEmailLog slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for send_email_log")
	}

	if len(sendEmailLogAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SendEmailLog) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSendEmailLog(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SendEmailLogSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SendEmailLogSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sendEmailLogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"send_email_log\".* FROM \"send_email_log\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sendEmailLogPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in SendEmailLogSlice")
	}

	*o = slice

	return nil
}

// SendEmailLogExists checks if the SendEmailLog row exists.
func SendEmailLogExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"send_email_log\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if send_email_log exists")
	}

	return exists, nil
}

// Exists checks if the SendEmailLog row exists.
func (o *SendEmailLog) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SendEmailLogExists(ctx, exec, o.ID)
}