package errorModel

import "errors"

type ErrorModel struct {
	Code                  int
	Error                 error
	FileName              string
	FuncName              string
	CausedBy              error
	ErrorParameter        []ErrorParameter
	AdditionalInformation []string
}

type ErrorParameter struct {
	ErrorParameterKey   string
	ErrorParameterValue string
}

func GenerateErrorModel(code int, err string, fileName string, funcName string, causedBy error) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.FileName = fileName
	errModel.FuncName = funcName
	errModel.CausedBy = causedBy
	return errModel
}

func GenerateNonErrorModel() ErrorModel {
	var errModel ErrorModel
	errModel.Code = 200
	errModel.Error = nil
	return errModel
}
