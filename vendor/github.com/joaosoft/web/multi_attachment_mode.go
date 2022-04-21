package web

type MultiAttachmentMode string

const (
	MultiAttachmentModeBoundary MultiAttachmentMode = "boundary"
	MultiAttachmentModeZip      MultiAttachmentMode = "zip"
)
