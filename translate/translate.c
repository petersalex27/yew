#include "translate.h"

uint64_t _extract_tag(Data data) {
  return data.tag;
}

Int _extract_Int(Data data) {
  return *(Int *)data.data;
}

Uint _extract_Uint(Data data) {
  return *(Uint *)data.data;
}

Float _extract_Float(Data data) {
  return *(Float *)data.data;
}

Char _extract_Char(Data data) {
  return *(Char *)data.data;
}

String _extract_String(Data data) {
  return *(String *)data.data;
}

String _new_string(Char *chars, Uint length) {
  String s;
  s.chars = chars;
  s.length = length;
  return s;
}



bool match(Data scrutinee, Data val) {
  do {
    // the types have been matched at compile time, so we know that the variables have the same type
    //
    // thus, if the scrutinee is a variable, it matches the value
    if (scrutinee.tag == VAR_TAG) {
      return true;
    }

    if (scrutinee.tag != val.tag) {
      return false;
    }
    
    switch (scrutinee.tag) {
      case INT_TAG:
        return _extract_Int(scrutinee) == _extract_Int(val);
      case UINT_TAG:
        return _extract_Uint(scrutinee) == _extract_Uint(val);
      case FLOAT_TAG:
        return _extract_Float(scrutinee) == _extract_Float(val);
      case CHAR_TAG:
        return _extract_Char(scrutinee) == _extract_Char(val);
      case STRING_TAG:
        String a = _extract_String(scrutinee);
        String b = _extract_String(val);
        if (a.length != b.length) {
          return false;
        }
        for (Uint i = 0; i < a.length; i++) {
          if (a.chars[i] != b.chars[i]) {
            return false;
          }
        }
        return true;
      default:
        scrutinee = *(Data *)scrutinee.data;
        val = *(Data *)val.data;
    }
  } while (1);
}