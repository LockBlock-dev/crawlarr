# Crawlarr

[![GitHub stars](https://img.shields.io/github/stars/LockBlock-dev/crawlarr.svg)](https://github.com/LockBlock-dev/crawlarr/stargazers)

Crawlarr is a web crawler built using the Go programming language. This tool allows users to input a base URL, and it will search through the HTML code to locate all anchor tags (`<a>`) on the page. Crawlarr will then follow these links and repeat the process, searching through each subsequent page for more anchor tags until either the end of the website or a user-defined maximum depth is reached. This tool leverages concurrency to significantly increase its speed.

See the [changelog](/CHANGELOG.md) for the latest updates.

## Table of content

-   [**Installation**](#installation)
-   [**Compiling from source**](#compiling-from-source)
-   [**Usage**](#usage)
-   [**Configuring Crawlarr**](#configuring-crawlarr)
-   [**Config details**](#config-details)
-   [**Matching types**](#matching-types)
-   [**License**](#copyright)

## Installation

-   Download [go](https://go.dev/dl/) (go 1.20 required).
-   Download or clone the project.
-   Download the binary from the [Releases](../../releases) or [build it](#compiling-from-source) yourself.
-   [Configure Crawlarr](#configuring-crawlarr).

## Compiling from source

-   Use [`build.sh`](/build.sh) or use `go build` in [`cmd/crawlarr/`](/cmd/crawlarr/)

## Usage

-   With a binary:
    -   Run `chmod +x crawlarr`.
    -   Start the tool with `./crawlarr`
-   Running from source:
    -   Start the tool with `go run ./cmd/crawlarr/main.go` or `cd ./cmd/crawlarr/ && go run .`

Find the results in `links.txt`.

## Configuring Crawlarr

The config can be found at the root of the project.

-   Open the [`config`](/config.json) in your favorite editor.
-   Enable the features you want to use. See [Config details](#config-details) for in-depth explanations.

## Config details

| Item       | Values    | Meaning                                  |
| ---------- | --------- | ---------------------------------------- |
| debug      | `boolean` | Enable debug logs                        |
| baseUrl    | `text`    | The URL to starts with                   |
| matchType  | `text`    | [Matching type](#matching-types) for URL |
| depthLimit | `number`  | Maximum crawling depth                   |
| delay      | `number`  | Delay in ms between crawls               |

## Matching types

-   `SAME_BASE`:  
    Match the same base URL, e.g:

    ```diff
    baseUrl: "http://example.com/this-page/"
    + valid match : http://example.com/this-page/random-page/
    - discarded match : http://example.com/another-page/
    - discarded match : http://random.site/a-third-page/
    ```

-   `SAME_HOST`:  
    Match the same host, e.g:

    ```diff
    baseUrl: "http://example.com/this-page/"
    + valid match : http://example.com/this-page/random-page/
    + valid match : http://example.com/another-page/
    - discarded match : http://random.site/another-page/
    ```

-   `DANGEROUS_NO_MATCH_TYPE_ONLY_ENABLE_IF_YOU_KNOW_WHAT_YOURE_DOING`:  
    Match any URL (this can go very far), e.g:
    ```diff
    baseUrl: "http://example.com/this-page/"
    + valid match : http://example.com/this-page/random-page/
    + valid match : http://example.com/another-page/
    + valid match : http://random.site/another-page/
    ```

## License

Copyright (c) 2023 LockBlock-dev

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with [project name]. If not, see <https://www.gnu.org/licenses/>.
