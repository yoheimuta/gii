# go-from-gist-to-issue #

[![GitHub release](http://img.shields.io/github/release/yoheimuta/go-from-gist-to-issue.svg?style=flat-square)][release]
[![Wercker](http://img.shields.io/wercker/ci/54393fe184570fc622001411.svg?style=flat-square)][wercker]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[release]: https://github.com/yoheimuta/go-from-gist-to-issue/releases
[wercker]: https://app.wercker.com/project/bykey/371feff8aaae40a8317fa0192a72803f
[license]: https://github.com/yoheimuta/go-from-gist-to-issue/blob/master/LICENSE

Easily import each gist content to github issue.

### Usage

Require a [personal access token](https://github.com/blog/1509-personal-api-tokens).

```ruby
$ go-from-gist-to-issue --gist example.txt --repo yourrepo --token yourtoken
```

### Installation

To install go-from-gist-to-issue, please use go get.

```ruby
$ go get github.com/yoheimuta/go-from-gist-to-issue
...
$ go-from-gist-to-issue help
...
```

Or you can download a binary from [github relases page](https://github.com/yoheimuta/go-from-gist-to-issue/releases) and place it in $PATH directory.

### Example

At first, prepare a file to list gist urls.

```ruby
$ cat valid_example.txt
https://gist.github.com/yoheimuta/a05cc8b41f161efd8e8c
https://gist.github.com/yoheimuta/c3d0be70ce9194c4556f
```

Then, dry-run to import these gist urls to a your repository.

```ruby
$ go-from-gist-to-issue --gist valid_example.txt --repo sample-go-from-gist-to-issue --token *** --no-gist-comment --dry-run
Downloading a gist and comments: a05cc8b41f161efd8e8c
Downloaded  a gist and comments: a05cc8b41f161efd8e8c
Dry-run to create an issue
Dry-run to create a comment
Dry-run to create a comment
Downloading a gist and comments: c3d0be70ce9194c4556f
Downloaded  a gist and comments: c3d0be70ce9194c4556f
Dry-run to create an issue
Dry-run to create a comment
Dry-run to create a comment
Completed to import from gists to issues: count=2
```

Finally, run to import these gist urls to a your repository.

```ruby
$ go-from-gist-to-issue --gist valid_example.txt --repo sample-go-from-gist-to-issue --token a8d9ab0c69d087a06161b4461bd95e4083c78611 --no-gist-comment
Downloading a gist and comments: c3d0be70ce9194c4556f
Downloaded  a gist and comments: c3d0be70ce9194c4556f
Created an issue: from https://gist.github.com/c3d0be70ce9194c4556f to https://github.com/yoheimuta/sample-go-from-gist-to-issue/issues/259
Created a comment: https://github.com/yoheimuta/sample-go-from-gist-to-issue/issues/259#issuecomment-99134645
Created a comment: https://github.com/yoheimuta/sample-go-from-gist-to-issue/issues/259#issuecomment-99134648
Downloading a gist and comments: a05cc8b41f161efd8e8c
Downloaded  a gist and comments: a05cc8b41f161efd8e8c
Created an issue: from https://gist.github.com/a05cc8b41f161efd8e8c to https://github.com/yoheimuta/sample-go-from-gist-to-issue/issues/260
Created a comment: https://github.com/yoheimuta/sample-go-from-gist-to-issue/issues/260#issuecomment-99134646
Created a comment: https://github.com/yoheimuta/sample-go-from-gist-to-issue/issues/260#issuecomment-99134650
Completed to import from gists to issues: count=2
```

### Options

```ruby
$ go-from-gist-to-issue help
NAME:
   go-from-gist-to-issue - importing each gist to github issue

USAGE:
   go-from-gist-to-issue [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --gist                       a text file to list up gist url
   --repo                       a repository name to be imported from gists
   --token                      a github personal access token
   --dry-run                    a flag to run without any changes
   --verbose                    a flag to log verbosely
   --sequence                   a flag to import sequentially
   --no-gist-comment            a flag not to create a gist comment after completing each import
   --help, -h                   show help
   --generate-bash-completion
   --version, -v                print the version
```
