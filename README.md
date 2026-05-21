# macfigure

Mac configuration in pkl. Simple alternative to nix-darwin

## Installation

Run the install script directly from GitHub:

```sh
curl -fsSL https://raw.githubusercontent.com/quintisimo/macfigure/main/install.sh | bash
```

The script will:

1. Install [Homebrew](https://brew.sh) if it is not already installed.
2. Download the latest [pkl](https://pkl-lang.org) release from GitHub and move it to `/usr/local/bin` (required to parse config files).
3. Download the latest `macfigure` release from GitHub and move it to `/usr/local/bin`.

## Configuration

Example `config.pkl`:

- Allowed values for each section can be found in the section pkl file in the [pkl folder](./pkl)

```pkl
amends "https://raw.githubusercontent.com/quintisimo/macfigure/pkl/config.pkl"

brew {
  casks = List("zed")
}
cron {
  new {
    schedule = "0 9 * * 1-5"
    source = "/Users/quintisimo/Github/personal/macfigure/hello.sh"
    target = "/tmp/hello.sh"
  }
}
nsglobaldomain {
  AppleIconAppearanceTheme = "RegularAutomatic"
}
dock {
  apps = List("/Applications/Zed.app", "spacer", "small-spacer")
  `show-recents` = false
}
secret {
  new {
    source = "./secrets.age"
    target = "/tmp/secrets.sh"
  }
}
```

Execute the config by running:

```bash
macfigure sync -c config.pkl
```
