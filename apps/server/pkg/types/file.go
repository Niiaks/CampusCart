package types

type File struct {
	Tags     []string          `json:"tags"`
	Context  map[string]string `json:"context"`
	MetaData map[string]string `json:"metadata"`
}
