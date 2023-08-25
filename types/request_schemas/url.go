package requestschemas

type CreateURL struct {
	FullUrl string `validate:"required,http_url"`
}
