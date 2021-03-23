package applicationModel

type ContextModel struct {
	LoggerModel        LoggerModel
	PermissionHave     string
	IsSignatureCheck   bool
	IsInternal         bool
	LimitedByCreatedBy int64
}
