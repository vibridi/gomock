package _example

type TestInterface interface {
	Get(param *TestParam) TestStruct
}

type TestParam struct {
}
