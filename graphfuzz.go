package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const fuzzChars = "!@#$%^&*()_+-=[]{}|;:',.<>?/"

var (
	endpoint    string
	gqlFile     string
	wordlistPath string
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "", "URL of the GraphQL endpoint")
	flag.StringVar(&gqlFile, "gqlfile", "", "Path to the GraphQL schema file in JSON format")
	flag.StringVar(&wordlistPath, "wordlist", "", "Path to the wordlist file (optional)")
}

type GraphQLType struct {
	Kind   string         `json:"kind"`
	Name   string         `json:"name"`
	Fields []GraphQLField `json:"fields"`
}

type GraphQLField struct {
	Name string              `json:"name"`
	Args []GraphQLFieldArg   `json:"args"`
	Type GraphQLFieldType    `json:"type"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
}

type GraphQLFieldArg struct {
	Name string           `json:"name"`
	Type GraphQLFieldType `json:"type"`
}

type GraphQLFieldType struct {
	Kind   string            `json:"kind"`
	Name   string            `json:"name"`
	OfType *GraphQLFieldType `json:"ofType"`
}

type GraphQLRequest struct {
	Query string `json:"query"`
}

func readSchemaFromFile(filePath string) ([]GraphQLType, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening schema file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading schema file: %w", err)
	}

	var schema struct {
		Data struct {
			Schema struct {
				Types []GraphQLType `json:"types"`
			} `json:"__schema"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bytes, &schema); err != nil {
		return nil, fmt.Errorf("error unmarshalling schema JSON: %w", err)
	}

	return schema.Data.Schema.Types, nil
}

func fuzzGraphQL(url string, types []GraphQLType, wordlist []string) {
	for _, t := range types {
		if t.Kind == "OBJECT" && t.Name != "Query" && t.Name != "Mutation" {
			for _, field := range t.Fields {
				fuzzList := []string{""}
				if len(wordlist) > 0 {
					fuzzList = wordlist
				}

				for _, fuzzString := range fuzzList {
					// Generating fuzzed query based on the field and type
					fuzzedArgs := make([]string, len(field.Args))
					for i, arg := range field.Args {
						fuzzedArgs[i] = fmt.Sprintf("%s: \"%s\"", arg.Name, fuzzString)
					}

					fuzzedQuery := fmt.Sprintf("{ %s(%s) { %s } }", t.Name, strings.Join(fuzzedArgs, ", "), field.Name)
					response, err := sendQuery(url, fuzzedQuery)
					if err != nil {
						fmt.Println("Error sending fuzzed request:", err)
					} else {
						fmt.Printf("Fuzzed Query: %s\nResponse: %v\n", fuzzedQuery, response)
					}
				}
			}
		}
	}
}

func sendQuery(url, query string) (*GraphQLResponse, error) {
		requestBody, err := json.Marshal(GraphQLRequest{Query: query})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Send the request to the GraphQL endpoint
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var graphQLResponse GraphQLResponse
	if err := json.Unmarshal(responseBody, &graphQLResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &graphQLResponse, nil
}

func readWordlist(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening wordlist file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var wordlist []string
	for scanner.Scan() {
		wordlist = append(wordlist, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading wordlist: %w", err)
	}

	return wordlist, nil
}

func main() {
	flag.Parse()

	if endpoint == "" || gqlFile == "" {
		fmt.Println("Both --endpoint and --gqlfile flags are required.")
		flag.Usage()
		return
	}

	var wordlist []string
	if wordlistPath != "" {
		wordlistData, err := readWordlist(wordlistPath)
		if err != nil {
			fmt.Printf("Error reading wordlist: %v\n", err)
			return
		}
		wordlist = wordlistData
	}

	types, err := readSchemaFromFile(gqlFile)
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		return
	}

	fuzzGraphQL(endpoint, types, wordlist)
}
