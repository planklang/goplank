package lexer

type TokenList struct {
	index int
	list  []*Lexer
}

func (list *TokenList) Peek() (*Lexer, bool) {
	if len(list.list) <= list.index {
		return nil, false
	}
	return list.list[list.index+1], true
}

func (list *TokenList) Current() *Lexer {
	return list.list[list.index]
}

func (list *TokenList) Next() (*Lexer, bool) {
	if len(list.list) <= list.index {
		return nil, false
	}
	list.index++
	return list.list[list.index], true
}

func (list *TokenList) Empty() bool {
	return len(list.list) <= list.index
}
