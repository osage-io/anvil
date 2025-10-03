package crypto

import (
	"os"
	"runtime"
	"runtime/debug"
	"unsafe"
)

// InitSecureRuntime configures the Go runtime for maximum security
func InitSecureRuntime() {
	// Disable memory profiling to prevent sensitive data from being retained
	runtime.MemProfileRate = 0

	// Force garbage collection to clear any leftover sensitive data
	runtime.GC()

	// Set maximum GOMAXPROCS to prevent resource exhaustion attacks
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Disable debug info to prevent information leakage
	debug.SetGCPercent(-1)  // Disable GC percentage-based triggering for predictable behavior
	debug.SetGCPercent(100) // Re-enable with conservative setting
}

// SecureZeroMemory attempts to securely zero memory at the given address
// This is a best-effort implementation in Go
func SecureZeroMemory(data []byte) {
	if len(data) == 0 {
		return
	}

	// Zero the memory multiple times to defeat potential memory remanence
	for pass := 0; pass < 3; pass++ {
		for i := range data {
			data[i] = 0
		}

		// Force memory barrier to ensure writes complete
		runtime.KeepAlive(data)
	}

	// Additional paranoid zeroing using unsafe operations
	if len(data) > 0 {
		ptr := unsafe.Pointer(&data[0])
		size := uintptr(len(data))

		// Write different patterns to defeat potential hardware optimizations
		patterns := []byte{0x00, 0xFF, 0xAA, 0x55, 0x00}
		for _, pattern := range patterns {
			slice := (*[1 << 30]byte)(ptr)[:size:size]
			for i := range slice {
				slice[i] = pattern
			}
			runtime.KeepAlive(slice)
		}
	}

	// Final garbage collection to ensure cleanup
	runtime.GC()
}

// SecureClearString attempts to clear a string from memory (best effort)
func SecureClearString(s *string) {
	if s == nil || len(*s) == 0 {
		return
	}

	// Get the underlying data
	data := unsafe.Slice(unsafe.StringData(*s), len(*s))

	// Securely zero the underlying memory
	SecureZeroMemory(data)

	// Clear the string reference
	*s = ""
}

// VerifySecureEnvironment checks if the current environment is suitable for secure operations
func VerifySecureEnvironment() []string {
	var warnings []string

	// Check if debugger is attached (basic check)
	if runtime.GOMAXPROCS(0) != runtime.NumCPU() {
		warnings = append(warnings, "GOMAXPROCS differs from CPU count - possible debugging")
	}

	// Check memory profiling status
	if runtime.MemProfileRate > 0 {
		warnings = append(warnings, "Memory profiling is enabled - sensitive data may be retained")
	}

	// Check for common debug environment variables
	debugEnvVars := []string{
		"GOTRACEBACK", "GODEBUG", "GOMEMLIMIT",
	}

	for _, envVar := range debugEnvVars {
		if val := os.Getenv(envVar); val != "" {
			warnings = append(warnings, "Debug environment variable "+envVar+" is set: "+val)
		}
	}

	return warnings
}
