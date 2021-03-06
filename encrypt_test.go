package kgo

import (
	"crypto/aes"
	"fmt"
	"strings"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	str := []byte("This is an string to encod")
	res := KEncr.Base64Encode(str)
	if !strings.HasSuffix(res, "=") {
		t.Error("Base64Encode fail")
		return
	}
}

func BenchmarkBase64Encode(b *testing.B) {
	b.ResetTimer()
	str := []byte("This is an string to encod")
	for i := 0; i < b.N; i++ {
		KEncr.Base64Encode(str)
	}
}

func TestBase64Decode(t *testing.T) {
	str := "VGhpcyBpcyBhbiBlbmNvZGVkIHN0cmluZw=="
	_, err := KEncr.Base64Decode(str)
	if err != nil {
		t.Error("Base64Decode fail")
		return
	}
	_, err = KEncr.Base64Decode("#iu3498r")
	if err == nil {
		t.Error("Base64Decode fail")
		return
	}
	_, err = KEncr.Base64Decode("VGhpcy")
	_, err = KEncr.Base64Decode("VGhpcyB")
}

func BenchmarkBase64Decode(b *testing.B) {
	b.ResetTimer()
	str := "VGhpcyBpcyBhbiBlbmNvZGVkIHN0cmluZw=="
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.Base64Decode(str)
	}
}

func TestBase64UrlEncodeDecode(t *testing.T) {
	str := []byte("This is an string to encod")
	res := KEncr.Base64UrlEncode(str)
	if strings.HasSuffix(res, "=") {
		t.Error("Base64UrlEncode fail")
		return
	}

	_, err := KEncr.Base64UrlDecode(res)
	if err != nil {
		t.Error("Base64UrlDecode fail")
		return
	}
}

func BenchmarkBase64UrlEncode(b *testing.B) {
	b.ResetTimer()
	str := []byte("This is an string to encod")
	for i := 0; i < b.N; i++ {
		KEncr.Base64UrlEncode(str)
	}
}

func BenchmarkBase64UrlDecode(b *testing.B) {
	b.ResetTimer()
	str := "VGhpcyBpcyBhbiBzdHJpbmcgdG8gZW5jb2Q"
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.Base64UrlDecode(str)
	}
}

func TestAuthCode(t *testing.T) {
	key := "123456"
	str := "hello world"

	res, _ := KEncr.AuthCode(str, key, true, 0)
	if res == "" {
		t.Error("AuthCode Encode fail")
		return
	}

	res2, _ := KEncr.AuthCode(res, key, false, 0)
	if res2 == "" {
		t.Error("AuthCode Decode fail")
		return
	}

	res, _ = KEncr.AuthCode(str, key, true, -3600)
	KEncr.AuthCode(res, key, false, 0)
	KEncr.AuthCode("", key, true, 0)
	KEncr.AuthCode("", "", true, 0)
	KEncr.AuthCode("7caeNfPt/N1zHdj5k/7i7pol6NHsVs0Cji7c15h4by1RYcrBoo7EEw==", key, false, 0)
	KEncr.AuthCode("7caeNfPt/N1zHdj5k/7i7pol6N", key, false, 0)
	KEncr.AuthCode("123456", "", false, 0)
	KEncr.AuthCode("1234#iu3498r", "", false, 0)
}

func BenchmarkAuthCodeEncode(b *testing.B) {
	b.ResetTimer()
	key := "123456"
	str := "hello world"
	for i := 0; i < b.N; i++ {
		KEncr.AuthCode(str, key, true, 0)
	}
}

func BenchmarkAuthCodeDecode(b *testing.B) {
	b.ResetTimer()
	key := "123456"
	str := "a79b5do3C9nbaZsAz5j3NQRj4e/L6N+y5fs2U9r1mO0LinOWtxmscg=="
	for i := 0; i < b.N; i++ {
		KEncr.AuthCode(str, key, false, 0)
	}
}

func TestPasswordHashVerify(t *testing.T) {
	pwd := []byte("123456")
	has, err := KEncr.PasswordHash(pwd)
	if err != nil {
		t.Error("PasswordHash fail")
		return
	}

	chk := KEncr.PasswordVerify(pwd, has)
	if !chk {
		t.Error("PasswordVerify fail")
		return
	}

	_, _ = KEncr.PasswordHash(pwd, 1)
	//慎用20以上,太耗时
	_, _ = KEncr.PasswordHash(pwd, 15)
	_, _ = KEncr.PasswordHash(pwd, 33)
}

