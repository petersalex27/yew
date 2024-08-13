#include "translate.h"

// returns shorter string (or second string if they are the same length)
//
// calling both _get_shorter_string(a, b) and _get_shorter_string(a, b) (with arguments in that order) will guarantee that each argument
// is returned by one of the two calls
String _get_shorter_string(String a, String b) {
  return a.length < b.length ? a : b;
}

// returns longer string (or first string if they are the same length)
//
// calling both _get_shorter_string(a, b) and _get_shorter_string(a, b) (with arguments in that order) will guarantee that each argument
// is returned by one of the two calls
String _get_longer_string(String a, String b) {
  String c = _get_shorter_string(a, b);
  // if c = a, 
  //  then it must be shorter than b, so return b
  // if c = b, 
  //  then either b is shorter than a or they are the same length, 
  //  in both cases a was the string not returned by _get_shorter_string
  return c.chars == a.chars ? b : a;
}

// compares two strings
//
// returns 0 if the strings are equal
// returns a negative number if a precedes b
// returns a positive number if a succeeds b
Int _string_compare(String a, String b) {
  String s1 = _get_shorter_string(a, b);
  String s2 = _get_longer_string(a, b);
  for (Uint32 i = 0; i < s1.length; i++) {
    if (s1.chars[i] != s2.chars[i]) {
      return (Int)s1.chars[i] - (Int)s2.chars[i];
    }
  }
  // all characters in the shorter string are equal to the corresponding characters in the longer string,
  // so the shorter string precedes the longer string or they are equal
  return (Int)s1.length - (Int)s2.length;
}