package main

import (
    "fmt"
    "os"

    "github.com/I-invincib1e/http-cli/internal/client"
    "github.com/I-invincib1e/http-cli/internal/config"
    "github.com/I-invincib1e/http-cli/internal/output"
    "github.com/I-invincib1e/http-cli/internal/styles"
)

func exitError(msg string) {
    fmt.Fprintf(os.Stderr, "Error: %v\n", msg)
    os.Exit(1)
}

func main() {
    cfg, err := config.ParseFlags()
    if err != nil {
        exitError(err.Error())
    }

    if err := cfg.Validate(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        config.PrintUsage()
        os.Exit(1)
    }

    st := styles.New()
    output.DisplayRequest(cfg, st)

    resp, err := client.ExecuteRequest(cfg)
    if err != nil {
        exitError(err.Error())
    }

    if err := output.DisplayResponse(cfg, resp, st); err != nil {
        exitError(err.Error())
    }
}