func BenchmarkPasswordHash(b *testing.B) {
	b.ResetTimer()
	pwd := []byte("123456")
	for i := 0; i < b.N; i++ {
		//太耗时,只测试少量的
		if i > 10 {
			break
		}
		_, _ = KEncr.PasswordHash(pwd)
	}
}

func BenchmarkPasswordVerify(b *testing.B) {
	b.ResetTimer()
	pwd := []byte("123456")
	has := []byte("$2a$10$kCv6ljsVuTSI54oPkWulreEmUNTW/zj0Dgh6qF4Vz0w4C3gVf/w7a")
	for i := 0; i < b.N; i++ {
		//太耗时,只测试少量的
		if i > 10 {
			break
		}
		KEncr.PasswordVerify(pwd, has)
	}
}

func TestEasyEncryptDecrypt(t *testing.T) {
	key := "123456"
	str := "hello world你好!hello world你好!hello world你好!hello world你好!"
	enc := KEncr.EasyEncrypt(str, key)
	if enc == "" {
		t.Error("EasyEncrypt fail")
		return
	}

	dec := KEncr.EasyDecrypt(enc, key)
	if dec != str {
		t.Error("EasyDecrypt fail")
		return
	}

	dec = KEncr.EasyDecrypt("你好，世界！", key)
	if dec != "" {
		t.Error("EasyDecrypt fail")
		return
	}

	KEncr.EasyEncrypt("", key)
	KEncr.EasyEncrypt("", "")
	KEncr.EasyDecrypt(enc, "1qwer")
	KEncr.EasyDecrypt("123", key)
	KEncr.EasyDecrypt("1234#iu3498r", key)
}

func BenchmarkEasyEncrypt(b *testing.B) {
	b.ResetTimer()
	key := "123456"
	str := "hello world你好"
	for i := 0; i < b.N; i++ {
		KEncr.EasyEncrypt(str, key)
	}
}

func BenchmarkEasyDecrypt(b *testing.B) {
	b.ResetTimer()
	key := "123456"
	str := "e10azZaczdODqqimpcY"
	for i := 0; i < b.N; i++ {
		KEncr.EasyDecrypt(str, key)
	}
}

func TestHmacShaX(t *testing.T) {
	str := []byte("hello world")
	key := []byte("123456")
	res1 := KEncr.HmacShaX(str, key, 1)
	res2 := KEncr.HmacShaX(str, key, 256)
	res3 := KEncr.HmacShaX(str, key, 512)
	if res1 == "" || res2 == "" || res3 == "" {
		t.Error("HmacShaX fail")
		return
	}
}

func TestHmacShaXPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover...:", r)
		}
	}()

	str := []byte("hello world")
	key := []byte("123456")
	KEncr.HmacShaX(str, key, 4)
}

func BenchmarkHmacShaX(b *testing.B) {
	b.ResetTimer()
	str := []byte("hello world")
	key := []byte("123456")
	for i := 0; i < b.N; i++ {
		KEncr.HmacShaX(str, key, 256)
	}
}

