package main

import _ "embed"

const (
	AsmTemplate = `; This assembly code was taken verbatim from https://github.com/blaquee/proxydll_template
; Commit attribution goes to Zeffy/Blaquee, but the license does not name anyone. 
;
; License Notice:
;  MIT License
;  
;  Copyright (c) 2017 
;  
;  Permission is hereby granted, free of charge, to any person obtaining a copy
;  of this software and associated documentation files (the "Software"), to deal
;  in the Software without restriction, including without limitation the rights
;  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
;  copies of the Software, and to permit persons to whom the Software is
;  furnished to do so, subject to the following conditions:
;  
;  The above copyright notice and this permission notice shall be included in all
;  copies or substantial portions of the Software.
;  
;  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
;  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
;  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
;  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
;  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
;  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
;  SOFTWARE.

ifndef x64
    .model flat
endif

.code

ifndef x64
    resolve_export_proc proto C, arg1:dword
else
    extern resolve_export_proc:proc
endif

M_EXPORT_PROC macro export, index
export proc
ifndef x64
    invoke resolve_export_proc, index
    jmp dword ptr eax
else
    push rcx
    push rdx
    push r8
    push r9
    push r10
    push r11
    sub rsp, 28h
if index eq 0
    xor rcx, rcx
else
    mov rcx, index
endif
    call resolve_export_proc
    add rsp, 28h
    pop r11
    pop r10
    pop r9
    pop r8
    pop rdx
    pop rcx
    jmp qword ptr rax
endif
export endp
endm

; Generated Export Definitions:
; x86-exclusive exports:
ifndef x64
%s
endif

; x86-64-exclusive exports:
ifdef x64
%s
endif

; Common exports:
%s

end`

	HeaderTemplate = `#pragma once

#define PROXYGEN_GENERATED

#define DLL_FNAME "\\%s.dll"

#ifdef _X86_
#define EXPORT_NAMES \
%s
#endif

#ifdef _X64_
#define EXPORT_NAMES \
%s
#endif`

	HeaderExportNoNameTemplate = `MAKEINTRESOURCEA(%d)`

	ExportDefTemplate = `EXPORTS
`

	RtRcDataTemplate = `
1 VERSIONINFO
FILEVERSION 6, 1, 7601, 23403
PRODUCTVERSION 6, 1, 7601, 23403
FILEOS 0x40004
FILETYPE 0x2
{
    BLOCK "StringFileInfo"
    {
        BLOCK "040904B0"
        {
            VALUE "CompanyName", "%s"
                VALUE "FileDescription", "%s"
                VALUE "FileVersion", "%s"
                VALUE "InternalName", "%s"
                VALUE "LegalCopyright", "%s"
                VALUE "OriginalFilename", "%s"
                VALUE "ProductName", "%s"
                VALUE "ProductVersion", "%s"
        }
    }

    BLOCK "VarFileInfo"
    {
        VALUE "Translation", 0x0409 0x04B0
    }
}
`
)
