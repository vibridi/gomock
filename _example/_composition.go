package _example

import "github.com/vibridi/gomock/v3/_example/sub"

type TestComposedInterface interface {
	testSubInterface
	sub.TestSubInterface
}

type testSubInterface interface {
	DoThings(s string, v int) error
}
