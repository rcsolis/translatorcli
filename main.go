package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	cli "github.com/rcsolis/translatorcli/internal/translator"
)

var sourceLang string = "es"
var targetLang string = "en"
var sourceText string = "Hola Mundo"

var wg sync.WaitGroup

/**
 * Initialize the command line flags
 */
func init() {
	flag.StringVar(&sourceLang, "sl", "es", "Source language")
	flag.StringVar(&targetLang, "tl", "en", "Target language")
	flag.StringVar(&sourceText, "txt", "Hola Mundo", "Source text")
}

/**
 * Main function
 */
func main() {
	// Declare a channel to receive the translation
	var ch chan string
	// Parse the command line flags
	flag.Parse()
	// Check if the user provided any flags
	if flag.NFlag() == 0 {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Initialize the channel
	ch = make(chan string)
	// Create the request body
	reqBody := &cli.RequestBody{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		SourceText: sourceText,
	}
	// Add 1 to the WaitGroup
	wg.Add(1)
	// Call the RequestTranslation function as goroutine
	go cli.RequestTranslation(reqBody, ch, &wg)
	// Receive the translation from the channel and replace the "+" with " "
	processedStr := strings.ReplaceAll(<-ch, "+", " ")
	// Print the translation
	printTranslation(processedStr)
	// Close the channel
	close(ch)
	// Wait for the goroutine to finish
	wg.Wait()
}

/**
 * Print the translation
 * @param translation string
 */
func printTranslation(translation string) {
	fmt.Println("Source language:", sourceLang)
	fmt.Println("Target language:", targetLang)
	fmt.Println("Source text:", sourceText)
	fmt.Println("Translation::", translation)
}
