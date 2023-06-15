; formats
@.fmtInt = private constant [4 x i8] c"%ld\00"
@.fmtChar = private constant [3 x i8] c"%c\00"
@.fmtDouble = private constant [3 x i8] c"%g\00"
@.fmtNewLine = private constant [2 x i8] c"\0a\00"

; integer builtins
define i64 @addInt(i64 %a, i64 %b) alwaysinline {
    %1 = add i64 %a, %b
    ret i64 %1
}
define i64 @subInt(i64 %a, i64 %b) alwaysinline {
    %1 = sub i64 %a, %b
    ret i64 %1
}
define i64 @mulInt(i64 %a, i64 %b) alwaysinline {
    %1 = mul i64 %a, %b
    ret i64 %1
}
define i64 @divInt(i64 %a, i64 %b) alwaysinline {
    %1 = sdiv i64 %a, %b
    ret i64 %1
}
define i64 @negInt(i64 %a) alwaysinline {
    %1 = sub i64 0, %a
    ret i64 %1
}
define i1 @equalsInt(i64 %a, i64 %b) alwaysinline {
    %1 = icmp eq i64 %a, %b
    ret i1 %1
}
define i1 @notEqualsInt(i64 %a, i64 %b) alwaysinline {
    %1 = icmp ne i64 %a, %b
    ret i1 %1
}
define i1 @greaterInt(i64 %a, i64 %b) alwaysinline {
    %1 = icmp sgt i64 %a, %b
    ret i1 %1
}
define i1 @lesserInt(i64 %a, i64 %b) alwaysinline {
    %1 = icmp slt i64 %a, %b
    ret i1 %1
}
define i1 @greaterOrEqualInt(i64 %a, i64 %b) alwaysinline {
    %1 = icmp sge i64 %a, %b
    ret i1 %1
}
define i1 @lesserOrEqualInt(i64 %a, i64 %b) alwaysinline {
    %1 = icmp sle i64 %a, %b
    ret i1 %1
}
define i64 @powInt(i64 %a, i64 %b) { ; for any non-negative integer %b
Entry:
    ; c code
    ; uint64_t *acc = alloca(8); uint64_t *out = alloca(8); *out = 1;
    ; for (*acc = 1; *acc < b; (*acc)++) { *out = (*out) * a; }
    ; return *out;
    
    ; allocations
    %acc = alloca i64 ; accumulator
    %out = alloca i64 ; output
    ; initializations
    store i64 1, i64* %acc 
    store i64 1, i64* %out
    ; start loop
    br label %TestLoop

TestLoop: ; check if output must be multiplied by %a
    ; if (*acc) > b then goto IfTrue else goto EndLoop
    %0 = load i64, i64* %acc
    %cond = icmp sgt i64 %0, %b
    br i1 %cond, label %IfTrue, label %EndLoop
IfTrue: ; loop body
    ; multiply output by %a
    %1 = load i64, i64* %out
    %2 = mul i64 %1, %a
    store i64 %2, i64* %out
    ; increment acumulator
    %3 = load i64, i64* %acc
    %4 = add i64 %3, 1
    store i64 %4, i64* %acc
    ; go back to loop test
    br label %TestLoop
EndLoop: ; end of loop

    %5 = load i64, i64* %out
    ret i64 %5
}
define i64 @remainderInt(i64 %a, i64 %b) {
    %0 = srem i64 %a, %b
    ret i64 %0
}

; float builtins
define double @addFloat(double %a, double %b) alwaysinline {
    %1 = fadd double %a, %b
    ret double %1
}
define double @subFloat(double %a, double %b) alwaysinline {
    %1 = fsub double %a, %b
    ret double %1
}
define double @mulFloat(double %a, double %b) alwaysinline {
    %1 = fmul double %a, %b
    ret double %1
}
define double @divFloat(double %a, double %b) alwaysinline {
    %1 = fdiv double %a, %b
    ret double %1
}
define double @negFloat(double %a) alwaysinline {
    %1 = fsub double 0.0, %a
    ret double %1
}
define i1 @equalsFloat(double %a, double %b) alwaysinline {
    %1 = fcmp oeq double %a, %b
    ret i1 %1
}
define i1 @notEqualsFloat(double %a, double %b) alwaysinline {
    %1 = fcmp one double %a, %b
    ret i1 %1
}
define i1 @greaterFloat(double %a, double %b) alwaysinline {
    %1 = fcmp ogt double %a, %b
    ret i1 %1
}
define i1 @lesserFloat(double %a, double %b) alwaysinline {
    %1 = fcmp olt double %a, %b
    ret i1 %1
}
define i1 @greaterOrEqualFloat(double %a, double %b) alwaysinline {
    %1 = fcmp oge double %a, %b
    ret i1 %1
}
define i1 @lesserOrEqualFloat(double %a, double %b) alwaysinline {
    %1 = fcmp ole double %a, %b
    ret i1 %1
}

; bool builtins
define i1 @equalsBool(i1 %a, i1 %b) alwaysinline {
    %1 = and i1 %a, %b
    ret i1 %1
}
define i1 @notEqualsBool(i1 %a, i1 %b) alwaysinline {
    %1 = xor i1 %a, %b
    ret i1 %1
}

; char builtins
define i1 @equalsChar(i8 %a, i8 %b) alwaysinline {
    %1 = icmp eq i8 %a, %b
    ret i1 %1
}
define i1 @notEqualsChar(i8 %a, i8 %b) alwaysinline {
    %1 = icmp ne i8 %a, %b
    ret i1 %1
}

declare i32 @printf(i8* noalias nocapture, ...)
define i32 @main() {
    %1 = call i64 @addInt(i64 1, i64 1)
    %2 = call i32 (i8*, ...) @printf(i8* @.fmtInt, i64 %1)
    %3 = call i32 (i8*, ...) @printf(i8* @.fmtNewLine)
    ret i32 0
}