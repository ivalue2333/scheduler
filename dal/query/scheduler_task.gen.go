// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/ivalue2333/scheduler/dal/rdm"
)

func newSchedulerTask(db *gorm.DB, opts ...gen.DOOption) schedulerTask {
	_schedulerTask := schedulerTask{}

	_schedulerTask.schedulerTaskDo.UseDB(db, opts...)
	_schedulerTask.schedulerTaskDo.UseModel(&rdm.SchedulerTask{})

	tableName := _schedulerTask.schedulerTaskDo.TableName()
	_schedulerTask.ALL = field.NewAsterisk(tableName)
	_schedulerTask.ID = field.NewInt64(tableName, "id")
	_schedulerTask.InstanceID = field.NewInt64(tableName, "instance_id")
	_schedulerTask.Name = field.NewString(tableName, "name")
	_schedulerTask.Namespace = field.NewString(tableName, "namespace")
	_schedulerTask.Category = field.NewString(tableName, "category")
	_schedulerTask.Priority = field.NewInt32(tableName, "priority")
	_schedulerTask.Status = field.NewInt32(tableName, "status")
	_schedulerTask.Extra = field.NewString(tableName, "extra")
	_schedulerTask.CreatedAt = field.NewTime(tableName, "created_at")
	_schedulerTask.UpdatedAt = field.NewTime(tableName, "updated_at")

	_schedulerTask.fillFieldMap()

	return _schedulerTask
}

// schedulerTask 任务表
type schedulerTask struct {
	schedulerTaskDo

	ALL        field.Asterisk
	ID         field.Int64  // 主键
	InstanceID field.Int64  // 执行实例
	Name       field.String // 任务名,绑定执行器
	Namespace  field.String // 实例空间,c0
	Category   field.String // 实例空间分类,c1,推荐拼接
	Priority   field.Int32  // 任务优先级, 1 - 100, 1:数值越小,优先级越高
	Status     field.Int32  // 状态; 1:等待中,2:执行中,3:暂停,4:停止,5:成功,6:过期
	Extra      field.String // 业务数据
	CreatedAt  field.Time   // 创建时间
	UpdatedAt  field.Time   // 更新时间

	fieldMap map[string]field.Expr
}

func (s schedulerTask) Table(newTableName string) *schedulerTask {
	s.schedulerTaskDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s schedulerTask) As(alias string) *schedulerTask {
	s.schedulerTaskDo.DO = *(s.schedulerTaskDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *schedulerTask) updateTableName(table string) *schedulerTask {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt64(table, "id")
	s.InstanceID = field.NewInt64(table, "instance_id")
	s.Name = field.NewString(table, "name")
	s.Namespace = field.NewString(table, "namespace")
	s.Category = field.NewString(table, "category")
	s.Priority = field.NewInt32(table, "priority")
	s.Status = field.NewInt32(table, "status")
	s.Extra = field.NewString(table, "extra")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")

	s.fillFieldMap()

	return s
}

func (s *schedulerTask) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *schedulerTask) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 10)
	s.fieldMap["id"] = s.ID
	s.fieldMap["instance_id"] = s.InstanceID
	s.fieldMap["name"] = s.Name
	s.fieldMap["namespace"] = s.Namespace
	s.fieldMap["category"] = s.Category
	s.fieldMap["priority"] = s.Priority
	s.fieldMap["status"] = s.Status
	s.fieldMap["extra"] = s.Extra
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
}

func (s schedulerTask) clone(db *gorm.DB) schedulerTask {
	s.schedulerTaskDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s schedulerTask) replaceDB(db *gorm.DB) schedulerTask {
	s.schedulerTaskDo.ReplaceDB(db)
	return s
}

type schedulerTaskDo struct{ gen.DO }

