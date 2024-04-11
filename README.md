# Locator

Go CLI tool that fetches `(lat,lng)` coordinates for search terms using the [Nominatim API from OpenStreetMap](https://nominatim.org/release-docs/latest/api/Overview/).

## Installation (requires Go)

```bash
$ go install github.com/marianogappa/locator@latest
```

## Usage

```bash
$ echo "Berlin,London" | tr ',' '\n' | locator
52.5170365	13.3888599	Berlin
51.4893335	-0.14405508452768728	London
$
```

Or say you have a file:

```
$ cat countries.csv
Germany
UK
$ cat countries.csv | locator
51.1638175	10.4478313	Germany
6.3110548	20.5447525	UK
$
```

For CSV output:

```
$ echo "Berlin,London" | tr ',' '\n' | locator --separator comma
52.5170365,13.3888599,Berlin
51.4893335,-0.14405508452768728,London
```

## Notes

- It caches results in the system's temp folder to be nice to OpenStreetMap. Your system should® automatically evict it.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
