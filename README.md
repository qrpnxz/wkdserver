# wkdserver

Command wkdserver serves a WKD as specified by
<https://tools.ietf.org/html/draft-koch-openpgp-webkey-service-11>

The first argument is the address on which the server will listen
for connections. It is optional.

Keys are taken from files in `pgpKeyDir`. For example, the keys for
the address `mister@example.org` are in the file named
`r3ptdiy83btqwgjkooeprx3udzwcr34a` in `pgpKeyDir`.

## Copyright
See [COPYRIGHT.md](COPYRIGHT.md).

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md).