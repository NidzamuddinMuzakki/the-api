package requests

type RequestAData struct {
	Id int `param:"id" validate:"required,number,min=1"`
}
