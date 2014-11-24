package edtd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "testing"
)

func TestHeaders(t *T) {

	test := `
        define header {
            DocType := "matroska";
            EBMLVersion := 1;
        }`

	m, err := parseAsRoot(bytes.NewBufferString(test))
	require.Nil(t, err)

	assert.Equal(t, m[0x282].def, []byte("matroska"))
	assert.Equal(t, m[0x282].mustMatchDef, true)

	assert.Equal(t, m[0x286].def, mustDefDataBytes(uint64(1)))
	assert.Equal(t, m[0x286].mustMatchDef, true)
}
