# dns-resolver

A simple DNS resolver using golang.

Original Article https://implement-dns.wizardzines.com/book/exercises.html

https://www.ietf.org/rfc/rfc1035.txt

                Local Host                         |  Foreign
                                                   |
    +---------+               +----------+         |  +--------+
    |         | user queries  |          |queries  |  |        |
    |  User   |-------------->|          |---------|->|Foreign |
    | Program |               | Resolver |         |  |  Name  |
    |         |<--------------|          |<--------|--| Server |
    |         | user responses|          |responses|  |        |
    +---------+               +----------+         |  +--------+
                                |     A            |
                cache additions |     | references |
                                V     |            |
                              +----------+         |
                              |  cache   |         |
                              +----------+         |


cache has not been implemented yet.

## Run

```shell
go run main.go

Querying 198.41.0.4 for twitter.com
Querying 192.12.94.30 for twitter.com
Querying 198.41.0.4 for a.r06.twtrdns.net
Querying 192.12.94.30 for a.r06.twtrdns.net
Querying 205.251.195.207 for a.r06.twtrdns.net
Querying 205.251.192.179 for twitter.com
Resolved IP for twitter.com is 104.244.42.65
```

## Root Nameservers

The authoritative name servers that serve the DNS root zone, commonly known as the “root servers”, are a network of hundreds of servers in many countries around the world. They are configured in the DNS root zone as 13 named authorities

[Root Servers](https://www.iana.org/domains/root/servers)
