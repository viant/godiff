package godiff

import "encoding/json"

type (
	//ChangeLog represents a change log
	ChangeLog struct {
		Changes []*Change
	}
)

//Size returns change log size
func (l *ChangeLog) Size() int {
	return len(l.Changes)
}

//Add adds change log
func (l *ChangeLog) Add(change *Change) {
	l.Changes = append(l.Changes, change)
}

//AddError adds an error
func (l *ChangeLog) AddError(path *Path, err error) {
	if err == nil {
		return
	}
	l.Add(&Change{Path: path, Error: err.Error()})
}

//AddCreate adds create change
func (l *ChangeLog) AddCreate(path *Path, value interface{}) {
	l.Add(&Change{Type: ChangeTypeCreate, Path: path, To: value})
}

//AddDelete adds delete change
func (l *ChangeLog) AddDelete(path *Path, value interface{}) {
	l.Add(&Change{Type: ChangeTypeDelete, Path: path, From: value})
}

//AddUpdate adds update change
func (l *ChangeLog) AddUpdate(path *Path, from, to interface{}) {
	l.Add(&Change{Type: ChangeTypeUpdate, Path: path, From: from, To: to})
}

//ToChangeRecords converts changeLog to change records
func (l *ChangeLog) ToChangeRecords(source, id string) []*ChangeRecord {
	var result []*ChangeRecord
	for _, change := range l.Changes {
		result = append(result, &ChangeRecord{
			Source: source,
			ID:     id,
			Path:   change.Path.String(),
			Change: string(change.Type),
			From:   change.From,
			To:     change.To,
		})
	}
	return result
}

//String stringify change
func (l *ChangeLog) String() string {
	if len(l.Changes) == 0 {
		return ""
	}
	records := l.ToChangeRecords("", "")
	result, _ := json.Marshal(records)
	return string(result)
}
