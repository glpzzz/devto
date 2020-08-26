package article

import (
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetImageLinks(t *testing.T) {
	type args struct {
		filename string
		images   map[string]string
	}

	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				filename: "./testdata/testdata.md",
				images: map[string]string{
					"./image.png":   "./a/image.png",
					"./image-2.png": "",
				},
			},
			want: `---
title: A title
published: false
description: A description
tags: tag-one, tag-two
---
![image](./a/image.png)
[Google](www.google.com)
![image](./image-2.png)
`,
			assertion: assert.NoError,
		},
		{
			name: "not found",
			args: args{
				filename: "./testdata/unknown.md",
				images:   map[string]string{"./image.png": "./a/image.png"},
			},
			want:      "",
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetImageLinks(tt.args.filename, tt.args.images)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetImageLinks_golden(t *testing.T) {
	content, err := SetImageLinks("./testdata/real_article.md", map[string]string{})
	assert.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "real_article", []byte(content))
}

func TestGetImageLinks(t *testing.T) {
	type args struct {
		filename string
	}

	tests := []struct {
		name      string
		args      args
		want      map[string]string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{filename: "./testdata/testdata.md"},
			want: map[string]string{
				"./image.png":   "",
				"./image-2.png": "",
			},
			assertion: assert.NoError,
		},
		{
			name:      "not found",
			args:      args{filename: "./testdata/unknown.md"},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetImageLinks(tt.args.filename)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrefixLinks(t *testing.T) {
	type args struct {
		links  map[string]string
		prefix string
		force  bool
	}

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "normal",
			args: args{
				links: map[string]string{
					"./image/image.png":   "image.png",
					"./image/picture.jpg": "",
				},
				prefix: "https://raw.githubusercontent.com/repo/user/",
			},
			want: map[string]string{
				"./image/image.png":   "image.png",
				"./image/picture.jpg": "https://raw.githubusercontent.com/repo/user/./image/picture.jpg",
			},
		},
		{
			name: "force",
			args: args{
				links: map[string]string{
					"./image/image.png":   "image.png",
					"./image/picture.jpg": "",
				},
				prefix: "https://raw.githubusercontent.com/repo/user/",
				force:  true,
			},
			want: map[string]string{
				"./image/image.png":   "https://raw.githubusercontent.com/repo/user/./image/image.png",
				"./image/picture.jpg": "https://raw.githubusercontent.com/repo/user/./image/picture.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, PrefixLinks(tt.args.links, tt.args.prefix, tt.args.force))
		})
	}
}
