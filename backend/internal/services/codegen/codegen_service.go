package codegen

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain"
)

// ProblemSignature is an alias for domain.ProblemSchema for local use if needed,
// but we'll use domain.ProblemSchema directly in methods.

type CodeGenService struct {
	typeImplRepo domain.TypeImplementationRepository
}

func NewCodeGenService(typeImplRepo domain.TypeImplementationRepository) *CodeGenService {
	return &CodeGenService{
		typeImplRepo: typeImplRepo,
	}
}

// GenerateStubCode generates starter code for users
func (s *CodeGenService) GenerateStubCode(signature domain.ProblemSchema, languageSlug string) (string, error) {
	// Validate inputs
	if signature.FunctionName == "" {
		return "", errors.New("function name is required")
	}
	if len(signature.Parameters) == 0 {
		return "", errors.New("at least one parameter is required")
	}

	// Identify custom types
	customTypes := s.identifyCustomTypes(signature)

	switch languageSlug {
	case "python":
		return s.generatePythonStub(signature, customTypes)
	case "javascript":
		return s.generateJavaScriptStub(signature, customTypes)
	case "java":
		// Java support for custom types todo
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for Java yet")
		}
		return s.generateJavaStub(signature)
	case "c++":
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for C++ yet")
		}
		return s.generateCppStub(signature)
	case "c":
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for C yet")
		}
		return s.generateCStub(signature)
	case "go":
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for Go yet")
		}
		return s.generateGoStub(signature)
	default:
		return "", fmt.Errorf("unsupported language: %s", languageSlug)
	}
}

func (s *CodeGenService) identifyCustomTypes(sig domain.ProblemSchema) []string {
	types := []string{}
	seen := make(map[domain.GenericType]bool)

	if s.isCustomType(sig.ReturnType) {
		types = append(types, string(sig.ReturnType))
		seen[sig.ReturnType] = true
	}

	for _, param := range sig.Parameters {
		if param.IsCustom && !seen[param.Type] {
			types = append(types, string(param.Type))
			seen[param.Type] = true
		}
	}

	return types
}

func (s *CodeGenService) isCustomType(typeName domain.GenericType) bool {
	// This could be dynamic, but checking seeded types for now
	validTypes := []domain.GenericType{"TreeNode", "ListNode", "GraphNode", "Node"}
	for _, vt := range validTypes {
		if typeName == vt {
			return true
		}
	}
	return false
}

// GenerateTestHarness generates wrapper code that runs test cases
func (s *CodeGenService) GenerateTestHarness(signature domain.ProblemSchema, userCode string, languageSlug string, testCases []domain.TestCase, validationType string) (string, error) {
	// Validate
	if userCode == "" {
		return "", errors.New("user code is required")
	}

	switch languageSlug {
	case "python":
		return s.GeneratePythonHarness(signature, userCode, testCases, validationType)
	case "javascript":
		return s.GenerateJavaScriptHarness(signature, userCode, testCases, validationType)
	case "java":
		return s.GenerateJavaHarness(signature, userCode, testCases, validationType)
	case "c++":
		return s.GenerateCppHarness(signature, userCode, testCases, validationType)
	case "c":
		return s.GenerateCHarness(signature, userCode, testCases, validationType)
	case "go":
		return s.GenerateGoHarness(signature, userCode, testCases, validationType)
	default:
		return "", fmt.Errorf("unsupported language: %s", languageSlug)
	}
}

