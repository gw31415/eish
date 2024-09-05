# eish

SSH command wrapper with [EC2 instance Connect Endpoint](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/create-ec2-instance-connect-endpoints.html).

## Installation

```sh
go install github.com/gw31415/eish@latest
```

## Usage

Same as standard `ssh` command.

### Tips

If you want to push git repository to the EC2 instance, you can set the env variable `GIT_SSH_COMMAND` (or `GIT_SSH`) to `eish`.

```sh
export GIT_SSH=eish
git remote add ec2 ec2-user@i-Ex4mp1e1n5tanc31d:/path/to/repo.git
```


## Spec

### Supported option types

- [x] `[-1] [-2] [-i] ...`: Single character options without its arguments (separated by space)
- [x] `[-B value] [-b value] ...` : Single character options with its arguments (separated by space)
- [x] `[-2v]` : Single character options without its arguments (not separated by space)
- [x] `[-vi identity_file]` : Single character options with its arguments before the option (not separated by space)
- [x] `[-ivalue]` : Single character options with its arguments just after the option (not separated by space)

### Supported options

- Single character options without value
  - `1246AaCfGgKkMNnqsTtVvXxYy`
- Single character options with value
  - `BbcDEeFIiJLlmOoPpRSWw`
- Option to call other functions without SSH connection
  - `Q`

## TODO

- [ ] logging and error message

## License

Apache-2.0
