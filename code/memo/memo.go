package memo

type Memo struct {
	hashes map[string]interface{}
}

func New() *Memo {
	return &Memo{
		hashes: make(map[string]interface{}),
	}
}
