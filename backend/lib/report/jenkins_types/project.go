package jenkins_types

type Project struct {
    URL    string `json:"url"`
    Name   string `json:"name"`
    Number int    `json:"number"`
}
