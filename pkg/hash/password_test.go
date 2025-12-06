package hash

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHash_Success(t *testing.T) {
	data := "password123"

	hashedData, err := Hash(data)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedData)
	assert.NotEqual(t, data, hashedData)
}

func TestHash_EmptyString(t *testing.T) {
	data := ""

	hashedData, err := Hash(data)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedData)
}

func TestHash_LongPassword(t *testing.T) {
	// bcrypt has a max length of 72 bytes, so this should return an error
	data := strings.Repeat("a", 100)

	hashedData, err := Hash(data)

	assert.Error(t, err)
	assert.Empty(t, hashedData)
	assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
}

func TestHash_DifferentPasswordsDifferentHashes(t *testing.T) {
	password1 := "password123"
	password2 := "password456"

	hash1, err1 := Hash(password1)
	hash2, err2 := Hash(password2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2)
}

func TestHash_SamePasswordDifferentHashes(t *testing.T) {
	// bcrypt generates different hashes for the same password due to salt
	password := "password123"

	hash1, err1 := Hash(password)
	hash2, err2 := Hash(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, hash2)
	assert.NotEqual(t, hash1, hash2)
}

func TestHash_SpecialCharacters(t *testing.T) {
	data := "p@ssw0rd!#$%^&*()"

	hashedData, err := Hash(data)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedData)
}

func TestHash_UnicodeCharacters(t *testing.T) {
	data := "пароль密码🔑"

	hashedData, err := Hash(data)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedData)
}

func TestHash_UsesDefaultCost(t *testing.T) {
	data := "password123"

	hashedData, err := Hash(data)

	assert.NoError(t, err)

	// Verify the cost is bcrypt.DefaultCost
	cost, err := bcrypt.Cost([]byte(hashedData))
	assert.NoError(t, err)
	assert.Equal(t, bcrypt.DefaultCost, cost)
}

func TestCheckHash_Success(t *testing.T) {
	password := "password123"

	hashedPassword, err := Hash(password)
	assert.NoError(t, err)

	err = CheckHash(password, hashedPassword)

	assert.NoError(t, err)
}

func TestCheckHash_WrongPassword(t *testing.T) {
	correctPassword := "password123"
	wrongPassword := "wrongpassword"

	hashedPassword, err := Hash(correctPassword)
	assert.NoError(t, err)

	err = CheckHash(wrongPassword, hashedPassword)

	assert.Error(t, err)
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
}

func TestCheckHash_EmptyPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := Hash(password)
	assert.NoError(t, err)

	err = CheckHash("", hashedPassword)

	assert.Error(t, err)
}

func TestCheckHash_EmptyHash(t *testing.T) {
	password := "password123"

	err := CheckHash(password, "")

	assert.Error(t, err)
}

func TestCheckHash_InvalidHash(t *testing.T) {
	password := "password123"
	invalidHash := "not-a-valid-bcrypt-hash"

	err := CheckHash(password, invalidHash)

	assert.Error(t, err)
}

func TestCheckHash_CaseSensitive(t *testing.T) {
	password := "Password123"

	hashedPassword, err := Hash(password)
	assert.NoError(t, err)

	// Should fail with lowercase
	err = CheckHash("password123", hashedPassword)
	assert.Error(t, err)

	// Should succeed with correct case
	err = CheckHash("Password123", hashedPassword)
	assert.NoError(t, err)
}

func TestCheckHash_WhitespaceMatters(t *testing.T) {
	password := "password123"

	hashedPassword, err := Hash(password)
	assert.NoError(t, err)

	// Should fail with extra whitespace
	err = CheckHash("password123 ", hashedPassword)
	assert.Error(t, err)

	err = CheckHash(" password123", hashedPassword)
	assert.Error(t, err)

	// Should succeed with exact match
	err = CheckHash("password123", hashedPassword)
	assert.NoError(t, err)
}

func TestCheckHash_SpecialCharacters(t *testing.T) {
	password := "p@ssw0rd!#$%"

	hashedPassword, err := Hash(password)
	assert.NoError(t, err)

	err = CheckHash(password, hashedPassword)

	assert.NoError(t, err)
}

func TestCheckHash_UnicodeCharacters(t *testing.T) {
	password := "пароль密码🔑"

	hashedPassword, err := Hash(password)
	assert.NoError(t, err)

	err = CheckHash(password, hashedPassword)

	assert.NoError(t, err)
}

func TestCheckHash_LongPassword(t *testing.T) {
	// Use a password within bcrypt's 72-byte limit
	password := strings.Repeat("a", 70)

	hashedPassword, err := Hash(password)
	assert.NoError(t, err)

	err = CheckHash(password, hashedPassword)

	assert.NoError(t, err)
}

func TestHashAndCheckHash_Integration(t *testing.T) {
	testCases := []string{
		"simple",
		"with spaces",
		"with!special@characters#",
		"1234567890",
		"MixedCase123",
		"very_long_" + strings.Repeat("pass", 10), // Within 72-byte limit
		"短い",
		"émojis🎉🔐",
	}

	for _, password := range testCases {
		t.Run(password, func(t *testing.T) {
			hashedPassword, err := Hash(password)
			assert.NoError(t, err)
			assert.NotEmpty(t, hashedPassword)

			err = CheckHash(password, hashedPassword)
			assert.NoError(t, err)

			// Verify wrong password fails
			err = CheckHash(password+"wrong", hashedPassword)
			assert.Error(t, err)
		})
	}
}

func TestHash_MultipleCallsForSamePassword(t *testing.T) {
	password := "testpassword"
	iterations := 5
	hashes := make([]string, iterations)

	for i := 0; i < iterations; i++ {
		hash, err := Hash(password)
		assert.NoError(t, err)
		hashes[i] = hash
	}

	// All hashes should be different
	for i := 0; i < iterations; i++ {
		for j := i + 1; j < iterations; j++ {
			assert.NotEqual(t, hashes[i], hashes[j])
		}
	}

	// All hashes should validate correctly
	for _, hash := range hashes {
		err := CheckHash(password, hash)
		assert.NoError(t, err)
	}
}

func TestCheckHash_DifferentHashFormats(t *testing.T) {
	password := "password123"

	// Test with manually created bcrypt hash
	manualHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = CheckHash(password, string(manualHash))
	assert.NoError(t, err)

	// Test with our Hash function
	ourHash, err := Hash(password)
	assert.NoError(t, err)

	err = CheckHash(password, ourHash)
	assert.NoError(t, err)
}