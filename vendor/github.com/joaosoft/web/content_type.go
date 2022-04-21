package web

import (
	"bytes"
	"encoding/binary"
)

// Mime type
type ContentType string

const (
	ContentTypeEmpty                  ContentType = ""
	ContentTypeApplicationJSON        ContentType = "application/json"
	ContentTypeApplicationJavaScript  ContentType = "application/javascript"
	ContentTypeApplicationXML         ContentType = "application/xml"
	ContentTypeTextXML                ContentType = "text/xml"
	ContentTypeTextCSV                ContentType = "text/csv"
	ContentTypeApplicationForm        ContentType = "application/x-www-form-urlencoded"
	ContentTypeApplicationProtobuf    ContentType = "application/protobuf"
	ContentTypeApplicationMsgpack     ContentType = "application/msgpack"
	ContentTypeTextHTML               ContentType = "text/html"
	ContentTypeTextPlain              ContentType = "text/plain"
	ContentTypeMultipartFormData      ContentType = "multipart/form-data"
	ContentTypeMultipartMixed         ContentType = "multipart/mixed"
	ContentTypeApplicationOctetStream ContentType = "application/octet-stream"
	ContentTypeApplicationZip         ContentType = "application/zip"
	ContentTypeApplication7z          ContentType = "application/x-7z-compressed"
	ContentTypeApplicationGzip        ContentType = "application/x-gzip"
	ContentTypeVideoMp4               ContentType = "video/mp4"
	ContentTypeApplicationPdf         ContentType = "application/pdf"
	ContentTypeApplicationPostScript  ContentType = "application/postscript"
	ContentTypeImageGif               ContentType = "image/gif"
	ContentTypeImagePng               ContentType = "image/png"
	ContentTypeImageJpeg              ContentType = "image/jpeg"
	ContentTypeImageBmp               ContentType = "image/bmp"
	ContentTypeImageWebp              ContentType = "image/webp"
	ContentTypeImageVnd               ContentType = "image/vnd.microsoft.icon"
	ContentTypeAudioWave              ContentType = "audio/wave"
	ContentTypeAudioAiff              ContentType = "audio/aiff"
	ContentTypeAudioBasic             ContentType = "audio/basic"
	ContentTypeApplicationOgg         ContentType = "application/ogg"
	ContentTypeAudioMidi              ContentType = "audio/midi"
	ContentTypeAudioMpeg              ContentType = "audio/mpeg"
	ContentTypeVideoAvi               ContentType = "video/avi"
	ContentTypeApplicationVnd         ContentType = "application/vnd.ms-fontobject"
	ContentTypeApplicationFontTtf     ContentType = "application/font-ttf"
	ContentTypeApplicationFontOff     ContentType = "application/font-off"
	ContentTypeApplicationFontCff     ContentType = "application/font-cff"
	ContentTypeApplicationFontWoff    ContentType = "application/font-woff"
	ContentTypeApplicationVideoWebm   ContentType = "video/webm"
	ContentTypeApplicationRar         ContentType = "application/x-rar-compressed"

	// the algorithm uses at most sniffLen bytes to make its decision.
	sniffLen = 512
)

type Content struct {
	contentType ContentType
	charset     Charset
}

var (
	extensions = map[string]*Content{
		"html": &Content{
			contentType: ContentTypeTextHTML,
			charset:     CharsetUTF8,
		},
		"js": &Content{
			contentType: ContentTypeApplicationJavaScript,
			charset:     CharsetUTF8,
		},
		"json": &Content{
			contentType: ContentTypeApplicationJSON,
			charset:     CharsetUTF8,
		},
		"xml": &Content{
			contentType: ContentTypeTextXML,
			charset:     CharsetUTF8,
		},
		"zip": &Content{
			contentType: ContentTypeApplicationZip,
		},
		"7z": &Content{
			contentType: ContentTypeApplication7z,
		},
		"rar": &Content{
			contentType: ContentTypeApplicationRar,
		},
		"txt": &Content{
			contentType: ContentTypeTextPlain,
		},
		"csv": &Content{
			contentType: ContentTypeTextPlain,
		},
		"pdf": &Content{
			contentType: ContentTypeApplicationPdf,
		},
		"ps": &Content{
			contentType: ContentTypeApplicationPostScript,
		},
		"gif": &Content{
			contentType: ContentTypeImageGif,
		},
		"png": &Content{
			contentType: ContentTypeImagePng,
		},
		"jpeg": &Content{
			contentType: ContentTypeImageJpeg,
		},
		"jpg": &Content{
			contentType: ContentTypeImageJpeg,
		},
		"bmp": &Content{
			contentType: ContentTypeImageBmp,
		},
		"vnd": &Content{
			contentType: ContentTypeImageVnd,
		},
		"mpeg": &Content{
			contentType: ContentTypeAudioMpeg,
		},
		"aiff": &Content{
			contentType: ContentTypeAudioAiff,
		},
		"midi": &Content{
			contentType: ContentTypeAudioMidi,
		},
		"wave": &Content{
			contentType: ContentTypeAudioWave,
		},
		"avi": &Content{
			contentType: ContentTypeVideoAvi,
		},
		"mp4": &Content{
			contentType: ContentTypeVideoMp4,
		},
		"webm": &Content{
			contentType: ContentTypeApplicationVideoWebm,
		},
		"ogg": &Content{
			contentType: ContentTypeApplicationOgg,
		},
		"ttf": &Content{
			contentType: ContentTypeApplicationFontTtf,
		},
		"cff": &Content{
			contentType: ContentTypeApplicationFontCff,
		},
		"off": &Content{
			contentType: ContentTypeApplicationFontOff,
		},
		"woff": &Content{
			contentType: ContentTypeApplicationFontWoff,
		},
	}
)