// Type mapping helpers
func (s *CodeGenService) mapTypeToPython(typ domain.GenericType) string {
	typeMap := map[domain.GenericType]string{
		domain.TypeInteger:      "int",
		domain.TypeIntegerArray: "List[int]",
		domain.TypeString:       "str",
		domain.TypeStringArray:  "List[str]",
		domain.TypeBoolean:      "bool",
		"double":                "float", // keep double for now for existing logic
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return string(typ)
}

func (s *CodeGenService) mapTypeToJavaScript(typ domain.GenericType) string {
	return ""
}

func (s *CodeGenService) mapTypeToJava(typ domain.GenericType) string {
	typeMap := map[domain.GenericType]string{
		domain.TypeInteger:      "int",
		domain.TypeIntegerArray: "int[]",
		domain.TypeString:       "String",
		domain.TypeStringArray:  "String[]",
		domain.TypeBoolean:      "boolean",
		"double":                "double",
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return string(typ)
}

func (s *CodeGenService) mapTypeToCpp(typ domain.GenericType) string {
	typeMap := map[domain.GenericType]string{
		domain.TypeInteger:      "int",
		domain.TypeIntegerArray: "vector<int>",
		domain.TypeString:       "string",
		domain.TypeStringArray:  "vector<string>",
		domain.TypeBoolean:      "bool",
		"double":                "double",
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return string(typ)
}

func (s *CodeGenService) mapTypeToC(typ domain.GenericType) string {
	typeMap := map[domain.GenericType]string{
		domain.TypeInteger:      "int",
		domain.TypeIntegerArray: "int*",
		domain.TypeString:       "char*",
		domain.TypeStringArray:  "char**",
		domain.TypeBoolean:      "bool",
		"double":                "double",
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return string(typ)
}

func (s *CodeGenService) mapTypeToGo(typ domain.GenericType) string {
	typeMap := map[domain.GenericType]string{
		domain.TypeInteger:      "int",
		domain.TypeIntegerArray: "[]int",
		domain.TypeString:       "string",
		domain.TypeStringArray:  "[]string",
		domain.TypeBoolean:      "bool",
		"double":                "float64",
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return string(typ)
}

func (s *CodeGenService) formatCppLiteral(val interface{}, typ domain.GenericType) string {
	if val == nil {
		return "{}"
	}

	switch typ {
	case domain.TypeInteger:
		return fmt.Sprintf("%v", val)
	case domain.TypeBoolean:
		return fmt.Sprintf("%v", val)
	case domain.TypeString:
		// Escape quotes for C++ string literal
		str := fmt.Sprintf("%v", val)
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(str, "\"", "\\\""))
	case domain.TypeIntegerArray:
		arr, ok := val.([]interface{})
		if !ok {
			return "{}"
		}
		items := []string{}
		for _, v := range arr {
			items = append(items, fmt.Sprintf("%v", v))
		}
		return "{" + strings.Join(items, ", ") + "}"
	case domain.TypeStringArray:
		arr, ok := val.([]interface{})
		if !ok {
			return "{}"
		}
		items := []string{}
		for _, v := range arr {
			str := fmt.Sprintf("%v", v)
			items = append(items, fmt.Sprintf("\"%s\"", strings.ReplaceAll(str, "\"", "\\\"")))
		}
		return "{" + strings.Join(items, ", ") + "}"
	}
	return fmt.Sprintf("%v", val)
}

func (s *CodeGenService) formatJavaLiteral(val interface{}, typ domain.GenericType) string {
	if val == nil {
		return "null"
	}

	switch typ {
	case domain.TypeInteger:
		return fmt.Sprintf("%v", val)
	case domain.TypeBoolean:
		return fmt.Sprintf("%v", val)
	case domain.TypeString:
		str := fmt.Sprintf("%v", val)
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(str, "\"", "\\\""))
	case domain.TypeIntegerArray:
		arr, ok := val.([]interface{})
		if !ok {
			return "new int[]{}"
		}
		items := []string{}
		for _, v := range arr {
			items = append(items, fmt.Sprintf("%v", v))
		}
		return "new int[]{" + strings.Join(items, ", ") + "}"
	case domain.TypeStringArray:
		arr, ok := val.([]interface{})
		if !ok {
			return "new String[]{}"
		}
		items := []string{}
		for _, v := range arr {
			str := fmt.Sprintf("%v", v)
			items = append(items, fmt.Sprintf("\"%s\"", strings.ReplaceAll(str, "\"", "\\\"")))
		}
		return "new String[]{" + strings.Join(items, ", ") + "}"
	}
	return fmt.Sprintf("%v", val)
}

func (s *CodeGenService) formatCLiteral(val interface{}, typ domain.GenericType) string {
	if val == nil {
		return "NULL"
	}

	switch typ {
	case domain.TypeInteger:
		return fmt.Sprintf("%v", val)
	case domain.TypeBoolean:
		if b, ok := val.(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
		return fmt.Sprintf("%v", val)
	case domain.TypeString:
		str := fmt.Sprintf("%v", val)
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(str, "\"", "\\\""))
	case domain.TypeIntegerArray:
		arr, ok := val.([]interface{})
		if !ok {
			return "{}"
		}
		items := []string{}
		for _, v := range arr {
			items = append(items, fmt.Sprintf("%v", v))
		}
		return "{" + strings.Join(items, ", ") + "}"
	case domain.TypeStringArray:
		arr, ok := val.([]interface{})
		if !ok {
			return "{NULL}"
		}
		items := []string{}
		for _, v := range arr {
			str := fmt.Sprintf("%v", v)
			items = append(items, fmt.Sprintf("\"%s\"", strings.ReplaceAll(str, "\"", "\\\"")))
		}
		return "{" + strings.Join(items, ", ") + "}"
	}
	return fmt.Sprintf("%v", val)
}

// Python stub generator
func (s *CodeGenService) generatePythonStub(sig domain.ProblemSchema, customTypes []string) (string, error) {
	var sb strings.Builder
	sb.WriteString("from typing import List, Optional\n\n")

	// Add custom type definitions
	for _, typeName := range customTypes {
		impl, err := s.typeImplRepo.GetByTypeAndLanguageSlug(typeName, "python")
		if err != nil {
			return "", fmt.Errorf("failed to get implementation for %s: %w", typeName, err)
		}
		sb.WriteString(impl.ClassDefinition)
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("def %s(", sig.FunctionName))
	params := []string{}
	for _, param := range sig.Parameters {
		paramType := s.mapTypeToPython(param.Type)
		params = append(params, fmt.Sprintf("%s: %s", param.Name, paramType))
	}
	sb.WriteString(strings.Join(params, ", "))
	returnType := s.mapTypeToPython(sig.ReturnType)
	sb.WriteString(fmt.Sprintf(") -> %s:\n", returnType))
	sb.WriteString("    # Write your code here\n")
	sb.WriteString("    pass\n")
	return sb.String(), nil
}

// Python harness generator
func (s *CodeGenService) GeneratePythonHarness(sig domain.ProblemSchema, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder
	sb.WriteString("import json\nimport sys\nimport time\nimport tracemalloc\nimport signal\nfrom typing import List, Optional, Any\n\n")

	customTypes := s.identifyCustomTypes(sig)

	// Add custom type definitions and helpers
	for _, typeName := range customTypes {
		impl, err := s.typeImplRepo.GetByTypeAndLanguageSlug(typeName, "python")
		if err != nil {
			return "", fmt.Errorf("failed to get implementation for %s: %w", typeName, err)
		}
		sb.WriteString(fmt.Sprintf("# Custom type: %s\n", typeName))
		sb.WriteString(impl.ClassDefinition)
		sb.WriteString("\n")
		sb.WriteString(impl.DeserializerCode)
		sb.WriteString("\n")
		sb.WriteString(impl.SerializerCode)
		sb.WriteString("\n\n")
	}

	sb.WriteString("# User's solution\n")
	sb.WriteString(userCode)
	sb.WriteString("\n\n")

	// Embed Test Cases
	sb.WriteString("# Embedded Test Cases\n")
	sb.WriteString("TEST_CASES = ")

	type EmbeddedTestCase struct {
		Input    interface{} `json:"input"`
		Expected interface{} `json:"expected"`
		IsSample bool        `json:"is_sample"`
	}

	embeddedCases := []EmbeddedTestCase{}
	for _, tc := range testCases {
		var input, expected interface{}
		json.Unmarshal([]byte(tc.Input), &input)
		json.Unmarshal([]byte(tc.ExpectedOutput), &expected)
		embeddedCases = append(embeddedCases, EmbeddedTestCase{
			Input:    input,
			Expected: expected,
			IsSample: tc.IsSample,
		})
	}
	casesJSON, _ := json.MarshalIndent(embeddedCases, "", "    ")
	sb.WriteString("json.loads('''")
	sb.WriteString(string(casesJSON))
	sb.WriteString("''')")
	sb.WriteString("\n\n")

	// Validation Logic
	sb.WriteString("def compare_outputs(actual, expected, val_type):\n")
	sb.WriteString("    if val_type == 'EXACT':\n")
	sb.WriteString("        return json.dumps(actual, sort_keys=True) == json.dumps(expected, sort_keys=True)\n")
	sb.WriteString("    elif val_type == 'UNORDERED':\n")
	sb.WriteString("        if not isinstance(actual, list) or not isinstance(expected, list):\n")
	sb.WriteString("            return json.dumps(actual, sort_keys=True) == json.dumps(expected, sort_keys=True)\n")
	sb.WriteString("        return sorted(json.dumps(i, sort_keys=True) for i in actual) == sorted(json.dumps(i, sort_keys=True) for i in expected)\n")
	sb.WriteString("    return json.dumps(actual, sort_keys=True) == json.dumps(expected, sort_keys=True)\n\n")

	sb.WriteString("def timeout_handler(signum, frame):\n")
	sb.WriteString("    raise TimeoutError(\"Test exceeded timeout\")\n\n")
	sb.WriteString("signal.signal(signal.SIGALRM, timeout_handler)\n\n")

	sb.WriteString("if __name__ == \"__main__\":\n")
	sb.WriteString("    results = []\n")
	sb.WriteString(fmt.Sprintf("    validation_type = '%s'\n\n", validationType))
	sb.WriteString("    for i, test in enumerate(TEST_CASES):\n")
	sb.WriteString("        status = \"passed\"\n")
	sb.WriteString("        output = None\n")
	sb.WriteString("        time_ms = 0\n")
	sb.WriteString("        memory_kb = 0\n")
	sb.WriteString("        error = None\n")
	sb.WriteString("        \n")
	sb.WriteString("        tracemalloc.start()\n")
	sb.WriteString("        start_time = time.perf_counter()\n")
	sb.WriteString("        signal.alarm(2)\n")
	sb.WriteString("        \n")
	sb.WriteString("        try:\n")
	// Parameter deserialization
	for j, param := range sig.Parameters {
		if param.IsCustom {
			deserializerFunc := fmt.Sprintf("deserialize_%s", strings.ToLower(string(param.Type)))
			sb.WriteString(fmt.Sprintf("            %s = %s(test['input'][%d])\n", param.Name, deserializerFunc, j))
		} else {
			sb.WriteString(fmt.Sprintf("            %s = test['input'][%d]\n", param.Name, j))
		}
	}
	sb.WriteString("\n")
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("            res = %s(%s)\n", sig.FunctionName, strings.Join(paramNames, ", ")))

	// Serializing actual result for comparison and output
	if s.isCustomType(sig.ReturnType) {
		serializerFunc := fmt.Sprintf("serialize_%s", strings.ToLower(string(sig.ReturnType)))
		sb.WriteString(fmt.Sprintf("            actual_res = %s(res)\n", serializerFunc))
	} else {
		sb.WriteString("            actual_res = res\n")
	}

	sb.WriteString("            output = actual_res\n")
	sb.WriteString("            if not compare_outputs(actual_res, test['expected'], validation_type):\n")
	sb.WriteString("                status = \"failed\"\n")

	sb.WriteString("        except TimeoutError:\n            status = \"timeout\"\n")
	sb.WriteString("        except Exception as e:\n            status = \"runtime_error\"\n            error = str(e)\n")
	sb.WriteString("        finally:\n            signal.alarm(0)\n")
	sb.WriteString("            end_time = time.perf_counter()\n")
	sb.WriteString("            current, peak = tracemalloc.get_traced_memory()\n")
	sb.WriteString("            tracemalloc.stop()\n")
	sb.WriteString("            time_ms = int((end_time - start_time) * 1000)\n")
	sb.WriteString("            memory_kb = int(peak / 1024)\n")
	sb.WriteString("        \n")
	sb.WriteString("        results.append({\n")
	sb.WriteString("            \"status\": status,\n")
	sb.WriteString("            \"time_ms\": time_ms,\n")
	sb.WriteString("            \"memory_kb\": memory_kb,\n")
	sb.WriteString("            \"output\": json.dumps(output, sort_keys=True) if output is not None else \"\",\n")
	sb.WriteString("            \"error\": error\n")
	sb.WriteString("        })\n")
	sb.WriteString("    \n")
	sb.WriteString("    # Standardized Verdict Aggregation\n")
	sb.WriteString("    final_verdict = \"ACCEPTED\"\n")
	sb.WriteString("    max_runtime = 0\n")
	sb.WriteString("    max_memory = 0\n")
	sb.WriteString("    test_results = []\n")
	sb.WriteString("    \n")
	sb.WriteString("    for i, res in enumerate(results):\n")
	sb.WriteString("        status = res[\"status\"]\n")
	sb.WriteString("        if status != \"passed\" and final_verdict == \"ACCEPTED\":\n")
	sb.WriteString("            if status == \"timeout\": final_verdict = \"TLE\"\n")
	sb.WriteString("            elif status == \"runtime_error\": final_verdict = \"RUNTIME_ERROR\"\n")
	sb.WriteString("            elif status == \"failed\": final_verdict = \"WRONG_ANSWER\"\n")
	sb.WriteString("            else: final_verdict = status.upper()\n")
	sb.WriteString("            \n")
	sb.WriteString("        if res[\"time_ms\"] > max_runtime: max_runtime = res[\"time_ms\"]\n")
	sb.WriteString("        if res[\"memory_kb\"] > max_memory: max_memory = res[\"memory_kb\"]\n")
	sb.WriteString("        \n")
	sb.WriteString("        test_results.append({\n")
	sb.WriteString("            \"passed\": status == \"passed\",\n")
	sb.WriteString("            \"input\": json.dumps(TEST_CASES[i][\"input\"]),\n")
	sb.WriteString("            \"actual\": res[\"output\"],\n")
	sb.WriteString("            \"error\": res[\"error\"]\n")
	sb.WriteString("        })\n")
	sb.WriteString("        \n")
	sb.WriteString("    print(json.dumps({\n")
	sb.WriteString("        \"verdict\": final_verdict,\n")
	sb.WriteString("        \"runtime\": max_runtime,\n")
	sb.WriteString("        \"memory\": max_memory,\n")
	sb.WriteString("        \"test_results\": test_results\n")
	sb.WriteString("    }))\n")
	sb.WriteString("    sys.exit(0)\n")
	return sb.String(), nil
}

// JavaScript stub generator
func (s *CodeGenService) generateJavaScriptStub(sig domain.ProblemSchema, customTypes []string) (string, error) {
	var sb strings.Builder

	// Add custom type definitions
	for _, typeName := range customTypes {
		impl, err := s.typeImplRepo.GetByTypeAndLanguageSlug(typeName, "javascript")
		if err != nil {
			return "", fmt.Errorf("failed to get implementation for %s: %w", typeName, err)
		}
		sb.WriteString(impl.ClassDefinition)
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("function %s(", sig.FunctionName))
	params := []string{}
	for _, param := range sig.Parameters {
		params = append(params, param.Name)
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") {\n")
	sb.WriteString("    // Write your code here\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}

// JavaScript harness generator
func (s *CodeGenService) GenerateJavaScriptHarness(sig domain.ProblemSchema, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder

	customTypes := s.identifyCustomTypes(sig)

	// Add custom type definitions and helpers
	for _, typeName := range customTypes {
		impl, err := s.typeImplRepo.GetByTypeAndLanguageSlug(typeName, "javascript")
		if err != nil {
			return "", fmt.Errorf("failed to get implementation for %s: %w", typeName, err)
		}
		sb.WriteString(fmt.Sprintf("// Custom type: %s\n", typeName))
		sb.WriteString(impl.ClassDefinition)
		sb.WriteString("\n")
		sb.WriteString(impl.DeserializerCode)
		sb.WriteString("\n")
		sb.WriteString(impl.SerializerCode)
		sb.WriteString("\n\n")
	}

	sb.WriteString("// User's solution\n")
	sb.WriteString(userCode)
	sb.WriteString("\n\n")

	// Embed Test Cases
	sb.WriteString("// Embedded Test Cases\n")
	sb.WriteString("const TEST_CASES = ")

	type EmbeddedTestCase struct {
		Input    interface{} `json:"input"`
		Expected interface{} `json:"expected"`
		IsSample bool        `json:"is_sample"`
	}

	embeddedCases := []EmbeddedTestCase{}
	for _, tc := range testCases {
		var input, expected interface{}
		json.Unmarshal([]byte(tc.Input), &input)
		json.Unmarshal([]byte(tc.ExpectedOutput), &expected)
		embeddedCases = append(embeddedCases, EmbeddedTestCase{
			Input:    input,
			Expected: expected,
			IsSample: tc.IsSample,
		})
	}
	casesJSON, _ := json.MarshalIndent(embeddedCases, "", "    ")
	sb.WriteString(string(casesJSON))
	sb.WriteString(";\n\n")

	// Validation Logic
	sb.WriteString("function compareOutputs(actual, expected, valType) {\n")
	sb.WriteString("    const actualStr = JSON.stringify(actual, Object.keys(actual || {}).sort());\n")
	sb.WriteString("    const expectedStr = JSON.stringify(expected, Object.keys(expected || {}).sort());\n")
	sb.WriteString("    if (valType === 'EXACT') {\n")
	sb.WriteString("        return actualStr === expectedStr;\n")
	sb.WriteString("    } else if (valType === 'UNORDERED') {\n")
	sb.WriteString("        if (!Array.isArray(actual) || !Array.isArray(expected)) return actualStr === expectedStr;\n")
	sb.WriteString("        const sortFn = (a, b) => JSON.stringify(a).localeCompare(JSON.stringify(b));\n")
	sb.WriteString("        return JSON.stringify([...actual].sort(sortFn)) === JSON.stringify([...expected].sort(sortFn));\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return actualStr === expectedStr;\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// Test harness\n")
	sb.WriteString("const results = [];\n")
	sb.WriteString(fmt.Sprintf("const validationType = '%s';\n\n", validationType))
	sb.WriteString("(async () => {\n")
	sb.WriteString("    for (let i = 0; i < TEST_CASES.length; i++) {\n")
	sb.WriteString("        const test = TEST_CASES[i];\n")
	sb.WriteString("        let status = \"passed\";\n")
	sb.WriteString("        let output = null;\n")
	sb.WriteString("        let error = null;\n")
	sb.WriteString("        const startTime = process.hrtime.bigint();\n")
	sb.WriteString("        const startMem = process.memoryUsage().heapUsed;\n\n")
	sb.WriteString("        const testPromise = (async () => {\n")
	// Parameter deserialization
	for j, param := range sig.Parameters {
		if param.IsCustom {
			deserializerFunc := fmt.Sprintf("deserialize%s", param.Type)
			sb.WriteString(fmt.Sprintf("            const %s = %s(test.input[%d]);\n", param.Name, deserializerFunc, j))
		} else {
			sb.WriteString(fmt.Sprintf("            const %s = test.input[%d];\n", param.Name, j))
		}
	}
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("            let res = await %s(%s);\n", sig.FunctionName, strings.Join(paramNames, ", ")))

	if s.isCustomType(sig.ReturnType) {
		serializerFunc := fmt.Sprintf("serialize%s", sig.ReturnType)
		sb.WriteString(fmt.Sprintf("            let actualRes = %s(res);\n", serializerFunc))
	} else {
		sb.WriteString("            let actualRes = res;\n")
	}

	sb.WriteString("            return actualRes;\n")
	sb.WriteString("        })();\n\n")

	sb.WriteString("        const timeoutPromise = new Promise((_, reject) => setTimeout(() => reject(new Error('timeout')), 2000));\n\n")
	sb.WriteString("        try {\n")
	sb.WriteString("            output = await Promise.race([testPromise, timeoutPromise]);\n")
	sb.WriteString("            if (!compareOutputs(output, test.expected, validationType)) {\n")
	sb.WriteString("                status = \"failed\";\n")
	sb.WriteString("            }\n")
	sb.WriteString("        } catch (e) {\n")
	sb.WriteString("            if (e.message === 'timeout') {\n")
	sb.WriteString("                status = \"timeout\";\n")
	sb.WriteString("            } else {\n")
	sb.WriteString("                status = \"runtime_error\";\n")
	sb.WriteString("                error = e.message;\n")
	sb.WriteString("            }\n")
	sb.WriteString("        }\n\n")

	sb.WriteString("        const endTime = process.hrtime.bigint();\n")
	sb.WriteString("        const endMem = process.memoryUsage().heapUsed;\n")
	sb.WriteString("        const timeMs = Number((endTime - startTime) / BigInt(1000000));\n")
	sb.WriteString("        const memoryKb = Math.max(0, Math.floor((endMem - startMem) / 1024));\n\n")
	sb.WriteString("        results.push({\n")
	sb.WriteString("            status: status,\n")
	sb.WriteString("            time_ms: timeMs,\n")
	sb.WriteString("            memory_kb: memoryKb,\n")
	sb.WriteString("            output: output !== null ? JSON.stringify(output) : \"\",\n")
	sb.WriteString("            error: error\n")
	sb.WriteString("        });\n")
	sb.WriteString("    }\n")
	sb.WriteString("\n")
	sb.WriteString("    // Standardized Verdict Aggregation\n")
	sb.WriteString("    let finalVerdict = \"ACCEPTED\";\n")
	sb.WriteString("    let maxRuntime = 0;\n")
	sb.WriteString("    let maxMemory = 0;\n")
	sb.WriteString("    const testResults = results.map((res, i) => {\n")
	sb.WriteString("        if (res.status !== \"passed\" && finalVerdict === \"ACCEPTED\") {\n")
	sb.WriteString("            if (res.status === \"timeout\") finalVerdict = \"TLE\";\n")
	sb.WriteString("            else if (res.status === \"runtime_error\") finalVerdict = \"RUNTIME_ERROR\";\n")
	sb.WriteString("            else if (res.status === \"failed\") finalVerdict = \"WRONG_ANSWER\";\n")
	sb.WriteString("            else finalVerdict = res.status.toUpperCase();\n")
	sb.WriteString("        }\n")
	sb.WriteString("        if (res.time_ms > maxRuntime) maxRuntime = res.time_ms;\n")
	sb.WriteString("        if (res.memory_kb > maxMemory) maxMemory = res.memory_kb;\n")
	sb.WriteString("\n")
	sb.WriteString("        return {\n")
	sb.WriteString("            passed: res.status === \"passed\",\n")
	sb.WriteString("            input: JSON.stringify(TEST_CASES[i].input),\n")
	sb.WriteString("            actual: res.output,\n")
	sb.WriteString("            error: res.error\n")
	sb.WriteString("        };\n")
	sb.WriteString("    });\n")
	sb.WriteString("\n")
	sb.WriteString("    console.log(JSON.stringify({\n")
	sb.WriteString("        verdict: finalVerdict,\n")
	sb.WriteString("        runtime: maxRuntime,\n")
	sb.WriteString("        memory: maxMemory,\n")
	sb.WriteString("        test_results: testResults\n")
	sb.WriteString("    }));\n")
	sb.WriteString("    process.exit(0);\n")
	sb.WriteString("})();\n")
	return sb.String(), nil
}

// Java stub generator
func (s *CodeGenService) generateJavaStub(sig domain.ProblemSchema) (string, error) {
	var sb strings.Builder
	sb.WriteString("public class Solution {\n")
	sb.WriteString(fmt.Sprintf("    public %s %s(", s.mapTypeToJava(sig.ReturnType), sig.FunctionName))
	params := []string{}
	for _, param := range sig.Parameters {
		params = append(params, fmt.Sprintf("%s %s", s.mapTypeToJava(param.Type), param.Name))
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") {\n")
	sb.WriteString("        // Write your code here\n")
	sb.WriteString("        return null;\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}

// Java harness generator (Library-free version)
func (s *CodeGenService) GenerateJavaHarness(sig domain.ProblemSchema, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder
	sb.WriteString("import java.util.*;\n")
	sb.WriteString("import java.util.concurrent.*;\n\n")

	// Start with the harness class (public class Solution with main method)
	sb.WriteString("public class Solution {\n")

	// Manual JSON helpers for Java
	sb.WriteString("    private static String escapeJSON(String s) {\n")
	sb.WriteString("        if (s == null) return \"null\";\n")
	sb.WriteString("        StringBuilder sb = new StringBuilder();\n")
	sb.WriteString("        for (char c : s.toCharArray()) {\n")
	sb.WriteString("            if (c == '\"') sb.append(\"\\\\\\\"\");\n")
	sb.WriteString("            else if (c == '\\\\') sb.append(\"\\\\\\\\\");\n")
	sb.WriteString("            else sb.append(c);\n")
	sb.WriteString("        }\n")
	sb.WriteString("        return sb.toString();\n")
	sb.WriteString("    }\n\n")

	sb.WriteString("    private static String toJson(Object v) {\n")
	sb.WriteString("        if (v == null) return \"null\";\n")
	sb.WriteString("        if (v instanceof String) return \"\\\"\" + escapeJSON((String)v) + \"\\\"\";\n")
	sb.WriteString("        if (v instanceof Integer || v instanceof Long || v instanceof Boolean) return v.toString();\n")
	sb.WriteString("        if (v instanceof int[]) {\n")
	sb.WriteString("            int[] arr = (int[])v;\n")
	sb.WriteString("            StringBuilder sb = new StringBuilder(\"[\");\n")
	sb.WriteString("            for (int i = 0; i < arr.length; i++) {\n")
	sb.WriteString("                sb.append(arr[i]);\n")
	sb.WriteString("                if (i < arr.length - 1) sb.append(\",\");\n")
	sb.WriteString("            }\n")
	sb.WriteString("            sb.append(\"]\");\n")
	sb.WriteString("            return sb.toString();\n")
	sb.WriteString("        }\n")
	sb.WriteString("        if (v instanceof String[]) {\n")
	sb.WriteString("            String[] arr = (String[])v;\n")
	sb.WriteString("            StringBuilder sb = new StringBuilder(\"[\");\n")
	sb.WriteString("            for (int i = 0; i < arr.length; i++) {\n")
	sb.WriteString("                sb.append(\"\\\"\").append(escapeJSON(arr[i])).append(\"\\\"\");\n")
	sb.WriteString("                if (i < arr.length - 1) sb.append(\",\");\n")
	sb.WriteString("            }\n")
	sb.WriteString("            sb.append(\"]\");\n")
	sb.WriteString("            return sb.toString();\n")
	sb.WriteString("        }\n")
	sb.WriteString("        return \"null\";\n")
	sb.WriteString("    }\n\n")

	// Validation Helper
	sb.WriteString("    private static boolean compareOutputs(Object actual, Object expected, String valType) {\n")
	sb.WriteString("        if (valType.equals(\"UNORDERED\")) {\n")
	sb.WriteString("            if (actual instanceof int[] && expected instanceof int[]) {\n")
	sb.WriteString("                int[] a = ((int[])actual).clone(); int[] e = ((int[])expected).clone();\n")
	sb.WriteString("                Arrays.sort(a); Arrays.sort(e);\n")
	sb.WriteString("                return Arrays.equals(a, e);\n")
	sb.WriteString("            }\n")
	sb.WriteString("            if (actual instanceof String[] && expected instanceof String[]) {\n")
	sb.WriteString("                String[] a = ((String[])actual).clone(); String[] e = ((String[])expected).clone();\n")
	sb.WriteString("                Arrays.sort(a); Arrays.sort(e);\n")
	sb.WriteString("                return Arrays.equals(a, e);\n")
	sb.WriteString("            }\n")
	sb.WriteString("        }\n")
	sb.WriteString("        if (actual instanceof int[]) return Arrays.equals((int[])actual, (int[])expected);\n")
	sb.WriteString("        if (actual instanceof String[]) return Arrays.equals((String[])actual, (String[])expected);\n")
	sb.WriteString("        return Objects.equals(actual, expected);\n")
	sb.WriteString("    }\n\n")

	sb.WriteString("    public static void main(String[] args) {\n")
	sb.WriteString("        List<Map<String, Object>> results = new ArrayList<>();\n")
	sb.WriteString("        UserSolution sol = new UserSolution();\n")
	sb.WriteString("        ExecutorService executor = Executors.newSingleThreadExecutor();\n\n")

	for _, tc := range testCases {
		var inputArr []interface{}
		json.Unmarshal([]byte(tc.Input), &inputArr)

		sb.WriteString("        {\n")
		sb.WriteString("            String status = \"passed\";\n")
		sb.WriteString("            String output_val = \"\";\n")
		sb.WriteString("            String error_msg = \"\";\n")
		sb.WriteString("            String input_desc = \"\";\n")

		// Assign inputs
		paramNames := []string{}
		for j, param := range sig.Parameters {
			val := interface{}(nil)
			if j < len(inputArr) {
				val = inputArr[j]
			}
			literal := s.formatJavaLiteral(val, param.Type)
			javaType := s.mapTypeToJava(param.Type)
			sb.WriteString(fmt.Sprintf("            %s %s = %s;\n", javaType, param.Name, literal))
			paramNames = append(paramNames, param.Name)
			sb.WriteString(fmt.Sprintf("            input_desc += (input_desc.isEmpty() ? \"\" : \", \") + toJson(%s);\n", param.Name))
		}

		// Expected
		var expected interface{}
		json.Unmarshal([]byte(tc.ExpectedOutput), &expected)
		expectedLiteral := s.formatJavaLiteral(expected, sig.ReturnType)
		sb.WriteString(fmt.Sprintf("            %s expected = %s;\n", s.mapTypeToJava(sig.ReturnType), expectedLiteral))

		sb.WriteString("            long startTime = System.nanoTime();\n")
		sb.WriteString("            long startMem = Runtime.getRuntime().totalMemory() - Runtime.getRuntime().freeMemory();\n")

		sb.WriteString("            Future<Object> future = executor.submit(() -> {\n")
		sb.WriteString(fmt.Sprintf("                return sol.%s(%s);\n", sig.FunctionName, strings.Join(paramNames, ", ")))
		sb.WriteString("            });\n")

		sb.WriteString("            try {\n")
		sb.WriteString("                Object res = future.get(2, TimeUnit.SECONDS);\n")
		sb.WriteString("                output_val = toJson(res);\n")
		sb.WriteString("                if (!compareOutputs(res, expected, \"" + validationType + "\")) status = \"failed\";\n")
		sb.WriteString("            } catch (TimeoutException e) {\n")
		sb.WriteString("                status = \"timeout\";\n")
		sb.WriteString("                future.cancel(true);\n")
		sb.WriteString("            } catch (Exception e) {\n")
		sb.WriteString("                status = \"runtime_error\";\n")
		sb.WriteString("                error_msg = e.getMessage() != null ? e.getMessage() : e.toString();\n")
		sb.WriteString("            }\n")

		sb.WriteString("            long endTime = System.nanoTime();\n")
		sb.WriteString("            long endMem = Runtime.getRuntime().totalMemory() - Runtime.getRuntime().freeMemory();\n")
		sb.WriteString("            long timeMs = (endTime - startTime) / 1000000;\n")
		sb.WriteString("            long memKb = Math.max(0, (endMem - startMem) / 1024);\n")

		sb.WriteString("            Map<String, Object> r = new HashMap<>();\n")
		sb.WriteString("            r.put(\"status\", status); r.put(\"time_ms\", timeMs); r.put(\"memory_kb\", memKb);\n")
		sb.WriteString("            r.put(\"output\", output_val); r.put(\"error\", error_msg); r.put(\"input\", \"[\" + input_desc + \"]\");\n")
		sb.WriteString("            results.add(r);\n")
		sb.WriteString("        }\n")
	}

	sb.WriteString("\n")
	sb.WriteString("        // Standardized Verdict Aggregation\n")
	sb.WriteString("        String final_verdict = \"ACCEPTED\";\n")
	sb.WriteString("        long max_runtime = 0; long max_memory = 0;\n")
	sb.WriteString("        for (Map<String, Object> r : results) {\n")
	sb.WriteString("            String status = (String)r.get(\"status\");\n")
	sb.WriteString("            if (!status.equals(\"passed\") && final_verdict.equals(\"ACCEPTED\")) {\n")
	sb.WriteString("                if (status.equals(\"timeout\")) final_verdict = \"TLE\";\n")
	sb.WriteString("                else if (status.equals(\"runtime_error\")) final_verdict = \"RUNTIME_ERROR\";\n")
	sb.WriteString("                else if (status.equals(\"failed\")) final_verdict = \"WRONG_ANSWER\";\n")
	sb.WriteString("                else final_verdict = status.toUpperCase();\n")
	sb.WriteString("            }\n")
	sb.WriteString("            max_runtime = Math.max(max_runtime, (long)r.get(\"time_ms\"));\n")
	sb.WriteString("            max_memory = Math.max(max_memory, (long)r.get(\"memory_kb\"));\n")
	sb.WriteString("        }\n\n")

	sb.WriteString("        // Manual JSON output\n")
	sb.WriteString("        System.out.print(\"{\");\n")
	sb.WriteString("        System.out.print(\"\\\"verdict\\\":\\\"\" + final_verdict + \"\\\",\");\n")
	sb.WriteString("        System.out.print(\"\\\"runtime\\\":\" + max_runtime + \",\");\n")
	sb.WriteString("        System.out.print(\"\\\"memory\\\":\" + max_memory + \",\");\n")
	sb.WriteString("        System.out.print(\"\\\"test_results\\\":[\");\n")
	sb.WriteString("        for (int i = 0; i < results.size(); i++) {\n")
	sb.WriteString("            Map<String, Object> r = results.get(i);\n")
	sb.WriteString("            System.out.print(\"{\");\n")
	sb.WriteString("            System.out.print(\"\\\"passed\\\":\" + r.get(\"status\").equals(\"passed\") + \",\");\n")
	sb.WriteString("            System.out.print(\"\\\"input\\\":\\\"\" + escapeJSON((String)r.get(\"input\")) + \"\\\",\");\n")
	sb.WriteString("            System.out.print(\"\\\"actual\\\":\\\"\" + escapeJSON((String)r.get(\"output\")) + \"\\\",\");\n")
	sb.WriteString("            System.out.print(\"\\\"error\\\":\\\"\" + escapeJSON((String)r.get(\"error\")) + \"\\\"\");\n")
	sb.WriteString("            System.out.print(\"}\");\n")
	sb.WriteString("            if (i < results.size() - 1) System.out.print(\",\");\n")
	sb.WriteString("        }\n")
	sb.WriteString("        System.out.println(\"]}\");\n")
	sb.WriteString("        executor.shutdownNow();\n")
	sb.WriteString("        System.exit(0);\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n\n")

	// Append user code at the end (renamed to UserSolution)
	re := regexp.MustCompile(`\bSolution\b`)
	userCodeData := re.ReplaceAllString(userCode, "UserSolution")
	// Also ensure we don't have "public class UserSolution" (only one public class allowed)
	rePublic := regexp.MustCompile(`public\s+class\s+UserSolution\b`)
	userCodeData = rePublic.ReplaceAllString(userCodeData, "class UserSolution")
	sb.WriteString(userCodeData)
	sb.WriteString("\n")

	return sb.String(), nil
}

// C++ stub generator
func (s *CodeGenService) generateCppStub(sig domain.ProblemSchema) (string, error) {
	var sb strings.Builder
	sb.WriteString("#include <iostream>\n#include <vector>\n#include <string>\n\nusing namespace std;\n\n")
	sb.WriteString("class Solution {\npublic:\n")
	sb.WriteString(fmt.Sprintf("    %s %s(", s.mapTypeToCpp(sig.ReturnType), sig.FunctionName))
	params := []string{}
	for _, param := range sig.Parameters {
		params = append(params, fmt.Sprintf("%s %s", s.mapTypeToCpp(param.Type), param.Name))
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") {\n")
	sb.WriteString("        // Write your code here\n")
	sb.WriteString("    }\n")
	sb.WriteString("};\n")
	return sb.String(), nil
}

// C++ harness generator (Library-free version)
func (s *CodeGenService) GenerateCppHarness(sig domain.ProblemSchema, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder
	sb.WriteString("#include <iostream>\n#include <vector>\n#include <string>\n#include <chrono>\n#include <sys/resource.h>\n#include <signal.h>\n#include <setjmp.h>\n#include <algorithm>\n#include <sstream>\n#include <unistd.h>\n\n")
	sb.WriteString("using namespace std;\n\n")
	sb.WriteString("jmp_buf jump_buffer;\n")
	sb.WriteString("void timeout_handler(int sig) { longjmp(jump_buffer, 1); }\n\n")

	// Pre-declarations for manual JSON building
	sb.WriteString("// Manual JSON serialization helpers\n")
	sb.WriteString("string escapeJSON(string s) {\n")
	sb.WriteString("    string res = \"\";\n")
	sb.WriteString("    for (char c : s) {\n")
	sb.WriteString("        if (c == '\"') res += \"\\\\\\\"\";\n")
	sb.WriteString("        else if (c == '\\\\') res += \"\\\\\\\\\";\n")
	sb.WriteString("        else res += c;\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return res;\n")
	sb.WriteString("}\n\n")

	sb.WriteString("string toJson(int v) { return to_string(v); }\n")
	sb.WriteString("string toJson(long v) { return to_string(v); }\n")
	sb.WriteString("string toJson(long long v) { return to_string(v); }\n")
	sb.WriteString("string toJson(bool v) { return v ? \"true\" : \"false\"; }\n")
	sb.WriteString("string toJson(string v) { return \"\\\"\" + escapeJSON(v) + \"\\\"\"; }\n\n")

	sb.WriteString("template<typename T>\n")
	sb.WriteString("string toJson(const vector<T>& v) {\n")
	sb.WriteString("    string res = \"[\";\n")
	sb.WriteString("    for (size_t i = 0; i < v.size(); ++i) {\n")
	sb.WriteString("        res += toJson(v[i]);\n")
	sb.WriteString("        if (i < v.size() - 1) res += \",\";\n")
	sb.WriteString("    }\n")
	sb.WriteString("    res += \"]\";\n")
	sb.WriteString("    return res;\n")
	sb.WriteString("}\n\n")

	sb.WriteString(userCode)
	sb.WriteString("\n\n")

	sb.WriteString("struct TestResult {\n")
	sb.WriteString("    string status;\n")
	sb.WriteString("    long time_ms;\n")
	sb.WriteString("    long memory_kb;\n")
	sb.WriteString("    string output;\n")
	sb.WriteString("    string error;\n")
	sb.WriteString("    string input_description;\n")
	sb.WriteString("};\n\n")

	sb.WriteString("int main() {\n")
	sb.WriteString("    vector<TestResult> results;\n")
	sb.WriteString("    Solution sol;\n")
	sb.WriteString("    signal(SIGALRM, timeout_handler);\n\n")

	for i, tc := range testCases {
		var inputArr []interface{}
		json.Unmarshal([]byte(tc.Input), &inputArr)

		sb.WriteString(fmt.Sprintf("    // Test Case %d\n", i+1))
		sb.WriteString("    {\n")
		sb.WriteString("        string status = \"passed\";\n")
		sb.WriteString("        string output_val = \"\";\n")
		sb.WriteString("        string error_msg = \"\";\n")
		sb.WriteString("        string input_desc = \"\";\n")

		// Assign inputs
		paramNames := []string{}
		for j, param := range sig.Parameters {
			val := interface{}(nil)
			if j < len(inputArr) {
				val = inputArr[j]
			}
			literal := s.formatCppLiteral(val, param.Type)
			cppType := s.mapTypeToCpp(param.Type)
			sb.WriteString(fmt.Sprintf("        %s %s = %s;\n", cppType, param.Name, literal))
			paramNames = append(paramNames, param.Name)
			sb.WriteString(fmt.Sprintf("        input_desc += (input_desc.empty() ? \"\" : \", \") + toJson(%s);\n", param.Name))
		}

		// Expected val
		var expected interface{}
		json.Unmarshal([]byte(tc.ExpectedOutput), &expected)
		expectedLiteral := s.formatCppLiteral(expected, sig.ReturnType)
		sb.WriteString(fmt.Sprintf("        %s expected = %s;\n", s.mapTypeToCpp(sig.ReturnType), expectedLiteral))

		sb.WriteString("        auto start_time = chrono::high_resolution_clock::now();\n")
		sb.WriteString("        struct rusage usage_start, usage_end;\n")
		sb.WriteString("        getrusage(RUSAGE_SELF, &usage_start);\n\n")

		sb.WriteString("        alarm(2);\n")
		sb.WriteString("        if (setjmp(jump_buffer) == 0) {\n")
		sb.WriteString("            try {\n")
		sb.WriteString(fmt.Sprintf("                auto res = sol.%s(%s);\n", sig.FunctionName, strings.Join(paramNames, ", ")))
		sb.WriteString("                output_val = toJson(res);\n")

		if validationType == "UNORDERED" && (sig.ReturnType == domain.TypeIntegerArray || sig.ReturnType == domain.TypeStringArray) {
			sb.WriteString("                auto actual_sorted = res; auto expected_sorted = expected;\n")
			sb.WriteString("                sort(actual_sorted.begin(), actual_sorted.end()); sort(expected_sorted.begin(), expected_sorted.end());\n")
			sb.WriteString("                if (actual_sorted != expected_sorted) status = \"failed\";\n")
		} else {
			sb.WriteString("                if (res != expected) status = \"failed\";\n")
		}

		sb.WriteString("            } catch (const exception& e) {\n")
		sb.WriteString("                status = \"runtime_error\";\n")
		sb.WriteString("                error_msg = e.what();\n")
		sb.WriteString("            } catch (...) {\n")
		sb.WriteString("                status = \"runtime_error\";\n")
		sb.WriteString("                error_msg = \"Unknown error\";\n")
		sb.WriteString("            }\n")
		sb.WriteString("            alarm(0);\n")
		sb.WriteString("        } else {\n")
		sb.WriteString("            status = \"timeout\";\n")
		sb.WriteString("        }\n\n")

		sb.WriteString("        auto end_time = chrono::high_resolution_clock::now();\n")
		sb.WriteString("        getrusage(RUSAGE_SELF, &usage_end);\n")
		sb.WriteString("        auto time_ms = chrono::duration_cast<chrono::milliseconds>(end_time - start_time).count();\n")
		sb.WriteString("        auto memory_kb = usage_end.ru_maxrss;\n\n")

		sb.WriteString("        results.push_back({status, (long)time_ms, (long)memory_kb, output_val, error_msg, \"[\" + input_desc + \"]\"});\n")
		sb.WriteString("    }\n\n")
	}

	// Final Aggregation
	sb.WriteString("    // Standardized Verdict Aggregation\n")
	sb.WriteString("    string final_verdict = \"ACCEPTED\";\n")
	sb.WriteString("    long max_runtime = 0;\n")
	sb.WriteString("    long max_memory = 0;\n")
	sb.WriteString("    \n")
	sb.WriteString("    for (const auto& res : results) {\n")
	sb.WriteString("        if (res.status != \"passed\" && final_verdict == \"ACCEPTED\") {\n")
	sb.WriteString("            if (res.status == \"timeout\") final_verdict = \"TLE\";\n")
	sb.WriteString("            else if (res.status == \"runtime_error\") final_verdict = \"RUNTIME_ERROR\";\n")
	sb.WriteString("            else if (res.status == \"failed\") final_verdict = \"WRONG_ANSWER\";\n")
	sb.WriteString("            else final_verdict = res.status;\n")
	sb.WriteString("        }\n")
	sb.WriteString("        if (res.time_ms > max_runtime) max_runtime = res.time_ms;\n")
	sb.WriteString("        if (res.memory_kb > max_memory) max_memory = res.memory_kb;\n")
	sb.WriteString("    }\n\n")

	sb.WriteString("    // Output JSON manually\n")
	sb.WriteString("    cout << \"{\";\n")
	sb.WriteString("    cout << \"\\\"verdict\\\":\\\"\" + final_verdict + \"\\\",\";\n")
	sb.WriteString("    cout << \"\\\"runtime\\\":\" + to_string(max_runtime) + \",\";\n")
	sb.WriteString("    cout << \"\\\"memory\\\":\" + to_string(max_memory) + \",\";\n")
	sb.WriteString("    cout << \"\\\"test_results\\\":[\";\n")
	sb.WriteString("    for (size_t i = 0; i < results.size(); ++i) {\n")
	sb.WriteString("        cout << \"{\";\n")
	sb.WriteString("        cout << \"\\\"passed\\\":\" << (results[i].status == \"passed\" ? \"true\" : \"false\") << \",\";\n")
	sb.WriteString("        cout << \"\\\"input\\\":\\\"\" << escapeJSON(results[i].input_description) << \"\\\",\";\n")
	sb.WriteString("        cout << \"\\\"actual\\\":\\\"\" << escapeJSON(results[i].output) << \"\\\",\";\n")
	sb.WriteString("        cout << \"\\\"error\\\":\\\"\" << escapeJSON(results[i].error) << \"\\\"\";\n")
	sb.WriteString("        cout << \"}\";\n")
	sb.WriteString("        if (i < results.size() - 1) cout << \",\";\n")
	sb.WriteString("    }\n")
	sb.WriteString("    cout << \"]}\" << endl;\n")
	sb.WriteString("    return 0;\n")
	sb.WriteString("}\n")

	return sb.String(), nil
}

// C stub generator
func (s *CodeGenService) generateCStub(sig domain.ProblemSchema) (string, error) {
	var sb strings.Builder
	sb.WriteString("#include <stdio.h>\n#include <stdlib.h>\n#include <stdbool.h>\n#include <string.h>\n\n")
	sb.WriteString(fmt.Sprintf("%s %s(", s.mapTypeToC(sig.ReturnType), sig.FunctionName))
	params := []string{}
	for _, param := range sig.Parameters {
		params = append(params, fmt.Sprintf("%s %s", s.mapTypeToC(param.Type), param.Name))
		if param.Type == domain.TypeIntegerArray || param.Type == domain.TypeStringArray {
			params = append(params, fmt.Sprintf("int %sSize", param.Name))
		}
	}
	if sig.ReturnType == domain.TypeIntegerArray || sig.ReturnType == domain.TypeStringArray {
		params = append(params, "int* returnSize")
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") {\n")
	sb.WriteString("    // Write your code here\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}

// C harness generator (Library-free)
func (s *CodeGenService) GenerateCHarness(sig domain.ProblemSchema, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder
	sb.WriteString("#define _POSIX_C_SOURCE 199309L\n")
	sb.WriteString("#include <stdio.h>\n#include <stdlib.h>\n#include <stdbool.h>\n#include <string.h>\n#include <time.h>\n#include <sys/resource.h>\n#include <signal.h>\n#include <setjmp.h>\n#include <unistd.h>\n\n")

	sb.WriteString("jmp_buf jump_buffer;\n")
	sb.WriteString("void timeout_handler(int sig) { longjmp(jump_buffer, 1); }\n\n")

	// Manual JSON helpers for C
	sb.WriteString("// Manual JSON serialization helpers\n")
	sb.WriteString("void escapeJSON(const char* s, char* dest) {\n")
	sb.WriteString("    if (!s) { strcpy(dest, \"null\"); return; }\n")
	sb.WriteString("    while (*s) {\n")
	sb.WriteString("        if (*s == '\"') { *dest++ = '\\\\'; *dest++ = '\"'; }\n")
	sb.WriteString("        else if (*s == '\\\\') { *dest++ = '\\\\'; *dest++ = '\\\\'; }\n")
	sb.WriteString("        else *dest++ = *s;\n")
	sb.WriteString("        s++;\n")
	sb.WriteString("    }\n")
	sb.WriteString("    *dest = '\\0';\n")
	sb.WriteString("}\n\n")

	sb.WriteString(userCode)
	sb.WriteString("\n\n")

	sb.WriteString("typedef struct {\n")
	sb.WriteString("    char status[20];\n")
	sb.WriteString("    long time_ms;\n")
	sb.WriteString("    long memory_kb;\n")
	sb.WriteString("    char output[4096];\n")
	sb.WriteString("    char error[4096];\n")
	sb.WriteString("    char input_desc[4096];\n")
	sb.WriteString("} TestResult;\n\n")

	// Helper to JSONize arrays for the harness
	sb.WriteString("char* arrayToJson(int* arr, int size) {\n")
	sb.WriteString("    if (!arr) { char* r = malloc(5); strcpy(r, \"null\"); return r; }\n")
	sb.WriteString("    char* res = malloc(size * 12 + 2);\n")
	sb.WriteString("    strcpy(res, \"[\");\n")
	sb.WriteString("    for (int i = 0; i < size; i++) {\n")
	sb.WriteString("        char buf[12]; sprintf(buf, \"%d\", arr[i]);\n")
	sb.WriteString("        strcat(res, buf);\n")
	sb.WriteString("        if (i < size - 1) strcat(res, \",\");\n")
	sb.WriteString("    }\n")
	sb.WriteString("    strcat(res, \"]\");\n")
	sb.WriteString("    return res;\n")
	sb.WriteString("}\n\n")

	sb.WriteString("char* stringArrayToJson(char** arr, int size) {\n")
	sb.WriteString("    if (!arr) { char* r = malloc(5); strcpy(r, \"null\"); return r; }\n")
	sb.WriteString("    char* res = malloc(size * 1024 + 2);\n")
	sb.WriteString("    strcpy(res, \"[\");\n")
	sb.WriteString("    for (int i = 0; i < size; i++) {\n")
	sb.WriteString("        strcat(res, \"\\\"\");\n")
	sb.WriteString("        char buf[1024]; escapeJSON(arr[i], buf);\n")
	sb.WriteString("        strcat(res, buf);\n")
	sb.WriteString("        strcat(res, \"\\\"\");\n")
	sb.WriteString("        if (i < size - 1) strcat(res, \",\");\n")
	sb.WriteString("    }\n")
	sb.WriteString("    strcat(res, \"]\");\n")
	sb.WriteString("    return res;\n")
	sb.WriteString("}\n\n")

	sb.WriteString("int main() {\n")
	sb.WriteString("    TestResult results[100];\n")
	sb.WriteString("    int test_count = 0;\n")
	sb.WriteString("    signal(SIGALRM, timeout_handler);\n\n")

	for i, tc := range testCases {
		var inputArr []interface{}
		json.Unmarshal([]byte(tc.Input), &inputArr)

		sb.WriteString(fmt.Sprintf("    // Test Case %d\n", i+1))
		sb.WriteString("    {\n")
		sb.WriteString("        TestResult* r = &results[test_count++];\n")
		sb.WriteString("        strcpy(r->status, \"passed\");\n")
		sb.WriteString("        strcpy(r->output, \"\");\n")
		sb.WriteString("        strcpy(r->error, \"\");\n")
		sb.WriteString("        strcpy(r->input_desc, \"[\");\n")

		// Assign inputs
		paramNames := []string{}
		for j, param := range sig.Parameters {
			val := interface{}(nil)
			if j < len(inputArr) {
				val = inputArr[j]
			}
			literal := s.formatCLiteral(val, param.Type)
			cType := s.mapTypeToC(param.Type)

			if param.Type == domain.TypeIntegerArray {
				items, _ := val.([]interface{})
				count := len(items)
				sb.WriteString(fmt.Sprintf("        int %s_arr[] = %s;\n", param.Name, literal))
				sb.WriteString(fmt.Sprintf("        %s %s = %s_arr;\n", cType, param.Name, param.Name))
				sb.WriteString(fmt.Sprintf("        int %sSize = %d;\n", param.Name, count))
				sb.WriteString(fmt.Sprintf("        char* %s_json = arrayToJson(%s_arr, %d);\n", param.Name, param.Name, count))
				sb.WriteString(fmt.Sprintf("        strcat(r->input_desc, %s_json); free(%s_json);\n", param.Name, param.Name))
				paramNames = append(paramNames, param.Name, fmt.Sprintf("%sSize", param.Name))
			} else if param.Type == domain.TypeStringArray {
				items, _ := val.([]interface{})
				count := len(items)
				sb.WriteString(fmt.Sprintf("        char* %s_arr[] = %s;\n", param.Name, literal))
				sb.WriteString(fmt.Sprintf("        %s %s = %s_arr;\n", cType, param.Name, param.Name))
				sb.WriteString(fmt.Sprintf("        int %sSize = %d;\n", param.Name, count))
				sb.WriteString(fmt.Sprintf("        char* %s_json = stringArrayToJson(%s_arr, %d);\n", param.Name, param.Name, count))
				sb.WriteString(fmt.Sprintf("        strcat(r->input_desc, %s_json); free(%s_json);\n", param.Name, param.Name))
				paramNames = append(paramNames, param.Name, fmt.Sprintf("%sSize", param.Name))
			} else {
				sb.WriteString(fmt.Sprintf("        %s %s = %s;\n", cType, param.Name, literal))
				if param.Type == domain.TypeInteger {
					sb.WriteString(fmt.Sprintf("        char %s_buf[12]; sprintf(%s_buf, \"%%d\", %s); strcat(r->input_desc, %s_buf);\n", param.Name, param.Name, param.Name, param.Name))
				} else if param.Type == domain.TypeString {
					sb.WriteString(fmt.Sprintf("        strcat(r->input_desc, \"\\\"\"); strcat(r->input_desc, %s); strcat(r->input_desc, \"\\\"\");\n", param.Name))
				} else if param.Type == domain.TypeBoolean {
					sb.WriteString(fmt.Sprintf("        strcat(r->input_desc, %s ? \"true\" : \"false\");\n", param.Name))
				}
				paramNames = append(paramNames, param.Name)
			}
			if j < len(sig.Parameters)-1 {
				sb.WriteString("        strcat(r->input_desc, \", \");\n")
			}
		}
		sb.WriteString("        strcat(r->input_desc, \"]\");\n")

		// Expected
		var expected interface{}
		json.Unmarshal([]byte(tc.ExpectedOutput), &expected)
		expectedLiteral := s.formatCLiteral(expected, sig.ReturnType)
		if sig.ReturnType == domain.TypeIntegerArray {
			sb.WriteString(fmt.Sprintf("        int expected_arr[] = %s;\n", expectedLiteral))
			items, _ := expected.([]interface{})
			sb.WriteString(fmt.Sprintf("        int expected_size = %d;\n", len(items)))
			sb.WriteString("        int returnSize = 0;\n")
			paramNames = append(paramNames, "&returnSize")
		} else if sig.ReturnType == domain.TypeStringArray {
			sb.WriteString(fmt.Sprintf("        char* expected_arr[] = %s;\n", expectedLiteral))
			items, _ := expected.([]interface{})
			sb.WriteString(fmt.Sprintf("        int expected_size = %d;\n", len(items)))
			sb.WriteString("        int returnSize = 0;\n")
			paramNames = append(paramNames, "&returnSize")
		} else {
			sb.WriteString(fmt.Sprintf("        %s expected = %s;\n", s.mapTypeToC(sig.ReturnType), expectedLiteral))
		}

		sb.WriteString("        struct timespec start, end;\n")
		sb.WriteString("        clock_gettime(CLOCK_MONOTONIC, &start);\n")
		sb.WriteString("        struct rusage usage_start, usage_end;\n")
		sb.WriteString("        getrusage(RUSAGE_SELF, &usage_start);\n\n")

		sb.WriteString("        alarm(2);\n")
		sb.WriteString("        if (setjmp(jump_buffer) == 0) {\n")
		sb.WriteString(fmt.Sprintf("            %s res = %s(%s);\n", s.mapTypeToC(sig.ReturnType), sig.FunctionName, strings.Join(paramNames, ", ")))
		sb.WriteString("            alarm(0);\n")

		// Comparison
		if sig.ReturnType == domain.TypeIntegerArray {
			sb.WriteString("            bool match = (returnSize == expected_size);\n")
			sb.WriteString("            if (match) {\n")
			sb.WriteString("                for(int k=0; k<expected_size; k++) if(res[k] != expected_arr[k]) match = false;\n")
			sb.WriteString("            }\n")
			sb.WriteString("            if (!match) strcpy(r->status, \"failed\");\n")
			sb.WriteString("            char* res_json = arrayToJson(res, returnSize);\n")
			sb.WriteString("            strcpy(r->output, res_json); free(res_json);\n")
		} else if sig.ReturnType == domain.TypeStringArray {
			sb.WriteString("            bool match = (returnSize == expected_size);\n")
			sb.WriteString("            if (match) {\n")
			sb.WriteString("                for(int k=0; k<expected_size; k++) if(strcmp(res[k], expected_arr[k]) != 0) match = false;\n")
			sb.WriteString("            }\n")
			sb.WriteString("            if (!match) strcpy(r->status, \"failed\");\n")
			sb.WriteString("            char* res_json = stringArrayToJson(res, returnSize);\n")
			sb.WriteString("            strcpy(r->output, res_json); free(res_json);\n")
		} else if sig.ReturnType == domain.TypeInteger {
			sb.WriteString("            if (res != expected) strcpy(r->status, \"failed\");\n")
			sb.WriteString("            sprintf(r->output, \"%d\", res);\n")
		} else if sig.ReturnType == domain.TypeString {
			sb.WriteString("            if (strcmp(res, expected) != 0) strcpy(r->status, \"failed\");\n")
			sb.WriteString("            strcpy(r->output, res);\n")
		} else if sig.ReturnType == domain.TypeBoolean {
			sb.WriteString("            if (res != expected) strcpy(r->status, \"failed\");\n")
			sb.WriteString("            strcpy(r->output, res ? \"true\" : \"false\");\n")
		}

		sb.WriteString("        } else {\n")
		sb.WriteString("            strcpy(r->status, \"timeout\");\n")
		sb.WriteString("        }\n\n")

		sb.WriteString("        clock_gettime(CLOCK_MONOTONIC, &end);\n")
		sb.WriteString("        getrusage(RUSAGE_SELF, &usage_end);\n")
		sb.WriteString("        r->time_ms = (end.tv_sec - start.tv_sec) * 1000 + (end.tv_nsec - start.tv_nsec) / 1000000;\n")
		sb.WriteString("        r->memory_kb = usage_end.ru_maxrss;\n")
		sb.WriteString("    }\n\n")
	}

	// Final Aggregation
	sb.WriteString("    char final_verdict[20] = \"ACCEPTED\";\n")
	sb.WriteString("    long max_runtime = 0; long max_memory = 0;\n")
	sb.WriteString("    for (int i = 0; i < test_count; i++) {\n")
	sb.WriteString("        if (strcmp(results[i].status, \"passed\") != 0 && strcmp(final_verdict, \"ACCEPTED\") == 0) {\n")
	sb.WriteString("            if (strcmp(results[i].status, \"timeout\") == 0) strcpy(final_verdict, \"TLE\");\n")
	sb.WriteString("            else if (strcmp(results[i].status, \"failed\") == 0) strcpy(final_verdict, \"WRONG_ANSWER\");\n")
	sb.WriteString("            else strcpy(final_verdict, \"RUNTIME_ERROR\");\n")
	sb.WriteString("        }\n")
	sb.WriteString("        if (results[i].time_ms > max_runtime) max_runtime = results[i].time_ms;\n")
	sb.WriteString("        if (results[i].memory_kb > max_memory) max_memory = results[i].memory_kb;\n")
	sb.WriteString("    }\n\n")

	sb.WriteString("    printf(\"{\");\n")
	sb.WriteString("    printf(\"\\\"verdict\\\":\\\"%s\\\",\", final_verdict);\n")
	sb.WriteString("    printf(\"\\\"runtime\\\":%ld,\", max_runtime);\n")
	sb.WriteString("    printf(\"\\\"memory\\\":%ld,\", max_memory);\n")
	sb.WriteString("    printf(\"\\\"test_results\\\":[\");\n")
	sb.WriteString("    for (int i = 0; i < test_count; i++) {\n")
	sb.WriteString("        char escaped_input[8192], escaped_actual[8192], escaped_error[8192];\n")
	sb.WriteString("        escapeJSON(results[i].input_desc, escaped_input);\n")
	sb.WriteString("        escapeJSON(results[i].output, escaped_actual);\n")
	sb.WriteString("        escapeJSON(results[i].error, escaped_error);\n")
	sb.WriteString("        printf(\"{\");\n")
	sb.WriteString("        printf(\"\\\"passed\\\":%s,\", strcmp(results[i].status, \"passed\") == 0 ? \"true\" : \"false\");\n")
	sb.WriteString("        printf(\"\\\"input\\\":\\\"%s\\\",\", escaped_input);\n")
	sb.WriteString("        printf(\"\\\"actual\\\":\\\"%s\\\",\", escaped_actual);\n")
	sb.WriteString("        printf(\"\\\"error\\\":\\\"%s\\\"\", escaped_error);\n")
	sb.WriteString("        printf(\"}\");\n")
	sb.WriteString("        if (i < test_count - 1) printf(\",\");\n")
	sb.WriteString("    }\n")
	sb.WriteString("    printf(\"]}\\n\");\n")
	sb.WriteString("    return 0;\n")
	sb.WriteString("}\n")

	return sb.String(), nil
}

// Go stub generator
func (s *CodeGenService) generateGoStub(sig domain.ProblemSchema) (string, error) {
	var sb strings.Builder
	sb.WriteString("func ")
	sb.WriteString(fmt.Sprintf("%s(", sig.FunctionName))
	params := []string{}
	for _, param := range sig.Parameters {
		params = append(params, fmt.Sprintf("%s %s", param.Name, s.mapTypeToGo(param.Type)))
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(fmt.Sprintf(") %s {\n", s.mapTypeToGo(sig.ReturnType)))
	sb.WriteString("    // Write your code here\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}

// Go harness generator
func (s *CodeGenService) GenerateGoHarness(sig domain.ProblemSchema, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder
	sb.WriteString("package main\n\nimport (\n	\"encoding/json\"\n	\"fmt\"\n	\"time\"\n	\"runtime\"\n	\"context\"\n	\"reflect\"\n	\"strings\"\n	\"sort\"\n)\n\n")
	sb.WriteString(userCode)
	sb.WriteString("\n\n")

	// Embed Test Cases
	sb.WriteString("// Embedded Test Cases\n")
	sb.WriteString("var TEST_CASES_JSON = []byte(`")

	type EmbeddedTestCase struct {
		Input    interface{} `json:"input"`
		Expected interface{} `json:"expected"`
		IsSample bool        `json:"is_sample"`
	}

	embeddedCases := []EmbeddedTestCase{}
	for _, tc := range testCases {
		var input, expected interface{}
		json.Unmarshal([]byte(tc.Input), &input)
		json.Unmarshal([]byte(tc.ExpectedOutput), &expected)
		embeddedCases = append(embeddedCases, EmbeddedTestCase{
			Input:    input,
			Expected: expected,
			IsSample: tc.IsSample,
		})
	}
	casesJSON, _ := json.Marshal(embeddedCases)
	sb.WriteString(string(casesJSON))
	sb.WriteString("`)\n\n")

	// Validation Helper
	sb.WriteString("func compareOutputs(actual, expected interface{}, valType string) bool {\n")
	sb.WriteString("    switch valType {\n")
	sb.WriteString("    case \"EXACT\":\n")
	sb.WriteString("        return reflect.DeepEqual(actual, expected)\n")
	sb.WriteString("    case \"UNORDERED\":\n")
	sb.WriteString("        aList, ok1 := actual.([]interface{})\n")
	sb.WriteString("        eList, ok2 := expected.([]interface{})\n")
	sb.WriteString("        if !ok1 || !ok2 {\n")
	sb.WriteString("            return reflect.DeepEqual(actual, expected)\n")
	sb.WriteString("        }\n")
	sb.WriteString("        if len(aList) != len(eList) { return false }\n")
	sb.WriteString("        aStrList := make([]string, len(aList))\n")
	sb.WriteString("        eStrList := make([]string, len(eList))\n")
	sb.WriteString("        for i, v := range aList { b, _ := json.Marshal(v); aStrList[i] = string(b) }\n")
	sb.WriteString("        for i, v := range eList { b, _ := json.Marshal(v); eStrList[i] = string(b) }\n")
	sb.WriteString("        sort.Strings(aStrList)\n")
	sb.WriteString("        sort.Strings(eStrList)\n")
	sb.WriteString("        return reflect.DeepEqual(aStrList, eStrList)\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return reflect.DeepEqual(actual, expected)\n")
	sb.WriteString("}\n\n")

	sb.WriteString("func main() {\n")
	sb.WriteString("    var testCases []map[string]interface{}\n")
	sb.WriteString("    json.Unmarshal(TEST_CASES_JSON, &testCases)\n")
	sb.WriteString("    results := []map[string]interface{}{}\n")
	sb.WriteString(fmt.Sprintf("    validationType := \"%s\"\n\n", validationType))
	sb.WriteString("    for _, test := range testCases {\n")
	sb.WriteString("        status := \"passed\"\n")
	sb.WriteString("        var output interface{}\n")
	sb.WriteString("        var errStr string\n")
	sb.WriteString("        \n")
	sb.WriteString("        start := time.Now()\n")
	sb.WriteString("        var ms runtime.MemStats\n")
	sb.WriteString("        runtime.ReadMemStats(&ms)\n")
	sb.WriteString("        startAlloc := ms.TotalAlloc\n\n")
	sb.WriteString("        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)\n")
	sb.WriteString("        resChan := make(chan interface{}, 1)\n")
	sb.WriteString("        errChan := make(chan error, 1)\n\n")
	sb.WriteString("        go func() {\n")
	sb.WriteString("            defer func() {\n                if r := recover(); r != nil {\n                    errChan <- fmt.Errorf(\"%v\", r)\n                }\n            }()\n")
	sb.WriteString("            input := test[\"input\"].([]interface{})\n")
	// Casting each input to its expected Go type
	for j, param := range sig.Parameters {
		typeName := s.mapTypeToGo(param.Type)
		// Correcting for the fact that json.Unmarshal uses float64 for all numbers
		if typeName == "int" {
			sb.WriteString(fmt.Sprintf("            %s := int(input[%d].(float64))\n", param.Name, j))
		} else if typeName == "[]int" {
			sb.WriteString(fmt.Sprintf("            %sRaw := input[%d].([]interface{})\n", param.Name, j))
			sb.WriteString(fmt.Sprintf("            %s := make([]int, len(%sRaw))\n", param.Name, param.Name))
			sb.WriteString(fmt.Sprintf("            for k, v := range %sRaw { %s[k] = int(v.(float64)) }\n", param.Name, param.Name))
		} else {
			sb.WriteString(fmt.Sprintf("            %s := input[%d].(%s)\n", param.Name, j, typeName))
		}
	}
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("            resChan <- %s(%s)\n", sig.FunctionName, strings.Join(paramNames, ", ")))
	sb.WriteString("        }()\n\n")
	sb.WriteString("        select {\n")
	sb.WriteString("        case res := <-resChan:\n            ints := res.([]int)\n")
	sb.WriteString("            normalized := make([]interface{}, len(ints))\n")
	sb.WriteString("            for i, v := range ints { normalized[i] = float64(v) }\n")
	sb.WriteString("            output = normalized\n")
	sb.WriteString("            if !compareOutputs(output, test[\"expected\"], validationType) {\n")
	sb.WriteString("                status = \"failed\"\n")
	sb.WriteString("            }\n")
	sb.WriteString("        case err := <-errChan:\n            status = \"runtime_error\"\n            errStr = err.Error()\n")
	sb.WriteString("        case <-ctx.Done():\n            status = \"timeout\"\n")
	sb.WriteString("        }\n")
	sb.WriteString("        cancel()\n\n")
	sb.WriteString("        duration := time.Since(start)\n")
	sb.WriteString("        runtime.ReadMemStats(&ms)\n")
	sb.WriteString("        memKb := int64((ms.TotalAlloc - startAlloc) / 1024)\n")
	// sb.WriteString("        if memKb < 0 { memKb = 0 }\n\n")
	sb.WriteString("        outStr, _ := json.Marshal(output)\n")
	// sb.WriteString("        expStr, _ := json.Marshal(test[\"expected\"])\n")
	sb.WriteString("        results = append(results, map[string]interface{}{\n")
	sb.WriteString("            \"status\": status,\n")
	sb.WriteString("            \"time_ms\": duration.Milliseconds(),\n")
	sb.WriteString("            \"memory_kb\": memKb,\n")
	sb.WriteString("            \"output\": string(outStr),\n")
	sb.WriteString("            \"error\": errStr,\n")
	sb.WriteString("        })\n")
	sb.WriteString("    }\n\n")
	sb.WriteString("    // Standardized Verdict Aggregation\n")
	sb.WriteString("    finalVerdict := \"ACCEPTED\"\n")
	sb.WriteString("    var maxRuntime int64\n")
	sb.WriteString("    var maxMemory int64\n")
	sb.WriteString("    testResults := make([]map[string]interface{}, len(results))\n\n")
	sb.WriteString("    for i, res := range results {\n")
	sb.WriteString("        status := res[\"status\"].(string)\n")
	sb.WriteString("        if status != \"passed\" && finalVerdict == \"ACCEPTED\" {\n")
	sb.WriteString("            if status == \"timeout\" {\n")
	sb.WriteString("                finalVerdict = \"TLE\"\n")
	sb.WriteString("            } else if status == \"runtime_error\" {\n")
	sb.WriteString("                finalVerdict = \"RUNTIME_ERROR\"\n")
	sb.WriteString("            } else if status == \"failed\" {\n")
	sb.WriteString("                finalVerdict = \"WRONG_ANSWER\"\n")
	sb.WriteString("                \n")
	sb.WriteString("            } else {\n")
	sb.WriteString("                finalVerdict = strings.ToUpper(status)\n")
	sb.WriteString("            }\n")
	sb.WriteString("        }\n")
	sb.WriteString("        if res[\"time_ms\"].(int64) > maxRuntime { maxRuntime = res[\"time_ms\"].(int64) }\n")
	sb.WriteString("        if res[\"memory_kb\"].(int64) > maxMemory { maxMemory = res[\"memory_kb\"].(int64) }\n\n")
	sb.WriteString("        inStr, _ := json.Marshal(testCases[i][\"input\"])\n")
	sb.WriteString("        testResults[i] = map[string]interface{}{\n")
	sb.WriteString("            \"passed\": status == \"passed\",\n")
	sb.WriteString("            \"input\": string(inStr),\n")
	sb.WriteString("            \"actual\": res[\"output\"],\n")
	sb.WriteString("            \"error\": res[\"error\"],\n")
	sb.WriteString("        }\n")
	sb.WriteString("    }\n\n")
	sb.WriteString("    verdictObj := map[string]interface{}{\n")
	sb.WriteString("        \"verdict\": finalVerdict,\n")
	sb.WriteString("        \"runtime\": maxRuntime,\n")
	sb.WriteString("        \"memory\": maxMemory,\n")
	sb.WriteString("        \"test_results\": testResults,\n")
	sb.WriteString("    }\n")
	sb.WriteString("    finalData, _ := json.Marshal(verdictObj)\n")
	sb.WriteString("    fmt.Println(string(finalData))\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}