func TestPkcs7PaddingUnPadding(t *testing.T) {
	var emp1 []byte
	var emp2 = []byte("")
	key1 := []byte("1234")
	dat1 := []byte{49, 50, 51, 52, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	dat2 := []byte{49, 50, 51, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var tests = []struct {
		cipher    []byte
		orig      []byte
		size      int
		zero      bool
		expected1 []byte
		expected2 []byte
	}{
		{nil, nil, aes.BlockSize, false, nil, nil},
		{emp1, emp1, aes.BlockSize, false, nil, nil},
		{emp2, emp2, aes.BlockSize, false, nil, nil},
		{key1, key1, 0, false, nil, nil},
		{key1, dat1, aes.BlockSize, false, dat1, key1},
		{key1, dat2, aes.BlockSize, true, dat2, nil},
		{key1, dat2, aes.BlockSize, false, dat1, emp1},
	}

	for _, test := range tests {
		actual1 := pkcs7Padding(test.cipher, test.size, test.zero)
		if !KArr.IsEqualArray(actual1, test.expected1) {
			t.Errorf("Expected pkcs7Padding(%v, %d, %t) to be %v, got %v", test.cipher, test.size, test.zero, test.expected1, actual1)
		}

		actual2 := pkcs7UnPadding(test.orig, test.size)
		if !KArr.IsEqualArray(actual2, test.expected2) {
			t.Errorf("Expected pkcs7UnPadding(%v, %d) to be %v, got %v", test.orig, test.size, test.expected2, actual2)
		}
	}
}

func BenchmarkPkcs7Padding(b *testing.B) {
	b.ResetTimer()
	str := []byte("1234")
	for i := 0; i < b.N; i++ {
		pkcs7Padding(str, 16, false)
	}
}

func BenchmarkPkcs7UnPadding(b *testing.B) {
	b.ResetTimer()
	data := []byte{49, 50, 51, 52, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	for i := 0; i < b.N; i++ {
		pkcs7UnPadding(data, 16)
	}
}

func TestZeroPaddingUnPadding(t *testing.T) {
	key := []byte("hello")
	ori := zeroPadding(key, 16)
	res := zeroUnPadding(ori)

	if ori == nil {
		t.Error("zeroPadding fail")
		return
	}

	if !KArr.IsEqualArray(key, res) {
		t.Error("zeroUnPadding fail")
		return
	}
}

func BenchmarkZeroPadding(b *testing.B) {
	b.ResetTimer()
	key := []byte("hello")
	for i := 0; i < b.N; i++ {
		zeroPadding(key, 16)
	}
}

func BenchmarkZeroUnPadding(b *testing.B) {
	b.ResetTimer()
	ori := []byte{104, 101, 108, 108, 111, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < b.N; i++ {
		zeroUnPadding(ori)
	}
}

func TestAesCBCEncryptDecrypt(t *testing.T) {
	ori := []byte("hello")
	key := []byte("1234567890123456")
	emp := []byte("")
	var err error
	var enc, des []byte

	_, err = KEncr.AesCBCEncrypt(ori, []byte("123"))
	if err == nil {
		t.Error("AesCBCEncrypt fail")
		return
	}

	enc, err = KEncr.AesCBCEncrypt(ori, key)
	des, err = KEncr.AesCBCDecrypt(enc, key)
	if !KArr.IsEqualArray(ori, des) {
		t.Error("AesCBCEncrypt fail")
		return
	}

	enc, err = KEncr.AesCBCEncrypt(ori, key, PKCS_SEVEN)
	des, err = KEncr.AesCBCDecrypt(enc, key, PKCS_SEVEN)
	if !KArr.IsEqualArray(ori, des) {
		t.Error("AesCBCEncrypt fail")
		return
	}

	enc, err = KEncr.AesCBCEncrypt(emp, key, PKCS_SEVEN)
	des, err = KEncr.AesCBCDecrypt(enc, key, PKCS_SEVEN)
	if !KArr.IsEqualArray(emp, des) {
		t.Error("AesCBCEncrypt fail")
		return
	}

	enc, err = KEncr.AesCBCEncrypt(ori, key, PKCS_ZERO)
	des, err = KEncr.AesCBCDecrypt(enc, key, PKCS_ZERO)
	if !KArr.IsEqualArray(ori, des) {
		t.Error("AesCBCEncrypt fail")
		return
	}

	enc = []byte{83, 28, 170, 254, 29, 174, 21, 129, 241, 233, 243, 84, 1, 250, 95, 122, 104, 101, 108, 108, 111, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	des, err = KEncr.AesCBCDecrypt(enc, key, PKCS_ZERO)
	if err == nil {
		t.Error("AesCBCDecrypt fail")
		return
	}

	_, err = KEncr.AesCBCDecrypt(enc, []byte("1"))
	if err == nil {
		t.Error("AesCBCDecrypt fail")
		return
	}

	_, err = KEncr.AesCBCDecrypt([]byte("1234"), key)
	if err == nil {
		t.Error("AesCBCDecrypt fail")
		return
	}

}

func BenchmarkAesCBCEncrypt(b *testing.B) {
	b.ResetTimer()
	ori := []byte("hello")
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesCBCEncrypt(ori, key)
	}
}

func BenchmarkAesCBCDecryptZero(b *testing.B) {
	b.ResetTimer()
	enc := []byte{214, 214, 97, 208, 185, 68, 246, 40, 124, 3, 155, 58, 5, 84, 136, 10, 104, 101, 108, 108, 111, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesCBCDecrypt(enc, key, PKCS_ZERO)
	}
}

func BenchmarkAesCBCDecryptSeven(b *testing.B) {
	b.ResetTimer()
	enc := []byte{17, 195, 8, 206, 231, 183, 143, 246, 244, 137, 216, 185, 120, 175, 90, 111, 104, 101, 108, 108, 111, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11}
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesCBCDecrypt(enc, key, PKCS_SEVEN)
	}
}

func TestAesCFBEncryptDecrypt(t *testing.T) {
	ori := []byte("hello")
	key := []byte("1234567890123456")
	emp := []byte("")
	var err error
	var enc, des []byte

	_, err = KEncr.AesCFBEncrypt(ori, []byte("123"))
	if err == nil {
		t.Error("AesCFBEncrypt fail")
		return
	}

	enc, err = KEncr.AesCFBEncrypt(ori, key)
	des, err = KEncr.AesCFBDecrypt(enc, key)
	if !KArr.IsEqualArray(ori, des) {
		t.Error("AesCFBEncrypt fail")
		return
	}

	enc, err = KEncr.AesCFBEncrypt(emp, key)
	des, err = KEncr.AesCFBDecrypt(enc, key)
	if !KArr.IsEqualArray(emp, des) {
		t.Error("AesCFBEncrypt fail")
		return
	}

	_, err = KEncr.AesCFBDecrypt(enc, []byte("1"))
	if err == nil {
		t.Error("AesCFBDecrypt fail")
		return
	}

	_, err = KEncr.AesCFBDecrypt([]byte("1234"), key)
	if err == nil {
		t.Error("AesCFBDecrypt fail")
		return
	}

}

func BenchmarkAesCFBEncrypt(b *testing.B) {
	b.ResetTimer()
	ori := []byte("hello")
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesCFBEncrypt(ori, key)
	}
}

func BenchmarkAesCFBDecrypt(b *testing.B) {
	b.ResetTimer()
	enc := []byte{150, 234, 226, 46, 34, 206, 171, 155, 186, 66, 116, 201, 63, 67, 227, 217, 104, 101, 108, 108, 111}
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesCFBDecrypt(enc, key)
	}
}

func TestAesCTREncryptDecrypt(t *testing.T) {
	ori := []byte("hello")
	key := []byte("1234567890123456")
	emp := []byte("")
	var err error
	var enc, des []byte

	_, err = KEncr.AesCTREncrypt(ori, []byte("123"))
	if err == nil {
		t.Error("AesCTREncrypt fail")
		return
	}

	enc, err = KEncr.AesCTREncrypt(ori, key)
	des, err = KEncr.AesCTRDecrypt(enc, key)
	if !KArr.IsEqualArray(ori, des) {
		t.Error("AesCTREncrypt fail")
		return
	}

	enc, err = KEncr.AesCTREncrypt(emp, key)
	des, err = KEncr.AesCTRDecrypt(enc, key)
	if !KArr.IsEqualArray(emp, des) {
		t.Error("AesCTREncrypt fail")
		return
	}

	_, err = KEncr.AesCTRDecrypt(enc, []byte("1"))
	if err == nil {
		t.Error("AesCTRDecrypt fail")
		return
	}

	_, err = KEncr.AesCTRDecrypt([]byte("1234"), key)
	if err == nil {
		t.Error("AesCTRDecrypt fail")
		return
	}

}

func BenchmarkAesCTREncrypt(b *testing.B) {
	b.ResetTimer()
	ori := []byte("hello")
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesCTREncrypt(ori, key)
	}
}

