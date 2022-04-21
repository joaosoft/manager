package errors

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func NewBlock(try func(), catch func(Exception), finally func()) *Block {
	return &Block{
		Try:     try,
		Catch:   catch,
		Finally: finally,
	}
}

func (block Block) Do() {
	if block.Finally != nil {
		defer block.Finally()
	}
	if block.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				block.Catch(r)
			}
		}()
	}
	block.Try()
}
