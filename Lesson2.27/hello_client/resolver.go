package main

import (
	"google.golang.org/grpc/resolver"
)

const (
	myScheme   = "q1mi"
	myEndpoint = "resolver.liwenzhou.com"
)

var addrs = []string{"127.0.0.1:8972", "127.0.0.1:8973"}

type q1miResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *q1miResolver) ResolveNow(o resolver.ResolveNowOptions) {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*q1miResolver) Close() {}

type q1miResolverBuilder struct{}

func (*q1miResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &q1miResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			myEndpoint: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*q1miResolverBuilder) Scheme() string { return myScheme }

// 在newclient中指定了自定义解析器
func init() {
	resolver.Register(&q1miResolverBuilder{})
}
