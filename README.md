# guard-dog 🔒🐕

A TUI to manage users, roles, and permissions in [Stardog](https://www.stardog.com/).

https://user-images.githubusercontent.com/23270779/176809760-984de87f-649d-4223-a86e-b6fa571ccccc.mov

`guard-dog` currently supports:
- viewing users, roles and their respective permissions
- granting and revoking permissions from users and roles
- assigning and unassigning users from/to roles
- deleting users and roles
- changing users' passwords
- enabling and disabling users

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

### Controls

- To switch between panes, `tab`. This is important as actions can only be performed if the pane is active. For example, in order to create a new user (`ctrl+n`), the `user list` tab must be active. In order to add an explicit permission to a user, the `user details` pane must be active. 
  
-  Controls are highlighted throughout the app. If needed,`?` displays a help menu on the right side of the app. The help menu lists the controls for the current active pane. 

<img width="1298" alt="help" src="https://user-images.githubusercontent.com/23270779/176811048-3c5879d7-4f28-40f7-9064-8c9f8e5df59b.png">

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

`guard-dog` will attempt to authenticate using the default superuser `admin` with password `admin` on `http://localhost:5820` if no credentials are provided.

### Environment Variables

| Environment Variable  |  Description |
|---|---|
| `GUARD_DOG_USERNAME`  | username |
| `GUARD_DOG_PASSWORD`  | password |
| `GUARD_DOG_SERVER`  | Stardog server to connect to |


## Built With

- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [bubbles](https://github.com/charmbracelet/bubbles)
- [bubble-table](https://github.com/Evertras/bubble-table)
- [teacup](https://github.com/knipferrc/teacup)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
- [go-stardog](https://github.com/noahgorstein/go-stardog)

## Notes

- 🚧 this project is under active development and is relatively unstable. Please file an issue if you see a bug or have suggestions on how to improve the app.

