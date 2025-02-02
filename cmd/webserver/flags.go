package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type flags struct {
	port   string
	dbname string
}

func parseArgs(args []string) flags {
	flagSet := flag.NewFlagSet("flags", flag.PanicOnError)

	port := flagSet.String("port", "", "port that app will listen on (required)")
	dbname := flagSet.String("dbname", "database.db", "name of created SQLite database")

	err := flagSet.Parse(args)
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	if *port == "" {
		fmt.Println("Error: -port flag is required")
		flagSet.Usage()
		os.Exit(1)
	}

	return flags{
		port:   *port,
		dbname: *dbname,
	}

}
