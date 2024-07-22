package dns

import (
	"fmt"
	"net"
	"strings"
)

const (
	RecordTypeA     = "A"
	RecordTypeCNAME = "CNAME"
	RecordTypeMX    = "MX"
	RecordTypeNS    = "NS"
	RecordTypeTXT   = "TXT"
)

type Record struct {
	Value string
	Type  string
}

func (me Record) String() string {
	return fmt.Sprintf("%s\t%s", me.Type, me.Value)
}

type Reverse struct {
	IP      string
	Domains []string
}

func (me Reverse) String() string {
	return fmt.Sprintf("%s\t%s", me.IP, strings.Join(me.Domains, ", "))
}

type Domain struct {
	Name    string
	Records []Record
	Reverse []Reverse
}

var lookedups = make(map[string]struct{})

func Lookup(name string) *Domain {
	if _, ok := lookedups[name]; ok {
		return nil
	}
	lookedups[name] = struct{}{}
	me := &Domain{Name: name}
	me.Lookup()
	if len(me.Records) == 0 {
		return nil
	}
	return me
}

func (me *Domain) String() string {
	builder := strings.Builder{}

	builder.WriteString("---\n")
	builder.WriteString("\n")
	builder.WriteString("Domain: " + me.Name + "\n")

	builder.WriteString("\n")

	builder.WriteString("Records:\n")
	for _, record := range me.Records {
		builder.WriteString("  " + record.String() + "\n")
	}

	if len(me.Reverse) > 0 {
		builder.WriteString("\n")

		builder.WriteString("Reverse Lookup:\n")
		for _, reverse := range me.Reverse {
			builder.WriteString("  " + reverse.String() + "\n")
		}
	}

	return builder.String()
}

func (me *Domain) Lookup() {
	domain := me.Name

	records := []Record{}

	cname, err := net.LookupCNAME(domain)
	cname, _ = strings.CutSuffix(cname, ".")
	if err == nil && cname != domain {
		records = append(records, Record{Value: cname, Type: RecordTypeCNAME})
		me.Records = records
		// A CNAME record is not allowed to coexist with any other data.
		return
	}

	ns, err := net.LookupNS(domain)
	if err == nil {
		for _, n := range ns {
			v, _ := strings.CutSuffix(n.Host, ".")
			records = append(records, Record{Value: v, Type: RecordTypeNS})
		}
	}

	a, err := net.LookupIP(domain)
	if err == nil {
		for _, r := range a {
			ip := r.String()
			records = append(records, Record{Value: ip, Type: RecordTypeA})

			// Reverse lookup
			names, err := net.LookupAddr(ip)
			if err == nil {
				me.Reverse = append(me.Reverse, Reverse{IP: ip, Domains: names})
			}
		}
	}

	mx, err := net.LookupMX(domain)
	if err == nil {
		for _, m := range mx {
			v, _ := strings.CutSuffix(m.Host, ".")
			records = append(records, Record{Value: v, Type: RecordTypeMX})
		}
	}

	txt, err := net.LookupTXT(domain)
	if err == nil {
		for _, t := range txt {
			records = append(records, Record{Value: t, Type: RecordTypeTXT})
		}
	}

	me.Records = records
}

type Config struct {
	Hosts     []string `json:"hosts"`
	Recursive bool     `json:"recursive"`
}

func Run(domain string, config *Config) []*Domain {
	d := Lookup(domain)
	if d == nil {
		return nil
	}

	domains := []*Domain{d}

	for _, host := range config.Hosts {
		domain := host + "." + domain
		d := Lookup(domain)
		if d != nil {
			domains = append(domains, d)
		}
	}

	// Recursive lookup
	// Lookup CNAME and MX records and add them to the list of domains to lookup

	if config.Recursive {
		for _, d := range domains {
			for _, record := range d.Records {
				if record.Type == RecordTypeCNAME || record.Type == RecordTypeMX {
					domain := record.Value
					d := Lookup(domain)
					if d != nil {
						domains = append(domains, d)
					}
				}
			}
		}
	}

	return domains
}
