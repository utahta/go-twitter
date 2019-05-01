# go-twitter

Go client library for the [Twitter REST APIs v1.1](https://dev.twitter.com/rest/public)

## Features

- [POST statuses/update](https://dev.twitter.com/rest/reference/post/statuses/update)
- [POST media/upload](https://dev.twitter.com/rest/reference/post/media/upload)
- [POST media/upload (Async)](https://dev.twitter.com/rest/reference/post/media/upload-init)

...Maybe

## Install

```
go get github.com/utahta/go-twitter
```

#### Enable Modules (go v1.11+)
```
export GO111MODULE=on
go mod init
```

## Example Usage

A basic file uploader can be found in:
```
example/file_uploader.go
```

## Contributing

1. Fork it ( https://github.com/utahta/go-twitter/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