type ISchedulerTaskDo interface {
	gen.SubQuery
	Debug() ISchedulerTaskDo
	WithContext(ctx context.Context) ISchedulerTaskDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ISchedulerTaskDo
	WriteDB() ISchedulerTaskDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ISchedulerTaskDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ISchedulerTaskDo
	Not(conds ...gen.Condition) ISchedulerTaskDo
	Or(conds ...gen.Condition) ISchedulerTaskDo
	Select(conds ...field.Expr) ISchedulerTaskDo
	Where(conds ...gen.Condition) ISchedulerTaskDo
	Order(conds ...field.Expr) ISchedulerTaskDo
	Distinct(cols ...field.Expr) ISchedulerTaskDo
	Omit(cols ...field.Expr) ISchedulerTaskDo
	Join(table schema.Tabler, on ...field.Expr) ISchedulerTaskDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ISchedulerTaskDo
	RightJoin(table schema.Tabler, on ...field.Expr) ISchedulerTaskDo
	Group(cols ...field.Expr) ISchedulerTaskDo
	Having(conds ...gen.Condition) ISchedulerTaskDo
	Limit(limit int) ISchedulerTaskDo
	Offset(offset int) ISchedulerTaskDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ISchedulerTaskDo
	Unscoped() ISchedulerTaskDo
	Create(values ...*rdm.SchedulerTask) error
	CreateInBatches(values []*rdm.SchedulerTask, batchSize int) error
	Save(values ...*rdm.SchedulerTask) error
	First() (*rdm.SchedulerTask, error)
	Take() (*rdm.SchedulerTask, error)
	Last() (*rdm.SchedulerTask, error)
	Find() ([]*rdm.SchedulerTask, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*rdm.SchedulerTask, err error)
	FindInBatches(result *[]*rdm.SchedulerTask, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*rdm.SchedulerTask) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ISchedulerTaskDo
	Assign(attrs ...field.AssignExpr) ISchedulerTaskDo
	Joins(fields ...field.RelationField) ISchedulerTaskDo
	Preload(fields ...field.RelationField) ISchedulerTaskDo
	FirstOrInit() (*rdm.SchedulerTask, error)
	FirstOrCreate() (*rdm.SchedulerTask, error)
	FindByPage(offset int, limit int) (result []*rdm.SchedulerTask, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ISchedulerTaskDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	UpdateStatus(ids []int64, instanceID int64, oldStatus int32, status int32) (rowsAffected int64, err error)
	PurgeInstanceTask(ids []int64, instanceID int64) (rowsAffected int64, err error)
	ListWaittingTasks(namespace string, category string, instanceID int64, id int64) (result []*rdm.SchedulerTask, err error)
	ListInstanceRunningTasks(id int64, instanceID int64) (result []*rdm.SchedulerTask, err error)
}

// update @@table
//
//	set status = @status
//	where id in @ids and status = @oldStatus and instance_id = @instanceID
func (s schedulerTaskDo) UpdateStatus(ids []int64, instanceID int64, oldStatus int32, status int32) (rowsAffected int64, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, status)
	params = append(params, ids)
	params = append(params, oldStatus)
	params = append(params, instanceID)
	generateSQL.WriteString("update scheduler_task set status = ? where id in ? and status = ? and instance_id = ? ")

	var executeSQL *gorm.DB
	executeSQL = s.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	rowsAffected = executeSQL.RowsAffected
	err = executeSQL.Error

	return
}

// update @@table
//
//	set status = 1, instance_id = 0
//	where id in @ids
//	and instance_id = @instanceID
//	and status = 2
func (s schedulerTaskDo) PurgeInstanceTask(ids []int64, instanceID int64) (rowsAffected int64, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, ids)
	params = append(params, instanceID)
	generateSQL.WriteString("update scheduler_task set status = 1, instance_id = 0 where id in ? and instance_id = ? and status = 2 ")

	var executeSQL *gorm.DB
	executeSQL = s.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	rowsAffected = executeSQL.RowsAffected
	err = executeSQL.Error

	return
}

// select * from @@table
//
//	where namespace = @namespace
//	and category = @category
//	and status = 1
//	and instance_id in (@instanceID, 0)
//	and id > @id
//	order by id limit 20
func (s schedulerTaskDo) ListWaittingTasks(namespace string, category string, instanceID int64, id int64) (result []*rdm.SchedulerTask, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, namespace)
	params = append(params, category)
	params = append(params, instanceID)
	params = append(params, id)
	generateSQL.WriteString("select * from scheduler_task where namespace = ? and category = ? and status = 1 and instance_id in (?, 0) and id > ? order by id limit 20 ")

	var executeSQL *gorm.DB
	executeSQL = s.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// select * from @@table
//
//	where instance_id = @instanceID
//	and status = 2
//	and id > @id
//	order by id limit 20
//
// todo order by id 是否OK
func (s schedulerTaskDo) ListInstanceRunningTasks(id int64, instanceID int64) (result []*rdm.SchedulerTask, err error) {
	var generateSQL strings.Builder
	generateSQL.WriteString("todo order by id 是否OK ")

	var executeSQL *gorm.DB
	executeSQL = s.UnderlyingDB().Raw(generateSQL.String()).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (s schedulerTaskDo) Debug() ISchedulerTaskDo {
	return s.withDO(s.DO.Debug())
}

func (s schedulerTaskDo) WithContext(ctx context.Context) ISchedulerTaskDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s schedulerTaskDo) ReadDB() ISchedulerTaskDo {
	return s.Clauses(dbresolver.Read)
}

