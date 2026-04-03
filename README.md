# stereoctl

stereoctl é uma pequena CLI para detectar e corrigir incompatibilidades de áudio/vídeo/contêiner (por exemplo para compatibilidade com DaVinci Resolve Free).

Funcionalidades principais
- `convert`: converte/remuxa áudio para AAC estéreo
- `check`: analisa um arquivo e sugere ações por perfil
- `fix`: aplica um perfil (ex.: `resolve-free`) e converte/remuxa quando necessário

Instalação (rápida)

- Requisitos: Go (para desenvolvimento), `ffmpeg` e `ffprobe` (para execução/integration tests).
- Para instalar hooks locais (opcional):

```bash
bash scripts/install-lefthook.sh
```

Instalando dependências de sistema (Ubuntu/Debian):

```bash
sudo apt-get update
sudo apt-get install -y ffmpeg
```

Uso (exemplos)

- `convert` — converte/remuxa um arquivo:

```bash
stereoctl convert movie.mkv
stereoctl convert movie.mkv --output fixed.mp4 --bitrate 256k
```

- `check` — somente checa e sugere ações:

```bash
stereoctl check movie.mkv
```

- `fix` — aplica perfil `resolve-free` (padrão):

```bash
# modo normal: converte/remuxa
stereoctl fix movie.mkv

# modo preview: mostra o comando ffmpeg sem executar
stereoctl fix --preview movie.mkv

# modo batch: aceita diretório ou glob e processa vários arquivos
stereoctl fix --batch "*.mkv"
stereoctl fix --batch /path/to/videos
```

Flags importantes
- `--output, -o`: especifica path de saída (por padrão usa o mesmo nome com `.mp4`)
- `--bitrate, -b`: taxa de bits de áudio para conversão (`convert`)
- `--profile, -p`: perfil a aplicar (`fix`)
- `--preview, -n`: apenas imprime o comando `ffmpeg` sem executar (`fix`)
- `--batch, -B`: trata o argumento como diretório/glob e processa vários arquivos (`fix`)

Testes

Unit + integração (requere `ffmpeg` no PATH para executar o teste de integração):

```bash
go test ./... -v
```

CI

Há um workflow `.github/workflows/integration.yml` que executa testes em `ubuntu-latest` e instala `ffmpeg` antes de rodar os testes de integração.

Troubleshooting
- `ffmpeg` ou `ffprobe` não encontrados: instale via gerenciador de pacotes (apt, brew, etc.) e verifique `ffmpeg -version`.
- `lefthook` com asdf shim: se o `lefthook install` falhar por causa de um shim (`No version is set ...`), rode o script de instalação incluído que tenta `go install` ou usar Homebrew; em último caso, execute o binário instalado diretamente:

```bash
# exemplo quando Homebrew instalou o binário em /home/linuxbrew/.linuxbrew/bin
/home/linuxbrew/.linuxbrew/bin/lefthook install
```

Contribuindo

- Abra issues/PRs para bugs ou melhorias. Siga o padrão de commits e execute os hooks locais antes de push.

Próximos passos sugeridos
- Adicionar `make` targets para `hooks-install`, `test`, `build` e para empacotar binários.
- Documentar perfis e heurísticas de decisão (`internal/profiles`).
