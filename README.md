# YouTube Video Summarizer

Este é um projeto em Go que permite buscar vídeos no YouTube por URL, extrair a transcrição do áudio utilizando o **Whisper**, e gerar um resumo de seu conteúdo.

[Demo](https://yt.savioaraujogomes.com)

![Animação do projeto](https://github.com/savio04/youtube-video-summarizer/blob/main/ui/public/gif.gif)

## Funcionalidades

- Extrai as legendas/transcrições de vídeos do YouTube utilizando **Whisper**.
- Gera um resumo das transcrições extraídas.
- Suporte a vídeos em múltiplos idiomas (desde que legendas ou transcrições estejam disponíveis).

## Execução do projeto 

### Pré-requisitos

Antes de começar, você vai precisar ter instalado em sua máquina as seguintes ferramentas:
- [Git](https://git-scm.com)
- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/ubuntu)
- [Docker Compose](https://docs.docker.com/compose/install)
- [Yt-dlp](https://github.com/yt-dlp/yt-dlp)
- [Ffmpeg](https://www.ffmpeg.org/download.html)

### 🎲 Executando o Projeto

```bash
# Clone este repositório
$ git clone git@github.com:savio04/youtube-video-summarizer.git

# Acesse a pasta do projeto no terminal/cmd
$ cd youtube-video-summarizer

# Preencha as envs baseado no exemplo .env.example
$ cp .env-example .env

# Crie os containers
$ sudo docker-compose up -d

# Veja o arquivo cookies.txt para instruções sobre como preencher os cookies

# Execute o projeto
$ make run
```
