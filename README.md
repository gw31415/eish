# eish

## Supported option types
- [x] `[-1] [-2] [-i] ...`: Single character options (separated by space)
- [x] `[-B value] [-b value] ...` : Single character options with arguments (separated by space)
- [x] `[-2v]` : Single character options (not separated by space)
- [x] `[-vi identity_file]` : Single character options with arguments before the option (not separated by space)

## Supported options

- Single character options without value
  - `1246AaCfGgKkMNnqsTtVvXxYy`
- Single character options with value
  - `BbcDEeFIiJLlmOoPpRSWw`
- Option to call other functions without SSH connection
  - `Q`
