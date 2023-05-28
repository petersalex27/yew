package demangler

import (
	"strconv"
	"sync"
)

var itoa = strconv.Itoa

type demanglePrefix string
const (
	typePrefix demanglePrefix = ".ty"
	classPrefix demanglePrefix = ".cls"
	functionPrefix demanglePrefix = ".fn"
	curryPrefix demanglePrefix = ".cur"
	labelPrefix demanglePrefix = ".l"
)

// Types of things that can be demangled
type DemangleType int
const (
	TYPE DemangleType = iota
	CLASS
	FUNCTION
	CURRY
	LABEL
)

// holds counters and their respective locks
var demangler = struct{
	typeCounter int
	typeLock sync.Mutex

	classCounter int
	classLock sync.Mutex
	
	functionCounter int
	functionLock sync.Mutex

	curryCounter int
	curryLock sync.Mutex

	labelCounter int
	labelLock sync.Mutex
} {
	typeCounter: 0,
	classCounter: 0,
	functionCounter: 0,
	curryCounter: 0,
	labelCounter: 0,
}

// demangle type -> prefix string
var prefixMap = map[DemangleType]string {
	TYPE: string(typePrefix),
	CLASS: string(classPrefix),
	FUNCTION: string(functionPrefix),
	CURRY: string(curryPrefix),
	LABEL: string(labelPrefix),
}

// Get `count` demangler prefixes for a given demangle type when the demangle type is valid,
// else get an empty string
func (d DemangleType) GetDemanglerPrefixes(count int) (demanglerPrefixes []string) {
	prefix := prefixMap[d]
	demanglerPrefixes = make([]string, count)
	switch d {
	case TYPE:
		demangler.typeLock.Lock()
		for i := range demanglerPrefixes {
			demanglerPrefixes[i] = prefix + itoa(demangler.typeCounter)
			demangler.typeCounter++
		}
		demangler.typeLock.Unlock()
		return
	case CLASS:
		demangler.classLock.Lock()
		for i := range demanglerPrefixes {
			demanglerPrefixes[i] = prefix + itoa(demangler.classCounter) 
			demangler.classCounter++
		}
		demangler.classLock.Unlock()
		return
	case FUNCTION:
		demangler.functionLock.Lock()
		for i := range demanglerPrefixes {
			demanglerPrefixes[i] = prefix + itoa(demangler.functionCounter)
			demangler.functionCounter++
		}
		demangler.functionLock.Unlock()
		return
	case CURRY:
		demangler.curryLock.Lock()
		for i := range demanglerPrefixes {
			demanglerPrefixes[i] = prefix + itoa(demangler.curryCounter)
			demangler.curryCounter++
		}
		demangler.curryLock.Unlock()
		return
	case LABEL:
		demangler.labelLock.Lock()
		for i := range demanglerPrefixes {
			demanglerPrefixes[i] = prefix + itoa(demangler.labelCounter)
			demangler.labelCounter++
		}
		demangler.labelLock.Unlock()
		return
	default:
		for i := range demanglerPrefixes {
			demanglerPrefixes[i] = ""
		}
		return // no prefixes
	}
}

// Get a single demangler prefix for a given demangle type when the type is valid, 
// else get an empty string
func (d DemangleType) GetDemanglerPrefix() (demanglerPrefix string) {
	return d.GetDemanglerPrefixes(1)[0]
}