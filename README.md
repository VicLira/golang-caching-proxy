# caching-proxy

Proxy HTTP em Go que encaminha requisições para um servidor de origem, cacheia as respostas e retorna `X-Cache: HIT` ou `X-Cache: MISS`. Permite limpar o cache via CLI e serve qualquer endpoint de forma genérica.

---

## Funcionalidades

* Encaminha requisições GET para o servidor de origem.
* Cacheia respostas em memória.
* Adiciona header `X-Cache: HIT` ou `X-Cache: MISS`.
* Limpeza do cache via CLI (`--clear-cache`).
* Funciona para qualquer URL de forma genérica.

---

## Passo a passo

1. **Clonar o repositório**

```bash
git clone <url-do-repo>
cd caching-proxy
```

2. **Rodar o proxy**

```bash
go run main.go --port 3000 --origin http://dummyjson.com
```

3. **Testar a primeira requisição**

```bash
curl -i http://localhost:3000/products
```

* Esperado: `X-Cache: MISS`

4. **Testar uma requisição repetida**

```bash
curl -i http://localhost:3000/products
```

* Esperado: `X-Cache: HIT`

5. **Limpar o cache**

```bash
go run main.go --clear-cache
```

6. **Testar novamente após limpar cache**

* A próxima requisição retornará `X-Cache: MISS` novamente.

---

## Observações

* O proxy é **genérico** e funciona para qualquer endpoint GET.
* Cache em memória → reiniciar o servidor limpa automaticamente o cache.
* Ideal para **dados que não mudam constantemente**, como listas de produtos ou categorias.
* Essa é a resolução do [projeto](https://roadmap.sh/projects/caching-server) do site roadmap.sh.