package iterm2_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/francescomari/iterm2"
	"github.com/google/go-cmp/cmp"
)

func TestInlineImageOptions(t *testing.T) {
	tests := []struct {
		name       string
		option     iterm2.InlineImageOption
		serialized string
	}{
		{
			name:       "name",
			option:     iterm2.WithName("foo"),
			serialized: "name=Zm9v",
		},
		{
			name:       "width-cells",
			option:     iterm2.WithWidthCells(3),
			serialized: "width=3",
		},
		{
			name:       "width-pixels",
			option:     iterm2.WithWidthPixels(3),
			serialized: "width=3px",
		},
		{
			name:       "width-percent",
			option:     iterm2.WithWidthPercent(3),
			serialized: "width=3%",
		},
		{
			name:       "width-auto",
			option:     iterm2.WithWidthAuto(),
			serialized: "width=auto",
		},
		{
			name:       "height-cells",
			option:     iterm2.WithHeightCells(3),
			serialized: "height=3",
		},
		{
			name:       "height-pixels",
			option:     iterm2.WithHeightPixels(3),
			serialized: "height=3px",
		},
		{
			name:       "height-percent",
			option:     iterm2.WithHeightPercent(3),
			serialized: "height=3%",
		},
		{
			name:       "height-auto",
			option:     iterm2.WithHeightAuto(),
			serialized: "height=auto",
		},
		{
			name:       "preserve-aspect-ratio-on",
			option:     iterm2.WithPreserveAspectRatio(true),
			serialized: "preserveAspectRatio=1",
		},
		{
			name:       "preserve-aspect-ratio-off",
			option:     iterm2.WithPreserveAspectRatio(false),
			serialized: "preserveAspectRatio=0",
		},
		{
			name:       "inline-on",
			option:     iterm2.WithInline(true),
			serialized: "inline=1",
		},
		{
			name:       "inline-off",
			option:     iterm2.WithInline(false),
			serialized: "inline=0",
		},
	}

	data, err := os.ReadFile(filepath.Join("testdata/image.jpeg"))
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	for _, test := range tests {
		var (
			option     = test.option
			serialized = test.serialized
		)

		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := iterm2.InlineImageTo(&buffer, data, option); err != nil {
				t.Fatalf("inline image: %v", err)
			}

			if diff := cmp.Diff(buffer.String(), buildInlineImage(data, serialized)); diff != "" {
				t.Fatalf("invalid serialization:\n%s", diff)
			}
		})
	}
}

func buildInlineImage(data []byte, serializedOptions string) string {
	return fmt.Sprintf("\033]1337;File=size=%d;%s:%s\a", len(data), serializedOptions, base64.StdEncoding.EncodeToString(data))
}
