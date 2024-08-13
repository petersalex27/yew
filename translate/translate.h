#ifndef TRANSLATE_H
#define TRANSLATE_H

#include <stdint.h>
#include <stdbool.h>

// each tag should be associated with an array of step values to get 
// from one value to the next
typedef struct {
  uint64_t tag;
  // will either point to a primitive value or another Data struct
  void *data;
} Data;

#define INT_TAG 0
#define UINT_TAG 1
#define FLOAT_TAG 3
#define CHAR_TAG 2
#define STRING_TAG 4
#define VAR_TAG 5

typedef int64_t Int;
typedef uint64_t Uint;
typedef uint32_t Uint32;
typedef double Float;
typedef int8_t Char;

typedef struct {
  Char *chars;
  Uint32 length;
} String;

uint64_t _extract_tag(Data data);
Int _extract_Int(Data data);
Uint _extract_Uint(Data data);
Float _extract_Float(Data data);
Char _extract_Char(Data data);
String _extract_String(Data data);

String _new_string(Char *chars, Uint length);
Int _string_compare(String a, String b);

bool match(Data a, Data b);

#endif // TRANSLATE_H