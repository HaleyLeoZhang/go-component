package db

import "github.com/jinzhu/gorm"

// DB常用的查询条件封装
type DBConditions struct {
	And       map[string]interface{}
	Or        map[string]interface{}
	Not       map[string]interface{}
	Limit     interface{}
	Offset    interface{}
	Order     interface{}
	Select    interface{}
	Group     string
	Having    interface{}
	NeedCount bool
	Count     int64
}

// 填充查询条件
func (d *DBConditions) Fill(db *gorm.DB) *gorm.DB {
	if d.Select != nil {
		db = db.Select(d.Select)
	}

	for cond, val := range d.And {
		db = db.Where(cond, val)
	}
	for cond, val := range d.Not {
		db = db.Not(cond, val)
	}
	for cond, val := range d.Or {
		db = db.Or(cond, val)
	}

	if d.NeedCount {
		db = db.Count(&d.Count)
	}
	if d.Order != nil {
		db = db.Order(d.Order)
	}
	if d.Limit != nil {
		db = db.Limit(d.Limit)
	}
	if d.Offset != nil {
		db = db.Offset(d.Offset)
	}
	if d.Group != "" {
		db = db.Group(d.Group)
	}
	if d.Having != nil {
		db = db.Having(d.Having)
	}

	return db
}

// 计算分页信息
func GetPageInfo(page int, size int) (int, int) {
	offset := (page - 1) * size
	return offset, size
}


/* Demo
cond := &base.DBConditions{
	Select = "id"
	And: map[string][]interface{}{
		"id IN (?)": {95,96,97},
	},
	Not: map[string][]interface{}{
		"id": {96},
	},
	Limit: 1,
	Offset: 1,
	Order: "id DESC",
}

// 查询列表
func (d *Dao) ComicListComplete(ctx context.Context, conditions *base.DBConditions) ([]*po.Comic, error) {
	db := d.DB.Context(ctx)

	var res []*po.Comic

	db = db.Table(po.Comic{}.TableName())

	db = conditions.Fill(db)

	err := db.Find(&res).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}
*/
