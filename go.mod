module github.com/adriancampos/ha-tray

go 1.14

require (
	github.com/akavel/rsrc v0.8.0 // indirect
	github.com/apenwarr/fixconsole v0.0.0-20191012055117-5a9f6489cc29
	github.com/cratonica/2goarray v0.0.0-20190331194516-514510793eaa // indirect
	github.com/gen2brain/beeep v0.0.0-20200526185328-e9c15c258e28
	github.com/gen2brain/dlgs v0.0.0-20200211102745-b9c2664df42f
	github.com/getlantern/appdir v0.0.0-20180320102544-7c0f9d241ea7 // indirect
	github.com/getlantern/ops v0.0.0-20200403153110-8476b16edcd6 // indirect
	github.com/getlantern/systray v1.0.4
	github.com/getlantern/uuid v1.2.0 // indirect
	github.com/google/go-cmp v0.5.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/lxn/walk v0.0.0-20191128110447-55ccb3a9f5c1 // indirect
	github.com/lxn/win v0.0.0-20191128105842-2da648fda5b4 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71 // indirect
	golang.org/x/sys v0.0.0-20200622214017-ed371f2e16b4 // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

// Until the PR is accepted, use a fork of systray that supports menu open/close events
replace github.com/getlantern/systray => github.com/adriancampos/systray v1.0.4-fork
