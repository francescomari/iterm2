package iterm2

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// InlineImageOption configures the behaviour of InlineImage.
type InlineImageOption func(opts *inlineImageOptions)

// WithName gives a name to the inlined or downloaded image.
func WithName(name string) InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.name = name
	}
}

// WithWidthCells configures the width of the image in character cells.
func WithWidthCells(cells int) InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.width = fmt.Sprintf("%d", cells)
	}
}

// WithWidthPixels configures the width of the image in pixels.
func WithWidthPixels(pixels int) InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.width = fmt.Sprintf("%dpx", pixels)
	}
}

// WithWidthPercent configures the width of the image as a percentage of the
// width of the terminal session.
func WithWidthPercent(percent int) InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.width = fmt.Sprintf("%d%%", percent)
	}
}

// WithWidthAuto will use the original image size to determine the inlined size.
func WithWidthAuto() InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.width = "auto"
	}
}

// WithHeightCells configures the height of the image in character cells.
func WithHeightCells(cells int) InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.height = fmt.Sprintf("%d", cells)
	}
}

// WithHeightPixels configures the height of the image in pixels.
func WithHeightPixels(pixels int) InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.height = fmt.Sprintf("%dpx", pixels)
	}
}

// WithHeightPercent configures the height of the image as a percentage of the
// height of the terminal session.
func WithHeightPercent(percent int) InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.height = fmt.Sprintf("%d%%", percent)
	}
}

// WithHeightAuto will use the original image size to determine the inlined
// size.
func WithHeightAuto() InlineImageOption {
	return func(opts *inlineImageOptions) {
		opts.height = "auto"
	}
}

// WithPreserveAspectRatio determines whether to respect the original image
// aspect ration.
func WithPreserveAspectRatio(flag bool) InlineImageOption {
	return func(opts *inlineImageOptions) {
		if flag {
			opts.preserveAspectRatio = "1"
		} else {
			opts.preserveAspectRatio = "0"
		}
	}
}

// WithInline determines whether the image should be inlined or downloaded.
func WithInline(flag bool) InlineImageOption {
	return func(opts *inlineImageOptions) {
		if flag {
			opts.inline = "1"
		} else {
			opts.inline = "0"
		}
	}
}

type inlineImageOptions struct {
	name                string
	width               string
	height              string
	preserveAspectRatio string
	inline              string
}

// InlineImage is equivalent to calling InlineImageTo with os.Stdout as the
// output io.Writer.
func InlineImage(data []byte, opts ...InlineImageOption) (int, error) {
	return InlineImageTo(os.Stdout, data, opts...)
}

// InlineImageTo implements iterm2's Inline Images Protocol. It writes the
// necessary escape sequences to the provided io.Writer to inline or download
// the provided image. InlineImageTo accepts zero or more optional configuration
// options. If not provided, the configuration options will use the default
// values documented by the iterm2's Inline Images Protocol.
//
// InlineImageTo keeps track of number of written bytes as reported by the
// provided io.Writer, and returns it to the caller. If the provided io.Writer
// returns an error, InlineImageTo returns immediately with that error and with
// the number of bytes written up to that point.
func InlineImageTo(w io.Writer, data []byte, opts ...InlineImageOption) (int, error) {
	var written int

	if err := inlineImageTo(newPrinter(w, &written), data, opts...); err != nil {
		return written, err
	}

	return written, nil
}

func inlineImageTo(print printer, data []byte, opts ...InlineImageOption) error {
	var options inlineImageOptions

	for _, opt := range opts {
		opt(&options)
	}

	if err := print("\033]1337;File=size=%d", len(data)); err != nil {
		return err
	}

	if options.name != "" {
		if err := print(";name=%s", base64.StdEncoding.EncodeToString([]byte(options.name))); err != nil {
			return err
		}
	}

	if options.height != "" {
		if err := print(";height=%s", options.height); err != nil {
			return err
		}
	}

	if options.width != "" {
		if err := print(";width=%s", options.width); err != nil {
			return err
		}
	}

	if options.preserveAspectRatio != "" {
		if err := print(";preserveAspectRatio=%s", options.preserveAspectRatio); err != nil {
			return err
		}
	}

	if options.inline != "" {
		if err := print(";inline=%s", options.inline); err != nil {
			return err
		}
	}

	if err := print(":%s\a", base64.StdEncoding.EncodeToString(data)); err != nil {
		return err
	}

	return nil
}

type printer func(format string, args ...interface{}) error

func newPrinter(w io.Writer, written *int) printer {
	return func(format string, args ...interface{}) error {
		n, err := fmt.Fprintf(w, format, args...)

		*written += n

		return err
	}
}
