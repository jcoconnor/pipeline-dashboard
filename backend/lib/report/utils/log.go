package utils

import (
	"fmt"
)

// LogTree - This is a null function for the moment.
func LogTree(shown string, hidden string, indent int) {
    hide := true
    // := viper.GetBool("hidetreelog") // retrieve values from viper instead of pflag

    if hide {
        return
    }

    // OK - this is dead code - we can add this in later just want to use the last
    // usage of viper flags until we can get a better way of handling configuration
    // parameters.

    hidden = "URL"

    for i := 0; i < indent; i++ {
        fmt.Print("  ")
    }
    fmt.Printf("\\_ [%s] %s\n", hidden, shown)

}
