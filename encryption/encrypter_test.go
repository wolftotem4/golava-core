package encryption

import (
	"encoding/base64"
	"testing"
)

func TestEncryption(t *testing.T) {
	plaintext := "Hello, World!"

	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Error generating key: %v", err)
	}

	encrypter := NewEncrypter(key)
	encrypted, err := encrypter.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Error encrypting value: %v", err)
	}

	decrypted, err := encrypter.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Error decrypting value: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("Expected %s, got %s", plaintext, string(decrypted))
	}
}

func TestGCMEncrypt(t *testing.T) {
	plaintext := []byte(`Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`)

	key := []byte("12345678901234567890123456789012")
	nonce := []byte("255e59e73ed3")

	encrypted, err := NewEncrypter(key).gcmEncrypt(plaintext, nonce)
	if err != nil {
		t.Fatalf("Error encrypting value: %v", err)
	}

	result := base64.StdEncoding.EncodeToString(encrypted)
	expectsStr := `IWDACTj/Lj9+DzxEQUdDwGwUga6RSCkE1mvU0Bxj2mpbX88fJteyOjz8QbjSE/t/0qjDONO/DM1Gob6Lb81iF7j/vGKn/RtgwRVrulzIk4k4Sv7qS+3j9aDuwNLxKtzCwUnA9plP3DkeMQzbE+oj1Lm+SjqQn8wgaP3p+P6FcVts3d+nl8CtWvLZCF9viJnMDnBr8Tyh0NKEnDJLmemMcrw7alKaZ1BX8FVDvPw3jtHX5lgoG2ebaeIS6obe+YXzrdVVxrk4DyU/Fl5sz+9qkGCMOnkFTq/8KzV7SdXwZwDd4b2kyUosY/ORE1hdCYhdbFqmGyd8NtfADQx8eIrw8BFFdXNlOYUqCZ3Tvz4r4oLtCOZ3o0k1COb/nsCwI7po3PBfONq1wo8IVeMM6e9+VRbwtiR8d8/MELtTIq38uaRBIp8R7J9ClG3iq9BxjnujWH3UBsPFWdVZAhNYJuVurfmvxWPf1dPCGIDwJw2kHmpox4TLcfSBJMkF3Mz+f9fSr72rXBF5VNqqtCtr9mhMshHkLcuCIJmM3uPbSs1hBQ5ngyXzN/Mr4fVNMuY/A5+gAt3sUK3sDOrtK4VlgUNntz0r6UFU3E6rdgG5vZ42NMOeHarHRFCWT2kEkW9WIISpeuLifoR5rs0CECQg3jVhIqzHPRcxcz9ib29t9IzTj1wXW6wC4NOdemy646sTLLQxsnyt9/BfjIJ8Zl/6iNTXf1drt2Fda+LyfmrsLt+zbgj/EM6nBY/GELoXWrsZaXMQ4vQMhSoQApLjfz7/+BOM/PQl7PSUQnsMmJxRg//5/YkgxZ9w4FavGgv7g5Taj2vlnU5ehSY42XwtSmfq+E7KsCtonsFR6GUQfv38UHxRmb6iR7DeK8x9CqQfWElNumll+JrxJ3ovKlWN3l9cNTUAN4ANJfK8xbMaBHUmeubmCJoEzD58jTd22njMGTzcJ9TBaWLID5gDYs2ZX4RAJknYcvpn1LKfyrRwAIfeF86azq5ekHYFRFn8/3apR3a7zWGc+0hON7K0fYgJEok=`

	if result != expectsStr {
		t.Errorf("Expected %s, got %s", expectsStr, result)
	}
}

