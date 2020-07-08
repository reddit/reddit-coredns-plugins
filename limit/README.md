# limit
## Name
*limit* - limits the number of answers returned.

## Description
*limit* limits a response's answers to the configured limit.

## Syntax
```txt
limit [LIMIT]
```

**[LIMIT]** is an int value for setting the number of records that 
can be returned in an answer. It must be set. Any integer > 0 is 
accepted.

## Examples
Enable limiting the number of responses from the resolver (172.31.0.10):
```corefile
. {
    limit 100
    forward . 172.31.0.10
    log
}
```

Enable limiting the number of answers as an authoritative nameserver:
```corefile
. {
    limit 50
    file db.example.org
    log
}
```

## Considerations

In environments where RFC 3484 Section 6 Rule 9 is implemented and 
enforced (i.e. DNS answers are always sorted and therefore never 
random), clients may need to set this value to 1 to preserve the 
expected randomized distribution behavior (note: RFC 3484 has been 
obsoleted by RFC 6724 and as a result it should be increasingly 
uncommon to need to change this value with modern resolvers).
