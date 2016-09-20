package githubdl

import (
    "net/http"
    "io/ioutil"
    "fmt"
    "encoding/json"
    "encoding/base64"
    "strip"
)

type Params struct {
    Language *string
    Search *string
    Sort *string
    Order *string
    Name *string
    RepoSize *string
    Fork *string
    User *string
    Stars *string
    In *string
    CloneFlags *string
}

func (p Params) String() string {
    s := fmt.Sprintf("https://api.github.com/search/repositories?q=%s",*p.Search)

    if(len(*p.Search) > 0 && len(*p.In) > 0){
        s += fmt.Sprintf("+in:%s", *p.In)
    }

    if(len(*p.User) > 0){
        s += fmt.Sprintf("+user:%s", *p.User)
    }

    if(len(*p.Language) > 0){
        s += fmt.Sprintf("+language:%s", *p.Language)
    }

    if(len(*p.RepoSize) > 0){
        s += fmt.Sprintf("+size:%s", *p.RepoSize)
    }

    if(len(*p.Fork) > 0){
        s += fmt.Sprintf("+fork:%s", *p.Fork)
    }

    if(len(*p.Fork) > 0){
        s += fmt.Sprintf("+stars:%s", *p.Stars)
    }


    if(len(*p.Sort) > 0){
        s += fmt.Sprintf("&sort=%s&order=%s", *p.Sort, *p.Order)
    }

    return fmt.Sprintf("%s", s)
}


type Owner struct {
    Name string `json:"login"`
}

type Repo struct {
    Name string `json:"name"`
    Fullname string `json:"full_name"`
    HtmlUrl string `json:"html_url"`
    CloneUrl string `json:"clone_url"`
    Description string `json:"description"`
    Language string `json:"language"`
    Stars int `json:"stargazers_count"`
    Watchers int `json:"watchers_count"`
    Forks int `json:"forks_count"`
    Size int `json:"size"`
    Owner Owner `json:"owner"`
    Readme string
    Lines []string
}

func initStringSlice(width int) func(string) []string {
    var stringSlice []string
    i := 0
    return func(s string) []string {
        lineWidth := width
        if(len(s) > width - 5){
        lineWidth -= 10
            for len(s) > 0 {
                var l int
                if(len(s) > lineWidth){
                    l = lineWidth
                } else {
                    l = len(s)
                }
                stringSlice = append(stringSlice, s[:l])
                s = s[l:]
                i += 1
            }
        } else {
            stringSlice = append(stringSlice, s)
            i += 1
        }
        return stringSlice
    }
}

func (repo *Repo) GenerateDisplay(lineWidth int){
    appendString := initStringSlice(lineWidth)
    appendString(fmt.Sprintf("%v stars, %v watchers, %v forks",repo.Stars, repo.Watchers, repo.Forks))
    appendString("")
    appendString(repo.Fullname)
    appendString(repo.Description)
    appendString("")
    appendString(repo.HtmlUrl)
    appendString(fmt.Sprintf("%v KB, %s", repo.Size, repo.Language))
    appendString("")
    stringSlice := appendString(repo.Readme)
    repo.Lines = stringSlice
}

func (repo *Repo) DisplaySlice(lineWidth int, start int) []string {
    if(len(repo.Lines) == 0){
        repo.GenerateDisplay(lineWidth)
    } 
    return repo.Lines[start:]
}

func (repo *Repo) GetReadme() {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", repo.Owner.Name, repo.Name)
    resp, err := http.Get(url)
    defer resp.Body.Close()
    if err != nil {
        // handle error
    } else {
        var r Readme
        body, _ := ioutil.ReadAll(resp.Body)
        json.Unmarshal(body, &r)
        data, _ := base64.StdEncoding.DecodeString(r.Content)

        repo.Readme = strip.Strip(string(data))
    }
}

type Response struct {
    Count int `json:"total_count"`
    Items []*Repo `json:"items"`
}

type Readme struct {
    Encoding string `json:"encoding"`
    Content string `json:"content"`
}