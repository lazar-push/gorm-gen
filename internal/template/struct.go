package template

const (
	// TableQueryStruct table query struct
	TableQueryStruct = createMethod + `
	{{.QueryStructComment}}
	type {{.QueryStructName}} struct {
		{{.QueryStructName}}Do
		` + fields + `
	}
	` + tableMethod + asMethond + updateFieldMethod + getFieldMethod + fillFieldMapMethod + cloneMethod + replaceMethod + relationship + defineMethodStruct + `

	type {{.QueryStructName}}AsMethod struct {
		{{.QueryStructName}}
	}

	func New{{.QueryStructName}}AsMethod({{.QueryStructName}} {{.QueryStructName}}) *{{.QueryStructName}}AsMethod {
		return &{{.QueryStructName}}AsMethod{ {{.QueryStructName}} }
	}

	func (c {{.QueryStructName}}AsMethod) As(alias string) gen.Dao {
		return c.{{.QueryStructName}}Do.As(alias)
	}

	`

	// TableQueryStructWithContext table query struct with context
	TableQueryStructWithContext = createMethod + `
	{{.QueryStructComment}}
	type {{.QueryStructName}} struct {
		{{.QueryStructName}}Do {{.QueryStructName}}Do
		` + fields + `
	}
	` + tableMethod + asMethond + updateFieldMethod + `
	
	func ({{.S}} *{{.QueryStructName}}) WithContext(ctx context.Context) {{.ReturnObject}} { return {{.S}}.{{.QueryStructName}}Do.WithContext(ctx)}

	func ({{.S}} {{.QueryStructName}}) TableName() string { return {{.S}}.{{.QueryStructName}}Do.TableName() } 

	func ({{.S}} {{.QueryStructName}}) Alias() string { return {{.S}}.{{.QueryStructName}}Do.Alias() }

	func ({{.S}} {{.QueryStructName}}) Columns(cols ...field.Expr) gen.Columns { return {{.S}}.{{.QueryStructName}}Do.Columns(cols...) }

	` + getFieldMethod + fillFieldMapMethod + cloneMethod + replaceMethod + relationship + defineMethodStruct + `

	type {{.QueryStructName}}AsMethod struct {
		{{.QueryStructName}}
	}

	func New{{.QueryStructName}}AsMethod({{.QueryStructName}} {{.QueryStructName}}) *{{.QueryStructName}}AsMethod {
		return &{{.QueryStructName}}AsMethod{{{.QueryStructName}}}
	}

	func (c {{.QueryStructName}}AsMethod) As(alias string) gen.Dao {
		return c.{{.QueryStructName}}Do.As(alias)
	}
	
	`

	// TableQueryIface table query interface
	TableQueryIface = defineDoInterface
)

