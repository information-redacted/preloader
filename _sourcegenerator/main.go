package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	peparser "github.com/saferwall/pe"
)

type LimitedExport struct {
	Name    string
	Ordinal uint32
	NoName  bool
}

var (
	exports32     []LimitedExport
	exports64     []LimitedExport
	commonExports []LimitedExport

	dllBaseName = ""
	rcData      map[string]string
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Please provide a DLL as the first input parameter.")
		return
	}

	dllBaseName = strings.TrimSuffix(args[1], ".dll")

	p32 := "C:\\Windows\\SysWOW64\\" + args[1]
	p64 := "C:\\Windows\\System32\\" + args[1]

	if _, err := os.Stat(p32); err != nil {
		fmt.Println("failed to open 32-bit version of DLL:", err)
		return
	}

	if _, err := os.Stat(p64); err != nil {
		fmt.Println("failed to open 64-bit version of DLL:", err)
		return
	}

	err := parseExports(p32, &exports32)
	if err != nil {
		log.Fatalf("could not parse 32-bit exports: %v", err)
	}

	err = parseExports(p64, &exports64)
	if err != nil {
		log.Fatalf("could not parse 64-bit exports: %v", err)
	}

	err = os.MkdirAll(fmt.Sprintf("generated_%s", dllBaseName), 0777)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatalf("could not create storage directory: %v", err)
		}
	}

	err = os.WriteFile(fmt.Sprintf("generated_%s/generated.rc", dllBaseName),
		[]byte(fmt.Sprintf(RtRcDataTemplate,
			rcData["CompanyName"],
			rcData["FileDescription"],
			rcData["FileVersion"],
			rcData["InternalName"],
			rcData["LegalCopyright"],
			rcData["OriginalFilename"],
			rcData["ProductName"],
			rcData["ProductVersion"])),
		0)
	if err != nil {
		log.Fatalf("could not write resource file: %v", err)
	}

	slices.SortFunc(exports32, SortSliceCmp)
	slices.SortFunc(exports64, SortSliceCmp)

	commonExports = make([]LimitedExport, 0)
	for _, ex := range exports64 {
		for _, ex32 := range exports32 {
			if ex.Name == ex32.Name && ex.Ordinal == ex32.Ordinal {
				commonExports = append(commonExports, ex)
				break
			}
		}
	}

	slices.SortFunc(commonExports, SortSliceCmp)

	exports64 = SubtractSlices(exports64, commonExports)
	exports32 = SubtractSlices(exports32, commonExports)

	sExports32 := ""
	for _, e := range exports32 {
		sExports32 += fmt.Sprintf("M_EXPORT_PROC %s, %d\r\n", e.Name, e.Ordinal)
	}

	sExports64 := ""
	for _, e := range exports64 {
		sExports64 += fmt.Sprintf("M_EXPORT_PROC %s, %d\r\n", e.Name, e.Ordinal)
	}

	sCommonExports := ""
	for _, e := range commonExports {
		sCommonExports += fmt.Sprintf("M_EXPORT_PROC %s, %d\r\n", e.Name, e.Ordinal)
	}

	sExports32 = strings.TrimSpace(sExports32)
	sExports64 = strings.TrimSpace(sExports64)
	sCommonExports = strings.TrimSpace(sCommonExports)

	err = os.WriteFile(fmt.Sprintf("generated_%s/generated.asm", dllBaseName),
		[]byte(fmt.Sprintf(AsmTemplate, sExports32, sExports64, sCommonExports)), 777)
	if err != nil {
		log.Fatalf("could not write assembly file: %v", err)
	}

	sExports32 = generateExportDef(append(commonExports, exports32...))
	sExports64 = generateExportDef(append(commonExports, exports64...))

	err = os.WriteFile(fmt.Sprintf("generated_%s/exports_x86.def", dllBaseName), []byte(sExports32), 777)
	if err != nil {
		log.Fatalf("could not write x86 export definitions file: %v", err)
	}

	err = os.WriteFile(fmt.Sprintf("generated_%s/exports_x64.def", dllBaseName), []byte(sExports64), 777)
	if err != nil {
		log.Fatalf("could not write x86-64 export definitions file: %v", err)
	}

	allExports32 := append(commonExports, exports32...)
	allExports64 := append(commonExports, exports64...)

	slices.SortFunc(allExports32, SortSliceCmp)
	slices.SortFunc(allExports64, SortSliceCmp)

	sExportsHeader32 := ""
	for _, e := range allExports32 {
		h := ""
		if e.NoName {
			h = fmt.Sprintf(HeaderExportNoNameTemplate, e.Ordinal)
		} else {
			h = fmt.Sprintf("\"%s\"", e.Name)
		}
		if sExportsHeader32 != "" {
			sExportsHeader32 += ", "
		}
		sExportsHeader32 += h
	}

	sExportsHeader64 := ""
	for _, e := range allExports64 {
		h := ""
		if e.NoName {
			h = fmt.Sprintf(HeaderExportNoNameTemplate, e.Ordinal)
		} else {
			h = fmt.Sprintf("\"%s\"", e.Name)
		}

		if sExportsHeader64 != "" {
			sExportsHeader64 += ", "
		}
		sExportsHeader64 += h
	}

	sExportsHeader32 = SplitDefine(sExportsHeader32, 80)
	sExportsHeader64 = SplitDefine(sExportsHeader64, 80)

	err = os.WriteFile(fmt.Sprintf("generated_%s/header.h", dllBaseName),
		[]byte(fmt.Sprintf(HeaderTemplate, dllBaseName, sExportsHeader32, sExportsHeader64)), 777)
	if err != nil {
		log.Fatalf("could not write generated header file: %v", err)
	}

	err = os.WriteFile(fmt.Sprintf("generated_%s/target_name.txt", dllBaseName), []byte(dllBaseName), 777)
	if err != nil {
		log.Fatalf("could not write target name file: %v", err)
	}
}

