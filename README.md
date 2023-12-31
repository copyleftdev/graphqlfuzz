# GraphQL Fuzzing Tool

The GraphQL Fuzzing Tool is a command-line utility for testing and fuzzing GraphQL endpoints. It allows you to generate a variety of GraphQL queries with fuzzed input to test the robustness of your GraphQL server. This tool can be used to identify potential vulnerabilities and issues in your GraphQL API.

## Features

- Fuzz GraphQL queries with various input data.
- Test different query types, including mutations and queries.
- Specify a GraphQL schema in JSON format.
- Optional wordlist support for custom fuzzing input.
- Detailed response logging for analysis.

## Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/copyleftdev/graphqlfuzz.git
   cd graphqlfuzz
   ```

2. Build the tool:

   ```bash
   go build
   ```

3. Run the tool with the following command:

   ```bash
   ./graphqlfuzz -endpoint <GraphQL endpoint URL> -gqlfile <Path to GraphQL schema in JSON format> -wordlist <Path to wordlist file (optional)>
   ```

   Replace `<GraphQL endpoint URL>` with the URL of your GraphQL endpoint, `<Path to GraphQL schema in JSON format>` with the path to your GraphQL schema file in JSON format, and `<Path to wordlist file>` with the path to an optional wordlist file for custom fuzzing input.

## Example

```bash
./graphqlfuzz -endpoint http://localhost:8080/graphql -gqlfile schema.json -wordlist wordlist.txt
```

## Dependencies

- Go (Golang)
- External dependencies are managed using Go Modules.
