package main 

import (
    "flag"
    "net/http"
    "fmt"
    "os/exec"
    "encoding/json"
    "encoding/base64"
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
        defer termbox.Close()
        termbox.SetInputMode(termbox.InputEsc)
        viewRepos(repos)
    } else {
        fmt.Println("Usage: github-dl [-search <repository>] [-in <name,description,readme>] [-user <user>] [-lang <language>] [-stars <min..max>] [-size <min..max>] [-showforks <true/only>] [-sort <field>] [-order <asc/desc>]")
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

func getRepos(p githubdl.Params) []githubdl.Repo {
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

func getReadme(repo *githubdl.Repo) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", repo.Owner.Name, repo.Name)
    resp, err := http.Get(url)
    defer resp.Body.Close()
    if err != nil {
        // handle error
    } else {
        var r githubdl.Readme
        body, _ := ioutil.ReadAll(resp.Body)
        json.Unmarshal(body, &r)
        data, _ := base64.StdEncoding.DecodeString(r.Content)
        repo.Readme = string(data)
    }
}

func viewRepos(repos []githubdl.Repo) {
    i := 0
    width, height := termbox.Size()

    printRepo(repos[0], width, height)

    mainloop:
    for {
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            switch ev.Key {
            case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlQ:
                break mainloop
            case termbox.KeyArrowRight:
                if(i < len(repos) - 1){
                    i += 1
                    go printRepo(repos[i], width, height)
                }
            case termbox.KeyArrowLeft:
                if(i > 0){
                    i -= 1
                    go printRepo(repos[i], width, height)
                }
            case termbox.KeyCtrlG:
                t := make(chan string)
                go gitClone(repos[i].CloneUrl, t)
                go func() {
                    for str := range t {
                        y := height - 5
                        clearLine(y, width)
                        // fmt.Printf(str)
                        printStringAtXY(str, width / 2, y, width, true)
                        termbox.Flush()
                    }
                }()
            case termbox.KeyCtrlR:
                getReadme(&repos[i])
                printRepo(repos[i], width, height)
            }
        case termbox.EventError:
            panic(ev.Err)
        }
    }
}

func printRepo(repo githubdl.Repo, width int, height int) {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    center := (width / 2) -1

    var y int

    y = 1

     _, y = printStringAtXY(fmt.Sprintf("%v stars, %v watchers, %v forks",repo.Stars, repo.Watchers, repo.Forks), center, y, width, false)

    y += 2

    _, y = printStringAtXY(repo.Fullname, center, y, width, false)

    y += 1

    _, y = printStringAtXY(repo.Description, center, y, width, false)

    y += 2

    _, y = printStringAtXY(repo.HtmlUrl, center, y, width, false)

    y += 2

    if(len(repo.Language) == 0){
        repo.Language = "Unknown"
    }
    
    _, y = printStringAtXY(fmt.Sprintf("%v KB, %s", repo.Size, repo.Language), center, y, width, false)

    y += 1

    if(len(repo.Readme) > 0){
        _, y = printStringAtXY(repo.Readme, center, y, width, false)
    }

    y = height - 2

    printStringAtXY("<- : Prev | -> : Next | Ctrl-G : Clone | Ctrl-C : Exit", center, y, width, false)

    termbox.Flush()
}

func printStringAtXY(s string, x int, y int, maxw int, override bool) (int, int) {
    var stringArr []string
    if(len(s) > maxw - 5){
        maxw -= 10
        for len(s) > 0 {
            var l int
            if(len(s) > maxw){
                l = maxw
            } else {
                l = len(s)
            }
            stringArr = append(stringArr, s[:l])
            s = s[l:]
        }
    } else {
        stringArr = append(stringArr, s)
    }

    var ax, ay int
    for j := 0; j < len(stringArr); j++ {
        sc := len(stringArr[j])
        ax = x - (sc / 2 - 1)
        if(override == false){
            ay = y + j
        } else {
            ay = y
        }
        for len(stringArr[j]) > 0 {
            r, size := utf8.DecodeRuneInString(stringArr[j])
            termbox.SetCell(ax, ay, r, termbox.ColorDefault, termbox.ColorDefault)
            stringArr[j] = stringArr[j][size:]
            ax += 1
        }
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
    // if err != nil {
    //     // log.Fatal(err)
    // }
}
