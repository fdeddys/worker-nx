package errorModel

var DefaultError map[string]ErrorClass

type ErrorClass struct {
	ErrorCode    string
	ErrorMessage string
}

func GenerateInternalDBServerError(fileName string, funcName string, causedBy error) ErrorModel {
	return GenerateErrorModel(500, "E-5-MAD-DBS-001", fileName, funcName, causedBy)
}
