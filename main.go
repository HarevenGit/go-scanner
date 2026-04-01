package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
)

func main() {
	// 1. Настраиваем "крутилки" (Флаги)
	// Имя флага, значение по умолчанию, описание
	target := flag.String("t", "scanme.nmap.org", "Target IP/Domain")
	threads := flag.Int("w", 100, "Number of workers (threads)")
	maxPort := flag.Int("p", 100, "Scan ports from 1 to P")

	flag.Parse() // Считываем ввод пользователя

	var openPorts []int
	var wg sync.WaitGroup
	var resultsWg sync.WaitGroup

	jobs := make(chan int, 100)
	results := make(chan int)

	// Запускаем воркеров (количество берем из флага -w)
	for i := 1; i <= *threads; i++ {
		go worker(i, *target, jobs, results, &wg)
	}

	// Сборщик результатов
	resultsWg.Add(1)
	go func() {
		for p := range results {
			openPorts = append(openPorts, p)
		}
		resultsWg.Done()
	}()

	// Раздаем задачи (до порта, указанного в -p)
	for i := 1; i <= *maxPort; i++ {
		wg.Add(1)
		jobs <- i
	}

	close(jobs)
	wg.Wait()
	close(results)
	resultsWg.Wait()

	// СОРТИРОВКА перед выводом
	sort.Ints(openPorts)

	fmt.Printf("\n--- Scan Results for %s ---\n", *target)
	fmt.Println("Open Ports:", openPorts)
}

func worker(id int, target string, jobs <-chan int, results chan int, wg *sync.WaitGroup) {
	for port := range jobs {
		address := fmt.Sprintf("%s:%d", target, port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			wg.Done()
			continue
		}
		conn.Close()
		results <- port
		wg.Done()
	}
}
