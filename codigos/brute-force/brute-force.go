package main

// referencia -> https://youtu.be/37uG6C7RmkI

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yeka/zip"
)

const (
	zipPath  = "./assets/bat.zip"
	passPath = "./assets/rockyou.txt"
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

		start := time.Now()

		for i := 0; i < len(listaSenhas); i++ {

			senhaAtual := listaSenhas[i]

			arquivo.SetPassword(senhaAtual)

			r, err := arquivo.Open()

			if err != nil {
				panic(err.Error())
			}

			defer r.Close()

			buf, err := io.ReadAll(r)

			if err != nil {
				continue
			}

			if buf != nil {
				fmt.Printf("\nSenha encontrada: %v\n", senhaAtual)
				fmt.Printf("localizada na linha: %v\n", i+1)
				fmt.Printf("\nTempo brute-force: %v\n", time.Since(start).Seconds())
				break
			}
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