func (s schedulerTaskDo) WriteDB() ISchedulerTaskDo {
	return s.Clauses(dbresolver.Write)
}

func (s schedulerTaskDo) Session(config *gorm.Session) ISchedulerTaskDo {
	return s.withDO(s.DO.Session(config))
}

func (s schedulerTaskDo) Clauses(conds ...clause.Expression) ISchedulerTaskDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s schedulerTaskDo) Returning(value interface{}, columns ...string) ISchedulerTaskDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s schedulerTaskDo) Not(conds ...gen.Condition) ISchedulerTaskDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s schedulerTaskDo) Or(conds ...gen.Condition) ISchedulerTaskDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s schedulerTaskDo) Select(conds ...field.Expr) ISchedulerTaskDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s schedulerTaskDo) Where(conds ...gen.Condition) ISchedulerTaskDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s schedulerTaskDo) Order(conds ...field.Expr) ISchedulerTaskDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s schedulerTaskDo) Distinct(cols ...field.Expr) ISchedulerTaskDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s schedulerTaskDo) Omit(cols ...field.Expr) ISchedulerTaskDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s schedulerTaskDo) Join(table schema.Tabler, on ...field.Expr) ISchedulerTaskDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s schedulerTaskDo) LeftJoin(table schema.Tabler, on ...field.Expr) ISchedulerTaskDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s schedulerTaskDo) RightJoin(table schema.Tabler, on ...field.Expr) ISchedulerTaskDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s schedulerTaskDo) Group(cols ...field.Expr) ISchedulerTaskDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s schedulerTaskDo) Having(conds ...gen.Condition) ISchedulerTaskDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s schedulerTaskDo) Limit(limit int) ISchedulerTaskDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s schedulerTaskDo) Offset(offset int) ISchedulerTaskDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s schedulerTaskDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ISchedulerTaskDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s schedulerTaskDo) Unscoped() ISchedulerTaskDo {
	return s.withDO(s.DO.Unscoped())
}

func (s schedulerTaskDo) Create(values ...*rdm.SchedulerTask) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s schedulerTaskDo) CreateInBatches(values []*rdm.SchedulerTask, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s schedulerTaskDo) Save(values ...*rdm.SchedulerTask) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s schedulerTaskDo) First() (*rdm.SchedulerTask, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*rdm.SchedulerTask), nil
	}
}

func (s schedulerTaskDo) Take() (*rdm.SchedulerTask, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*rdm.SchedulerTask), nil
	}
}

func (s schedulerTaskDo) Last() (*rdm.SchedulerTask, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*rdm.SchedulerTask), nil
	}
}

func (s schedulerTaskDo) Find() ([]*rdm.SchedulerTask, error) {
	result, err := s.DO.Find()
	return result.([]*rdm.SchedulerTask), err
}

func (s schedulerTaskDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*rdm.SchedulerTask, err error) {
	buf := make([]*rdm.SchedulerTask, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s schedulerTaskDo) FindInBatches(result *[]*rdm.SchedulerTask, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s schedulerTaskDo) Attrs(attrs ...field.AssignExpr) ISchedulerTaskDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s schedulerTaskDo) Assign(attrs ...field.AssignExpr) ISchedulerTaskDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s schedulerTaskDo) Joins(fields ...field.RelationField) ISchedulerTaskDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s schedulerTaskDo) Preload(fields ...field.RelationField) ISchedulerTaskDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s schedulerTaskDo) FirstOrInit() (*rdm.SchedulerTask, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*rdm.SchedulerTask), nil
	}
}

func (s schedulerTaskDo) FirstOrCreate() (*rdm.SchedulerTask, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*rdm.SchedulerTask), nil
	}
}

func (s schedulerTaskDo) FindByPage(offset int, limit int) (result []*rdm.SchedulerTask, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s schedulerTaskDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s schedulerTaskDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s schedulerTaskDo) Delete(models ...*rdm.SchedulerTask) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *schedulerTaskDo) withDO(do gen.Dao) *schedulerTaskDo {
	s.DO = *do.(*gen.DO)
	return s
}
