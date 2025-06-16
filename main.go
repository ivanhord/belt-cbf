package main

func main() {
	err := DecryptFile("ciphertext.txt", "decrypted_output.txt")
	if err != nil {
		panic(err)
	}

	err = VerifyMessagesFromFile("decrypted_output.txt")
	if err != nil {
		panic(err)
	}
}
