package response

type FieldType int

const (
	Tbool FieldType = iota
	Toption
	Tnumber
	Ttext
)

type FieldSchema struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Desc  string    `json:"desc"`
	Type  FieldType `json:"type"`
	Extra any       `json:"extra"`
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type OptionExtra struct {
	Option []Option `json:"option"`
}

func NewOptionSchema(
	id, title, desc string,
	opts []Option,
) *FieldSchema {
	return &FieldSchema{
		ID:    id,
		Title: title,
		Desc:  desc,
		Type:  Toption,
		Extra: opts,
	}
}

type BoolExtra struct {
	TrueLabel  string `json:"trueLabel"`
	FalseLabel string `json:"falseLabel"`
}

func NewBoolSchema(
	id, title, desc string,
	trueL, falseL string,
) *FieldSchema {
	return &FieldSchema{
		ID:    id,
		Title: title,
		Desc:  desc,
		Type:  Toption,
		Extra: BoolExtra{
			TrueLabel:  trueL,
			FalseLabel: falseL,
		},
	}
}

type NumberExtra struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Step float64 `json:"step"`
}

func NewNumberSchema(
	id, title, desc string,
	max, min, step float64,
) *FieldSchema {
	return &FieldSchema{
		ID:    id,
		Title: title,
		Desc:  desc,
		Type:  Toption,
		Extra: NumberExtra{
			Max:  max,
			Min:  min,
			Step: step,
		},
	}
}

type ModuleSchema struct {
	Name   string
	Fields []*FieldSchema
}

var registry = map[string]*ModuleSchema{}

func Register(schema *ModuleSchema) {
	// 同名校验
	if _, ok := registry[schema.Name]; ok {
		return
	}

	// 是否需要重复id检查

	registry[schema.Name] = schema
}
