package schemas

type GetTokenResponseSchema struct {
	Signature      string
	SignatureKeyId string
	Timestamp      int
	Result         GetTokenResultSchema
}

type GetTokenResultSchema struct {
	Uid             string
	PacketType      string
	Data            GetTokenDataSchema
	EncryptionKeyId string
	SymmetricKey    string
	Iv              string
}

type GetTokenDataSchema struct {
	Token     string
	ExpiresIn int
}

type SendInvoiceSchema struct {
	Signature      string
	SignatureKeyId string
	Timestamp      int
	Result         []SendInvoiceResultSchema
}

type SendInvoiceResultSchema struct {
	Uid             string
	ReferenceNumber string
	ErrorCode       string
	ErrorDetail     string
}

type InquiryResponseSchema struct {
	Signature      string
	SignatureKeyId string
	Timestamp      int
	Result         InquiryResultSchema
}

type InquiryResultSchema struct {
	Uid             string
	PacketType      string
	Data            []InquiryDataSchema
	EncryptionKeyId string
	SymmetricKey    string
	Iv              string
}

type InquiryDataSchema struct {
	ReferenceNumber string
	Uid             string
	Status          string
	Data            InquiryInnerDataSchema
	PacketType      string
	FiscalId        string
}

type InquiryInnerDataSchema struct {
	ConfirmationReferenceId string
	Error                   []InquiryInnerDataErrorSchema
	Warning                 []InquiryInnerDataWarningSchema
	Success                 bool
}

type InquiryInnerDataErrorSchema struct {
	Code    string
	Message string
}

type InquiryInnerDataWarningSchema struct {
	Code    string
	Message string
}