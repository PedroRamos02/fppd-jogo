// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"fmt"
	"os"
	"time"
	"github.com/nsf/termbox-go"
)


// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
const (
	CorPadrao     Cor = termbox.ColorDefault
	CorCinzaEscuro    = termbox.ColorDarkGray
	CorVermelho       = termbox.ColorRed
	CorVerde          = termbox.ColorGreen
	CorParede         = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede    = termbox.ColorDarkGray
	CorTexto          = termbox.ColorDarkGray
	CorBranco = termbox.ColorWhite
)

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Key == termbox.KeyEnter {
		return EventoTeclado{Tipo: "atirar"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
} 

// Renderiza todo o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()
	interfaceDesenharTexto(0, len(jogo.Mapa)+1, fmt.Sprintf("⏱ Tempo restante: %ds", tempoRestante), CorBranco, CorPadrao)


	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// Desenha o personagem sobre o mapa
	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)

	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo)

	for _, inimigo := range jogo.Inimigos {
		interfaceDesenharElemento(inimigo.X, inimigo.Y, Inimigo)
	}
	// Força a atualização do terminal
	interfaceAtualizarTela()
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal com os dados desenhados
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

func verificarVitoria(jogo *Jogo) {
	for _, inimigo := range jogo.Inimigos {
		if inimigo.Ativo {
			return // Ainda há inimigos ativos
		}
	}

	interfaceLimparTela()
	interfaceDesenharVitoria()
	os.Exit(0) // Finaliza o programa
}

func interfaceDesenharVitoria() {
	println()
	println("██      ██  ██████   ██████  ██████      ██      ██ ███████ ███    ██  ██████  ███████  ██    ██ ")
	println(" ██    ██  ██    ██ ██       ██           ██    ██  ██      ████   ██ ██       ██       ██    ██ ")
	println("  ██  ██   ██    ██ ██       █████         ██  ██   █████   ██ ██  ██ ██       █████    ██    ██ ")
	println("   █  █    ██    ██ ██       ██             █  █    ██      ██  ██ ██ ██       ██       ██    ██ ")
	println("    ██      ██████   ██████  ██████          ██     ███████ ██   ████  ██████  ███████   ██████  ")
	println()
	println("                       Parabéns! Todos os inimigos foram derrotados.                             ")
}


// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	// Linha de status dinâmica
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}

	// Instruções fixas
	msg := "Use WASD para mover e Enter para disparar. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
	}
}

var tempoRestante = 60

func iniciarTimerVisual(jogo *Jogo) {
	go func() {
		for tempoRestante >= 0 {
			interfaceDesenharJogo(jogo) // redesenha tudo com tempo
			time.Sleep(1 * time.Second)
			tempoRestante--

			if tempoRestante < 0 {
				interfaceLimparTela()
				fmt.Println(`
   _____                        ____                 
  / ____|                      / __ \                
 | |  __  __ _ _ __ ___   ___ | |  | |_   _____ _ __ 
 | | |_ |/ _` + "`" + ` | '_ ` + "`" + ` _ \ / _ \| |  | \ \ / / _ \ '__|
 | |__| | (_| | | | | | | (_) | |__| |\ V /  __/ |   
  \_____|\__,_|_| |_| |_|\___/ \____/  \_/ \___|_|   

                 TEMPO ESGOTADO!
`)
				os.Exit(0)
			}
		}
	}()
}

func interfaceDesenharTexto(x, y int, texto string, corTexto, corFundo Cor) {
	for i, c := range texto {
		termbox.SetCell(x+i, y, c, corTexto, corFundo)
	}
}

