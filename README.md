# guard-dog

```                                      
  ____ _   _   _    ____  ____    ____   ___   ____ 
 / ___| | | | / \  |  _ \|  _ \  |  _ \ / _ \ / ___|
| |  _| | | |/ _ \ | |_) | | | | | | | | | | | |  _ 
| |_| | |_| / ___ \|  _ <| |_| | | |_| | |_| | |_| |
 \____|\___/_/   \_|_| \_|____/  |____/ \___/ \____|                                                 
```

A TUI to manage users, roles, and permissions in [Stardog](https://www.stardog.com/).


## Installation

### homebrew

```bash
brew install noahgorstein/tap/guard-dog
```

### Github releases

Download the relevant asset for your operating system from the latest Github release. Unpack it, then move the binary to somewhere accessible in your `PATH`, e.g. `mv ./guard-dog /usr/local/bin`.

### Build from source

Clone this repo, build from source with `cd guard-dog && go build`, then move the binary to somewhere accessible in your `PATH`, e.g. `mv ./guard-dog /usr/local/bin`.


## Usage

Run the app by running `guard-dog` in a terminal. See `guard-dog --help` and [configuration](#configuration) section below for details.

## Configuration

`guard-dog` can be configured in a yaml file at `$HOME/.guard-dog.yaml`.

Example yaml file:

```yaml
# .guard-dog.yaml
username: "admin"
password: "admin"
server: "http://localhost:5820"
```

Alternatively, `guard-dog` can be configured via environment variables, or via command line args visible by running `guard-dog --help`.

> Command line args take precedence over both the configuation file and environment variables. Environment variables take precedence over the configuration file.

### Environment Variables

| Environment Variable  |  Description |
|---|---|
| `GUARD_DOG_USERNAME`  | username |
| `GUARD_DOG_PASSWORD`  | password |
| `GUARD_DOG_SERVER`  | Stardog server to connect to |


### Built With

- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [bubbles](https://github.com/charmbracelet/bubbles)
- [bubble-table](https://github.com/Evertras/bubble-table)
- [teacup](https://github.com/knipferrc/teacup)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
- [Cobra](https://github.com/spf13/cobra)
- [go-stardog](https://github.com/noahgorstein/go-stardog)
