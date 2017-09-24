package srp

import (
	"encoding/base64"
	"math/big"
	"testing"
)

var r5054username = "alice"
var r5054password = "password123"

type kdfTestVector struct {
	sk           string
	alg          string
	iterations   int
	method       string
	saltB64      string
	mp           string
	email        string
	expectedXhex string
}

var kh = kdfTestVector{
	sk:      "A3-X5ZKSD-673PE8-CHH5Z-NAJMJ-5YFPF-3N5JP",
	alg:     "PBES2-HS256",
	method:  "SRP-4096",
	saltB64: "BW28BWLu6w9y2unDFwHnKg",
	mp:      "snowdrop ax organism pam",
	email:   "kevinhayestest@icloud.com",

	expectedXhex: "5b401cd715a53a0f2bb27de5554c2dde94d72680ac924094c5adbb74b8355b24",
}

/* from B5Book:/server/api/v1/srp.html

email address: kevinhayestest@icloud.com
password (hashed with SHA256, base64URLEncoded (with trailing '=' characters trimmed): tWkOikZrlNbt4r6CwzJP8EBLDaTEfNj6nRMKZ6k2UKI
    (original password unhashed: snowdrop ax organism pam)
Secret Key: A3-X5ZKSD-673PE8-CHH5Z-NAJMJ-5YFPF-3N5JP
auth params:
{
    alg = "PBES2-HS256";
    iterations = 100000;
    method = "SRP-4096";
    salt = BW28BWLu6w9y2unDFwHnKg;
}

SRP Verifier =
facbcd929a8aa1273cd59b05e10e81dfff4fa188cdf8f9e65b0013b35fe261f180794598faa20ca0da4a96e69026beb109c85ea565ba95fb67c7287f4c2254a120166758220ff64634fba1671798055163ee1669fca1808901c747320717b30b6aea97141de027b4728187ae54bcca6cef656b670e178af9e121678a2f891073f2775c1073825c826f02908f2343dc404b9b67ccd4da72ca732f99d3bb349cdb0e8b4d78bfa79a88bb7175a55e14845bc4ddbdba5b84c736fca94012fecbe1caf4396c7a0dcd447b4346676f817cfdf65d3aaf8a9b2190455c78f62556343b8dcf49809a65ac0052e4039a22218c178c502203cf84d386b18f30a93789d2d96004afa652967bc4426cf26b01a1a5ff7bfa80353e78d24376e250d08cbf1ce8b41fa020ee8e293af0a57384d2a529f2ee601f7b0a0339301dfc716ae07bde740dc95c29d18c6f6991845ef1eb4bfe2db666c4d54bffddef13170ee6c5b41feb2291edb5f2ce2d6fa7d2c9efcdd5123b40fb7a618498d2821341ff251ba03c19ecdb9dacd47b845e8131aab4a624bf8453943153b0a3d586d6a5f2ac8890a463987ec2b2d347ef2a7abec2a26094325fac8e07dd45879ae19bfaf33ab149ee86103a97e35fa27c776e117a29cbd032a46c67fb46adec4f5e903a087548ec7f70d0d3b92fb1bfcf9773cf69c63f1b21b3baebb33f1a3114d0847048f37065f406b2

PBKDF2 Derived = a3361b5791bc9f0bb874cefb7d0c6420224d6fd76ae801d7885ec294ce653aa4
Derived Combined Secret Key = 5b401cd715a53a0f2bb27de5554c2dde94d72680ac924094c5adbb74b8355b24

*/

func TestComputeX(t *testing.T) {

	var ak *AccountKey
	var err error
	var x, expectedX *big.Int

	if ak, err = NewAccountKeyFromString(kh.sk); ak == nil || err != nil {
		t.Errorf("Couldn't construct Secret Key: %s", err)
	}

	salt, _ := base64.RawURLEncoding.DecodeString(kh.saltB64)
	if x, err = CalculateX(kh.method, kh.alg, kh.email, kh.mp, salt, kh.iterations, ak); x == nil || err != nil {
		t.Errorf("Couldn't calculate x: %s", err)
	}

	expectedX, _ = new(big.Int).SetString(kh.expectedXhex, 16)

	if x.Cmp(expectedX) != 0 {
		t.Errorf("Calcuated x: %s\nExpected x: %s\n", x, expectedX)
	}

}
