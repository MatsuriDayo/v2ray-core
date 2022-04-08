module github.com/v2fly/v2ray-core/v5

go 1.18

require (
	github.com/adrg/xdg v0.4.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.7
	github.com/gorilla/websocket v1.5.0
	github.com/jhump/protoreflect v1.12.0
	github.com/lucas-clemente/quic-go v0.27.0
	github.com/marten-seemann/qtls-go1-16 v0.1.5
	github.com/marten-seemann/qtls-go1-17 v0.1.1
	github.com/marten-seemann/qtls-go1-18 v0.1.1
	github.com/miekg/dns v1.1.48
	github.com/pires/go-proxyproto v0.6.2
	github.com/seiflotfy/cuckoofilter v0.0.0-20220312154859-af7fbb8e765b
	github.com/stretchr/testify v1.7.1
	github.com/v2fly/BrowserBridge v0.0.0-20210430233438-0570fc1d7d08
	github.com/v2fly/ss-bloomring v0.0.0-20210312155135-28617310f63e
	go.starlark.net v0.0.0-20220302181546-5411bad688d1
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064
	golang.org/x/net v0.0.0-20220325170049-de3da57026de
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220325203850-36772127a21f
	golang.zx2c4.com/wireguard v0.0.0-20211209221555-9c9e7e272434 //SagerNet
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
	gvisor.dev/gvisor v0.0.0 //SagerNet
	h12.io/socks v1.0.3
	inet.af/netaddr v0.0.0-20211027220019-c74959edd3b6
)

replace gvisor.dev/gvisor => github.com/sagernet/gvisor v0.0.0-20211227140739-33ed11d8e732

//SagerNet
require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da
	github.com/dgryski/go-camellia v0.0.0-20191119043421-69a8a13fb23d
	github.com/dgryski/go-idea v0.0.0-20170306091226-d2fb45a411fb
	github.com/dgryski/go-rc2 v0.0.0-20150621095337-8a9021637152
	github.com/geeksbaek/seed v0.0.0-20180909040025-2a7f5fb92e22
	github.com/kierdavis/cfb8 v0.0.0-20180105024805-3a17c36ee2f8
)

require (
	github.com/cheekybits/genny v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-metro v0.0.0-20200812162917-85c65e2d0165 // indirect
	github.com/ebfe/rc2 v0.0.0-20131011165748-24b9757f5521 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/riobard/go-bloom v0.0.0-20200614022211-cdc8013cb5b3 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/xtaci/smux v1.5.15 // indirect
	go4.org/intern v0.0.0-20211027215823-ae77deb06f29 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20211027215541-db492cf91b37 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.1.6-0.20210726203631-07bc1bf47fb2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	golang.zx2c4.com/go118/netip v0.0.0-20211111135330-a4a02eeacf9d // indirect
	golang.zx2c4.com/wintun v0.0.0-20211104114900-415007cec224 // indirect
	google.golang.org/genproto v0.0.0-20210722135532-667f2b7c528f // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
