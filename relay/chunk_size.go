package relay

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding"
	"reflect"
	"unsafe"
)

// computeInMemoryFrameSize computes the size of a blob's chunks in memory.
func computeInMemoryFrameSize(frames []*encoding.Frame) (uint64, error) {

	if len(frames) == 0 {
		return 0, fmt.Errorf("no frames provided")
	}

	firstFrame := frames[0]
	firstFrameSize, err := SizeOf(firstFrame)
	if err != nil {
		return 0, fmt.Errorf("error calculating size of first frame: %w", err)
	}

	// all frames for a particular blob are the same size
	size := firstFrameSize * uint64(len(frames))

	return size, nil
}

// SizeOf calculates the size of an object in memory using reflection. Includes the memory
// referenced by the object. This function assumes that there are no circular references
// in the object graph. If there are, then this function will enter an infinite loop
// (likely ending with a stack overflow).
//
// This has non-trivial performance implications and should be used carefully.
func SizeOf(object any) (uint64, error) {
	return recursiveSizeOf(object, true)
}

// recursiveSizeOf recursively walks through a data structure and calculates the memory it uses.
//
// If the indirect flag is true, then the provided object was referenced by its
// parent (e.g. it was a pointer or in a map). If the indirect flag is false,
// then the provided object was directly embedded in its parent.
func recursiveSizeOf(object any, indirect bool) (uint64, error) {

	size := uint64(0)
	if object == nil {
		return size, nil
	}

	///////////////////////////////////////////////////////////////////////////////////
	//                   Determine the size of this object.                          //
	///////////////////////////////////////////////////////////////////////////////////

	// If indirect is false, then this object's size will have been counted by its parent.
	// If indirect is true, then we need to count this object's size.
	if indirect {
		// SizeOf is actually safe, it's just that the creators of golang decided that
		// software engineers weren't smart enough to use it properly.
		// Well, ok. That's a little insulting. I'm going to use it anyway.
		size = uint64(unsafe.Sizeof(object))
	}

	///////////////////////////////////////////////////////////////////////////////////
	//          Determine the size of the memory referenced by this object.          //
	///////////////////////////////////////////////////////////////////////////////////

	val := reflect.ValueOf(object)
	objectType := val.Type().Kind()

	switch objectType {
	case reflect.Pointer:
		// Although the bytes for the pointer itself will have been counted,
		// the thing being pointed to will not have been counted.

		referencedObject := val.Elem().Interface()
		referencedSize, err := recursiveSizeOf(referencedObject, true)
		if err != nil {
			return 0, fmt.Errorf("error calculating size of referenced object: %w", err)
		}
		size += referencedSize

	case reflect.Struct:
		// iterate over the fields in the struct

		fieldCount := val.NumField()
		for index := 0; index < fieldCount; index++ {
			field := val.Field(index)
			fieldSize, err := recursiveSizeOf(field.Interface(), false)
			if err != nil {
				return 0, fmt.Errorf("error calculating size of field: %w", err)
			}

			size += fieldSize
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		// The slice/array header will have been counted, but the memory it references will not have been.
		// This is a little tricky because slices are pointers to arrays, so we need to get the size of the array.
		length := val.Len()
		for i := 0; i < length; i++ {
			fieldSize, err := recursiveSizeOf(val.Index(i).Interface(), true)
			if err != nil {
				return 0, fmt.Errorf("error calculating size of field: %w", err)
			}
			size += fieldSize
		}
	case reflect.String:
		// The string header will have been counted, but not the data contained in the string.
		size += uint64(len(val.String()))
	case reflect.Map:
		// The map header will have been counted, but not the map's keys and values.
		keys := val.MapKeys()
		for _, key := range keys {
			keySize, err := recursiveSizeOf(key.Interface(), true)
			if err != nil {
				return 0, fmt.Errorf("error calculating size of map key: %w", err)
			}
			size += keySize

			valueSize, err := recursiveSizeOf(val.MapIndex(key).Interface(), true)
			if err != nil {
				return 0, fmt.Errorf("error calculating size of map value: %w", err)
			}
			size += valueSize
		}
	case reflect.Bool:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		// There is no memory referenced by these types.
	default:
		// This utility was created to calculate the size of simple object types, not as a general purpose
		// memory calculator. If you're seeing this error, then you're trying to calculate the size of
		// an object with some fancy type in it that I didn't bother with because I didn't need it.
		// Take your unsafe pointers, functions, and other hoo haa and go calculate the size yourself, thank you.
		return 0, fmt.Errorf("unsupported object type: %v", objectType)
	}

	return size, nil
}
