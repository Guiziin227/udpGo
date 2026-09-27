//package main
//
//import "testing"
//
//type Objeto struct {
//	ID int
//}
//
//func BenchmarkLoopStack(b *testing.B) {
//	for i := 0; i < b.N; i++ { // b.N será controlado pelo Go, mas faremos o loop interno de 10M
//		for j := 0; j < 10_000_000_0; j++ {
//			obj := Objeto{ID: j} // Alocado na Stack
//			_ = obj.ID           // Apenas usa a variável
//		}
//	}
//}
//
//func BenchmarkLoopHeap(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		for j := 0; j < 10_000_000_0; j++ {
//			obj := &Objeto{ID: j} // Alocado na Heap por causa do ponteiro '&'
//			_ = obj.ID
//		}
//	}
//}

package main

import "testing"

type ObjetoPequeno struct {
	ID int
}

var ResultadoGlobal *ObjetoPequeno
var ResultadoValor int

func BenchmarkLoopStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10_000_000; j++ {
			obj := ObjetoPequeno{ID: j}
			ResultadoValor = obj.ID // Joga para a global de valor
		}
	}
}

func BenchmarkLoopHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10_000_000; j++ {
			obj := &ObjetoPequeno{ID: j}
			ResultadoGlobal = obj // GARANTE O ESCAPE
		}
	}
}