func detectDataContentType(data []byte) (ContentType, Charset) {
	if len(data) > sniffLen {
		data = data[:sniffLen]
	}

	// index of the first non-whitespace byte in data.
	firstNonWhiteSpace := 0
	for ; firstNonWhiteSpace < len(data) && isWhiteSpace(data[firstNonWhiteSpace]); firstNonWhiteSpace++ {
	}

	for _, signature := range sniffSignatures {
		if contentType, charset := signature.match(data, firstNonWhiteSpace); contentType != "" {
			return contentType, charset
		}
	}

	return ContentTypeApplicationOctetStream, ""
}

func DetectContentType(extension string, data []byte) (ContentType, Charset) {
	if content, ok := extensions[extension]; ok {
		return content.contentType, content.charset
	} else {
		return detectDataContentType(data)
	}
}

func isWhiteSpace(b byte) bool {
	switch b {
	case '\t', '\n', '\x0c', '\r', ' ':
		return true
	}
	return false
}

type sniffSignature interface {
	match(data []byte, firstNonWS int) (ContentType, Charset)
}

var sniffSignatures = []sniffSignature{
	html{
		signature:   []byte(`<!DOCTYPE HTML`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<HTML`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<HEAD`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<SCRIPT`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<IFRAME`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<H1`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<DIV`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<FONT`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<TABLE`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<A`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<STYLE`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<TITLE`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<B`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<BODY`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<BR`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<P`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<!DOCTYPE HTML`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	html{
		signature:   []byte(`<!--`),
		contentType: ContentTypeTextHTML,
		charset:     CharsetUTF8,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\xFF\xFF"),
		pat:         []byte("<?xml"),
		skipWS:      true,
		contentType: ContentTypeTextXML,
		charset:     CharsetUTF8,
	},
	&exact{
		signature:   []byte("%PDF-"),
		contentType: ContentTypeApplicationPdf,
	},
	&exact{
		signature:   []byte("%!PS-Adobe-"),
		contentType: ContentTypeApplicationPostScript,
	},
	&masked{
		mask:        []byte("\xFF\xFF\x00\x00"),
		pat:         []byte("\xFE\xFF\x00\x00"),
		contentType: ContentTypeTextPlain,
		charset:     CharsetUTF16be,
	},
	&masked{
		mask:        []byte("\xFF\xFF\x00\x00"),
		pat:         []byte("\xFF\xFE\x00\x00"),
		contentType: ContentTypeTextPlain,
		charset:     CharsetUTF16le,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\x00"),
		pat:         []byte("\xEF\xBB\xBF\x00"),
		contentType: ContentTypeTextPlain,
		charset:     CharsetUTF8,
	},
	&exact{
		signature:   []byte("GIF87a"),
		contentType: ContentTypeImageGif,
	},
	&exact{
		signature:   []byte("GIF89a"),
		contentType: ContentTypeImageGif,
	},
	&exact{
		signature:   []byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"),
		contentType: ContentTypeImagePng,
	},
	&exact{
		signature:   []byte("\xFF\xD8\xFF"),
		contentType: ContentTypeImageJpeg,
	},
	&exact{
		signature:   []byte("BM"),
		contentType: ContentTypeImageBmp,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF\xFF\xFF"),
		pat:         []byte("RIFF\x00\x00\x00\x00WEBPVP"),
		contentType: ContentTypeImageWebp,
	},
	&exact{
		signature:   []byte("\x00\x00\x01\x00"),
		contentType: ContentTypeImageVnd,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),
		pat:         []byte("RIFF\x00\x00\x00\x00WAVE"),
		contentType: ContentTypeAudioWave,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),
		pat:         []byte("FORM\x00\x00\x00\x00AIFF"),
		contentType: ContentTypeAudioAiff,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\xFF"),
		pat:         []byte(".snd"),
		contentType: ContentTypeAudioBasic,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\xFF\xFF"),
		pat:         []byte("OggS\x00"),
		contentType: ContentTypeApplicationOgg,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"),
		pat:         []byte("MThd\x00\x00\x00\x06"),
		contentType: ContentTypeAudioMidi,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF"),
		pat:         []byte("ID3"),
		contentType: ContentTypeAudioMpeg,
	},
	&masked{
		mask:        []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),
		pat:         []byte("RIFF\x00\x00\x00\x00AVI "),
		contentType: ContentTypeVideoAvi,
	},
	&masked{
		pat:         []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x4C\x50"),
		mask:        []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xFF\xFF"),
		contentType: ContentTypeApplicationVnd,
	},
	&exact{
		signature:   []byte("\x00\x01\x00\x00"),
		contentType: ContentTypeApplicationFontTtf,
	},
	&exact{
		signature:   []byte("OTTO"),
		contentType: ContentTypeApplicationFontOff,
	},
	&exact{
		signature:   []byte("ttcf"),
		contentType: ContentTypeApplicationFontCff,
	},
	&exact{
		signature:   []byte("wOFF"),
		contentType: ContentTypeApplicationFontWoff,
	},
	&exact{
		signature:   []byte("\x1A\x45\xDF\xA3"),
		contentType: ContentTypeApplicationVideoWebm},
	&exact{
		signature:   []byte("\x52\x61\x72\x20\x1A\x07\x00"),
		contentType: ContentTypeApplicationRar},
	&exact{
		signature:   []byte("\x50\x4B\x03\x04"),
		contentType: ContentTypeApplicationZip,
	},
	&exact{
		signature:   []byte("\x1F\x8B\x08"),
		contentType: ContentTypeApplicationGzip,
	},
	mp4Signature{
		contentType: ContentTypeVideoMp4,
	},
	textSignature{
		contentType: ContentTypeTextPlain,
	},
}

