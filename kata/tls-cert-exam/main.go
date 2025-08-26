package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// ---- Helpers ----

func mustRSA(bits int) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	return key
}

func mustCreateCert(template, parent *x509.Certificate, pub *rsa.PublicKey, parentKey *rsa.PrivateKey) []byte {
	der, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentKey)
	if err != nil {
		panic(err)
	}
	return der
}

func pemEncodeCert(der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

func pemEncodeKey(key *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
}

func mustSerial() *big.Int {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		panic(err)
	}
	return serial
}

func subjectKeyID(pub *rsa.PublicKey) []byte {
	spki, _ := x509.MarshalPKIXPublicKey(pub)
	sum := sha1.Sum(spki)
	return sum[:]
}

func printPEM(title string, b []byte) {
	fmt.Printf("\n===== %s =====\n%s\n", title, string(b))
}

func days(n int) time.Time { return time.Now().Add(time.Duration(n) * 24 * time.Hour) }

// ---- Main flow ----

func main() {
	// 1) 生成 Root CA（自签）
	rootKey := mustRSA(2048)
	rootTmpl := &x509.Certificate{
		SerialNumber:          mustSerial(),
		Subject:               pkix.Name{CommonName: "Demo Root CA", Organization: []string{"Acme Corp"}},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              days(3650), // 10 年
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		SubjectKeyId:          subjectKeyID(&rootKey.PublicKey),
	}
	rootDER := mustCreateCert(rootTmpl, rootTmpl, &rootKey.PublicKey, rootKey)
	rootCert, _ := x509.ParseCertificate(rootDER)

	printPEM("ROOT CA CERT", pemEncodeCert(rootDER))
	printPEM("ROOT CA KEY", pemEncodeKey(rootKey))

	// 2) 生成 Intermediate CA，由 Root 签发
	intKey := mustRSA(2048)
	intTmpl := &x509.Certificate{
		SerialNumber:          mustSerial(),
		Subject:               pkix.Name{CommonName: "Demo Intermediate CA", Organization: []string{"Acme Corp"}},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              days(3650),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		SubjectKeyId:          subjectKeyID(&intKey.PublicKey),
		AuthorityKeyId:        rootCert.SubjectKeyId,
	}
	intDER := mustCreateCert(intTmpl, rootCert, &intKey.PublicKey, rootKey)
	intCert, _ := x509.ParseCertificate(intDER)

	printPEM("INTERMEDIATE CA CERT", pemEncodeCert(intDER))
	printPEM("INTERMEDIATE CA KEY", pemEncodeKey(intKey))

	// 3) 生成服务器私钥与 CSR（包含 SAN）
	leafKey := mustRSA(2048)
	csrTmpl := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: "fanyamin.com", Organization: []string{"Fanyamin Ltd"}},
		DNSNames: []string{"fanyamin.com", "www.fanyamin.com", "api.fanyamin.com"}, // SAN
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, csrTmpl, leafKey)
	if err != nil {
		panic(err)
	}
	csr, err := x509.ParseCertificateRequest(csrDER)
	if err != nil {
		panic(err)
	}
	if err := csr.CheckSignature(); err != nil {
		panic(fmt.Errorf("CSR signature invalid: %w", err))
	}
	printPEM("SERVER CSR", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER}))
	printPEM("SERVER KEY", pemEncodeKey(leafKey))

	// 4) 中级 CA 根据 CSR 签发服务器证书（服务器证书不是 CA）
	leafTmpl := &x509.Certificate{
		SerialNumber:          mustSerial(),
		Subject:               csr.Subject, // 沿用 CSR 主题
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              days(825), // 约 27 个月，符合常见限制
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		DNSNames:              csr.DNSNames, // 从 CSR 注入 SAN
		SubjectKeyId:          subjectKeyID(&leafKey.PublicKey),
		AuthorityKeyId:        intCert.SubjectKeyId,
	}
	leafDER := mustCreateCert(leafTmpl, intCert, &leafKey.PublicKey, intKey)
	leafCert, _ := x509.ParseCertificate(leafDER)
	printPEM("SERVER CERT", pemEncodeCert(leafDER))

	// 5) 构建验证所需的根与中级池，并验证证书链 + 主机名
	roots := x509.NewCertPool()
	roots.AddCert(rootCert)

	inters := x509.NewCertPool()
	inters.AddCert(intCert)

	// 验证：主机名 www.fanyamin.com 是否被证书覆盖（SAN 中包含）
	if _, err := leafCert.Verify(x509.VerifyOptions{
		DNSName:       "www.fanyamin.com",
		Roots:         roots,
		Intermediates: inters,
	}); err != nil {
		panic(fmt.Errorf("verify failed (www.fanyamin.com): %w", err))
	}
	fmt.Println("✓ 证书链与主机名验证通过：www.fanyamin.com")

	// 故意用一个未在 SAN 中的主机名，看看失败情况
	if _, err := leafCert.Verify(x509.VerifyOptions{
		DNSName:       "not-in-san.fanyamin.com",
		Roots:         roots,
		Intermediates: inters,
	}); err != nil {
		fmt.Println("✗ 预期的验证失败（主机名不在 SAN）：", err)
	} else {
		panic("unexpected: verification should have failed for not-in-san.fanyamin.com")
	}

	// 6) 补充：说明数字签名在链路验证中的体现
	fmt.Println("\n说明：以上 x509.Verify 实际做了“用上级证书公钥验证下级证书签名”的工作，")
	fmt.Println("即：Root 公钥 验证 Intermediate 的签名；Intermediate 公钥 验证 Leaf 的签名；")
	fmt.Println("再结合 DNSName 检查 SAN 与证书有效期、用途等扩展，最终确定是否可信。")
}
