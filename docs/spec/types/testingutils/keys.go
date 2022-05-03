package testingutils

import (
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestingSK1 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("1b87e7b14d99020a1f8c9381f67907b4fd4ec47db3c32ecde9ec7b23532c107b")
	return ret
}()
var TestingSK2 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("73251a39ad1f6ead1fc2f49a2b460b5caac62b3de5097b590376249e6f3c2925")
	return ret
}()
var TestingSK3 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("67bdecf889644b44ab8508a11447f8bf44c1c87e7e810dc4e905ff2c10261b40")
	return ret
}()
var TestingSK4 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("02f56a2937ba7b8009a88d5a85382baebd03544ae76bc293e0f3dc683c7956ad")
	return ret
}()
var TestingSK5 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("5e88c85c6295ec2f3ee2e7cd9f45f00d27522e45bd87a117e1fea0ae044ede83")
	return ret
}()
var TestingSK6 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("58f6a059510e27f7e7098a9b1e58ec93fe031acf50d338b78e7492a4abf14e5d")
	return ret
}()
var TestingSK7 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("24ce20e9b5ae09ad17c8995374a5d4b2f720effa22c0451d95319612eb2dc04d")
	return ret
}()
var TestingSK8 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("36defc4b41350c9a31949bc327415eff11a2653fa9be2498b164a0ce4687441d")
	return ret
}()
var TestingSK9 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("4ec648dd23907cafd285c5fff82c1f57a99ddf95de65ae5fae94c349818af891")
	return ret
}()
var TestingSK10 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("13acf18295807f00d2fdc27be5f13caf79162db35c36850502e5291d2387373c")
	return ret
}()
var TestingSK11 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("654003bcc8c59c0c4e9f03ed97d5033a29208a83247947fbbe1e80bb01eca35b")
	return ret
}()
var TestingSK12 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("57683dbb39a6314009db5d4835ebfc7c86530df05683dda9b88b4457bda44c06")
	return ret
}()
var TestingSK13 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("1eae6dd5f701f62d6cd59fc0b2fe285feb7412da750059f2193b68c41aaa720f")
	return ret
}()
