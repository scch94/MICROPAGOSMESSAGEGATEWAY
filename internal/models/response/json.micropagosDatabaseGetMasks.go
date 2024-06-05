package response

type MaskResponse struct {
	Result  uint8  `json:"result"`
	Message string `json:"message"`
	Masks   []Mask `json:"masks"`
}
type Mask struct {
	ID             string `json:"id"`
	ShortNumber    string `json:"short_number"`
	MaskPattern    string `json:"mask_pattern"`
	MinLength      string `json:"min_length"`
	MaxLength      string `json:"max_length"`
	ExcludePattern string `json:"exclude_pattern"`
	Direction      string `json:"direction"`
	ApplicationID  string `json:"application_id"`
}
