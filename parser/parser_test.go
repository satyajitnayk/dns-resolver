package parser_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/satyajitnayk/dns-resolver/parser"
	"github.com/stretchr/testify/require"
)

func TestEncodeDomainName(t *testing.T) {
	t.Run("test domain name encoder", func(t *testing.T) {
		out := parser.EncodeDomainName("google.com")
		fmt.Println(out)
		require.Equal(t, "06676f6f676c6503636f6d00", hex.EncodeToString(out))
	})
}
