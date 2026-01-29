package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	// 01 Ler flags
	port := flag.Int("port", 3000, "Port to run the caching proxy")
	origin := flag.String("origin", "", "Origin server URL")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache")
	flag.Parse()

	// 02. Criar cache em memória
	type CachedResponse struct {
		Body	[]byte
		Header	http.Header
		Status int
	}
	var cache = make(map[string]CachedResponse)
	var cacheMutex sync.RWMutex // para acessar cache de forma segura em concorrência

	if *clearCache {
		cacheMutex.Lock()
		cache = make(map[string]CachedResponse)
		cacheMutex.Unlock()
		fmt.Println("Cache cleared!")
		return
	}

	fmt.Printf("Starting caching proxy on port %d forwarding to %s\n", *port, *origin)

	// 03. Função para lidar com todas as requisições
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cacheKey := r.Method + " " + r.URL.RequestURI()

		// 01. Tentar responder do cache
		cacheMutex.RLock()
		cached, found := cache[cacheKey]
		cacheMutex.RUnlock()

		if found {
			for key, values := range cached.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(cached.Status)
			w.Write(cached.Body)
			return
		}

		// 02. Se não está no cache, encaminhar para origin
		fullURL := *origin + r.URL.Path
		if r.URL.RawQuery != "" {
			fullURL += "?" + r.URL.RawQuery
		}

		resp, err := http.Get(fullURL)
		if err != nil {
			http.Error(w, "Failed to reach origin", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read origin response", http.StatusInternalServerError)
			return
		}

		// Adicionar X-Cahe: MISS
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)

		cacheMutex.Lock()
		cache[cacheKey] = CachedResponse{
			Body:	body,
			Header: resp.Header.Clone(),
			Status: resp.StatusCode,
		}
		cacheMutex.Unlock()
	})

	// 04. Iniciar servidor
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}