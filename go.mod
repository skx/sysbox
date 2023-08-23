module github.com/skx/sysbox

go 1.16

require (
	github.com/armon/go-metrics v0.5.1 // indirect
	github.com/creack/pty v1.1.18
	github.com/gdamore/tcell/v2 v2.6.0
	github.com/google/btree v1.1.2 // indirect
	github.com/google/goexpect v0.0.0-20210430020637-ab937bf7fd6f
	github.com/google/goterm v0.0.0-20200907032337-555d40f16ae2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/memberlist v0.5.0
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/miekg/dns v1.1.55 // indirect
	github.com/nightlyone/lockfile v1.0.0
	github.com/peterh/liner v1.2.2
	github.com/rivo/tview v0.0.0-20230814110005-ccc2c8119703
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/skx/subcommands v0.9.2
	golang.org/x/net v0.14.0
	golang.org/x/term v0.11.0
	golang.org/x/tools v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.5.1
