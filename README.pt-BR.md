# 🚀 Zomboid Workshop API

[![English README](https://img.shields.io/badge/Language-English-blue)](README.md)
[![Português README](https://img.shields.io/badge/Idioma-Português-green)](README.pt-BR.md)

Microserviço de API ultrarrápido construído em **Go** e **Redis** para indexar, armazenar e consultar o mapeamento entre **Mod IDs** e **Workshop IDs** da Oficina do Steam para o Project Zomboid.

---

## 📌 Funcionalidades

- **Consultas Ultrarrápidas**: Desenvolvido em Go 1.22 com armazenamento em memória via Redis.
- **Suporte a CORS**: Pronto para integração com aplicações desktop e web, como o [Project Zomboid Mod Manager (PZMM)](https://github.com/WiliamMelo01/ProjectZomboidModManager).
- **Ingestão em Lote (Bulk)**: Endpoint `/mappings/bulk` otimizado para scanners multithreaded de mods.
- **Pronto para Servidores Linux**: Inclui arquivo de serviço `systemd` configurado para produção no Ubuntu/Debian.

---

## 🛠️ Endpoints da API

### 1. `GET /mappings`
Retorna todos os mapeamentos de Mod ID para Workshop ID cadastrados no formato JSON.

**Exemplo de Resposta:**
```json
{
  "Hydrocraft": "514427485",
  "Base.Vehicle": "2948504608",
  "Brita": "2460154811"
}
```

---

### 2. `POST /mappings/bulk`
Envia um lote de novos mapeamentos para salvar no Redis Hash em uma única requisição.

**Exemplo de Requisição:**
```json
{
  "mappings": [
    { "modId": "Hydrocraft", "workshopId": "514427485" },
    { "modId": "Brita", "workshopId": "2460154811" }
  ]
}
```

**Exemplo de Resposta:**
```json
{
  "status": "ok",
  "processed": 2
}
```

---

## 🚀 Como Executar Localmente

### 1. Requisitos
- **Go**: 1.22 ou superior
- **Redis**: 6.0 ou superior

### 2. Rodar a Aplicação
```bash
go run main.go
```

Por padrão, o serviço rodará na porta `8080` (configurável via variável de ambiente `PORT`).

---

## ⚙️ Variáveis de Ambiente

| Variável | Padrão | Descrição |
|---|---|---|
| `PORT` | `8080` | Porta onde o servidor HTTP escutará |
| `REDIS_ADDR` | `localhost:6379` | Endereço e porta do servidor Redis |

---

## 📄 Licença
Licença MIT
