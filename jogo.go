package main

import (
	"bufio"
	"os"
	"time"
)

type Elemento struct {
	simbolo   rune
	cor       Cor
	corFundo  Cor
	tangivel  bool
}

type InimigoMovel struct {
	X, Y     int
	Direita  bool
	Ativo    bool
	Cor      Cor  // Cor agora pertence ao inimigo
}

type Tiro struct {
	X, Y int
}

type Jogo struct {
	Mapa           [][]Elemento
	PosX, PosY     int
	UltimoVisitado Elemento
	StatusMsg      string
	Inimigos       []InimigoMovel
	Tiros          []Tiro
}

var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
	Projeteis  = Elemento{'*', CorVerde, CorPadrao, false}
)

func jogoNovo() Jogo {
	return Jogo{UltimoVisitado: Vazio}
}

func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				jogo.Inimigos = append(jogo.Inimigos, InimigoMovel{
					X: x, Y: y, Direita: true, Ativo: true, Cor: CorVermelho, // Inicializa com a cor vermelha
				})
				e = Vazio
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}
	if jogo.Mapa[y][x].tangivel {
		return false
	}
	return true
}

func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy
	elemento := jogo.Mapa[y][x]
	jogo.Mapa[y][x] = jogo.UltimoVisitado
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]
	jogo.Mapa[ny][nx] = elemento
}

func removerInimigoNaPosicao(jogo *Jogo, x, y int) {
	for i := 0; i < len(jogo.Inimigos); i++ {
		if jogo.Inimigos[i].X == x && jogo.Inimigos[i].Y == y {
			jogo.Inimigos = append(jogo.Inimigos[:i], jogo.Inimigos[i+1:]...)
			break
		}
	}
	jogo.Mapa[y][x] = Vazio
}

func moverInimigo(inimigo *InimigoMovel, jogo *Jogo) {
	if !inimigo.Ativo {
		// Garante que a posição no mapa seja esvaziada caso o inimigo tenha sido desativado
		if jogo.Mapa[inimigo.Y][inimigo.X].simbolo == Inimigo.simbolo {
			jogo.Mapa[inimigo.Y][inimigo.X] = Vazio
		}
		return
	}

	dx := 1
	if !inimigo.Direita {
		dx = -1
	}
	nx := inimigo.X + dx
	ny := inimigo.Y

	if nx < 0 || nx >= len(jogo.Mapa[0]) {
		inimigo.Direita = !inimigo.Direita
		return
	}
	destino := jogo.Mapa[ny][nx]
	if destino.tangivel || (jogo.PosX == nx && jogo.PosY == ny) {
		inimigo.Direita = !inimigo.Direita
		return
	}

	jogo.Mapa[inimigo.Y][inimigo.X] = Vazio
	jogo.Mapa[ny][nx] = Inimigo
	inimigo.X = nx
	inimigo.Y = ny
}

func atirar(jogo *Jogo) {
	tiro := Tiro{X: jogo.PosX, Y: jogo.PosY - 1}
	jogo.Tiros = append(jogo.Tiros, tiro)

	go func() {
		for {
			if tiro.Y < 0 || tiro.Y >= len(jogo.Mapa) {
				return
			}

			// Verifica colisão com inimigo
			for i := range jogo.Inimigos {
				if jogo.Inimigos[i].X == tiro.X && jogo.Inimigos[i].Y == tiro.Y && jogo.Inimigos[i].Ativo {
					jogo.Inimigos[i].Ativo = false
					jogo.Inimigos[i].Cor = CorCinzaEscuro // Muda a cor do inimigo para preto

					// Atualiza o inimigo no mapa
					jogo.Mapa[jogo.Inimigos[i].Y][jogo.Inimigos[i].X].cor = CorCinzaEscuro

					// Remove o tiro do mapa
					jogo.Mapa[tiro.Y][tiro.X] = Vazio

					interfaceDesenharJogo(jogo)
					verificarVitoria(jogo)
					return
				}
			}

			if jogo.Mapa[tiro.Y][tiro.X].tangivel {
				return
			}

			jogo.Mapa[tiro.Y][tiro.X] = Projeteis
			interfaceDesenharJogo(jogo)
			time.Sleep(100 * time.Millisecond)

			jogo.Mapa[tiro.Y][tiro.X] = Vazio
			tiro.Y--

			if tiro.Y < 0 {
				return
			}
		}
	}()
}
