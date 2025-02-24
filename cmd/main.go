package main

import (
	"flag"
	"github.com/valms/stress-test/service"
	"log"
)

func main() {
	// Leitura das flags
	url := flag.String("url", "",
		"URL do endpoint a ser testado (obrigatório)")

	requests := flag.Int("requests", 0,
		"Número total de requisições a serem executadas (obrigatório)")

	threads := flag.Int("threads", 0,
		"Número de threads concorrentes para execução do teste. (obrigatório)\n"+
			"\tSe maior que o número de requisições, será limitado ao número de requisições.")

	flag.Parse()

	benchmark, err := service.NewBenchmark(*url, *requests, *threads)

	if err != nil {
		log.Fatal(err)
	}

	if err := benchmark.Run(); err != nil {
		log.Fatal(err)
	}

	benchmark.PrintReport()
}
