# DNSBL Check

Checks if the given IP address(es) are listed at the given dnsrbl provider(s).

DNSRBL lists are used to fight against e-mail spams. They detect and list IP addresses that are used to send unwanted e-mails. When an e-mail sent from a blacklisted address, the anti-spam software will mark the message as Spam.

This software is useful for sysadmins who operates MTA (Mail Transfer Agent) softwares and want to see if their IP address they use to send e-mail is listed on an RBL list - therefore they can see a reason why some clients doesn't receive e-mail.

## Providers

Providers must be listed in a file, one line should be one provider.
Empty lines are ignored.
Lines started with `#` are ignored.

The file name must be passed with the `-p` parameter to the command or `-` for standard input.

## Addresses

IP addresses or domain names which needs to be tested against the providers.
The addresses must be passed with the `-i` parameter to the command. Multiple address can be listed, separate them by comma (`,`).

## Build

```sh
go build
```

## Execute

Test if the IP address `1.2.3.4` is blacklisted in providers listed in `providers` file (see [how to get provider list](#getting-provider-list) section):

Variations for the same operation:

```sh
./dnsbl-check -i 1.2.3.4 -p ipv4providers
```

```sh
./dnsbl-check -i 1.2.3.4 -p - < ipv4providers
```

```sh
cat ipv4providers | ./dnsbl-check -i 1.2.3.4 -p -
```

## Output

The program returns every result in a new line, fields are separated by TAB character `\t`.

The line starts with the status: `OK` or `FAIL` or `ERR`

- `OK` returned if no listing found for the address
- `FAIL` returned if listing found for the address
- `ERR` returned if the address lookup failed

Second field is the address
Third field is the provider
Fourth field is filled only if the status is either `FAIL` or `ERR`. If the status is `FAIL` and no reason returned from te provider, the `unknown reason` text will be shown. If the status is `ERR` the error message will be shown here.

```text
OK	127.0.0.2	dyn.rbl.polspam.pl
FAIL	127.0.0.2	bl.spamcop.net	Blocked - see https://www.spamcop.net/bl.shtml?127.0.0.2
ERR	127.0.0.2	spam.dnsbl.anonmails.de	lookup 2.0.0.127.spam.dnsbl.anonmails.de on 127.0.0.53:53: server misbehaving
```

## Getting provider list

List of providers coming from [http://multirbl.valli.org/list/](http://multirbl.valli.org/list/)

### IPv4 providers

To get ipv4 blacklist providers run the following command:

```sh
awk '$5 == "b" && $2 == "ipv4" && $1 != "(hidden)" { print $1 }' < providers > ipv4providers
```

Then you can test if a provider is working - responds to a test query (query the address [127.0.0.2](https://datatracker.ietf.org/doc/html/rfc5782#section-5)):

```sh
./dnsbl-check -p ipv4providers -i 127.0.0.2 | awk '$1 == "FAIL" { print $3 }' > ipv4verified
```

Then with that list you can check if your IP address (1.2.3.4) is blacklisted:

```sh
./dnsbl-check -p ip4verified -i 1.2.3.4
```

It can be piped into one command:

```sh
awk '$5 == "b" && $2 == "ipv4" && $1 != "(hidden)" { print $1 }' < providers | \
./dnsbl-check -p - -i 127.0.0.2 | awk '$1 == "FAIL" { print $3 }' | \
./dnsbl-check -p - -i 1.2.3.4
```

However it's recommended to keep the used provider list separately, to save the resources of the providers.

### Domain providers

Similar to the IPv4 providers, but we filter the multirbl list to items which are maintaining balck lists of domains:

```sh
awk '$5 == "b" && $4 == "dom" && $1 != "(hidden)" { print $1 }' < providers > domain_providers
```

Then check if the provider is looking good query the address [TEST](https://datatracker.ietf.org/doc/html/rfc5782#section-5)

```sh
./dnsbl-check -p domain_providers -i TEST | awk '$1 == "FAIL" { print$3 }' > domain_providers_verified
```

You can check the address `INVALID` as well, which should return `OK`:

```sh
./dnsbl-check -p domain_providers -i TEST | awk '$1 == "FAIL" { print $3 }' | \
./dnsbl-check -p - -i INVALID  | awk '$1 == "OK" { print $3 }' > domain_providers_verified

```

Then you can query your domain:

```sh
./dnsbl-check -p domain_providers_verified -i mail.example.com
```
