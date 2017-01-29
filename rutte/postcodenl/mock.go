package postcodenl

type mockAPI struct {
	key    string
	secret string
}

// NewMock creates a new API instance that returns fake data
func NewMock() API {
	return &mockAPI{}
}

// Check wraps the postcode.nl API and returns either a filled Data structure or one of the erors defined in this package.
func (a *mockAPI) Check(postcode string, housenumber uint64) (*Data, error) {
	return &Data{
		Street:                 "Peilstraat",
		HouseNumber:            housenumber,
		HouseNumberAddition:    "A",
		Postcode:               postcode,
		City:                   "Peilstad",
		Municipality:           "Gemeente Peilstad",
		Province:               "Noord-Peilland",
		RdX:                    123,
		RdY:                    456,
		Latitude:               4.12314,
		Longitude:              51.12334,
		BagNumberDesignationID: "",
		BagAddressableObjectID: "",
		AddressType:            "home",
		Purposes:               []string{},
		SurfaceArea:            42,
		HouseNumberAdditions:   []string{"A"},
	}, nil
}
