# Tical

A beautiful, minimalist **T**erminal **U**ser **I**nterface **Cal**culator,
themed with [Tokyo Night](https://github.com/folke/tokyonight.nvim) colors and
built on [Bubble Tea](https://github.com/charmbracelet/bubbletea),
[Bubbles](https://github.com/charmbracelet/bubbles) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss).

Operate it entirely with the **mouse** or the **keyboard** — hover and click the
keys, or type the numbers and operators directly. The calculator centres itself
in the terminal and fills the whole screen.

```
╭────────────────────────────────────────────╮
│    Tical  · terminal calculator            │
│  ╭──────────────────────────────────────╮  │
│  │                                      │  │
│  │                                   0  │  │
│  ╰──────────────────────────────────────╯  │
│     C         ⌫         %         ÷        │
│     7         8         9         ×        │
│     4         5         6         −        │
│     1         2         3         +        │
│     ±         0         .         =        │
│  ? help • q quit                           │
╰────────────────────────────────────────────╯
```

Keys are color-coded in the Tokyo Night palette: **teal** functions
(`C` `⌫` `±`), **blue** operators, **green** equals, and a **magenta**
highlight that follows the focused/hovered key.

## Operations

All basic operations are supported: addition `+`, subtraction `-`,
multiplication `*`, division `/`, and modulo `%`. Expressions evaluate
left-to-right, classic-calculator style.

## Install & run

From the [AUR](https://aur.archlinux.org/packages/tical):

```sh
yay -S tical      # or: paru -S tical
```

From source:

```sh
go run .          # run from source
# or
go build -o tical # build a binary
./tical
```

### Build the Arch package locally

```sh
makepkg -si       # build and install with pacman
```

### Releasing (maintainer notes)

Publishing to the AUR is automated by
[`.github/workflows/publish-aur.yml`](.github/workflows/publish-aur.yml).
To cut a new version:

1. Create a **GitHub Release** tagged `vX.Y.Z`.
2. The workflow updates `pkgver`, recomputes the source checksum, regenerates
   `.SRCINFO`, and prepares the AUR push — then **pauses for manual approval**.
3. Approve the run in the GitHub UI; the package is pushed to the AUR.

The approval gate is a GitHub *Environment* (`aur`) with the maintainer set as a
required reviewer, so no release reaches the AUR without an explicit approval.

**One-time setup** (repo *Settings*):

- *Settings → Environments → New environment* → name it `aur`, enable
  **Required reviewers**, and add yourself.
- Add an SSH **private** key as the environment secret `AUR_SSH_PRIVATE_KEY`
  (its matching public key must be registered on your AUR account — a dedicated
  CI key is recommended over your personal one).

## Controls

| Action            | Keyboard                          | Mouse              |
| ----------------- | --------------------------------- | ------------------ |
| Enter digits      | `0`–`9`, `.`                      | click a key        |
| Operators         | `+` `-` `*` `/` `%`               | click a key        |
| Evaluate          | `Enter` or `=`                    | click `=`          |
| Copy result       | `y`                               | —                  |
| Move focus        | arrow keys / `h` `j` `k` `l`      | hover              |
| Press focused key | `Space`                           | click              |
| Delete last digit | `Backspace`                       | click `⌫`          |
| Clear             | `c` / `C`                         | click `C`          |
| Toggle help       | `?`                               | —                  |
| Quit              | `q`, `Esc`, `Ctrl-C`              | —                  |

## Layout

- `internal/calc` — the calculator engine (pure logic, fully unit-tested)
- `internal/ui`   — the Bubble Tea model, Tokyo Night styles, and rendering
- `main.go`       — program entry point
