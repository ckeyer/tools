package tests

import (
	"testing"
)

func Test(t *testing.T) {
	ok, err := LocalCheck("sha1", "hash.txt")

	t.Errorf("%v %v", ok, err)
}
