package godiff

const (
	//ChangeTypeCreate defines create change type
	ChangeTypeCreate = ChangeType("create")
	//ChangeTypeUpdate defines update change type
	ChangeTypeUpdate = ChangeType("update")
	//ChangeTypeDelete defines delete change type
	ChangeTypeDelete = ChangeType("delete")
)

type (
	//ChangeType defines a change type
	ChangeType string

	//ChangeRecord defines a change record
	ChangeRecord struct {
		Source   string `json:",omitempty"`
		SourceID string `json:",omitempty"`
		UserID   string `json:",omitempty"`
		Path     string
		Change   string
		From     interface{} `json:",omitempty"`
		To       interface{} `json:",omitempty"`
		Error    string      `json:",omitempty"`
	}

	//Change represents a change
	Change struct {
		Type  ChangeType
		Path  *Path
		From  interface{}
		To    interface{}
		Error string `json:",omitempty"`
	}
)
