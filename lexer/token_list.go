package lexer

type TokenList struct {
	index int
	list  []*Lexer
}

func (list *TokenList) Current() *Lexer {
	if list.index < 0 || list.Empty() {
		return nil
	}
	return list.list[list.index]
}

func (list *TokenList) Next() bool {
	list.index++
	return list.index < len(list.list)
}

func (list *TokenList) Empty() bool {
	return list.index >= len(list.list)
}
