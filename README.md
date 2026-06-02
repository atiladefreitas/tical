# Tical

A beautiful, minimalist **T**erminal **U**ser **I**nterface **Cal**culator,
themed with [Tokyo Night](https://github.com/folke/tokyonight.nvim) colors and
built on [Bubble Tea](https://github.com/charmbracelet/bubbletea),
[Bubbles](https://github.com/charmbracelet/bubbles) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss).

Operate it entirely with the **mouse** or the **keyboard** — hover and click the
keys, or type the numbers and operators directly.

```
╭────────────────────────────────────╮
│    Tical  · terminal calculator    │
│  ╭──────────────────────────────╮  │
│  │                        42 ×  │  │
│  │                           3  │  │
│  ╰──────────────────────────────╯  │
│     C       ⌫       %       ÷      │
│     7       8       9       ×      │
│     4       5       6       −      │
│     1       2       3       +      │
│     ±       0       .       =      │
│  ? help • q quit                   │
╰────────────────────────────────────╯
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

The PKGBUILD pulls a tagged release tarball from GitHub. To cut a new version,
bump `pkgver` in `PKGBUILD`, then:

```sh
git tag v$pkgver && git push origin v$pkgver   # create the GitHub release/tag
updpkgsums                                      # refresh the source sha256 (needs pacman-contrib)
makepkg --printsrcinfo > .SRCINFO               # refresh .SRCINFO
```

Then push the updated `PKGBUILD` + `.SRCINFO` to the `tical` AUR repository.

## Controls

| Action            | Keyboard                          | Mouse              |
| ----------------- | --------------------------------- | ------------------ |
| Enter digits      | `0`–`9`, `.`                      | click a key        |
| Operators         | `+` `-` `*` `/` `%`               | click a key        |
| Evaluate          | `Enter` or `=`                    | click `=`          |
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
