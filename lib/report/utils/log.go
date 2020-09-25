package utils

import (
    "fmt"

    "github.com/spf13/viper"
)

func LogHeading(shown string, hidden string) {
    hidden = "URL"
    fmt.Println("\n\n#==================================================#")
    fmt.Printf("# [%s] %s\n", hidden, shown)
    fmt.Println("#==================================================#")
}

func Log(shown string, hidden string) {
    hidden = "URL"
    fmt.Printf("# [%s] %s\n", hidden, shown)
}

func LogTree(shown string, hidden string, indent int) {
    hide := viper.GetBool("hidetreelog") // retrieve values from viper instead of pflag

    if hide {
        return
    }

    hidden = "URL"

    for i := 0; i < indent; i++ {
        fmt.Print("  ")
    }
    fmt.Printf("\\_ [%s] %s\n", hidden, shown)

}
