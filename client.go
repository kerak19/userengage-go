package userengage

// Client is an client used for communication with userengage api
type Client struct {
	apiKey    string
	apiPrefix string
}

// New creates new Client with provided api key
func New(APIKey string) Client {
	return Client{
		apiKey:    APIKey,
		apiPrefix: "https://app.userengage.com/api/public",
	}
}

// This is for the pull request
