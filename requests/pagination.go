package requests

type Pagination struct {
	Limit  int `query:"limit" validate:"required,numeric,min=1"`
	Offset int `query:"offset" validate:"required,numeric,min=1"`
}