func BenchmarkAesCTRDecrypt(b *testing.B) {
	b.ResetTimer()
	enc := []byte{225, 187, 161, 145, 117, 191, 229, 20, 164, 43, 242, 23, 138, 241, 74, 27, 104, 101, 108, 108, 111}
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesCTRDecrypt(enc, key)
	}
}

func TestAesOFBEncryptDecrypt(t *testing.T) {
	ori := []byte("hello")
	key := []byte("1234567890123456")
	emp := []byte("")
	var err error
	var enc, des []byte

	_, err = KEncr.AesOFBEncrypt(ori, []byte("123"))
	if err == nil {
		t.Error("AesOFBEncrypt fail")
		return
	}

	enc, err = KEncr.AesOFBEncrypt(ori, key)
	des, err = KEncr.AesOFBDecrypt(enc, key)
	if !KArr.IsEqualArray(ori, des) {
		t.Error("AesOFBEncrypt fail")
		return
	}

	enc, err = KEncr.AesOFBEncrypt(emp, key)
	des, err = KEncr.AesOFBDecrypt(enc, key)
	if !KArr.IsEqualArray(emp, des) {
		t.Error("AesOFBEncrypt fail")
		return
	}

	_, err = KEncr.AesOFBDecrypt(enc, []byte("1"))
	if err == nil {
		t.Error("AesOFBDecrypt fail")
		return
	}

	_, err = KEncr.AesOFBDecrypt([]byte("1234"), key)
	if err == nil {
		t.Error("AesOFBDecrypt fail")
		return
	}

}

