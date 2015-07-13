# gii #

[![GitHub release](http://img.shields.io/github/release/yoheimuta/gii.svg?style=flat-square)][release]
[![Wercker](http://img.shields.io/wercker/ci/54393fe184570fc622001411.svg?style=flat-square)][wercker]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[release]: https://github.com/yoheimuta/gii/releases
[wercker]: https://app.wercker.com/project/bykey/371feff8aaae40a8317fa0192a72803f
[license]: https://github.com/yoheimuta/gii/blob/master/LICENSE

Gist Issue Importer.

`gii` enables you to import multiple Gists to GitHub Issues. `gii` can parallelize to import multiple gists.

### Usage

Require a [personal access token](https://github.com/blog/1509-personal-api-tokens).

```ruby
$ gii --gist example.txt --repo yourrepo --token yourtoken
```

### Installation

To install gii, please use go get.

```ruby
$ go get github.com/yoheimuta/gii
...
$ gii help
...
```

Or you can download a binary from [github relases page](https://github.com/yoheimuta/gii/releases) and place it in $PATH directory.

### Example

At first, prepare a file to list gist urls.

```ruby
$ cat example/valid_example.txt
https://gist.github.com/yoheimuta/a05cc8b41f161efd8e8c
https://gist.github.com/yoheimuta/c3d0be70ce9194c4556f
```

Then, dry-run to import these gist urls to your repository.

```ruby
$ gii --gist example/valid_example.txt --repo sample-go-from-gist-to-issue --token *** --no-gist-comment --dry-run
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
$ gii --gist example/valid_example.txt --repo sample-go-from-gist-to-issue --token *** --no-gist-comment
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

### Screenshots

<table style="width:100%">
  <tr>
    <td>gist</td>
    <td>to</td>
    <td>issue</td>
  </tr>
  <tr>
    <td><img src="https://raw.githubusercontent.com/yoheimuta/gii/master/screenshot/gist.png" /></td>
    <td>to</td>
    <td><img src="https://raw.githubusercontent.com/yoheimuta/gii/master/screenshot/issue.png" /></td>
  </tr>
</table>

### Options

```ruby
$ gii help
NAME:
   gii - CLI tool to bulk import each gist to github issue

USAGE:
   gii [global options] command [command options] [arguments...]

VERSION:
   0.0.2

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
