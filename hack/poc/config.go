package poc

type Poc struct {
	FileName   string
	Urls       string
	ReverseUrl string

	ErrChannel chan string `json:"-"`
	Stop       int32       `json:"-"`
}
