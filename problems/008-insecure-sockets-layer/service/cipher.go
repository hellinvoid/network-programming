package service

const (
	END          = 0x00
	REVERSE_BITS = 0x01
	XOR_N        = 0x02
	XOR_POS      = 0x03
	ADD_N        = 0x04
	ADD_POS      = 0x05
)

func reverseBits(in byte) byte {
	var out byte
	for i := 0; i < 8; i++ {
		out <<= 1
		out |= in & 1
		in >>= 1
	}
	return out
}

func xor(in byte, pos byte) byte {
	return in ^ pos
}

func add(in byte, pos byte) byte {
	var a, b int = int(in), int(pos)
	c := (a + b) % 256
	return byte(c)
}

func sub(in byte, pos byte) byte {
	var a, b int = int(in), int(pos)
	c := (a - b + 256) % 256
	return byte(c)
}

func Encode(original byte, cipher []byte, idx int, pos int) byte {
	if idx >= len(cipher) {
		return original
	}
	encoded := original

	next := idx + 1
	// jump 2 places if any of below codes
	if cipher[idx] == XOR_N || cipher[idx] == ADD_N {
		next++
	}

	// Encode sequentially and pass on
	switch cipher[idx] {
	case REVERSE_BITS:
		encoded = reverseBits(encoded)
	case XOR_N:
		n := cipher[idx+1]
		encoded = xor(encoded, n)
	case XOR_POS:
		encoded = xor(encoded, byte(pos))
	case ADD_N:
		n := cipher[idx+1]
		encoded = add(encoded, n)
	case ADD_POS:
		encoded = add(encoded, byte(pos))
	}

	return Encode(encoded, cipher, next, pos)
}

func Decode(encoded byte, cipher []byte, idx int, pos int) byte {
	if idx >= len(cipher) {
		return encoded
	}
	decoded := encoded

	next := idx + 1
	// jump 2 places if any of below codes
	if cipher[idx] == XOR_N || cipher[idx] == ADD_N {
		next++
	}
	// decode next one first recursively till end to start from end then pass on
	decoded = Decode(encoded, cipher, next, pos)

	switch cipher[idx] {
	case REVERSE_BITS:
		decoded = reverseBits(decoded)
	case XOR_N:
		n := cipher[idx+1]
		decoded = xor(decoded, n)
	case XOR_POS:
		decoded = xor(decoded, byte(pos))
	case ADD_N:
		n := cipher[idx+1]
		decoded = sub(decoded, n)
	case ADD_POS:
		decoded = sub(decoded, byte(pos))
	}

	return decoded
}
