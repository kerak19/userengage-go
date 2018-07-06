package userengage

// Client is an client used for communication with userengage api
type Client struct {
	apikey string
}

// New creates new Client with provided api key
func New(apikey string) Client {
	return Client{
		apikey: apikey,
	}
}
