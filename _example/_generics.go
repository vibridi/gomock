package _example

type TestInterface[T any, R ~int] interface {
	Get() R
	Foo(v T)
}
