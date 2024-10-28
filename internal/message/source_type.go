package message

type sourceType string

const (
	TELEGRAM sourceType = "TELEGRAM"
)

type SourceTypeEnum interface {
	SourceType() sourceType
}

func (receiver sourceType) SourceType() sourceType {
	return receiver
}
