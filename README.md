# go-tms

go-tms is a tmux session manager. It allows you to easily save, load, and switch between tmux sessions.

## Usage

### Launching and Initial Setup

The `-b` flag is used as a tmux launcher. It performs the following actions:

1.  **Server Launch:** Starts a tmux server if one isn't already running.
2.  **Temporary Session Creation:** Creates a temporary tmux session.
3.  **UI for Session Selection:** Attaches to the temporary session and presents a user interface (UI) for selecting a session to load.
4.  **Session Replacement:** Replaces the temporary session with the selected session, effectively loading the chosen session.
5.  **Session Attach:** If the server is already running, attaches to a session.

The `-b` flag also **starts a background daemon** that monitors the tmux server. The daemon will automatically shut down when the tmux server is no longer active.

This flag is typically used in the shell config as an alias.

### Automatic Session Saving (Daemon)

The `-d` flag starts a background daemon that automatically saves your tmux sessions at regular intervals. **This daemon is automatically managed by the `-b` command and will exit when the tmux server is closed.** The interval is configurable in the configuration file (see Configuration below).

The daemon does not need to be configured manually in your `~/.tmux.conf`. It is started automatically when you launch tmux with `go-tms -b`.

### Switching Sessions (Switcher)

The `-s` flag is used to launch the session switcher, allowing you to quickly switch between your saved and active tmux sessions.

1.  **Popup:** The `-s` flag is designed to be used within a tmux popup.
2.  **Keybinding:** You bind a key combination to the command `go-tms -s` in your tmux configuration.
3.  **Session Selection:** When the keybinding is triggered, a popup appears, allowing you to select a session to switch to.

Example tmux configuration (in your `~/.tmux.conf`):

```tmux
bind-key C-s display-popup \
  -d '#{pane_current_path}' \
  -w 80% \
  -h 80% \
  -E "/Users/aleksejbagno/projects/go-tms/release/darwin/go-tms -s"

```

## Configuration

go-tms uses a YAML configuration file to customize its behavior.

### Configuration File Location

The configuration file should be placed at:
`$HOME/.config/go-tms/config.yaml`

### Available Options and Defaults

Here are the configurable options and their default values:

| Option                  | Description                                                              | Default Value |
| :---------------------- | :----------------------------------------------------------------------- | :------------ |
| `auto-save-interval-minutes` | Interval in minutes for the daemon to automatically save tmux sessions.    | `10`          |
| `fzf-bind-new`          | Keybinding for creating a new session in the fzf switcher.               | `ctrl-n`      |
| `fzf-bind-delete`       | Keybinding for deleting a session in the fzf switcher.                   | `ctrl-d`      |
| `fzf-bind-interactive`  | Keybinding for interactive mode in the fzf switcher.                     | `ctrl-i`      |
| `fzf-bind-save`         | Keybinding for saving sessions in the fzf switcher.                      | `ctrl-s`      |
| `fzf-prompt`            | The prompt displayed by fzf in the session switcher.                     | ` Sessions>  `  |
| `fzf-env`               | Additional options passed to fzf.                                        | `--no-sort --reverse` |
| `zoxide-env`            | Options for zoxide integration (if used for directory suggestions).      | `--layout=reverse --style=full --border=bold --border=rounded --margin=3%` |

### Example `config.yaml`

```yaml
auto-save-interval-minutes: 15
fzf-bind-new: ctrl-x
fzf-prompt: Choose Session:
```

## Why go-tms?

   **Simplified Session Management:** Easily save, load, and switch between tmux sessions.
   **Automated Backups:** The daemon automatically saves your sessions, preventing data loss.
   **User-Friendly Interface:** The session switcher provides an intuitive way to navigate and select sessions.
   **Seamless Integration:** Integrates directly with your tmux setup.
