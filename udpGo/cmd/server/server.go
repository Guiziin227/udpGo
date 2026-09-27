package main

import (
	"fmt"
	"log"
	"strings"
	"udpGo/internal/comunicador"
	"udpGo/internal/service"
)

func main() {
	com, err := comunicador.NovoServidor(8080)
	if err != nil {
		log.Fatalf("Falha ao abrir servidor UDP: %v", err)
	}
	defer com.Fechar()

	authService := service.NovoAutenticadorService()
	authService.IniciarRelogio()

	fmt.Println("=== Servidor UDP de Autenticador Ativo na Porta 8080 ===")

	for {
		msg, clientAddr, err := com.Receber()
		if err != nil {
			fmt.Printf("Erro ao receber: %v\n", err)
			continue
		}

		partes := strings.Split(strings.TrimSpace(msg), ";")
		comando := strings.ToUpper(partes[0])
		var resposta string

		switch comando {
		case "CADASTRAR":
			// Formato esperado: CADASTRAR;Nome;email@email.com
			if len(partes) < 3 {
				resposta = "ERRO: formato invalido. Use: CADASTRAR;nome;email"
			} else {
				resposta = authService.Cadastrar(partes[1], partes[2])
			}

		case "TOKEN":
			// Formato esperado: TOKEN;email@email.com
			if len(partes) < 2 {
				resposta = "ERRO: formato invalido. Use: TOKEN;email"
			} else {
				resposta = authService.ObterToken(partes[1])
			}

		default:
			resposta = "ERRO: comando desconhecido. Comandos validos: CADASTRAR, TOKEN"
		}

		if err := com.EnviarPara(resposta, clientAddr); err != nil {
			fmt.Printf("Erro ao responder para %s: %v\n", clientAddr, err)
		}
	}
}
