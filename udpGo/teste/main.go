package main

type Dados struct {
	Valor int
}

func naStack() Dados {
	d := Dados{Valor: 42}
	return d
}

func naHeap() *Dados {
	d := &Dados{Valor: 42}
	return d
}

func main() {
	_ = naStack()
	_ = naHeap()
}
