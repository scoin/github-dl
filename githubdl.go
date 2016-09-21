package main 

import (
    "flag"
    "net/http"
    "fmt"
    "os/exec"
    "encoding/json"
    "bufio"
    "io/ioutil"
    "unicode/utf8"
    "github.com/nsf/termbox-go"
    "githubdl"
)

func main() {
    p := getFlagParams()
    repos := getRepos(p)

    if(len(repos) > 0){
        termbox.Init()
        termbox.Sync()
        defer termbox.Close()
        termbox.SetInputMode(termbox.InputEsc)
        viewRepos(repos, p)
    } else {
        fmt.Println("Usage: github-dl [-search <terms>] [-in <name,description,readme>] [-user <user>] [-lang <language>] [-stars <min..max>] [-size <min..max>] [-showforks <true/only>] [-sort <field>] [-order <asc/desc>]")
    }
}

func getFlagParams() githubdl.Params {
    var p githubdl.Params
    p.Language = flag.String("lang", "", "")
    p.Search = flag.String("search", "", "")
    p.Sort = flag.String("sort", "", "")
    p.Order = flag.String("order", "desc", "")
    p.Name = flag.String("name", "", "")
    p.RepoSize = flag.String("size", "", "")
    p.Fork = flag.String("showforks", "", "")
    p.User = flag.String("user", "", "")
    p.Stars = flag.String("stars", "", "")
    p.In = flag.String("in", "", "")
    flag.Parse()
    return p
}

func getRepos(p githubdl.Params) []*githubdl.Repo {
    resp, err := http.Get(p.String())
    if err != nil {
        panic("HTTP Request Error - are you connected to the internet?")
    }
    defer resp.Body.Close()

    var r githubdl.Response
    body, _ := ioutil.ReadAll(resp.Body)
    json.Unmarshal(body, &r)
    return r.Items
}

func viewRepos(repos []*githubdl.Repo, p githubdl.Params) {
    i := 0
    line := 0
    width, height := termbox.Size()

    printRepo(repos[0], width, height, line)

    mainloop:
    for {
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            switch ev.Key {
            case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlQ:
                break mainloop

            case termbox.KeyArrowRight:
                if(i < len(repos) - 1){
                    line = 0
                    i += 1
                    printRepo(repos[i], width, height, line)
                }

            case termbox.KeyArrowLeft:
                if(i > 0){
                    line = 0
                    i -= 1
                    printRepo(repos[i], width, height, line)
                }

            case termbox.KeyArrowDown, termbox.MouseWheelDown:
                if(line + (height -4) < len(repos[i].Lines)){
                    line += 1
                    printRepo(repos[i], width, height, line)
                }

            case termbox.KeyArrowUp, termbox.MouseWheelUp:
                if(line > 0){
                    line -= 1
                    printRepo(repos[i], width, height, line)
                }

            case termbox.KeyCtrlG:
                t := make(chan string)
                go gitClone(repos[i].CloneUrl, t)
                go func() {
                    for str := range t {
                        clearLine(height - 1, width)
                        printStringAtXY(str, 0, height - 1, termbox.ColorBlack, termbox.ColorWhite, true)
                        termbox.Flush()
                    }
                }()

            case termbox.KeyCtrlR:
                go func(){
                    repos[i].GetReadme()
                    repos[i].GenerateDisplay(width)
                    printRepo(repos[i], width, height, line)
                }()
            }
        case termbox.EventError:
            panic(ev.Err)
        }
    }
}

func printRepo(repo *githubdl.Repo, width int, height int, startLine int) {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    center := (width / 2)

    stringArr := repo.DisplaySlice(width, startLine, (startLine + height - 4))

    for i, s := range stringArr {
        printStringAtXY(s, center - (len(s) / 2), i + 1, termbox.ColorDefault, termbox.ColorDefault, false)
    }

    menu := `<- : Prev | -> : Next | ^R : Readme | ^G : Clone | ^C : Exit`
    printStringAtXY(menu, center - (len(menu) / 2), height - 2, termbox.ColorBlack, termbox.ColorWhite, false)

    termbox.Flush()
}

func printStringAtXY(s string, x int, y int, fg termbox.Attribute, bg termbox.Attribute, override bool) (int, int) {
    ax := x
    ay := y
    for len(s) > 0 {
        r, size := utf8.DecodeRuneInString(s)
        termbox.SetCell(ax, ay, r, fg, bg)
        s = s[size:]
        ax += 1
    }
    return ax, ay
}

func clearLine(y int, width int){
    for i := 0; i < width; i++ {
        termbox.SetCell(i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
    }
    termbox.Flush()
}

func gitClone(url string, read chan string) {
    cmd := exec.Command("git", "clone", "--progress", url)
    stderr, _ := cmd.StderrPipe()
    scanner := bufio.NewScanner(stderr)
    cmd.Start()


    for scanner.Scan() {
        read <- string(scanner.Text())
    }
    cmd.Wait()
    close(read)
}
