package githubdl

import "fmt"

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
}

type Response struct {
    Count int `json:"total_count"`
    Items []Repo `json:"items"`
}

type Readme struct {
    Encoding string `json:"encoding"`
    Content string `json:"content"`
}