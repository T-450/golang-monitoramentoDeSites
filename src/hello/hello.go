package main

import (
	"net/http"
	"fmt"
	"os"
	"time"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
)
import "io"

const monitoramentos = 3
const delay = 5 * time.Second

func main() {
	exibeIntroducao()
	registraLog("site-falso", false)
	for {
		exibeMenu()
		comandoLido := leComando()
		
		switch comandoLido {
			case 0: 
				fmt.Println("Saindo do programa...")
				os.Exit(0)
			case 1:
				iniciarMonitoramento()
			case 2:
				imprimeLogs()
			default: 
				fmt.Println("Comando não reconhecido")
				os.Exit(-1);
		}
	}

}

func leComando() int {
	var comando int
	fmt.Scan(&comando)
	return comando;
}

func exibeIntroducao() {
	nome := "Edward"
	versao := 1.1

	fmt.Println("Olá, sr.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs");
	fmt.Println("0 - Sair do Programa")
}

func iniciarMonitoramento() {
	sites := leSitesDoArquivo()
	fmt.Println("Monitorando...")
	for i := 0; i < monitoramentos; i++ {
		for i, site := range sites {
			fmt.Println("Testando site", i, ":", site);
			testaSite(site);
		}
		time.Sleep(delay)
		fmt.Println("");
	}
}

func testaSite(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site ", site," foi carregado com sucesso")
		registraLog(site, true);
	} else {
		fmt.Println("Site ", site," esta com problema. Status Code", resp.StatusCode);
		registraLog(site, false)
	}	
}

func registraLog(site string, status bool) {
	// var str_status string
	// if (status) { 
	// 	str_status = "- online"
	// } else { 
	// 		str_status = "- offline"
	// }
	// do not work in golang:
	// str_status := status ? "online" : "offline"
	arquivo, errOpen := os.OpenFile(
		"log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	
		// https://golang.org/src/time/format.go
	_, err := arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + "- online " + strconv.FormatBool(status) + "\n")

	if err != nil {
		fmt.Println(err)
	} 
	if errOpen != nil {
		fmt.Println(errOpen)
	}
	arquivo.Close()
	fmt.Println(arquivo)
}

func leSitesDoArquivo() []string {
	var sites []string
	arquivo, _ := os.Open("sites.txt")
	// arquivo, err := ioutil.ReadFile("sites.txt")
	// if err != nil {
	// 	fmt.Println("Ocorreu um erro", err)
	// }
	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Ocorreu um erro", err)
		}
		sites = append(sites, strings.TrimSpace(linha))
	}
	arquivo.Close()
	return sites;
}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro ao ler os logs", err)
		os.Exit(-1)
	}
	
	fmt.Println(string(arquivo));
}
