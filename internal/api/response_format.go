package api

type ResponseFormat string

const (
	ResponseFormatCSV  ResponseFormat = "csv"
	ResponseFormatJSON ResponseFormat = "json"
	ResponseFormatNone ResponseFormat = ""
)

func (f ResponseFormat) IsValid() bool {
	return f == ResponseFormatCSV || f == ResponseFormatJSON || f == ResponseFormatNone
}
