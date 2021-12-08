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
		name             string
		option           iterm2.InlineImageOption
		serializedOption string
	}{
		{
			name:             "name",
			option:           iterm2.WithName("foo"),
			serializedOption: "name=Zm9v",
		},
		{
			name:             "width-cells",
			option:           iterm2.WithWidthCells(3),
			serializedOption: "width=3",
		},
		{
			name:             "width-pixels",
			option:           iterm2.WithWidthPixels(3),
			serializedOption: "width=3px",
		},
		{
			name:             "width-percent",
			option:           iterm2.WithWidthPercent(3),
			serializedOption: "width=3%",
		},
		{
			name:             "width-auto",
			option:           iterm2.WithWidthAuto(),
			serializedOption: "width=auto",
		},
		{
			name:             "height-cells",
			option:           iterm2.WithHeightCells(3),
			serializedOption: "height=3",
		},
		{
			name:             "height-pixels",
			option:           iterm2.WithHeightPixels(3),
			serializedOption: "height=3px",
		},
		{
			name:             "height-percent",
			option:           iterm2.WithHeightPercent(3),
			serializedOption: "height=3%",
		},
		{
			name:             "height-auto",
			option:           iterm2.WithHeightAuto(),
			serializedOption: "height=auto",
		},
		{
			name:             "preserve-aspect-ratio-on",
			option:           iterm2.WithPreserveAspectRatio(true),
			serializedOption: "preserveAspectRatio=1",
		},
		{
			name:             "preserve-aspect-ratio-off",
			option:           iterm2.WithPreserveAspectRatio(false),
			serializedOption: "preserveAspectRatio=0",
		},
		{
			name:             "inline-on",
			option:           iterm2.WithInline(true),
			serializedOption: "inline=1",
		},
		{
			name:             "inline-off",
			option:           iterm2.WithInline(false),
			serializedOption: "inline=0",
		},
	}

	imageData, err := os.ReadFile(filepath.Join("testdata/image.jpeg"))
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	for _, test := range tests {
		var (
			option           = test.option
			serializedOption = test.serializedOption
		)

		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := iterm2.InlineImageTo(&buffer, imageData, option); err != nil {
				t.Fatalf("inline image: %v", err)
			}

			if diff := cmp.Diff(buffer.String(), buildInlineImage(imageData, serializedOption)); diff != "" {
				t.Fatalf("invalid serialization:\n%s", diff)
			}
		})
	}
}

func buildInlineImage(imageData []byte, serializedOption string) string {
	return fmt.Sprintf("\033]1337;File=size=%d;%s:%s\a", len(imageData), serializedOption, base64.StdEncoding.EncodeToString(imageData))
}
