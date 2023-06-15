	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 11, 0
	.globl	_addInt                         ; -- Begin function addInt
	.p2align	2
_addInt:                                ; @addInt
	.cfi_startproc
; %bb.0:
	add	x0, x0, x1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_subInt                         ; -- Begin function subInt
	.p2align	2
_subInt:                                ; @subInt
	.cfi_startproc
; %bb.0:
	sub	x0, x0, x1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_mulInt                         ; -- Begin function mulInt
	.p2align	2
_mulInt:                                ; @mulInt
	.cfi_startproc
; %bb.0:
	mul	x0, x0, x1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_divInt                         ; -- Begin function divInt
	.p2align	2
_divInt:                                ; @divInt
	.cfi_startproc
; %bb.0:
	sdiv	x0, x0, x1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_negInt                         ; -- Begin function negInt
	.p2align	2
_negInt:                                ; @negInt
	.cfi_startproc
; %bb.0:
	neg	x0, x0
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_equalsInt                      ; -- Begin function equalsInt
	.p2align	2
_equalsInt:                             ; @equalsInt
	.cfi_startproc
; %bb.0:
	cmp	x0, x1
	cset	w0, eq
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_notEqualsInt                   ; -- Begin function notEqualsInt
	.p2align	2
_notEqualsInt:                          ; @notEqualsInt
	.cfi_startproc
; %bb.0:
	cmp	x0, x1
	cset	w0, ne
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_greaterInt                     ; -- Begin function greaterInt
	.p2align	2
_greaterInt:                            ; @greaterInt
	.cfi_startproc
; %bb.0:
	cmp	x0, x1
	cset	w0, gt
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_lesserInt                      ; -- Begin function lesserInt
	.p2align	2
_lesserInt:                             ; @lesserInt
	.cfi_startproc
; %bb.0:
	cmp	x0, x1
	cset	w0, lt
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_greaterOrEqualInt              ; -- Begin function greaterOrEqualInt
	.p2align	2
_greaterOrEqualInt:                     ; @greaterOrEqualInt
	.cfi_startproc
; %bb.0:
	cmp	x0, x1
	cset	w0, ge
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_lesserOrEqualInt               ; -- Begin function lesserOrEqualInt
	.p2align	2
_lesserOrEqualInt:                      ; @lesserOrEqualInt
	.cfi_startproc
; %bb.0:
	cmp	x0, x1
	cset	w0, le
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_powerInt                       ; -- Begin function powerInt
	.p2align	2
_powerInt:                              ; @powerInt
	.cfi_startproc
