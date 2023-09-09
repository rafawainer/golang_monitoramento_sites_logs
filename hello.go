package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)
import "os"

const monitoramentos = 3
const delay = 3 * time.Second

func exibeVazio() {
	fmt.Println("")
}

func checaEof(err error) bool {
	if err == io.EOF {
		//fmt.Println("Final do arquivo", err)
		return true
	}
	return false
}

func checaErro(err error) {
	if err != nil {
		fmt.Println("Error", err)
		sairProgramaError()
	}
}

func exibeIntroducao() {
	nome := "Wainer"
	versao := 1.1
	fmt.Println("Olá, sr.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func opcoeIniciais() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir logs")
	fmt.Println("0 - Sair")
}

func leComando() int {
	var comandoLido int
	scan, err := fmt.Scan(&comandoLido)
	checaErro(err)
	fmt.Println("\nO comandoLido escolhido foi:", scan)

	return comandoLido
}

func sairProgramaSucesso() {
	os.Exit(0)
}

func sairProgramaError() {
	os.Exit(500)
}

func testaSite(logFile *os.File, site string) int {
	resp, err := http.Get(site)
	//fmt.Println("Resultado do resp:", resp)
	checaErro(err)

	switch true {
	case resp.StatusCode == 200:
		println("Site", site, "carregado com sucesso !!")
		registraLog(logFile, site, true)
	default:
		println("Site", site, "com problema. Status code:",
			resp.StatusCode)
		registraLog(logFile, site, false)
	}
	return resp.StatusCode
}

func iniciarMonitoramento() {
	fmt.Println("\nIniciando monitoramento...")

	arquivo := "hello/logs.txt"
	logFile := createReadFile(arquivo)

	sites := leSitesArquivo()

	for i := 0; i < monitoramentos; i++ {
		exibeVazio()
		fmt.Println(i+1, "execucao de teste de Monitoramento")
		exibeVazio()
		for i, v := range sites {
			fmt.Println("Iniciando o", i+1, "site:", v)

			testaSite(logFile, v)
		}
		time.Sleep(delay)
	}
	exibeVazio()
	closeFile(logFile)
}

func leSitesArquivo() []string {
	var sites []string
	arquivo := "hello/sites.txt"
	reader := createReadFile(arquivo)
	//arquivo, err := ioutil.ReadFile("hello/sites.txt") //Exibir todo o conteudo do arquivo

	leitor := bufio.NewReader(reader)

	for {
		res, err := leitor.ReadString('\n')
		if checaEof(err) {
			break
		}
		res = strings.TrimSpace(res)
		sites = append(sites, res)
	}

	closeFile(reader)
	//fmt.Println("Registros do arquivo:", sites)
	return sites
}

func createReadFile(arquivo string) *os.File {
	open, err := os.OpenFile(arquivo, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checaErro(err)
	return open
}

func closeFile(arquivo *os.File) {
	err := arquivo.Close()
	checaErro(err)
	println(arquivo.Name(), "fechado com sucesso")
}

func registraLog(reader *os.File, site string, status bool) bool {

	horarioAtual := time.Now().Format("02/01/2006 15:04:05")
	_, err := reader.WriteString(horarioAtual + " " + site + " - online: " + strconv.FormatBool(status) + "\n")
	checaErro(err)

	return true
}

func main() {

	//Exibe introducao com nome
	exibeIntroducao()

	for {
		// Funcao para exibir as opcoes iniciais
		opcoeIniciais()

		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 0:
			fmt.Println("\nSaindo do programa...")
			sairProgramaSucesso()
		default:
			fmt.Println("\nNão reconheço a opçao selecionada")
			os.Exit(-1)
		}
	}
}

func imprimeLogs() {
	fmt.Println("\nExibindo logs.....")
	arquivo := "hello/logs.txt"

	content, err := ioutil.ReadFile(arquivo)
	checaErro(err)

	println(string(content))
}