const (
	createMethod = `
	func new{{.ModelStructName}}(db *gorm.DB, opts ...gen.DOOption) {{.QueryStructName}} {
		_{{.QueryStructName}} := {{.QueryStructName}}{}
	
		_{{.QueryStructName}}.{{.QueryStructName}}Do.UseDB(db,opts...)
		_{{.QueryStructName}}.{{.QueryStructName}}Do.UseModel(&{{.StructInfo.Package}}.{{.StructInfo.Type}}{})
	
		tableName := _{{.QueryStructName}}.{{.QueryStructName}}Do.TableName()
		_{{$.QueryStructName}}.ALL = field.NewAsterisk(tableName)
		{{range .Fields -}}
		{{if not .IsRelation -}}
			{{- if .ColumnName -}}_{{$.QueryStructName}}.{{.Name}} = field.New{{.GenType}}(tableName, "{{.ColumnName}}"){{- end -}}
		{{- else -}}
			_{{$.QueryStructName}}.{{.Relation.Name}} = {{$.QueryStructName}}{{.Relation.RelationshipName}}{{.Relation.Name}}{
				db: db.Session(&gorm.Session{}),

				{{.Relation.StructFieldInit}}
			}
		{{end}}
		{{end}}

		_{{$.QueryStructName}}.fillFieldMap()
		
		return _{{.QueryStructName}}
	}
	`
	fields = `
	ALL field.Asterisk
	{{range .Fields -}}
		{{if not .IsRelation -}}
			{{if .MultilineComment -}}
			/*
{{.ColumnComment}}
    		*/
			{{end -}}
			{{- if .ColumnName -}}{{.Name}} field.{{.GenType}}{{if not .MultilineComment}}{{if .ColumnComment}}// {{.ColumnComment}}{{end}}{{end}}{{- end -}}
		{{- else -}}
			{{.Relation.Name}} {{$.QueryStructName}}{{.Relation.RelationshipName}}{{.Relation.Name}}
		{{end}}
	{{end}}

	fieldMap  map[string]field.Expr
`
	tableMethod = `
func ({{.S}} {{.QueryStructName}}) Table(newTableName string) *{{.QueryStructName}} { 
	{{.S}}.{{.QueryStructName}}Do.UseTable(newTableName)
	return {{.S}}.updateTableName(newTableName)
}
`

	asMethond = `	
func ({{.S}} {{.QueryStructName}}) As(alias string) *{{.QueryStructName}} { 
	{{.S}}.{{.QueryStructName}}Do.DO = *({{.S}}.{{.QueryStructName}}Do.As(alias).(*gen.DO))
	return {{.S}}.updateTableName(alias)
}
`
	updateFieldMethod = `
func ({{.S}} *{{.QueryStructName}}) updateTableName(table string) *{{.QueryStructName}} { 
	{{.S}}.ALL = field.NewAsterisk(table)
	{{range .Fields -}}
	{{if not .IsRelation -}}
		{{- if .ColumnName -}}{{$.S}}.{{.Name}} = field.New{{.GenType}}(table, "{{.ColumnName}}"){{- end -}}
	{{end}}
	{{end}}
	
	{{.S}}.fillFieldMap()

	return {{.S}}
}
`

	cloneMethod = `
func ({{.S}} {{.QueryStructName}}) clone(db *gorm.DB) {{.QueryStructName}} {
	{{.S}}.{{.QueryStructName}}Do.ReplaceConnPool(db.Statement.ConnPool)
	return {{.S}}
}
`
	replaceMethod = `
func ({{.S}} {{.QueryStructName}}) replaceDB(db *gorm.DB) {{.QueryStructName}} {
	{{.S}}.{{.QueryStructName}}Do.ReplaceDB(db)
	return {{.S}}
}
`
	getFieldMethod = `
func ({{.S}} *{{.QueryStructName}}) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := {{.S}}.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe,ok := _f.(field.OrderExpr)
	return _oe,ok
}
`
	relationship = `{{range .Fields}}{{if .IsRelation}}` +
		`{{- $relation := .Relation }}{{- $relationship := $relation.RelationshipName}}` +
		relationStruct + relationTx +
		`{{end}}{{end}}`
	defineMethodStruct = `type {{.QueryStructName}}Do struct { gen.DO }`

	fillFieldMapMethod = `
func ({{.S}} *{{.QueryStructName}}) fillFieldMap() {
	{{.S}}.fieldMap =  make(map[string]field.Expr, {{len .Fields}})
	{{range .Fields -}}
	{{if not .IsRelation -}}
		{{- if .ColumnName -}}{{$.S}}.fieldMap["{{.ColumnName}}"] = {{$.S}}.{{.Name}}{{- end -}}
	{{end}}
	{{end -}}
}
`

	defineDoInterface = `

type I{{.ModelStructName}}Base interface {
	Debug() I{{.ModelStructName}}Exported
	WithContext(ctx context.Context) I{{.ModelStructName}}Exported
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() I{{.ModelStructName}}Exported
	WriteDB() I{{.ModelStructName}}Exported
	As(alias string) gen.Dao
	Session(config *gorm.Session) I{{.ModelStructName}}Exported
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) I{{.ModelStructName}}Exported
	Not(conds ...gen.Condition) I{{.ModelStructName}}Exported
	Or(conds ...gen.Condition) I{{.ModelStructName}}Exported
	Select(conds ...field.Expr) I{{.ModelStructName}}Exported
	Where(conds ...gen.Condition) I{{.ModelStructName}}Exported
	Order(conds ...field.Expr) I{{.ModelStructName}}Exported
	Distinct(cols ...field.Expr) I{{.ModelStructName}}Exported
	Omit(cols ...field.Expr) I{{.ModelStructName}}Exported
	Join(table schema.Tabler, on ...field.Expr) I{{.ModelStructName}}Exported
	LeftJoin(table schema.Tabler, on ...field.Expr) I{{.ModelStructName}}Exported
	RightJoin(table schema.Tabler, on ...field.Expr) I{{.ModelStructName}}Exported
	Group(cols ...field.Expr) I{{.ModelStructName}}Exported
	Having(conds ...gen.Condition) I{{.ModelStructName}}Exported
	Limit(limit int) I{{.ModelStructName}}Exported
	Offset(offset int) I{{.ModelStructName}}Exported
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) I{{.ModelStructName}}Exported
	Unscoped() I{{.ModelStructName}}Exported
	Create(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) error
	CreateInBatches(values []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int) error
	Save(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) error
	First() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	Take() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	Last() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	Find() ([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, err error)
	FindInBatches(result *[]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) I{{.ModelStructName}}Exported
	Assign(attrs ...field.AssignExpr) I{{.ModelStructName}}Exported
	Joins(fields ...field.RelationField) I{{.ModelStructName}}Exported
	Preload(fields ...field.RelationField) I{{.ModelStructName}}Exported
	FirstOrInit() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	FirstOrCreate() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	FindByPage(offset int, limit int) (result []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) I{{.ModelStructName}}Exported
	UnderlyingDB() *gorm.DB
}

type I{{.ModelStructName}}Specific interface {
	{{range .Interfaces -}}
	{{.FuncSign}}
	{{end}}
}

type I{{.ModelStructName}}Exported interface {
	I{{.ModelStructName}}Base
	I{{.ModelStructName}}Specific
	gen.Condition
}


type I{{.ModelStructName}}Do interface {
	gen.SubQuery
	schema.Tabler

	I{{.ModelStructName}}Exported
}
`
)

