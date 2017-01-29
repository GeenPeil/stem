package postcodenl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var regexpPostcode = regexp.MustCompile(`^([0-9]{4}) ?([A-Za-z]{2})$`)

var (
	ErrInvalidPostcode    = errors.New("invalid postcode")
	ErrPostcodenlApiError = errors.New("postcode API error")
)

// Data provided from postcode.nl API's
type Data struct {
	Street                 string   `json:"street"`
	HouseNumber            uint64   `json:"houseNumber"`
	HouseNumberAddition    string   `json:"houseNumberAddition"`
	Postcode               string   `json:"postcode"`
	City                   string   `json:"city"`
	Municipality           string   `json:"municipality"`
	Province               string   `json:"province"`
	RdX                    int      `json:"rdX"`
	RdY                    int      `json:"rdY"`
	Latitude               float64  `json:"latitude"`
	Longitude              float64  `json:"longitude"`
	BagNumberDesignationID string   `json:"bagNumberDesignationId"`
	BagAddressableObjectID string   `json:"bagAddressableObjectId"`
	AddressType            string   `json:"addressType"`
	Purposes               []string `json:"purposes"`
	SurfaceArea            int      `json:"surfaceArea"`
	HouseNumberAdditions   []string `json:"houseNumberAdditions"`
}

type Exception struct {
	Exception   string `json:"exception"`
	ExceptionID string `json:"exceptionId"`
}

func (e Exception) Error() string {
	return e.Exception
}

type API interface {
	Check(postcode string, housenumber uint64) (*Data, error)
}

type realAPI struct {
	key    string
	secret string
}

// New creates a new API instance that wraps postcode.nl using provided key and secret
func New(key, secret string) API {
	return &realAPI{
		key:    key,
		secret: secret,
	}
}
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Check wraps the postcode.nl API and returns either a filled Data structure or one of the erors defined in this package, or an Exception.
func (a *realAPI) Check(postcode string, housenumber uint64) (*Data, error) {
	postcodeParts := regexpPostcode.FindStringSubmatch(postcode)
	if len(postcodeParts) != 3 {
		return nil, ErrInvalidPostcode
	}
	url := fmt.Sprintf("https://api.postcode.nl/rest/addresses/%s%s/%d/", postcodeParts[1], strings.ToUpper(postcodeParts[2]), housenumber)
	request, err := http.NewRequest(`GET`, url, nil)
	if err != nil {
		log.Fatalf("error creating postcode api request: %v", err)
	}
	request.Header.Set(`Authorization`, `Basic `+basicAuth(a.key, a.secret))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		var e Exception
		err = json.NewDecoder(response.Body).Decode(&e)
		response.Body.Close()
		if err != nil || e.ExceptionID == `` {
			return nil, ErrPostcodenlApiError
		}
		return nil, e
	}

	data := &Data{}
	err = json.NewDecoder(response.Body).Decode(data)
	response.Body.Close()
	if err != nil {
		log.Fatalf("error decoding json: %v", err)
	}
	return data, nil
}