func BenchmarkAesOFBEncrypt(b *testing.B) {
	b.ResetTimer()
	ori := []byte("hello")
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesOFBEncrypt(ori, key)
	}
}

func BenchmarkAesOFBDecrypt(b *testing.B) {
	b.ResetTimer()
	enc := []byte{66, 87, 29, 157, 2, 128, 196, 94, 141, 224, 221, 41, 162, 41, 159, 207, 104, 101, 108, 108, 111}
	key := []byte("1234567890123456")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.AesOFBDecrypt(enc, key)
	}
}

func TestGenerateRsaKeys(t *testing.T) {
	_, _, err := KEncr.GenerateRsaKeys(1)
	if err == nil {
		t.Error("GenerateRsaKeys fail")
		return
	}

	pri, pub, err := KEncr.GenerateRsaKeys(2048)
	if len(pri) < 100 || len(pub) < 100 {
		t.Error("GenerateRsaKeys fail")
		return
	}

	if err != nil {
		println(err.Error())
	}
}

func BenchmarkGenerateRsaKeys(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = KEncr.GenerateRsaKeys(1024)
	}
}

func TestRsaPublicEncryptPrivateDecrypt(t *testing.T) {
	var enc, des []byte
	var err error
	word := []byte("hello world")
	pubkey1, _ := KFile.ReadFile("testdata/rsa/public_key.pem")
	prikey1, _ := KFile.ReadFile("testdata/rsa/private_key.pem")
	pubkey2 := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDteXRcRyppm5sOVvteo37Dmaidbx6YrV6QWZ0L9mGfCmSW1a/A
d61kT6OoU0Z3DyId7vA9TtvULucEUpywPpSoP/r+820UHFihdyhcb1iy8Z3v6KUc
arWzUOZpo0mc+o4hW2O1VnzNxLcXmhQOA9NdEOV/M+zxubFKo4VsY0ti9QIDAQAB
AoGAZuD/MBsEnMv02LmGHPHnsQWYrtu8/ZfeJ9sq1kve7u+ptE7O3Sr7y0FVPU8W
b+32cdFZ8rV/NuU63/yKNTBnZcbPwwGV9DmNpXy9YCdjwXkxfjYiDqUX9Fsxth1M
EqMb0PRO85akxCKxxtMagHDHNWkQaVThLagG31sh5d38SwECQQDuVsbRTbEz/H/j
Ip1NNU+8XERwMv1ac0LE9GhSRlqzUWDhukQ1gp9DmoKic8QMr6DS+JYvTCq38J8t
LHMNmzcpAkEA/xJHH/MwRlUSHsfP+DGXBuue2cAyw3NVLgusNV222kIgDOLcVxLl
8YOAgnheD5iI8+/GIVB4cXIfXKgqvzMC7QJAPUg8uMaEQLy02V8mGRsTFHiY9Ex4
DlDCo0fApx8F5UOQaJnvPd8HOme5HTIs/6IM9RIL879e4IrTMtdSAfad+QJBANAc
Opmv0mBgAnPItT8cPsvvrGCfdwuO6x2xemTkPE9hikLZSctlaOUfVNeem6f/3SWi
-----END RSA PRIVATE KEY-----`
	prikey2 := `-----BEGIN RSA PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDteXRcRyppm5sOVvteo37Dmaid
bx6YrV6QWZ0L9mGfCmSW1a/Ad61kT6OoU0Z3DyId7vA9TtvULucEUpywPpSoP/r+
820UHFihdyhcb1iy8Z3v6KUcarWzUOZpo0mc+o4hW2O1VnzNxLcXmhQOA9NdEOV/
-----END RSA PUBLIC KEY-----`

	enc, err = KEncr.RsaPublicEncrypt(word, pubkey1)
	if err != nil {
		t.Error("RsaPublicEncrypt fail")
		return
	}

	des, err = KEncr.RsaPrivateDecrypt(enc, prikey1)
	if err != nil {
		t.Error("RsaPrivateDecrypt fail")
		return
	}

	if !KArr.IsEqualArray(word, des) {
		t.Error("RsaPrivateDecrypt fail")
		return
	}

	_, err = KEncr.RsaPublicEncrypt(word, []byte("123"))
	if err == nil {
		t.Error("RsaPublicEncrypt fail")
		return
	}

	_, err = KEncr.RsaPublicEncrypt(word, []byte(pubkey2))
	if err == nil {
		t.Error("RsaPublicEncrypt fail")
		return
	}

	_, err = KEncr.RsaPrivateDecrypt(enc, []byte("123"))
	if err == nil {
		t.Error("RsaPrivateDecrypt fail")
		return
	}

	_, err = KEncr.RsaPrivateDecrypt(enc, []byte(prikey2))
	if err == nil {
		t.Error("RsaPrivateDecrypt fail")
		return
	}

}

func BenchmarkRsaPublicEncrypt(b *testing.B) {
	b.ResetTimer()
	word := []byte("hello world")
	pubkey, _ := KFile.ReadFile("testdata/rsa/public_key.pem")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.RsaPublicEncrypt(word, pubkey)
	}
}

func BenchmarkRsaPrivateDecrypt(b *testing.B) {
	b.ResetTimer()
	data := []byte{143, 167, 230, 243, 173, 106, 253, 203, 191, 77, 142, 78, 116, 8, 81, 120, 197, 206, 141, 219, 255, 210, 42, 71, 202, 47, 153, 60, 152, 163, 160, 226, 110, 102, 50, 20, 165, 181, 236, 160, 109, 229, 1, 11, 80, 164, 9, 56, 188, 66, 199, 227, 69, 88, 88, 143, 159, 211, 41, 169, 231, 215, 241, 35, 79, 208, 44, 43, 143, 163, 64, 107, 166, 128, 101, 106, 73, 248, 161, 36, 201, 161, 171, 241, 227, 114, 137, 28, 156, 63, 147, 52, 189, 230, 136, 90, 123, 21, 73, 172, 188, 8, 53, 98, 36, 185, 131, 171, 222, 52, 124, 48, 207, 82, 123, 234, 5, 97, 53, 47, 234, 6, 81, 118, 81, 161, 130, 172}
	prikey, _ := KFile.ReadFile("testdata/rsa/private_key.pem")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.RsaPrivateDecrypt(data, prikey)
	}
}

func TestRsaPrivateEncryptPublicDecrypt(t *testing.T) {
	var enc, des []byte
	var err error
	str := strings.Repeat("hello world", 10)
	word := []byte(str)
	pubkey1, _ := KFile.ReadFile("testdata/rsa/public_key.pem")
	prikey1, _ := KFile.ReadFile("testdata/rsa/private_key.pem")
	pubkey2 := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDteXRcRyppm5sOVvteo37Dmaidbx6YrV6QWZ0L9mGfCmSW1a/A
d61kT6OoU0Z3DyId7vA9TtvULucEUpywPpSoP/r+820UHFihdyhcb1iy8Z3v6KUc
arWzUOZpo0mc+o4hW2O1VnzNxLcXmhQOA9NdEOV/M+zxubFKo4VsY0ti9QIDAQAB
AoGAZuD/MBsEnMv02LmGHPHnsQWYrtu8/ZfeJ9sq1kve7u+ptE7O3Sr7y0FVPU8W
b+32cdFZ8rV/NuU63/yKNTBnZcbPwwGV9DmNpXy9YCdjwXkxfjYiDqUX9Fsxth1M
EqMb0PRO85akxCKxxtMagHDHNWkQaVThLagG31sh5d38SwECQQDuVsbRTbEz/H/j
Ip1NNU+8XERwMv1ac0LE9GhSRlqzUWDhukQ1gp9DmoKic8QMr6DS+JYvTCq38J8t
LHMNmzcpAkEA/xJHH/MwRlUSHsfP+DGXBuue2cAyw3NVLgusNV222kIgDOLcVxLl
8YOAgnheD5iI8+/GIVB4cXIfXKgqvzMC7QJAPUg8uMaEQLy02V8mGRsTFHiY9Ex4
DlDCo0fApx8F5UOQaJnvPd8HOme5HTIs/6IM9RIL879e4IrTMtdSAfad+QJBANAc
Opmv0mBgAnPItT8cPsvvrGCfdwuO6x2xemTkPE9hikLZSctlaOUfVNeem6f/3SWi
-----END RSA PRIVATE KEY-----`
	prikey2 := `-----BEGIN RSA PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDteXRcRyppm5sOVvteo37Dmaid
bx6YrV6QWZ0L9mGfCmSW1a/Ad61kT6OoU0Z3DyId7vA9TtvULucEUpywPpSoP/r+
820UHFihdyhcb1iy8Z3v6KUcarWzUOZpo0mc+o4hW2O1VnzNxLcXmhQOA9NdEOV/
-----END RSA PUBLIC KEY-----`

	enc, err = KEncr.RsaPrivateEncrypt(word, prikey1)
	if err != nil {
		println(err.Error())
		t.Error("RsaPrivateEncrypt fail")
		return
	}

	des, err = KEncr.RsaPublicDecrypt(enc, pubkey1)
	if err != nil {
		t.Error("RsaPublicDecrypt fail")
		return
	}

	if !KArr.IsEqualArray(word, des) {
		t.Error("RsaPublicDecrypt fail")
		return
	}

	_, err = KEncr.RsaPrivateEncrypt(word, []byte("123"))
	if err == nil {
		t.Error("RsaPrivateEncrypt fail")
		return
	}

	_, err = KEncr.RsaPrivateEncrypt(word, []byte(prikey2))
	if err == nil {
		t.Error("RsaPrivateEncrypt fail")
		return
	}

	_, err = KEncr.RsaPublicDecrypt(enc, []byte("123"))
	if err == nil {
		t.Error("RsaPublicDecrypt fail")
		return
	}

	_, err = KEncr.RsaPublicDecrypt(enc, []byte(pubkey2))
	if err == nil {
		t.Error("RsaPublicDecrypt fail")
		return
	}

}

