package requestschemas

type CreateEntity struct {
	Name        string `validate:"required,min=3,max=6"`
	Description string `validate:"min=0,max=64"`
}

type CreateURL struct {
	FullUrl   string `validate:"required,http_url"`
	Entity    int64
	ShortCode string `validate:"min=0,max=10"`
}
