# VK Album Downloader

A script that downloads all users public albums from [VK](https://vk.com) social network in best quality available.

## Installation

Install the application using the following command:
`go get github.com/ZelJin/vk_album_downloader`

## Usage

Script requires two flags: 

- `uid`: user ID of the VK user. Default is `1`
- `path`: A local path where albums will be downloaded. Default is `./`

### Usage Example

`vk_album_downloader --uid=6447964 --path=/tmp`
