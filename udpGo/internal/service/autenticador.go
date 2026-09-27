package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Usuario struct {
	Nome       string
	Email      string
	TokenAtual string
}

type AutenticadorService struct {
	mu       sync.RWMutex
	usuarios map[string]*Usuario // Chave: Email
}

func NovoAutenticadorService() *AutenticadorService {
	return &AutenticadorService{
		usuarios: make(map[string]*Usuario),
	}
}

// IniciarRelogio inicia o loop em background que gera tokens a cada 15s
func (s *AutenticadorService) IniciarRelogio() {
	ticker := time.NewTicker(15 * time.Second)

	// Goroutine para atualizar tokens periodicamente
	go func() {
		for agora := range ticker.C {
			s.mu.Lock()
			if len(s.usuarios) > 0 {
				fmt.Printf("\n[RELOGIO %s] Atualizando tokens de %d usuario(s)...\n",
					agora.Format("15:04:05"), len(s.usuarios))

				for _, u := range s.usuarios {
					u.TokenAtual = s.gerarHash(u.Email, agora)
					fmt.Printf(" -> [%s] %s: %s\n", u.Nome, u.Email, u.TokenAtual)
				}
			}
			s.mu.Unlock()
		}
	}()
}

func (s *AutenticadorService) Cadastrar(nome, email string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, existe := s.usuarios[email]; existe {
		return "ERRO: usuario ja cadastrado"
	}

	agora := time.Now()
	token := s.gerarHash(email, agora)

	s.usuarios[email] = &Usuario{
		Nome:       nome,
		Email:      email,
		TokenAtual: token,
	}

	return fmt.Sprintf("SUCESSO: usuario %s cadastrado! Token inicial: %s", nome, token)
}

func (s *AutenticadorService) ObterToken(email string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, existe := s.usuarios[email]
	if !existe {
		return "usuário sem cadastro"
	}
	return fmt.Sprintf("TOKEN:%s", u.TokenAtual)
}

func (s *AutenticadorService) gerarHash(email string, t time.Time) string {
	intervalo := t.Unix() / 15
	dados := fmt.Sprintf("%s:%d", email, intervalo)
	h := sha256.New()
	h.Write([]byte(dados))
	return hex.EncodeToString(h.Sum(nil))[:8] // Primeiros 8 caracteres
}
