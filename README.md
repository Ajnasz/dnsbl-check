# DNSBL Check

Checks if the given IP address(es) are listed at the given dnsrbl provider(s).

## Providers

Providers must be listed in a file, one line should be one provider.
Empty lines are ignored.
Lines started with `#` are ignored.

The file name must be passed with the `-p` parameter to the command.

## Addresses

IP addresses or domain names which needs to be tested against the providers.
The addresses must be passed with the `-i` parameter to the command. Multipla address can be listed, separate them by comma (`,`).

## Build

```
go build
```

## Execute

```sh
./dnsbl-check -i 1.2.3.4 -p providers
```

## Output

The program returns every result in a new line, fields are separated by TAB character `\t`.

The line starts with the status: `OK` or `FAIL` or `ERR`
- `OK` returned if no listing found for the address
- `FAIL` returned if listing found for the address
- `ERR` returned if the address lookup failed
Second field is the address
Third field is the provider
Fourth field is filled only if the statis is either `FAIL` or `ERR`. If the status is `FAIL` and no reason returned from te provider, the `unknown reason` text will be shown. If the status is `ERR` the error message will be shown here.

```
OK	127.0.0.2	dyn.rbl.polspam.pl
FAIL	127.0.0.2	bl.spamcop.net	Blocked - see https://www.spamcop.net/bl.shtml?127.0.0.2
ERR	127.0.0.2	spam.dnsbl.anonmails.de	lookup 2.0.0.127.spam.dnsbl.anonmails.de on 127.0.0.53:53: server misbehaving
```
