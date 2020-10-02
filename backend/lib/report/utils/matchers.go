package utils

func CharactersTheSame(url string, project string) int {
    maxSame := 0
    currSame := 0

    for u, _ := range url {
        for i := 0; i <= len(project)-1; i++ {
            if u+i >= len(url) {
                continue
            }
            if url[u+i] == project[i] {
                currSame = currSame + 1
            } else {
                currSame = 0
            }
            if maxSame < currSame {
                maxSame = currSame
            }
        }
    }

    return maxSame
}
