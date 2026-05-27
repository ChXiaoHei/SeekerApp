//go:build windows

package everything

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

// Everything SDK request flags
const (
	requestFileName     = 0x00000001
	requestPath         = 0x00000002
	requestFullPathName = 0x00000004
	requestSize         = 0x00000010
	requestDateModified = 0x00000040
)

// Everything SDK sort constants
const (
	sortNameAsc = 1
)

// Everything SDK error codes
const (
	errorOK                = 0
	errorCreateIPC         = 1
	errorCreateWindow      = 2
	errorRegisterWindow    = 3
	errorIPCNotRunning     = 4  // Everything not running
	errorMemory            = 5
	errorQuery             = 6
	errorInvalidCall       = 7
	errorInvalidIndex      = 8
	errorInvalidRequest    = 9
	errorIPCTimeout        = 10
)

// SDK wraps the Everything64.dll for file searching.
type SDK struct {
	dll                          *syscall.LazyDLL
	setSearchW                   *syscall.LazyProc
	setRequestFlags              *syscall.LazyProc
	setSort                      *syscall.LazyProc
	setMax                       *syscall.LazyProc
	setMatchCase                 *syscall.LazyProc
	setMatchWholeWord            *syscall.LazyProc
	setMatchPath                 *syscall.LazyProc
	setRegex                     *syscall.LazyProc
	queryW                       *syscall.LazyProc
	getNumResults                *syscall.LazyProc
	getResultFullPathNameW       *syscall.LazyProc
	getResultFileNameW           *syscall.LazyProc
	getResultSize                *syscall.LazyProc
	getResultDateModified        *syscall.LazyProc
	getLastError                 *syscall.LazyProc
	isResultFolder               *syscall.LazyProc
	cleanUp                      *syscall.LazyProc
}

// NewSDK loads the Everything SDK DLL from the given path.
// If dllPath is empty, it defaults to "Everything64.dll" (must be in PATH or working dir).
func NewSDK(dllPath string) (*SDK, error) {
	if dllPath == "" {
		dllPath = "Everything64.dll"
	}

	dll := syscall.NewLazyDLL(dllPath)
	if err := dll.Load(); err != nil {
		return nil, fmt.Errorf("failed to load Everything SDK DLL from %q: %w", dllPath, err)
	}

	sdk := &SDK{
		dll:                          dll,
		setSearchW:                   dll.NewProc("Everything_SetSearchW"),
		setRequestFlags:              dll.NewProc("Everything_SetRequestFlags"),
		setSort:                      dll.NewProc("Everything_SetSort"),
		setMax:                       dll.NewProc("Everything_SetMax"),
		setMatchCase:                 dll.NewProc("Everything_SetMatchCase"),
		setMatchWholeWord:            dll.NewProc("Everything_SetMatchWholeWord"),
		setMatchPath:                 dll.NewProc("Everything_SetMatchPath"),
		setRegex:                     dll.NewProc("Everything_SetRegex"),
		queryW:                       dll.NewProc("Everything_QueryW"),
		getNumResults:                dll.NewProc("Everything_GetNumResults"),
		getResultFullPathNameW:       dll.NewProc("Everything_GetResultFullPathNameW"),
		getResultFileNameW:           dll.NewProc("Everything_GetResultFileNameW"),
		getResultSize:                dll.NewProc("Everything_GetResultSize"),
		getResultDateModified:        dll.NewProc("Everything_GetResultDateModified"),
		getLastError:                 dll.NewProc("Everything_GetLastError"),
		isResultFolder:               dll.NewProc("Everything_IsFolderResult"),
		cleanUp:                      dll.NewProc("Everything_CleanUp"),
	}

	return sdk, nil
}

