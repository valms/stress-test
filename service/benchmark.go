package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

type Benchmark struct {
	url      string
	requests int
	threads  int
	client   *http.Client
	results  chan Result
	wg       sync.WaitGroup
	metrics  *Metrics
}

type Metrics struct {
	totalRequests int64
	statusCodes   map[int]int64
	errorsByCode  map[int]map[string]int64
	errors        int64
	totalTime     time.Duration
	mu            sync.RWMutex
}

func NewBenchmark(url string, requests, threads int) (*Benchmark, error) {
	if url == "" {
		return nil, errors.New("a URL é obrigatória")
	}
	if requests <= 0 {
		return nil, errors.New("o número de requisições deve ser maior que zero")
	}
	if threads <= 0 {
		return nil, errors.New("o número de threads deve ser maior que zero")
	}

	if threads > requests {
		threads = requests
		fmt.Printf("\nAjuste automático: número de threads reduzido para %d para corresponder ao número de requisições\n",
			threads)
	}

	return &Benchmark{
		url:      url,
		requests: requests,
		threads:  threads,
		results:  make(chan Result, requests),
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		metrics: &Metrics{
			statusCodes:  make(map[int]int64),
			errorsByCode: make(map[int]map[string]int64),
		},
	}, nil
}

func (b *Benchmark) worker(ctx context.Context) {
	defer b.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if atomic.LoadInt64(&b.metrics.totalRequests) >= int64(b.requests) {
				return
			}

			atomic.AddInt64(&b.metrics.totalRequests, 1)
			start := time.Now()

			resp, err := b.client.Get(b.url)
			result := Result{
				Duration: time.Since(start),
			}

			if err != nil {
				result.Error = err
				b.results <- result
				continue
			}

			if resp.Body != nil {
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
			}

			result.StatusCode = resp.StatusCode
			b.results <- result
		}
	}
}

func (b *Benchmark) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	fmt.Printf("\n🚀 Iniciando Teste de Performance\n")
	fmt.Printf("================================\n")
	fmt.Printf("URL: %s\n", b.url)
	fmt.Printf("Total de Requisições: %d\n", b.requests)
	fmt.Printf("Threads Concorrentes: %d\n\n", b.threads)

	b.wg.Add(b.threads)
	for i := 0; i < b.threads; i++ {
		go b.worker(ctx)
	}

	go func() {
		b.wg.Wait()
		close(b.results)
	}()

	for result := range b.results {
		b.metrics.mu.Lock()
		if result.Error != nil {
			b.metrics.errors++
			// Agrupa erros por tipo
			errorMsg := result.Error.Error()
			if b.metrics.errorsByCode[0] == nil {
				b.metrics.errorsByCode[0] = make(map[string]int64)
			}
			b.metrics.errorsByCode[0][errorMsg]++
		} else {
			b.metrics.statusCodes[result.StatusCode]++
			if result.Error != nil {
				if b.metrics.errorsByCode[result.StatusCode] == nil {
					b.metrics.errorsByCode[result.StatusCode] = make(map[string]int64)
				}
				b.metrics.errorsByCode[result.StatusCode][result.Error.Error()]++
			}
		}
		b.metrics.mu.Unlock()
	}

	b.metrics.totalTime = time.Since(start)
	return nil
}

func (b *Benchmark) PrintReport() {
	fmt.Printf("\n📊 Relatório de Performance\n")
	fmt.Printf("==========================\n\n")

	fmt.Printf("⏱️  Duração Total: %v\n", b.metrics.totalTime)
	fmt.Printf("📈 Requisições Executadas: %d\n", b.metrics.totalRequests)

	b.metrics.mu.RLock()
	defer b.metrics.mu.RUnlock()

	fmt.Printf("\n✅ Distribuição por Status HTTP:\n")
	for code, count := range b.metrics.statusCodes {
		fmt.Printf("   Status %d: %d requisições\n", code, count)

		if errMap, exists := b.metrics.errorsByCode[code]; exists && len(errMap) > 0 {
			fmt.Printf("   ⚠️  Erros associados ao Status %d:\n", code)
			for errMsg, errCount := range errMap {
				fmt.Printf("      • %s: %d ocorrência(s)\n", errMsg, errCount)
			}
		}
	}

	if errMap, exists := b.metrics.errorsByCode[0]; exists && len(errMap) > 0 {
		fmt.Printf("\n❌ Erros de Conexão/Rede:\n")
		for errMsg, errCount := range errMap {
			fmt.Printf("   • %s: %d ocorrência(s)\n", errMsg, errCount)
		}
	}

	fmt.Printf("\n📈 Resumo Final:\n")
	fmt.Printf("   • Requisições bem-sucedidas (2xx): %d\n", b.metrics.statusCodes[200])
	fmt.Printf("   • Total de falhas: %d\n", b.metrics.errors)

	successRate := float64(b.metrics.statusCodes[200]) / float64(b.metrics.totalRequests) * 100
	fmt.Printf("   • Taxa de sucesso: %.2f%%\n", successRate)

	if successRate < 95 {
		fmt.Printf("\n⚠️  Recomendações:\n")
		if b.metrics.errors > 0 {
			fmt.Printf("   • Verifique a conectividade com o endpoint\n")
			fmt.Printf("   • Considere aumentar o timeout das requisições\n")
		}
		if len(b.metrics.statusCodes) > 1 {
			fmt.Printf("   • Analise os status HTTP não-200 retornados\n")
		}
	}

	fmt.Printf("\nTeste concluído em %v\n", b.metrics.totalTime)
}
