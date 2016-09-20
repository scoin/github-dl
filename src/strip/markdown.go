package strip

import "regexp"
import "fmt"

type StringReplacer []byte

func(s *StringReplacer) Replace(exp string, repl string){
    regx := regexp.MustCompile(exp)
    *s = regx.ReplaceAll(*s, []byte(repl))
    fmt.Println(*s)
}

func Strip(src string) string {
    output := StringReplacer(src)
    output.Replace("\n={2,}", "\n")
    output.Replace("~~", "")
    output.Replace("`{3}.*\n", "")
    output.Replace("<(.*?)>", " ")
    output.Replace(`^[=-]{2,}\s*$`, "")
    // output.Replace(`[^.+?](: .*?$)?`, "")
    output.Replace(`\s{0,2}[.*?]: .*?$`, "")
    output.Replace(`![.*?][[(].*?[])]`, "")
    output.Replace(`[(.*?)][[(].*?[])]`, "$1")
    output.Replace(">", "")
    output.Replace(`^\s{1,2}[(.*?)]: (\S+)( ".*?")?\s*$`, "")
    output.Replace(`^#{1,6}\s*([^#]*)\s*(#{1,6})?`, "$1")
    output.Replace(`([*_]{1,3})(\S.*?\S)`, "$2")
    output.Replace("(`{3,})(.*?)", "$2")
    output.Replace(`^-{3,}\s*$`, "")
    output.Replace("`(.+?)`", "$1")
    output.Replace("\n{2,}", "  ")
    return string(output)
}

// output.Replace(`\s{0,2}[.*?]: .*?$`, "")
//     output.Replace(`![.*?][[(].*?[])]`, "")
//     output.Replace(`[(.*?)][[(].*?[])]`, "$1")
//     output.Replace(`^\s{1,2}[(.*?)]: (\S+)( ".*?")?\s*$`, "")
//     output.Replace(`^#{1,6}\s*`, "")
//     output.Replace(`([*_]{1,2})(\S.*?\S)$1`, "$2")
//     output.Replace("(`{3,})(.*?)$1", "$2")
//     output.Replace(`^-{3,}\s*$`, "")
//     output.Replace("`(.+)`", "$1")
//     output.Replace(`n{2,}`, "  ")