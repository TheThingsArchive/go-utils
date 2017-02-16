package entries

import "github.com/apex/log"

type Entries struct {
	Entries []*log.Entry
}

func New() *Entries {
	return &Entries{
		Entries: make([]*log.Entry, 0, 10),
	}
}

func (e *Entries) HandleLog(entry *log.Entry) error {
	e.Entries = append(e.Entries, entry)
	return nil
}
