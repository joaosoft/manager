package web

// Mime type
type ContentDisposition string

const (
	ContentDispositionInline     ContentDisposition = "inline"
	ContentDispositionAttachment ContentDisposition = "attachment"
	ContentDispositionFormData   ContentDisposition = "form-data"
)
