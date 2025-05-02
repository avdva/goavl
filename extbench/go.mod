module bench

go 1.24

replace github.com/avdva/goavl => ../

require (
	github.com/avdva/goavl v1.5.0
	github.com/karask/go-avltree v0.0.0-20210208113816-739fb7fa1601
	github.com/tidwall/btree v1.7.0
)

require golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
