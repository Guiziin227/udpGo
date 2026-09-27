package comunicador

import (
	"fmt"
	"net"
	"time"
)

// ComunicadorUDP gerencia o socket UDP
type ComunicadorUDP struct {
	conn *net.UDPConn
}

// NovoServidor inicializa o comunicador escutando em uma porta local
func NovoServidor(porta int) (*ComunicadorUDP, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", porta))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &ComunicadorUDP{conn: conn}, nil
}

// NovoCliente inicializa o comunicador apontando para um servidor remoto
func NovoCliente(host string, porta int) (*ComunicadorUDP, error) {
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, porta))
	if err != nil {
		return nil, err
	}
	// No cliente, laddr é nil para o SO escolher a porta local
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return &ComunicadorUDP{conn: conn}, nil
}

// Receber lê um datagrama (usado pelo Servidor)
// Retorna a mensagem em string, o endereço de quem enviou e eventual erro
func (c *ComunicadorUDP) Receber() (string, *net.UDPAddr, error) {
	buffer := make([]byte, 2048) // Buffer seguro para evitar estouros
	n, clientAddr, err := c.conn.ReadFromUDP(buffer)
	if err != nil {
		return "", nil, err
	}
	return string(buffer[:n]), clientAddr, nil
}

// EnviarPara envia resposta a um endereço específico (usado pelo Servidor)
func (c *ComunicadorUDP) EnviarPara(msg string, destino *net.UDPAddr) error {
	_, err := c.conn.WriteToUDP([]byte(msg), destino)
	return err
}

// Enviar envia mensagem para o endereço configurado no Dial (usado pelo Cliente)
func (c *ComunicadorUDP) Enviar(msg string) error {
	_, err := c.conn.Write([]byte(msg))
	return err
}

// ReceberComTimeout espera resposta com prazo máximo (usado pelo Cliente)
func (c *ComunicadorUDP) ReceberComTimeout(timeout time.Duration) (string, error) {
	c.conn.SetReadDeadline(time.Now().Add(timeout))
	buffer := make([]byte, 2048)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

// Fechar encerra o socket
func (c *ComunicadorUDP) Fechar() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