const (
	relationStruct = `
type {{$.QueryStructName}}{{$relationship}}{{$relation.Name}} struct{
	db *gorm.DB
	
	field.RelationField
	
	{{$relation.StructField}}
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}) Where(conds ...field.Expr) *{{$.QueryStructName}}{{$relationship}}{{$relation.Name}} {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}) WithContext(ctx context.Context) *{{$.QueryStructName}}{{$relationship}}{{$relation.Name}} {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}) Session(session *gorm.Session) *{{$.QueryStructName}}{{$relationship}}{{$relation.Name}} {
	a.db = a.db.Session(session)
	return &a
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}) Model(m *{{$.StructInfo.Package}}.{{$.StructInfo.Type}}) *{{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx {
	return &{{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx{a.db.Model(m).Association(a.Name())}
}

`
	relationTx = `
type {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx struct{ tx *gorm.Association }

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx) Find() (result {{if eq $relationship "HasMany" "ManyToMany"}}[]{{end}}*{{$relation.Type}}, err error) {
	return result, a.tx.Find(&result)
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx) Append(values ...*{{$relation.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx) Replace(values ...*{{$relation.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx) Delete(values ...*{{$relation.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx) Clear() error {
	return a.tx.Clear()
}

func (a {{$.QueryStructName}}{{$relationship}}{{$relation.Name}}Tx) Count() int64 {
	return a.tx.Count()
}
`
)
