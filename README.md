# simple-resolver

A simple toy dns resolver(command line tool) to resolve domain names to IP addressess. I have used [miekg/dns](https://github.com/miekg/dns) library to create and parse DNS packets.

## How to run

```
go run resolve.go google.com.
```

The `Msg` struct from the `miekg/dns` library, which lists all the sections.
```
type Msg struct {
        MsgHdr
        Compress bool       `json:"-"` // If true, the message will be compressed when converted to wire format.
        Question []Question // Holds the RR(s) of the question section.
        Answer   []RR       // Holds the RR(s) of the answer section.
        Ns       []RR       // Holds the RR(s) of the authority section.
        Extra    []RR       // Holds the RR(s) of the additional section.
}
```

## Resources
- [miekg/dns](https://pkg.go.dev/github.com/miekg/dns@v1.1.47) Documentation.
