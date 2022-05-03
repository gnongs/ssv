package testingutils

import (
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
)

/**
SK: 3515c7d08e5affd729e9579f7588d30f2342ee6f6a9334acf006345262162c6f
Validator PK: 8e80066551a81b318258709edaf7dd1f63cd686a0e4db8b29bbb7acfe65608677af5a527d9448ee47835485e02b50bc0
*/

var TestingSK1 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("5bbc105ff365486a14e9ff2ff56f49bf3965d17aeb0da41171cb26cbb0cf06bc")
	return ret
}()
var TestingSK2 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("61bd3b5b32bb2bf20c563f22dd156ca2bb64364587d6e31c2cac1591b9eb65cf")
	return ret
}()
var TestingSK3 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("700fd2d5dd1505e3dc90afb995e952bf2e4f30c02f873beaf3e22f05c8002bfa")
	return ret
}()
var TestingSK4 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("4c3b608e678daf173725aebb21adb722591a7250ed257440c1c93033826b1ea2")
	return ret
}()
var TestingSK5 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("272cbe168dd33a3b927533d2b0bef91dc1e61d4a8d8010657201af0f23f15e9f")
	return ret
}()
var TestingSK6 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("4086e00b8d52bdfa9c713955ab68913d6dee47d2f3fd913ba2397d91bf293f0c")
	return ret
}()
var TestingSK7 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("6e09e1002f5e2e62cf155905339f3920758b003aacf00af1cac06b8cfdc4949c")
	return ret
}()
var TestingSK8 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("4d4c0336350d292d84c5be8aaad8c9a503b3ecafa4460dea989486213b3f0286")
	return ret
}()
var TestingSK9 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("28b4ab6d2abdfd2487133b92447dc01f0ac5441940421bb08b15ccd34ee23516")
	return ret
}()
var TestingSK10 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("0ba704e875caa47ae1ca936a03abe8815b97f1c478b7b204d0162d55cc8064e9")
	return ret
}()
var TestingSK11 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("1b7d708328c6a8902dae8e9f869e67f112a85194dc053e43517f701df6b3e8e4")
	return ret
}()
var TestingSK12 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("0662a290fd9865742249a0ea1e7fe04bbbe92e35cbe1f397e04d98f910e107e8")
	return ret
}()
var TestingSK13 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("0aed765a3b4667934149b632a2c472ae2d0efecf94d20eac612117dd32a5a64d")
	return ret
}()
var TestingWrongSK = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("344de4028bc9a073185bf5bb8e091f93da0a8f914a2392083ea327248324d256")
	return ret
}()