func BenchmarkRsaPrivateEncrypt(b *testing.B) {
	b.ResetTimer()
	word := []byte("hello world")
	prikey, _ := KFile.ReadFile("testdata/rsa/private_key.pem")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.RsaPrivateEncrypt(word, prikey)
	}
}

func BenchmarkRsaPublicDecrypt(b *testing.B) {
	b.ResetTimer()
	data := []byte{134, 85, 170, 196, 249, 255, 241, 73, 245, 105, 254, 226, 205, 183, 69, 1, 214, 60, 209, 162, 8, 50, 87, 148, 215, 2, 198, 212, 82, 5, 49, 39, 219, 182, 194, 12, 198, 23, 0, 99, 99, 145, 32, 138, 182, 104, 0, 190, 69, 46, 213, 2, 243, 139, 161, 15, 0, 69, 242, 145, 240, 86, 173, 242, 7, 71, 151, 160, 145, 21, 15, 117, 7, 202, 243, 70, 11, 105, 247, 198, 192, 213, 152, 56, 85, 76, 237, 38, 155, 78, 81, 212, 160, 223, 41, 54, 143, 110, 214, 97, 138, 180, 139, 240, 178, 14, 67, 77, 19, 169, 103, 222, 34, 172, 5, 141, 64, 8, 63, 17, 72, 180, 54, 59, 20, 105, 124, 221}
	pubkey, _ := KFile.ReadFile("testdata/rsa/public_key.pem")
	for i := 0; i < b.N; i++ {
		_, _ = KEncr.RsaPublicDecrypt(data, pubkey)
	}
}