func parseExports(peFile string, sl *[]LimitedExport) error {
	pe, err := peparser.New(peFile, &peparser.Options{})
	if err != nil {
		return err
	}
	defer pe.Close()

	err = pe.Parse()
	if err != nil {
		return err
	}

	*sl = make([]LimitedExport, len(pe.Export.Functions))
	for i, e := range pe.Export.Functions {
		n := e.Name
		noname := false
		if n == "" {
			n = strings.ToUpper(dllBaseName) + fmt.Sprintf("_%d", e.Ordinal)
			noname = true
		}

		(*sl)[i].Name = n
		(*sl)[i].Ordinal = e.Ordinal
		(*sl)[i].NoName = noname
	}

	rcData, _ = pe.ParseVersionResources()

	return nil
}

func generateExportDef(exports []LimitedExport) string {
	res := ExportDefTemplate
	for _, e := range exports {
		res += fmt.Sprintf("    %s @%d", e.Name, e.Ordinal)
		if e.NoName {
			res += " NONAME"
		}
		res += "\r\n"
	}

	return res
}

// SubtractSlices returns a new slice containing elements from slice1 that are not present in slice2
func SubtractSlices(slice1, slice2 []LimitedExport) []LimitedExport {
	exists := make(map[LimitedExport]bool)
	for _, item := range slice2 {
		exists[item] = true
	}

	var result []LimitedExport
	for _, item := range slice1 {
		if !exists[item] {
			result = append(result, item)
		}
	}
	return result
}

func SortSliceCmp(val1, val2 LimitedExport) int {
	if val1.Ordinal == val2.Ordinal {
		return 0
	}

	if val1.Ordinal < val2.Ordinal {
		return -1
	}

	if val1.Ordinal > val2.Ordinal {
		return 1
	}

	return 1
}

func SplitDefine(str string, maxLen int) string {
	// FIXME(redacted): MAKEINTRESOURCEA(int) will cause invalid splits.
	//   This won't ever happen with the usual suspects (version, winmm, etc)
	//   but if your header's fucked because of many unnamed export ordinals, this'll be fun to debug
	var builder strings.Builder
	lineLength := 0
	inQuotes := false

	for i, char := range str {
		if char == '"' {
			inQuotes = !inQuotes
		}

		if char == ',' && !inQuotes && lineLength > maxLen {
			builder.WriteString(", \\\r\n    ")
			lineLength = 0
			continue
		}

		builder.WriteRune(char)
		lineLength++

		if lineLength == maxLen && !inQuotes {
			if i+1 < len(str) && str[i+1] == ',' {
				lineLength--
				continue
			}
			builder.WriteString(" \\\r\n    ")
			lineLength = 0
		}
	}
	return builder.String()
}
