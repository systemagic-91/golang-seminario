package main

// referencia -> https://youtu.be/37uG6C7RmkI

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/yeka/zip"
)

const (
	zipPath  = "./assets/juda.zip"
	passPath = "./assets/rockyou.txt"

	numeroDeThreads = 2
	linhasPorThread = 7172196

	// numeroDeThreads = 4
	// linhasPorThread = 3586098

	// numeroDeThreads = 8
	// linhasPorThread = 1793049
)

func main() {

	abrirArquivoZip()
}

func abrirArquivoZip() {

	arquivos, err := zip.OpenReader(zipPath)

	if err != nil {
		panic(err.Error())
	}

	defer arquivos.Close()

	arquivo := arquivos.File[0]

	if arquivo.IsEncrypted() {

		listaSenhas := obterListaDeSenhas(passPath)

		// canais sao usados para comunicações entre as threads
		// nesse canal cabe 1 resultado
		// quando ele receber o primeiro resultado todas as outras threads vão parar.
		canal := make(chan string, 1)

		start := time.Now()

		linhaInicial := 0

		for i := 0; i < numeroDeThreads; i++ {

			linhaFinal := linhasPorThread * (i + 1)

			fmt.Printf("\nIniciando a Thread: %d ", (i + 1))
			fmt.Printf("Lendo da linha %d até a linha %d", linhaInicial, linhaFinal)
			go bruteForce(zipPath, listaSenhas[linhaInicial:linhaFinal], canal)

			linhaInicial = linhaFinal + 1
		}

		// esperando até o canal receber uma msg no canal
		select {
		case senha := <-canal:
			fmt.Printf("\n\nSenha encontrada: %v\n", senha)
			fmt.Printf("\nTempo brute-force: %vseg\n", time.Since(start).Seconds())
		case <-time.After(time.Duration(35) * time.Second):
			fmt.Println("\nSenha NÃO encontrada! \nTimeout.")
		}
	}
}

func obterListaDeSenhas(caminhoArquivo string) (senhas []string) {

	arquivo, err := os.Open(caminhoArquivo)

	if err != nil {
		panic(err.Error())
	}

	start := time.Now()

	scanner := bufio.NewScanner(arquivo)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		senhas = append(senhas, scanner.Text())
	}

	fmt.Printf("\nTempo de leitura arquivo de senhas: %v\n", time.Since(start))

	defer arquivo.Close()

	fmt.Printf("\nQuantidade de senhas no arquivo: %v\n", len(senhas))

	return
}

func bruteForce(caminhoZip string, listaSenhas []string, canal chan<- string) {

	arquivos, err := zip.OpenReader(caminhoZip)

	if err != nil {
		panic(err.Error())
	}

	defer arquivos.Close()

	arquivo := arquivos.File[0]

	for _, senha := range listaSenhas {

		arquivo.SetPassword(string(senha))

		reader, err := arquivo.Open()

		if err != nil {
			panic(err.Error())
		}

		buf, err := ioutil.ReadAll(reader)

		if err != nil {
			continue
		}

		defer reader.Close()

		if buf != nil {
			canal <- string(senha)
			break
		}
	}
}