; %bb.0:                                ; %Entry
	mov	w8, #1
	stp	x8, x8, [sp, #-16]!
	.cfi_def_cfa_offset 16
LBB11_1:                                ; %TestLoop
                                        ; =>This Inner Loop Header: Depth=1
	ldp	x8, x9, [sp]
	cmp	x9, x1
	b.le	LBB11_3
; %bb.2:                                ; %IfTrue
                                        ;   in Loop: Header=BB11_1 Depth=1
	ldr	x9, [sp, #8]
	mul	x8, x8, x0
	add	x9, x9, #1
	stp	x8, x9, [sp]
	b	LBB11_1
LBB11_3:                                ; %EndLoop
	mov	x0, x8
	add	sp, sp, #16
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_addDouble                      ; -- Begin function addDouble
	.p2align	2
_addDouble:                             ; @addDouble
	.cfi_startproc
; %bb.0:
	fadd	d0, d0, d1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_subDouble                      ; -- Begin function subDouble
	.p2align	2
_subDouble:                             ; @subDouble
	.cfi_startproc
; %bb.0:
	fsub	d0, d0, d1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_mulDouble                      ; -- Begin function mulDouble
	.p2align	2
_mulDouble:                             ; @mulDouble
	.cfi_startproc
; %bb.0:
	fmul	d0, d0, d1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_divDouble                      ; -- Begin function divDouble
	.p2align	2
_divDouble:                             ; @divDouble
	.cfi_startproc
; %bb.0:
	fdiv	d0, d0, d1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_negDouble                      ; -- Begin function negDouble
	.p2align	2
_negDouble:                             ; @negDouble
	.cfi_startproc
; %bb.0:
	movi	d1, #0000000000000000
	fsub	d0, d1, d0
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_equalsDouble                   ; -- Begin function equalsDouble
	.p2align	2
_equalsDouble:                          ; @equalsDouble
	.cfi_startproc
; %bb.0:
	fcmp	d0, d1
	cset	w0, eq
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_notEqualsDouble                ; -- Begin function notEqualsDouble
	.p2align	2
_notEqualsDouble:                       ; @notEqualsDouble
	.cfi_startproc
; %bb.0:
	fcmp	d0, d1
	cset	w8, mi
	csinc	w0, w8, wzr, le
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_greaterDouble                  ; -- Begin function greaterDouble
	.p2align	2
_greaterDouble:                         ; @greaterDouble
	.cfi_startproc
; %bb.0:
	fcmp	d0, d1
	cset	w0, gt
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_lesserDouble                   ; -- Begin function lesserDouble
	.p2align	2
_lesserDouble:                          ; @lesserDouble
	.cfi_startproc
; %bb.0:
	fcmp	d0, d1
	cset	w0, mi
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_greaterOrEqualDouble           ; -- Begin function greaterOrEqualDouble
	.p2align	2
_greaterOrEqualDouble:                  ; @greaterOrEqualDouble
	.cfi_startproc
; %bb.0:
	fcmp	d0, d1
	cset	w0, ge
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_lesserOrEqualDouble            ; -- Begin function lesserOrEqualDouble
	.p2align	2
_lesserOrEqualDouble:                   ; @lesserOrEqualDouble
	.cfi_startproc
; %bb.0:
	fcmp	d0, d1
	cset	w0, ls
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_equalsBool                     ; -- Begin function equalsBool
	.p2align	2
_equalsBool:                            ; @equalsBool
	.cfi_startproc
; %bb.0:
	and	w8, w0, w1
	and	w0, w8, #0x1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_notEqualsBool                  ; -- Begin function notEqualsBool
	.p2align	2
_notEqualsBool:                         ; @notEqualsBool
	.cfi_startproc
; %bb.0:
	eor	w8, w0, w1
	and	w0, w8, #0x1
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_equalsChar                     ; -- Begin function equalsChar
	.p2align	2
_equalsChar:                            ; @equalsChar
	.cfi_startproc
; %bb.0:
	and	w8, w0, #0xff
	cmp	w8, w1, uxtb
	cset	w0, eq
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_notEqualsChar                  ; -- Begin function notEqualsChar
	.p2align	2
_notEqualsChar:                         ; @notEqualsChar
	.cfi_startproc
; %bb.0:
	and	w8, w0, #0xff
	cmp	w8, w1, uxtb
	cset	w0, ne
	ret
	.cfi_endproc
                                        ; -- End function
	.globl	_main                           ; -- Begin function main
	.p2align	2
_main:                                  ; @main
	.cfi_startproc
; %bb.0:
	sub	sp, sp, #32
	.cfi_def_cfa_offset 32
	stp	x29, x30, [sp, #16]             ; 16-byte Folded Spill
	.cfi_offset w30, -8
	.cfi_offset w29, -16
	mov	w0, #1
	mov	w1, #1
	bl	_addInt
	mov	x8, x0
Lloh0:
	adrp	x0, l_.fmtInt@PAGE
Lloh1:
	add	x0, x0, l_.fmtInt@PAGEOFF
	str	x8, [sp]
	bl	_printf
Lloh2:
	adrp	x0, l_.fmtNewLine@PAGE
Lloh3:
	add	x0, x0, l_.fmtNewLine@PAGEOFF
	bl	_printf
	ldp	x29, x30, [sp, #16]             ; 16-byte Folded Reload
	mov	w0, wzr
	add	sp, sp, #32
	ret
	.loh AdrpAdd	Lloh2, Lloh3
	.loh AdrpAdd	Lloh0, Lloh1
	.cfi_endproc
                                        ; -- End function
	.section	__TEXT,__const
l_.fmtInt:                              ; @.fmtInt
	.asciz	"%ld"

l_.fmtChar:                             ; @.fmtChar
	.asciz	"%c"

l_.fmtDouble:                           ; @.fmtDouble
	.asciz	"%g"

l_.fmtNewLine:                          ; @.fmtNewLine
	.asciz	"\n"

.subsections_via_symbols