type exact struct {
	signature   []byte
	contentType ContentType
	charset     Charset
}

func (e *exact) match(data []byte, firstNonWS int) (ContentType, Charset) {
	if bytes.HasPrefix(data, e.signature) {
		return e.contentType, e.charset
	}
	return "", ""
}

type masked struct {
	mask, pat   []byte
	skipWS      bool
	contentType ContentType
	charset     Charset
}

func (m *masked) match(data []byte, firstNonWS int) (ContentType, Charset) {
	// pattern matching algorithm section 6
	// https://mimesniff.spec.whatwg.org/#pattern-matching-algorithm

	if m.skipWS {
		data = data[firstNonWS:]
	}
	if len(m.pat) != len(m.mask) {
		return "", ""
	}
	if len(data) < len(m.mask) {
		return "", ""
	}
	for i, mask := range m.mask {
		db := data[i] & mask
		if db != m.pat[i] {
			return "", ""
		}
	}
	return m.contentType, m.charset
}

type html struct {
	signature   []byte
	contentType ContentType
	charset     Charset
}

func (h html) match(data []byte, firstNonWS int) (ContentType, Charset) {
	data = data[firstNonWS:]
	if len(data) < len(h.signature)+1 {
		return "", ""
	}
	for i, b := range h.signature {
		db := data[i]
		if 'A' <= b && b <= 'Z' {
			db &= 0xDF
		}
		if b != db {
			return "", ""
		}
	}
	// Next byte must be space or right angle bracket.
	if db := data[len(h.signature)]; db != ' ' && db != '>' {
		return "", ""
	}
	return h.contentType, h.charset
}

var mp4ftype = []byte("ftyp")
var mp4 = []byte("mp4")

type mp4Signature struct {
	contentType ContentType
	charset     Charset
}

func (m mp4Signature) match(data []byte, firstNonWS int) (ContentType, Charset) {
	// https://mimesniff.spec.whatwg.org/#signature-for-mp4
	// c.f. section 6.2.1
	if len(data) < 12 {
		return "", ""
	}
	boxSize := int(binary.BigEndian.Uint32(data[:4]))
	if boxSize%4 != 0 || len(data) < boxSize {
		return "", ""
	}
	if !bytes.Equal(data[4:8], mp4ftype) {
		return "", ""
	}
	for st := 8; st < boxSize; st += 4 {
		if st == 12 {
			// minor version number
			continue
		}
		if bytes.Equal(data[st:st+3], mp4) {
			return m.contentType, ""
		}
	}
	return "", ""
}

type textSignature struct {
	contentType ContentType
	charset     Charset
}

func (t textSignature) match(data []byte, firstNonWS int) (ContentType, Charset) {
	// c.f. section 5, step 4.
	for _, b := range data[firstNonWS:] {
		switch {
		case b <= 0x08,
			b == 0x0B,
			0x0E <= b && b <= 0x1A,
			0x1C <= b && b <= 0x1F:
			return "", ""
		}
	}
	return t.contentType, ""
}
