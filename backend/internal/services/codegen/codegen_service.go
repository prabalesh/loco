package codegen

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain"
)

type Parameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsCustom bool   `json:"is_custom"`
}

type ProblemSignature struct {
	FunctionName string      `json:"function_name"`
	ReturnType   string      `json:"return_type"`
	Parameters   []Parameter `json:"parameters"`
}

type CodeGenService struct {
	typeImplRepo domain.TypeImplementationRepository
}

func NewCodeGenService(typeImplRepo domain.TypeImplementationRepository) *CodeGenService {
	return &CodeGenService{
		typeImplRepo: typeImplRepo,
	}
}

// GenerateStubCode generates starter code for users
func (s *CodeGenService) GenerateStubCode(signature ProblemSignature, languageSlug string) (string, error) {
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
	case "cpp":
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for C++ yet")
		}
		return s.generateCppStub(signature)
	case "go":
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for Go yet")
		}
		return s.generateGoStub(signature)
	default:
		return "", fmt.Errorf("unsupported language: %s", languageSlug)
	}
}

func (s *CodeGenService) identifyCustomTypes(sig ProblemSignature) []string {
	types := []string{}
	seen := make(map[string]bool)

	if s.isCustomType(sig.ReturnType) {
		types = append(types, sig.ReturnType)
		seen[sig.ReturnType] = true
	}

	for _, param := range sig.Parameters {
		if param.IsCustom && !seen[param.Type] {
			types = append(types, param.Type)
			seen[param.Type] = true
		}
	}

	return types
}

func (s *CodeGenService) isCustomType(typeName string) bool {
	// This could be dynamic, but checking seeded types for now
	validTypes := []string{"TreeNode", "ListNode", "GraphNode", "Node"}
	for _, vt := range validTypes {
		if typeName == vt {
			return true
		}
	}
	return false
}

// GenerateTestHarness generates wrapper code that runs test cases
func (s *CodeGenService) GenerateTestHarness(signature ProblemSignature, userCode string, languageSlug string, testCases []domain.TestCase, validationType string) (string, error) {
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
	case "cpp":
		return s.GenerateCppHarness(signature, userCode, testCases, validationType)
	case "go":
		return s.GenerateGoHarness(signature, userCode, testCases, validationType)
	default:
		return "", fmt.Errorf("unsupported language: %s", languageSlug)
	}
}

