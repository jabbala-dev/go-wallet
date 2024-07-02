# A Blockchain Wallet Using Go
Creating a basic blockchain wallet in Go involves generating pair of cryptographic keys (Public and Private Keys), stroing keys, and providing functionality to sign transactions. Below is a simplified example.

## 1. Install Dependiences
First, install some dependencies for cryptographic functions. We can use the `go-ethereum` package for this.

```sh
go get github.com/ethereum/go-ethereum
```

## 2. Basic Structure
Here's is a basic structure for a blockchain wallet: <br/>
`main.go`
```go
package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	fmt.Printf("Address: %s\n", address)
	saveKey(privateKey)
}

func saveKey(privateKey *ecdsa.PrivateKey) {
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	file, err := os.Create("private_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
}

func loadKey() (*ecdsa.PrivateKey, error) {
	privateKeyHex, err := os.ReadFile("private_key.txt")
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := hex.DecodeString(string(privateKeyHex))
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func signMessage(message string, privateKey *ecdsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature), nil
}

func verifyMessage(message, signatureHex string, publicKey *ecdsa.PublicKey) bool {
	hash := sha256.Sum256([]byte(message))
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}

	return crypto.VerifySignature(crypto.FromECDSAPub(publicKey), hash[:], signature[:len(signature)-1])
}

```

## 3. Explanation
1. **Generate Key Pair**: Generates a new ECDSA key pair using crypto.GenerateKey.
2. **Address Generation**: Uses the public key to generate an Ethereum-style address with crypto.PubkeyToAddress.
3. **Key Storage**: Saves the private key to a file (private_key.txt).
4. **Load Key**: Reads the private key from the file and converts it back to an ECDSA key.
5. **Sign Message**: Signs a message using the private key and SHA-256 hashing.
6. **Verify Message**: Verifies the message signature using the public key.

## 4. Running the Code
1. Generating Keys
```sh
go run main.go
```
This will generate a key pair and save the private key to `private_key.txt`.

2. Load and Use the Key
Modify `main.go` to load the key and sign a message
```go
func main() {
    privateKey, err := loadKey()
    if err != nil {
        log.Fatal(err)
    }

    message := "Hello, Blockchain!"
    signature, err := signMessage(message, privateKey)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Message: %s\n", message)
    fmt.Printf("Signature: %s\n", signature)

    publicKey := privateKey.Public().(*ecdsa.PublicKey)
    isValid := verifyMessage(message, signature, publicKey)
    fmt.Printf("Signature valid: %t\n", isValid)
}
```

## Enhanced Version of it
The project is enhanced to make it a complete application capable of integrating with other platforms
1. **Add RESTful APIs**: Allow external platforms to interact with the wallet
2. **Create and manage transactions**: Include functionality to create, sign, and broadcast transactions

We'll use `gin-gonic` framework for the RESTful API and `go-ethereum` for blockchain interactions. 

### install dependencies
```sh
go get github.com/gin-gonic/gin
go get github.com/ethereum/go-ethereum

```

### Project Structure
```go
blockchain-wallet/
│
├── main.go
├── handlers/
│   └── handlers.go
├── models/
│   └── wallet.go
├── services/
│   └── wallet.go
├── utils/
│   └── crypto.go
├── go.mod
└── go.sum
```

### 1. `main.go`
```go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jabbala-dev/go-wallet/handlers"
)

func main() {
	r := gin.Default()

	// Define routes
	r.GET("/generate", handlers.GenerateKeyPair)
	r.GET("/address", handlers.GetAddress)
	r.POST("/sign", handlers.SignMessage)
	r.POST("/verify", handlers.VerifyMessage)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
```

