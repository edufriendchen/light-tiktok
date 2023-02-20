package code

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type ErrCode int64

const (
	ErrCode_SuccessCode                ErrCode = 0
	ErrCode_ServiceErrCode             ErrCode = 10001
	ErrCode_ParamErrCode               ErrCode = 10002
	ErrCode_UserAlreadyExistErrCode    ErrCode = 10003
	ErrCode_AuthorizationFailedErrCode ErrCode = 10004
	ErrCode_Neo4jColumnFailedErr       ErrCode = 10005
)

func (p ErrCode) String() string {
	switch p {
	case ErrCode_SuccessCode:
		return "SuccessCode"
	case ErrCode_ServiceErrCode:
		return "ServiceErrCode"
	case ErrCode_ParamErrCode:
		return "ParamErrCode"
	case ErrCode_UserAlreadyExistErrCode:
		return "UserAlreadyExistErrCode"
	case ErrCode_AuthorizationFailedErrCode:
		return "AuthorizationFailedErrCode"
	}
	return "<UNSET>"
}

func ErrCodeFromString(s string) (ErrCode, error) {
	switch s {
	case "SuccessCode":
		return ErrCode_SuccessCode, nil
	case "ServiceErrCode":
		return ErrCode_ServiceErrCode, nil
	case "ParamErrCode":
		return ErrCode_ParamErrCode, nil
	case "UserAlreadyExistErrCode":
		return ErrCode_UserAlreadyExistErrCode, nil
	case "AuthorizationFailedErrCode":
		return ErrCode_AuthorizationFailedErrCode, nil
	}
	return ErrCode(0), fmt.Errorf("not a valid ErrCode string")
}

func ErrCodePtr(v ErrCode) *ErrCode { return &v }
func (p *ErrCode) Scan(value interface{}) (err error) {
	var result sql.NullInt64
	err = result.Scan(value)
	*p = ErrCode(result.Int64)
	return
}

func (p *ErrCode) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}
