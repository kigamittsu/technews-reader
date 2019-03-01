# technews-reader

## Project setup

```
go build
```

### Prerequisites

To run the application, you need to get the following Libraries.

1. [GopherJS](https://github.com/gopherjs/gopherjs)

```
$ go get -u github.com/gopherjs/gopherjs
```

2. [Vecty](https://github.com/gopherjs/vecty)

```
$ go get -u github.com/gopherjs/vecty
```

### Compiles and hot-reloads for development

Front-end

```
gopherjs serve
```

Back-end

```
go run server.go
```

### Compiles for production

```
go build
```

### Run the application

```
go run start
```

### Run your tests

```
go test
```

## Feature Plan

1. Rethink how to gather the main article
2. Change the HTTP request to using goroutine

## Contributing

When contributing to this repository, please first discuss the change you wish to make via issue,
email, or any other method with the owners of this repository before making a change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