### 2. `handlers/handlers.go`
```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jabbala-dev/go-wallet/services"
)

func GenerateKeyPair(c *gin.Context) {
	privateKey, address, err := services.GenerateKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"private_key": privateKey, "address": address})
}

func GetAddress(c *gin.Context) {
	address, err := services.GetAddress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

func SignMessage(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	signature, err := services.SignMessage(request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"signature": signature})
}

func VerifyMessage(c *gin.Context) {
	var request struct {
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	isValid, err := services.VerifyMessage(request.Message, request.Signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": isValid})
}
```
### 3. `services/wallet.go`
```go
package services

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jabbala-dev/go-wallet/utils"
)

var privateKeyFile = "private_key.txt"

func GenerateKeyPair() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	err = ioutil.WriteFile(privateKeyFile, []byte(privateKeyHex), 0600)
	if err != nil {
		return "", "", err
	}

	return privateKeyHex, address, nil
}

func GetAddress() (string, error) {
	privateKey, err := loadKey()
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	return address, nil
}

func SignMessage(message string) (string, error) {
	privateKey, err := loadKey()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(message))
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature), nil
}

func VerifyMessage(message, signatureHex string) (bool, error) {
	privateKey, err := loadKey()
	if err != nil {
		return false, err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	hash := sha256.Sum256([]byte(message))
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}

	return crypto.VerifySignature(crypto.FromECDSAPub(publicKey), hash[:], signature[:len(signature)-1]), nil
}

func loadKey() (*ecdsa.PrivateKey, error) {
	if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
		return nil, errors.New("private key file does not exist")
	}

	privateKeyHex, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := hex.DecodeString(string(privateKeyHex))
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
```
### 4. `utils/crypto.go`
```go
package utils

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func PublicKeyToAddress(publicKey *ecdsa.PublicKey) string {
	return crypto.PubkeyToAddress(*publicKey).Hex()
}

func PrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(crypto.FromECDSA(privateKey))
}

func HexToPrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	return crypto.ToECDSA(privateKeyBytes)
}

func SignMessage(privateKey *ecdsa.PrivateKey, message string) (string, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

func VerifyMessage(publicKey *ecdsa.PublicKey, message, signatureHex string) (bool, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}
	return crypto.VerifySignature(crypto.FromECDSAPub(publicKey), hash[:], signature[:len(signature)-1]), nil
}
```
### Running the Application
#### 1. Initialize project
```sh
go mod init github.com/yourusername/blockchain-wallet
go mod tidy
```

#### 2. Run the Application
```sh
go run main.go
```

### Endpoint Tests
#### 1. generate keys
```sh
 curl http://localhost:8080/generate
```

#### 2. Address
```sh
http://localhost:8080/address
```

#### 3. sign
```sh
curl -X POST -H "Content-Type: application/json" -d '{"message":"Hello, Go Wallet"}' http://localhost:8080/sign
```

#### 4. Verify the transaction
```sh
curl -X POST -H "Content-Type: application/json" -d '{"message":"Hello, Go Wallet", "signature":"c71274f99339471fdfa73a4ea82d842982bfedf07f5c5b1df6108e9878494c021d1de0955fcc8869e770d83da2d6cc4ae6c5dddd30fa6fb4d313bd28a088650d00"}' http://localhost:8080/verify
```

#### 5. Create and Send Transaction
```sh
curl -X POST http://localhost:8080/transaction -H "Content-Type: application/json" -d '{"to_address": "0xRecipientAddress", "value": 1000000000000000000}'
```


## Conclusion
This project demonstrates the creation of a basic blockchain wallet application in Go, providing key functionalities such as generating key pairs, retrieving wallet addresses, signing messages, and verifying signatures. By structuring the project into well-defined modules and utilizing the gin-gonic framework for RESTful APIs, we have laid a strong foundation for further development and integration with other platforms.

### Areas for Improvement
#### 1. Transaction Creation and Broadcasting:
* Enhancement: Implement functionality for creating, signing, and broadcasting transactions to the blockchain network.
* Benefit: Enables the wallet to perform actual transfers of cryptocurrency, making it fully functional.

#### 2. Blockchain Interaction:
* Enhancement: Integrate with blockchain nodes (e.g., Ethereum, Bitcoin) to fetch and display wallet balances and transaction history.
* Benefit: Provides real-time information and enhances user interaction with the blockchain.

#### 3. Security Enhancements:
* Enhancement: Implement secure storage solutions for private keys (e.g., hardware wallets, encrypted storage).
* Benefit: Increases the security of private keys, protecting users from potential threats and attacks.

#### 4. User Authentication and Authorization:
* Enhancement: Add user authentication and authorization mechanisms to secure API endpoints.
* Benefit: Ensures that only authorized users can access and manage their wallets, enhancing security.

#### 5. Improved Error Handling and Logging:
* Enhancement: Implement comprehensive error handling and logging throughout the application.
* Benefit: Helps in debugging, monitoring, and maintaining the application more effectively.

#### 6. Unit and Integration Testing:
* Enhancement: Develop unit and integration tests to ensure the reliability and correctness of the application.
* Benefit: Improves code quality and reduces the likelihood of bugs and regressions.

#### 7. Web and Mobile Interfaces:
* Enhancement: Develop user-friendly web and mobile interfaces for interacting with the wallet.
* Benefit: Enhances accessibility and user experience, making the wallet more appealing to a broader audience.

#### 8.Multi-currency Support:
* Enhancement: Extend the wallet to support multiple cryptocurrencies.
* Benefit: Increases the utility of the wallet by allowing users to manage various assets within a single application.

By addressing these areas for improvement, we can transform this basic blockchain wallet into a robust, secure, and user-friendly application capable of competing with existing solutions in the market.


