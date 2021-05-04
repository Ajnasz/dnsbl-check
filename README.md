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

./dnsbl-check -i 1.2.3.4 -p providers

## Output

The program returns every result in a new line, fields are separated by TAB character `\t`.

The line starts with the status: `OK` or `ERR`
Second field is the address
Third field is the provider
Fourth field is filled only if the address listed at the provider. If no reason returned from te provider, the `unknown reason` text will be shown.


```
OK	45.95.168.196	zen.spamhaus.org
ERR	45.95.168.196	dnsbl-1.uceprotect.net	IP 45.95.168.196 is UCEPROTECT-Level 1 listed. See http://www.uceprotect.net/rblcheck.php?ipr=45.95.168.196
```
