package libs

import (
	"strings"
	"testing"
)

func TestEncodeFile(t *testing.T) {
	password := "testpassword"
	longText := "This is a test message."

	encryptedData, err := EncodeFile(password, longText)
	if err != nil {
		t.Fatalf("EncodeFile failed: %v", err)
	}

	if !strings.HasPrefix(encryptedData, "@") {
		t.Errorf("Encoded data does not start with '@': %s", encryptedData)
	}

	parts := strings.SplitN(encryptedData[1:], "@", 2)
	if len(parts) != 2 {
		t.Errorf("Encoded data format is incorrect: %s", encryptedData)
	}
	if !IsEncryptedFile(encryptedData) {
		t.Errorf("IsEncryptedFile returned false for valid encrypted data: %s", encryptedData)
	}
}

func TestDecodeFile(t *testing.T) {
	password := "96e79218965eb72c92a549dd5a330112"
	longText := "This is a test message."

	encryptedData, err := EncodeFile(password, longText)
	if err != nil {
		t.Fatalf("EncodeFile failed: %v", err)
	}

	decryptedText, err := DecodeFile(password, encryptedData)
	if err != nil {
		t.Fatalf("DecodeFile failed: %v", err)
	}

	if decryptedText != longText {
		t.Errorf("Decrypted text does not match original: expected '%s', got '%s'", longText, decryptedText)
	}
}

func TestIsEncryptedFile(t *testing.T) {
	password := "testpassword"
	longText := "This is a test message."

	encryptedData, err := EncodeFile(password, longText)
	if err != nil {
		t.Fatalf("EncodeFile failed: %v", err)
	}

	if !IsEncryptedFile(encryptedData) {
		t.Errorf("IsEncryptedFile returned false for valid encrypted data: %s", encryptedData)
	}

	invalidData := "Invalid@Data"
	if IsEncryptedFile(invalidData) {
		t.Errorf("IsEncryptedFile returned true for invalid encrypted data: %s", invalidData)
	}
}
