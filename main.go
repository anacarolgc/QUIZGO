package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Text          string
	Options       []string
	AnswerCorrect int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Seja bem-vindo(a) ao Quiz!")
	fmt.Print("Digite seu nome: ")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		panic("Erro ao ler a entrada")
	}
	g.Name = strings.TrimSpace(name)
	fmt.Printf("Vamos ao jogo, %s!\n", g.Name)
}

func (g *GameState) ProcessCsv() {
	file, err := os.Open("quizgo.csv")
	if err != nil {
		panic("Erro ao abrir o arquivo CSV")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Erro ao ler o CSV")
	}

	for index, record := range records {
		if index == 0 {
			continue // Ignorar cabeçalho
		}
		correctAnswer, err := strconv.Atoi(record[5])
		if err != nil {
			fmt.Println("Erro ao converter resposta correta:", err)
			continue
		}
		question := Question{
			Text:          record[0],
			Options:       record[1:5],
			AnswerCorrect: correctAnswer,
		}
		g.Questions = append(g.Questions, question)
	}
}

func (g *GameState) Run() {
	for i, question := range g.Questions {
		fmt.Printf("\n\033[33m%d. %s\033[0m\n", i+1, question.Text)
		for j, option := range question.Options {
			fmt.Printf("[%d] %s\n", j+1, option)
		}
		fmt.Print("Digite uma alternativa: ")

		answerCh := make(chan int)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			answer, err := strconv.Atoi(input)
			if err != nil {
				answerCh <- -1
			} else {
				answerCh <- answer
			}
		}()

		select {
		case answer := <-answerCh:
			if answer == question.AnswerCorrect {
				fmt.Println("Parabéns, você acertou!")
				g.Points += 10
			} else {
				fmt.Println("Você errou!")
			}
		case <-time.After(10 * time.Second):
			fmt.Println("\nTempo esgotado! O jogo foi encerrado.")
			return
		}
	}
}

func main() {
	game := &GameState{Points: 0}
	game.ProcessCsv()
	game.Init()
	game.Run()
	fmt.Printf("\nFim de Jogo! Você fez %d pontos.\n", game.Points)
}
