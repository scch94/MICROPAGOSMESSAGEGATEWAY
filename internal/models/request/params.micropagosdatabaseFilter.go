package request

type FilterRequest struct {
	MobileNumber string
	ShortNUmber  string
}

func NewFilterRequest(mobileNumber string, shortNUmber string) *FilterRequest {
	return &FilterRequest{
		MobileNumber: mobileNumber,
		ShortNUmber:  shortNUmber,
	}
}
