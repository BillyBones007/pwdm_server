// Models package describes models objects for working with storage.
package models

// UserModel - model for authentication/registration users.
type UserModel struct {
	Login    string
	Password string
}

// ListRecordsModel - model list records id. Used for delete all user records.
type ListRecordsModel struct {
	ListID []int32
}

// DataRecordModel - model for show information by one record.
type DataRecordModel struct {
	Title   string // record title
	Tag     string // tag for record
	Comment string // comment for record
	Type    int32  // data type (for example: 1 - login/password, 2 - card, 3 - text, 4 - binary)
	ID      int32  // record id in database
}

// ReqLogPwdModel - model login/password pair for request.
type ReqLogPwdModel struct {
	LogPwdModel
	ReqTechDataModel
}

// RespLogPwdModel - model login/password pair for response.
type RespLogPwdModel struct {
	LogPwdModel
	RespTechDataModel
}

// ReqCardModel - model card data for request.
type ReqCardModel struct {
	CardModel
	ReqTechDataModel
}

// RespCardModel - model card data for response.
type RespCardModel struct {
	CardModel
	RespTechDataModel
}

// ReqTextModel - model text data for request.
type ReqTextModel struct {
	TextDataModel
	ReqTechDataModel
}

// RespTextModel - model text data for response.
type RespTextModel struct {
	TextDataModel
	RespTechDataModel
}

// ReqBinaryModel - model binary data for request.
type ReqBinaryModel struct {
	BinaryDataModel
	ReqTechDataModel
}

// RespBinaryModel - model binary data for response.
type RespBinaryModel struct {
	BinaryDataModel
	RespTechDataModel
}

// ReqTechDataModel - model with general information for requests.
type ReqTechDataModel struct {
	Title   string
	Tag     string
	Comment string
	Type    int32
}

// RespTechDataModel - model with general information for responses.
type RespTechDataModel struct {
	Title   string
	Tag     string
	Comment string
	Error   error
	ID      int32
}

// InsertRespModel - model insert data response.
type InsertRespModel struct {
	Title string
	Error error
	ID    int32
}

// IDModel - model id record in database.
type IDModel struct {
	ID int32 // id record in database
}

// LogPassModel - model login/password pair.
type LogPwdModel struct {
	Login    string
	Password string
}

// CardModel - model card data.
type CardModel struct {
	Num       string // card number
	Data      string // validity period
	CVC       string // cvc card code
	FirstName string
	LastName  string
}

// TextDataModel - model for text data.
type TextDataModel struct {
	Data string // some text data
}

// BinaryDataModel - model for binary data.
type BinaryDataModel struct {
	Data []byte // some binary data
}
