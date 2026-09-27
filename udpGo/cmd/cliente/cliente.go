package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"udpGo/internal/comunicador"
)

func main() {
	cliente, err := comunicador.NovoCliente("127.0.0.1", 8080)
	if err != nil {
		log.Fatalf("Falha ao conectar cliente: %v", err)
	}
	defer cliente.Fechar()

	fmt.Println(cliente)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Cliente UDP Conectado ao Servidor.")
	fmt.Println("Exemplos de entrada:")
	fmt.Println("  1. CADASTRAR;usuario;email")
	fmt.Println("  2. TOKEN;email")

	for {
		fmt.Print("\nEsperando comando > ")
		if !scanner.Scan() {
			break
		}
		texto := strings.TrimSpace(scanner.Text())
		if texto == "sair" {
			break
		}
		if texto == "" {
			continue
		}

		// Envia datagrama
		if err := cliente.Enviar(texto); err != nil {
			fmt.Printf("Erro ao enviar: %v\n", err)
			continue
		}

		resposta, err := cliente.ReceberComTimeout(2 * time.Second)
		if err != nil {
			fmt.Printf("Sem resposta do servidor (Timeout ou erro): %v\n", err)
			continue
		}

		fmt.Printf("Servidor respondeu: %s\n", resposta)
	}
}
