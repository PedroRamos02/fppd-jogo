package main

import (
	"os"
	"time"
	"math/rand"
)

func main() {
	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	interfaceDesenharJogo(&jogo)

	// Goroutines para inimigos
	for i := range jogo.Inimigos {
		inimigo := &jogo.Inimigos[i]
		go func(inimigo *InimigoMovel) {
			rand.Seed(time.Now().UnixNano())
			for {
				if rand.Float64() < 0.3 {
					inimigo.Direita = !inimigo.Direita
				}
				moverInimigo(inimigo, &jogo)
				interfaceDesenharJogo(&jogo)
				time.Sleep(300 * time.Millisecond)
			}
		}(inimigo)
	}

	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}
