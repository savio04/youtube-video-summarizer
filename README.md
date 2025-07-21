# 📺 YouTube Video Summarizer

[![Último Commit](https://img.shields.io/github/last-commit/savio04/youtube-video-summarizer?style=for-the-badge)](https://github.com/savio04/youtube-video-summarizer/commits/main)
[![Go](https://img.shields.io/badge/Go-1.21-blue?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Docker](https://img.shields.io/badge/Docker-%2300f.svg?logo=docker&logoColor=white&style=for-the-badge)](https://www.docker.com/)
[![Whisper](https://img.shields.io/badge/Whisper-OpenAI-blue?style=for-the-badge)](https://openai.com/research/whisper)
[![Live Demo](https://img.shields.io/badge/Demo-Online-brightgreen?logo=google-chrome&logoColor=white&style=for-the-badge)](https://yt.sistema-solar.fun)

---

## ✨ **Resumo do Projeto**  
O **YouTube Video Summarizer** é uma aplicação desenvolvida em **Go** que permite ao usuário inserir a URL de um vídeo do YouTube e obter um resumo textual do conteúdo.  

🔊 Utiliza o modelo **Whisper (OpenAI)** para transcrever o áudio automaticamente, mesmo para vídeos em outros idiomas, e gera um resumo inteligente da transcrição. É ideal para:  
- 📖 Consumir conteúdo longo de forma rápida.  
- 🧑‍💻 Pesquisas acadêmicas ou profissionais.  
- 📰 Produção de insights a partir de vídeos informativos.  

🌐 [Acesse a aplicação online](https://yt.sistema-solar.fun)  

---

## 📖 Sumário

- [📺 Visão Geral](#-visão-geral)
- [🚀 Funcionalidades](#-funcionalidades)
- [🛠 Tecnologias Utilizadas](#-tecnologias-utilizadas)
- [🏁 Como Começar](#-como-começar)
  - [📋 Pré-requisitos](#-pré-requisitos)
  - [⚡ Instalação](#-instalação)
  - [🚀 Execução](#-execução)

---

## 📺 Visão Geral

![Demo GIF](https://github.com/savio04/youtube-video-summarizer/blob/main/ui/public/gif.gif)

O **YouTube Video Summarizer** permite:  
✅ Inserir uma URL de vídeo do YouTube.  
✅ Extrair automaticamente a transcrição do áudio com suporte a múltiplos idiomas.  
✅ Gerar um resumo textual do vídeo, facilitando o consumo de conteúdos longos e densos.

---

## 🚀 Funcionalidades

✔️ Extração de legendas/transcrições usando o **Whisper (OpenAI)**  
✔️ Geração automática de resumo do conteúdo transcrito  
✔️ Suporte a múltiplos idiomas (dependendo das legendas disponíveis)  
✔️ Interface web amigável e responsiva  
✔️ Fácil configuração e execução com **Docker Compose**

---

## 🛠 Tecnologias Utilizadas

- [Go](https://go.dev)
- [Gin Framework](https://gin-gonic.com)
- [Whisper (OpenAI)](https://openai.com/research/whisper)
- [Docker & Docker Compose](https://www.docker.com/)
- [Redis](https://redis.io/)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- [ffmpeg](https://ffmpeg.org/)

---

## 🏁 Como Começar

### 📋 Pré-requisitos

Para rodar o projeto localmente, você precisa ter instalado:  

- [Git](https://git-scm.com)  
- [Go](https://go.dev/doc/install)  
- [Docker](https://docs.docker.com/engine/install)  
- [Docker Compose](https://docs.docker.com/compose/install)  
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)  
- [ffmpeg](https://ffmpeg.org/download.html)  

---

### ⚡ Instalação

Clone o repositório e configure o ambiente:  

```bash
# Clone o repositório
git clone git@github.com:savio04/youtube-video-summarizer.git

# Acesse a pasta do projeto
cd youtube-video-summarizer

# Copie o arquivo de variáveis de ambiente
cp .env-example .env

# Crie os containers Docker
sudo docker-compose up -d
```

### 🚀 Execução

Preencha o `cookies.txt` conforme instruções no arquivo para vídeos privados ou com restrição de idade.  

Inicie a API (backend):  

```bash
make run
```

### 🌐 Interface Web (Frontend)

A interface web é um HTML estático que pode ser aberto diretamente no navegador:

1. Navegue até a pasta `ui/public`.

2. Abra o arquivo `index.html` com o navegador de sua preferência:

```bash
# Exemplo com Google Chrome
google-chrome ui/public/index.html
```

Ou simplesmente dê dois cliques no arquivo index.html para abrir no navegador padrão.
