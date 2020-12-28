module github.com/patrickdevivo/emsql

go 1.15

require (
	github.com/go-openapi/strfmt v0.19.11 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/packethost/packngo v0.5.1
	golang.org/x/crypto v0.0.0-20200420201142-3c4aac89819a
)

replace github.com/mattn/go-sqlite3 v1.14.5 => github.com/patrickdevivo/go-sqlite3 v1.14.6-0.20201202222432-05dee20a1df6
