package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_split_flags_visitor(t *testing.T) {
	asserter := assert.New(t)
	info := newCommandInfo("-l=axyn=xxx1", "--selector", "axyn=xxx2")
	err := splitFlagsVisitor(info, nil)
	asserter.NoError(err)
}