// Type mapping helpers
func (s *CodeGenService) mapTypeToPython(typ string) string {
	typeMap := map[string]string{
		"int":      "int",
		"int[]":    "List[int]",
		"string":   "str",
		"string[]": "List[str]",
		"bool":     "bool",
		"double":   "float",
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return typ
}

func (s *CodeGenService) mapTypeToJavaScript(typ string) string {
	return ""
}

func (s *CodeGenService) mapTypeToJava(typ string) string {
	typeMap := map[string]string{
		"int":      "int",
		"int[]":    "int[]",
		"string":   "String",
		"string[]": "String[]",
		"bool":     "boolean",
		"double":   "double",
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return typ
}

func (s *CodeGenService) mapTypeToCpp(typ string) string {
	typeMap := map[string]string{
		"int":      "int",
		"int[]":    "vector<int>",
		"string":   "string",
		"string[]": "vector<string>",
		"bool":     "bool",
		"double":   "double",
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return typ
}

func (s *CodeGenService) mapTypeToGo(typ string) string {
	typeMap := map[string]string{
		"int":      "int",
		"int[]":    "[]int",
		"string":   "string",
		"string[]": "[]string",
		"bool":     "bool",
		"double":   "float64",
	}
	if mapped, ok := typeMap[typ]; ok {
		return mapped
	}
	return typ
}

// Python stub generator
func (s *CodeGenService) generatePythonStub(sig ProblemSignature, customTypes []string) (string, error) {
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
func (s *CodeGenService) GeneratePythonHarness(sig ProblemSignature, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
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
	sb.WriteString(string(casesJSON))
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
			deserializerFunc := fmt.Sprintf("deserialize_%s", strings.ToLower(param.Type))
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
		serializerFunc := fmt.Sprintf("serialize_%s", strings.ToLower(sig.ReturnType))
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
	sb.WriteString("            \"test_id\": i + 1,\n")
	sb.WriteString("            \"status\": status,\n")
	sb.WriteString("            \"time_ms\": time_ms,\n")
	sb.WriteString("            \"memory_kb\": memory_kb,\n")
	sb.WriteString("            \"output\": json.dumps(output, sort_keys=True) if output is not None else \"\",\n")
	sb.WriteString("            \"expected\": json.dumps(test['expected'], sort_keys=True),\n")
	sb.WriteString("            \"is_sample\": test['is_sample'],\n")
	sb.WriteString("            \"error\": error\n")
	sb.WriteString("        })\n")
	sb.WriteString("    \n")
	sb.WriteString("    print(json.dumps(results))\n")
	return sb.String(), nil
}

// JavaScript stub generator
func (s *CodeGenService) generateJavaScriptStub(sig ProblemSignature, customTypes []string) (string, error) {
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
func (s *CodeGenService) GenerateJavaScriptHarness(sig ProblemSignature, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
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
	sb.WriteString("            test_id: i + 1,\n")
	sb.WriteString("            status: status,\n")
	sb.WriteString("            time_ms: timeMs,\n")
	sb.WriteString("            memory_kb: memoryKb,\n")
	sb.WriteString("            output: output !== null ? JSON.stringify(output) : \"\",\n")
	sb.WriteString("            expected: JSON.stringify(test.expected),\n")
	sb.WriteString("            is_sample: test.is_sample,\n")
	sb.WriteString("            error: error\n")
	sb.WriteString("        });\n")
	sb.WriteString("    }\n")
	sb.WriteString("    console.log(JSON.stringify(results));\n")
	sb.WriteString("    process.exit(0);\n")
	sb.WriteString("})();\n")
	return sb.String(), nil
}

// Java stub generator
func (s *CodeGenService) generateJavaStub(sig ProblemSignature) (string, error) {
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

// Java harness generator
func (s *CodeGenService) GenerateJavaHarness(sig ProblemSignature, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder
	sb.WriteString("import java.util.*;\n")
	sb.WriteString("import com.fasterxml.jackson.databind.ObjectMapper;\n")
	sb.WriteString("import com.fasterxml.jackson.databind.JsonNode;\n")
	sb.WriteString("import java.util.concurrent.*;\n\n")

	sb.WriteString(userCode)
	sb.WriteString("\n\n")

	sb.WriteString("public class Main {\n")

	// Validation Helper
	sb.WriteString("    private static boolean compareOutputs(Object actual, Object expected, String valType, ObjectMapper mapper) {\n")
	sb.WriteString("        try {\n")
	sb.WriteString("            JsonNode actualNode = mapper.valueToTree(actual);\n")
	sb.WriteString("            JsonNode expectedNode = mapper.valueToTree(expected);\n")
	sb.WriteString("            if (valType.equals(\"EXACT\")) {\n")
	sb.WriteString("                return actualNode.equals(expectedNode);\n")
	sb.WriteString("            } else if (valType.equals(\"UNORDERED\")) {\n")
	sb.WriteString("                if (!actualNode.isArray() || !expectedNode.isArray()) return actualNode.equals(expectedNode);\n")
	sb.WriteString("                List<String> actualList = new ArrayList<>();\n")
	sb.WriteString("                actualNode.forEach(n -> actualList.add(n.toString()));\n")
	sb.WriteString("                List<String> expectedList = new ArrayList<>();\n")
	sb.WriteString("                expectedNode.forEach(n -> expectedList.add(n.toString()));\n")
	sb.WriteString("                Collections.sort(actualList);\n")
	sb.WriteString("                Collections.sort(expectedList);\n")
	sb.WriteString("                return actualList.equals(expectedList);\n")
	sb.WriteString("            }\n")
	sb.WriteString("            return actualNode.equals(expectedNode);\n")
	sb.WriteString("        } catch (Exception e) { return false; }\n")
	sb.WriteString("    }\n\n")

	sb.WriteString("    public static void main(String[] args) throws Exception {\n")
	sb.WriteString("        ObjectMapper mapper = new ObjectMapper();\n")

	// Embed Test Cases
	sb.WriteString("        String testCasesJson = \"")

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
	// Escape double quotes for Java string
	escapedJSON := strings.ReplaceAll(string(casesJSON), "\"", "\\\"")
	sb.WriteString(escapedJSON)
	sb.WriteString("\";\n\n")

	sb.WriteString("        List<Map<String, Object>> testCases = mapper.readValue(testCasesJson, List.class);\n")
	sb.WriteString("        List<Map<String, Object>> results = new ArrayList<>();\n")
	sb.WriteString("        Solution sol = new Solution();\n")
	sb.WriteString(fmt.Sprintf("        String validationType = \"%s\";\n", validationType))
	sb.WriteString("        ExecutorService executor = Executors.newSingleThreadExecutor();\n")
	sb.WriteString("        \n")
	sb.WriteString("        for (int i = 0; i < testCases.size(); i++) {\n")
	sb.WriteString("            Map<String, Object> test = testCases.get(i);\n")
	sb.WriteString("            final int index = i;\n")
	sb.WriteString("            String status = \"passed\";\n")
	sb.WriteString("            Object output = null;\n")
	sb.WriteString("            String error = null;\n")
	sb.WriteString("            long startTime = System.nanoTime();\n")
	sb.WriteString("            long startMem = Runtime.getRuntime().totalMemory() - Runtime.getRuntime().freeMemory();\n")
	sb.WriteString("            \n")
	sb.WriteString("            Future<Object> future = executor.submit(() -> {\n")
	sb.WriteString("                List<Object> inputList = (List<Object>) test.get(\"input\");\n")
	for j, param := range sig.Parameters {
		sb.WriteString(fmt.Sprintf("                %s %s = mapper.convertValue(inputList.get(%d), %s.class);\n", s.mapTypeToJava(param.Type), param.Name, j, s.mapTypeToJava(param.Type)))
	}
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("                return sol.%s(%s);\n", sig.FunctionName, strings.Join(paramNames, ", ")))
	sb.WriteString("            });\n")
	sb.WriteString("            \n")
	sb.WriteString("            try {\n")
	sb.WriteString("                output = future.get(2, TimeUnit.SECONDS);\n")
	sb.WriteString("                if (!compareOutputs(output, test.get(\"expected\"), validationType, mapper)) {\n")
	sb.WriteString("                    status = \"failed\";\n")
	sb.WriteString("                }\n")
	sb.WriteString("            } catch (TimeoutException e) {\n")
	sb.WriteString("                status = \"timeout\";\n")
	sb.WriteString("                future.cancel(true);\n")
	sb.WriteString("            } catch (Exception e) {\n")
	sb.WriteString("                status = \"runtime_error\";\n")
	sb.WriteString("                error = e.getMessage() != null ? e.getMessage() : e.toString();\n")
	sb.WriteString("            }\n")
	sb.WriteString("            \n")
	sb.WriteString("            long endTime = System.nanoTime();\n")
	sb.WriteString("            long endMem = Runtime.getRuntime().totalMemory() - Runtime.getRuntime().freeMemory();\n")
	sb.WriteString("            long timeMs = (endTime - startTime) / 1000000;\n")
	sb.WriteString("            long memoryKb = Math.max(0, (endMem - startMem) / 1024);\n")
	sb.WriteString("            \n")
	sb.WriteString("            Map<String, Object> res = new HashMap<>();\n")
	sb.WriteString("            res.put(\"test_id\", index + 1);\n")
	sb.WriteString("            res.put(\"status\", status);\n")
	sb.WriteString("            res.put(\"time_ms\", timeMs);\n")
	sb.WriteString("            res.put(\"memory_kb\", memoryKb);\n")
	sb.WriteString("            res.put(\"output\", output != null ? mapper.writeValueAsString(output) : \"\");\n")
	sb.WriteString("            res.put(\"expected\", mapper.writeValueAsString(test.get(\"expected\")));\n")
	sb.WriteString("            res.put(\"is_sample\", test.get(\"is_sample\"));\n")
	sb.WriteString("            res.put(\"error\", error);\n")
	sb.WriteString("            results.add(res);\n")
	sb.WriteString("        }\n")
	sb.WriteString("        System.out.println(mapper.writeValueAsString(results));\n")
	sb.WriteString("        executor.shutdownNow();\n")
	sb.WriteString("        System.exit(0);\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}

// C++ stub generator
func (s *CodeGenService) generateCppStub(sig ProblemSignature) (string, error) {
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

// C++ harness generator
func (s *CodeGenService) GenerateCppHarness(sig ProblemSignature, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder
	sb.WriteString("#include <iostream>\n#include <vector>\n#include <string>\n#include <chrono>\n#include <sys/resource.h>\n#include <signal.h>\n#include <setjmp.h>\n#include <algorithm>\n#include <nlohmann/json.hpp>\n\n")
	sb.WriteString("using namespace std;\nusing json = nlohmann::json;\n\n")
	sb.WriteString("jmp_buf jump_buffer;\n")
	sb.WriteString("void timeout_handler(int sig) { longjmp(jump_buffer, 1); }\n\n")

	sb.WriteString(userCode)
	sb.WriteString("\n\n")

	// Embed Test Cases
	sb.WriteString("// Embedded Test Cases\n")
	sb.WriteString("const char* TEST_CASES_JSON = R\"(")

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
	sb.WriteString(")\";\n\n")

	// Validation Logic
	sb.WriteString("bool compareOutputs(json actual, json expected, string valType) {\n")
	sb.WriteString("    if (valType == \"EXACT\") {\n")
	sb.WriteString("        return actual == expected;\n")
	sb.WriteString("    } else if (valType == \"UNORDERED\") {\n")
	sb.WriteString("        if (!actual.is_array() || !expected.is_array()) return actual == expected;\n")
	sb.WriteString("        json actual_sorted = actual;\n")
	sb.WriteString("        json expected_sorted = expected;\n")
	sb.WriteString("        sort(actual_sorted.begin(), actual_sorted.end());\n")
	sb.WriteString("        sort(expected_sorted.begin(), expected_sorted.end());\n")
	sb.WriteString("        return actual_sorted == expected_sorted;\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return actual == expected;\n")
	sb.WriteString("}\n\n")

	sb.WriteString("int main() {\n")
	sb.WriteString("    json test_cases = json::parse(TEST_CASES_JSON);\n")
	sb.WriteString("    json results = json::array();\n")
	sb.WriteString("    Solution sol;\n")
	sb.WriteString(fmt.Sprintf("    string validation_type = \"%s\";\n", validationType))
	sb.WriteString("    signal(SIGALRM, timeout_handler);\n\n")
	sb.WriteString("    for (int i = 0; i < test_cases.size(); ++i) {\n")
	sb.WriteString("        auto& test = test_cases[i];\n")
	sb.WriteString("        string status = \"passed\";\n")
	sb.WriteString("        json output = nullptr;\n")
	sb.WriteString("        string error = \"\";\n")
	sb.WriteString("        auto start_time = chrono::high_resolution_clock::now();\n")
	sb.WriteString("        struct rusage usage_start, usage_end;\n")
	sb.WriteString("        getrusage(RUSAGE_SELF, &usage_start);\n\n")
	sb.WriteString("        alarm(2);\n")
	sb.WriteString("        if (setjmp(jump_buffer) == 0) {\n")
	sb.WriteString("            try {\n")
	for j, param := range sig.Parameters {
		sb.WriteString(fmt.Sprintf("                auto %s = test[\"input\"][%d].get<%s>();\n", param.Name, j, s.mapTypeToCpp(param.Type)))
	}
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("                auto res = sol.%s(%s);\n", sig.FunctionName, strings.Join(paramNames, ", ")))
	sb.WriteString("                output = res;\n")
	sb.WriteString("                if (!compareOutputs(output, test[\"expected\"], validation_type)) {\n")
	sb.WriteString("                    status = \"failed\";\n")
	sb.WriteString("                }\n")
	sb.WriteString("            } catch (const exception& e) {\n")
	sb.WriteString("                status = \"runtime_error\";\n")
	sb.WriteString("                error = e.what();\n")
	sb.WriteString("            }\n")
	sb.WriteString("            alarm(0);\n")
	sb.WriteString("        } else {\n")
	sb.WriteString("            status = \"timeout\";\n")
	sb.WriteString("        }\n\n")
	sb.WriteString("        auto end_time = chrono::high_resolution_clock::now();\n")
	sb.WriteString("        getrusage(RUSAGE_SELF, &usage_end);\n")
	sb.WriteString("        auto time_ms = chrono::duration_cast<chrono::milliseconds>(end_time - start_time).count();\n")
	sb.WriteString("        auto memory_kb = usage_end.ru_maxrss;\n\n")
	sb.WriteString("        results.push_back({\n")
	sb.WriteString("            {\"test_id\", i + 1},\n")
	sb.WriteString("            {\"status\", status},\n")
	sb.WriteString("            {\"time_ms\", time_ms},\n")
	sb.WriteString("            {\"memory_kb\", memory_kb},\n")
	sb.WriteString("            {\"output\", output.is_null() ? \"\" : output.dump()},\n")
	sb.WriteString("            {\"expected\", test[\"expected\"].dump()},\n")
	sb.WriteString("            {\"is_sample\", test[\"is_sample\"]},\n")
	sb.WriteString("            {\"error\", error.empty() ? nullptr : error}\n")
	sb.WriteString("        });\n")
	sb.WriteString("    }\n")
	sb.WriteString("    cout << results.dump() << endl;\n")
	sb.WriteString("    return 0;\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}

// Go stub generator
func (s *CodeGenService) generateGoStub(sig ProblemSignature) (string, error) {
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
func (s *CodeGenService) GenerateGoHarness(sig ProblemSignature, userCode string, testCases []domain.TestCase, validationType string) (string, error) {
	var sb strings.Builder
	sb.WriteString("package main\n\nimport (\n	\"encoding/json\"\n	\"fmt\"\n	\"time\"\n	\"runtime\"\n	\"context\"\n	\"reflect\"\n	\"sort\"\n)\n\n")
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
	sb.WriteString("    if valType == \"EXACT\" {\n")
	sb.WriteString("        return reflect.DeepEqual(actual, expected)\n")
	sb.WriteString("    } else if valType == \"UNORDERED\" {\n")
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
	sb.WriteString("    for i, test := range testCases {\n")
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
	sb.WriteString("        case res := <-resChan:\n            output = res\n")
	sb.WriteString("            if !compareOutputs(output, test[\"expected\"], validationType) {\n")
	sb.WriteString("                status = \"failed\"\n")
	sb.WriteString("            }\n")
	sb.WriteString("        case err := <-errChan:\n            status = \"runtime_error\"\n            errStr = err.Error()\n")
	sb.WriteString("        case <-ctx.Done():\n            status = \"timeout\"\n")
	sb.WriteString("        }\n")
	sb.WriteString("        cancel()\n\n")
	sb.WriteString("        duration := time.Since(start)\n")
	sb.WriteString("        runtime.ReadMemStats(&ms)\n")
	sb.WriteString("        memKb := (ms.TotalAlloc - startAlloc) / 1024\n")
	sb.WriteString("        if memKb < 0 { memKb = 0 }\n\n")
	sb.WriteString("        outStr, _ := json.Marshal(output)\n")
	sb.WriteString("        expStr, _ := json.Marshal(test[\"expected\"])\n")
	sb.WriteString("        results = append(results, map[string]interface{}{\n")
	sb.WriteString("            \"test_id\": i + 1,\n")
	sb.WriteString("            \"status\": status,\n")
	sb.WriteString("            \"time_ms\": duration.Milliseconds(),\n")
	sb.WriteString("            \"memory_kb\": memKb,\n")
	sb.WriteString("            \"output\": string(outStr),\n")
	sb.WriteString("            \"expected\": string(expStr),\n")
	sb.WriteString("            \"is_sample\": test[\"is_sample\"],\n")
	sb.WriteString("            \"error\": errStr,\n")
	sb.WriteString("        })\n")
	sb.WriteString("    }\n")
	sb.WriteString("    finalData, _ := json.Marshal(results)\n")
	sb.WriteString("    fmt.Println(string(finalData))\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}
