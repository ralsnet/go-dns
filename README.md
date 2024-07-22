# DNS: DNS Request Tool

`dns` is a dns request tool like `dig`, `nslookup`.

## Installing

Install `dns` by running:

`go install github.com/ralsnet/cmd/dns@latest`

and ensuring that `$GOPATH/bin` is added to your `$PATH`

## CLI Options

| option         | short | description                                              | e.g.               |
| :------------- | :---- | :------------------------------------------------------- | :----------------- |
| `--config`     | `-c`  | Config file path (default: `$HOME/.config/dns/dns.json`) | `-c ./config.json` |
| `--subdomains` | `-s`  | Subdomain hosts to lookup (comma separated)              | `-s www,mail,ftp`  |
| `--recursive`  | `-r`  | Recursive lookup MX, CNAME records                       | `-r`               |
| `--format`     | `-f`  | Output format (`json`, `text`)                           | `-f json`          |

## Config file

Default config filepath is `$HOME/.config/dns/dns.json`.

```jsonc
{
  "hosts": ["www", "smtp", "ftp"],
  "recursive": false
}
```