// Search performs a file search using Everything and returns matching results.
func (s *SDK) Search(query string, opts SearchOptions) ([]SearchResult, error) {
	if opts.MaxResults <= 0 {
		opts.MaxResults = 100
	}

	queryPtr, err := syscall.UTF16PtrFromString(query)
	if err != nil {
		return nil, fmt.Errorf("invalid query string: %w", err)
	}

	s.setSearchW.Call(uintptr(unsafe.Pointer(queryPtr)))
	s.setRequestFlags.Call(uintptr(requestFileName | requestPath | requestFullPathName | requestSize | requestDateModified))
	s.setSort.Call(uintptr(sortNameAsc))
	s.setMax.Call(uintptr(opts.MaxResults))
	s.setMatchCase.Call(boolToUintptr(opts.MatchCase))
	s.setMatchWholeWord.Call(boolToUintptr(opts.MatchWholeWord))
	s.setMatchPath.Call(boolToUintptr(opts.MatchPath))
	s.setRegex.Call(boolToUintptr(opts.UseRegex))

	ret, _, _ := s.queryW.Call(uintptr(1)) // TRUE = wait for results
	if ret == 0 {
		errCode, _, _ := s.getLastError.Call()
		return nil, newSDKError(int(errCode))
	}

	numResults, _, _ := s.getNumResults.Call()
	count := int(numResults)

	results := make([]SearchResult, 0, count)
	pathBuf := make([]uint16, 4096)

	for i := 0; i < count; i++ {
		result := SearchResult{}

		// Get full path
		s.getResultFullPathNameW.Call(
			uintptr(i),
			uintptr(unsafe.Pointer(&pathBuf[0])),
			uintptr(len(pathBuf)),
		)
		result.FullPath = syscall.UTF16ToString(pathBuf)

		// Get file name
		ret, _, _ := s.getResultFileNameW.Call(uintptr(i))
		if ret != 0 {
			result.FileName = utf16PtrToString(ret)
		}

		// Check if folder
		isFolder, _, _ := s.isResultFolder.Call(uintptr(i))
		result.IsFolder = isFolder != 0

		// Get file size
		var fileSize int64
		s.getResultSize.Call(uintptr(i), uintptr(unsafe.Pointer(&fileSize)))
		result.Size = fileSize

		// Get date modified
		var ft syscall.Filetime
		s.getResultDateModified.Call(uintptr(i), uintptr(unsafe.Pointer(&ft)))
		result.DateModified = filetimeToTime(ft)

		results = append(results, result)
	}

	return results, nil
}

// CleanUp releases Everything SDK resources.
func (s *SDK) CleanUp() {
	if s.cleanUp != nil {
		s.cleanUp.Call()
	}
}

func boolToUintptr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func utf16PtrToString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	// Walk the memory until null terminator
	var utf16Chars []uint16
	for i := uintptr(0); ; i += 2 {
		char := *(*uint16)(unsafe.Pointer(ptr + i))
		if char == 0 {
			break
		}
		utf16Chars = append(utf16Chars, char)
	}
	return syscall.UTF16ToString(utf16Chars)
}

func filetimeToTime(ft syscall.Filetime) time.Time {
	// Windows FILETIME: 100-nanosecond intervals since January 1, 1601
	nsec := int64(ft.HighDateTime)<<32 + int64(ft.LowDateTime)
	// Difference between 1601 and 1970 in 100-nanosecond intervals
	const epochDiff = 116444736000000000
	if nsec <= epochDiff {
		return time.Time{}
	}
	return time.Unix(0, (nsec-epochDiff)*100)
}

// newSDKError creates an error from an Everything SDK error code.
func newSDKError(code int) error {
	var msg string
	switch code {
	case errorOK:
		return nil
	case errorCreateIPC:
		msg = "failed to create IPC"
	case errorCreateWindow:
		msg = "failed to create window"
	case errorRegisterWindow:
		msg = "failed to register window"
	case errorIPCNotRunning:
		msg = "Everything is not running. Please start Everything first."
	case errorMemory:
		msg = "memory allocation failed"
	case errorQuery:
		msg = "query failed"
	case errorInvalidCall:
		msg = "invalid call"
	case errorInvalidIndex:
		msg = "invalid result index"
	case errorInvalidRequest:
		msg = "invalid request flags"
	case errorIPCTimeout:
		msg = "IPC timeout"
	default:
		msg = fmt.Sprintf("unknown error code: %d", code)
	}
	return fmt.Errorf("Everything SDK error: %s", msg)
}