func TestCBCEncrypt(t *testing.T) {
	plaintext := []byte(`Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`)

	key := []byte("12345678901234567890123456789012")
	iv := []byte("255e59e73ed37fc0")

	encrypted, err := NewEncrypter(key).cbcEncrypt(plaintext, iv)
	if err != nil {
		t.Fatalf("Error encrypting value: %v", err)
	}

	result := base64.StdEncoding.EncodeToString(encrypted)
	expectsStr := `+mL9otPqjUbNM/CTksY0XM4bwEhW1QmUFQsCoks+3XlB2IFms3dkEbHn+XqL6wHEquGjYqwRdBN0jVrsprqDeQAbAo3DawpqhBtUzCruA+aGVdHO5ZfLT+oL6aad2uF5LWqHiJzzicwbF3TXT1O1LbwMYeyb8iWB/T9bQcsnQGB6sASc78I/6WjqgkLBXgAUKzFQboTYF3lQSfp3t0ei248B0mU8Rn5uPueR0a+OsGZQK190tx90FnWKNX1gIED91wLMqASBJYIMxOxpT5IYAaKv5vcA2Ea2/RWRfWJgH9udXGf36v7WkTInjN8bqsXWdYc41GAzzrPT93pUvTnZAz7YnbY8Fe/oLFrAqSZRzQ9Zkv2Bxwczktp1d2KJE8VyOhjYyPUtsKykMAknNOyxA9GQrC5WhztS0F+9Sc2pK0NR66ECn6pQybJNzptSBoAPQCLLbUfpyN+FiIyCURjkUVbfh0qH+spbHz7ugzOY+8F3L4psoF3Fr0lZhedvVW0PXIbdIefWeh7Ym0SWZ7dlj5Lyked/Xs1fJuc6EoaN7hWUHkGjBrC2tzts5HL5HwYZXQqVzYytplO9VZ4pA0Vtrz2etHIBAIKyqaQE5jfzPoWA3Q5f6cSm5yu4+JWlaJJxOKH5Na4AAgQ2A/mt0yDhkTac+CB3QYyBrJKYpRdZQARpec58tLQDPR437Ve3FQI7hE+Tg7g9PyXyeVSUikagpVlLoPYRSa5CZiuctE9LofUhy/Ex1bzA81V4LvmUGW012tY0DU8yvPDbVUNIJI1xzwGTePsCOFJ2om5e3Y5w/be8aN7soKyyqupzfksh7z4/2+a6u56CJC/5JrXWoDCR9lJXYgB+H0i5dcGNboi2x5AJXjYw1PS4ZsXxbFeSd0AIOXUtOjIfql8X4AIla0S65smoq4Rnsvfla7WQ8e3OlD03jP1HcdZ4Tpzus6Y/6pTFm+qjO96jZYiAu6geAUZBmugfXPcJEtZ7lZCdjzre/17jkFO1mjxB0OUpCzX9kHfl`

	if result != expectsStr {
		t.Errorf("Expected %s, got %s", expectsStr, result)
	}
}

func TestGCMDecrypt(t *testing.T) {
	ciphertext, err := base64.StdEncoding.DecodeString(`IWDACTj/Lj9+DzxEQUdDwGwUga6RSCkE1mvU0Bxj2mpbX88fJteyOjz8QbjSE/t/0qjDONO/DM1Gob6Lb81iF7j/vGKn/RtgwRVrulzIk4k4Sv7qS+3j9aDuwNLxKtzCwUnA9plP3DkeMQzbE+oj1Lm+SjqQn8wgaP3p+P6FcVts3d+nl8CtWvLZCF9viJnMDnBr8Tyh0NKEnDJLmemMcrw7alKaZ1BX8FVDvPw3jtHX5lgoG2ebaeIS6obe+YXzrdVVxrk4DyU/Fl5sz+9qkGCMOnkFTq/8KzV7SdXwZwDd4b2kyUosY/ORE1hdCYhdbFqmGyd8NtfADQx8eIrw8BFFdXNlOYUqCZ3Tvz4r4oLtCOZ3o0k1COb/nsCwI7po3PBfONq1wo8IVeMM6e9+VRbwtiR8d8/MELtTIq38uaRBIp8R7J9ClG3iq9BxjnujWH3UBsPFWdVZAhNYJuVurfmvxWPf1dPCGIDwJw2kHmpox4TLcfSBJMkF3Mz+f9fSr72rXBF5VNqqtCtr9mhMshHkLcuCIJmM3uPbSs1hBQ5ngyXzN/Mr4fVNMuY/A5+gAt3sUK3sDOrtK4VlgUNntz0r6UFU3E6rdgG5vZ42NMOeHarHRFCWT2kEkW9WIISpeuLifoR5rs0CECQg3jVhIqzHPRcxcz9ib29t9IzTj1wXW6wC4NOdemy646sTLLQxsnyt9/BfjIJ8Zl/6iNTXf1drt2Fda+LyfmrsLt+zbgj/EM6nBY/GELoXWrsZaXMQ4vQMhSoQApLjfz7/+BOM/PQl7PSUQnsMmJxRg//5/YkgxZ9w4FavGgv7g5Taj2vlnU5ehSY42XwtSmfq+E7KsCtonsFR6GUQfv38UHxRmb6iR7DeK8x9CqQfWElNumll+JrxJ3ovKlWN3l9cNTUAN4ANJfK8xbMaBHUmeubmCJoEzD58jTd22njMGTzcJ9TBaWLID5gDYs2ZX4RAJknYcvpn1LKfyrRwAIfeF86azq5ekHYFRFn8/3apR3a7zWGc+0hON7K0fYgJEok=`)
	if err != nil {
		t.Fatalf("Error decoding ciphertext: %v", err)
	}

	key := []byte("12345678901234567890123456789012")
	nonce := []byte("255e59e73ed3")

	decrypted, err := NewEncrypter(key).gcmDecrypt(ciphertext, nonce)
	if err != nil {
		t.Fatalf("Error decrypting value: %v", err)
	}

	result := string(decrypted)
	expectsStr := `Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`

	if result != expectsStr {
		t.Errorf("Expected %s, got %s", expectsStr, result)
	}
}
