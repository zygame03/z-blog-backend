package config

type FieldType int

const (
	TBool FieldType = iota
	TOption
	TNumber
	TText
)

type BaseSchema struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  FieldType
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type OptionSchema struct {
	BaseSchema
	Option []Option
}

type BoolSchema struct {
	BaseSchema
	TrueLabel  string `json:"trueLabel"`
	FalseLabel string `json:"falseLabel"`
}

type NumberSchema struct {
	BaseSchema
	Min  *float64
	Max  *float64
	Step *float64
}
