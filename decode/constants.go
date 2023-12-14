package decode

const START_OBJECT byte = 0xFA
const END_OBJECT byte = 0xFb

const START_ARRAY byte = 0xF8
const END_ARRAY byte = 0xF9

const LONG_VARIABLE_ASCII = 0xE0
const LONG_UTF8 = 0xE4
const SHARED_STRING_REFERENCE_LONG_1 = 0xEC
const SHARED_STRING_REFERENCE_LONG_2 = 0xED
const SHARED_STRING_REFERENCE_LONG_3 = 0xEE
const SHARED_STRING_REFERENCE_LONG_4 = 0xEF
const STRING_END byte = 0xFC
const START_BINARY byte = 0xFD
const END_BINARY byte = 0xFF

const EMPTY_STRING byte = 0x20
const NULL byte = 0x21
const FALSE byte = 0x22
const TRUE byte = 0x23
const INT_32 byte = 0x24
const INT_64 byte = 0x25
const BIG_INT byte = 0x26
const FLOAT_32 byte = 0x28
const FLOAT_64 byte = 0x29
const BIG_DECIMAL byte = 0x2A // TODO: Use
