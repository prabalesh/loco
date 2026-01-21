package codegen

import (
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
func (s *CodeGenService) GenerateTestHarness(signature ProblemSignature, userCode string, languageSlug string) (string, error) {
	// Validate
	if userCode == "" {
		return "", errors.New("user code is required")
	}

	customTypes := s.identifyCustomTypes(signature)

	switch languageSlug {
	case "python":
		return s.generatePythonHarness(signature, userCode, customTypes)
	case "javascript":
		return s.generateJavaScriptHarness(signature, userCode, customTypes)
	case "java":
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for Java yet")
		}
		return s.generateJavaHarness(signature, userCode)
	case "cpp":
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for C++ yet")
		}
		return s.generateCppHarness(signature, userCode)
	case "go":
		if len(customTypes) > 0 {
			return "", errors.New("custom types not supported for Go yet")
		}
		return s.generateGoHarness(signature, userCode)
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
func (s *CodeGenService) generatePythonHarness(sig ProblemSignature, userCode string, customTypes []string) (string, error) {
	var sb strings.Builder
	sb.WriteString("import json\nimport sys\n\n")

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
	sb.WriteString("# Test harness\n")
	sb.WriteString("if __name__ == \"__main__\":\n")
	sb.WriteString("    test_cases = json.loads(sys.stdin.read())\n")
	sb.WriteString("    results = []\n    \n")
	sb.WriteString("    for test in test_cases:\n")
	sb.WriteString("        try:\n")
	sb.WriteString("            # Unpack parameters\n")

	// Parameter deserialization
	for i, param := range sig.Parameters {
		if param.IsCustom {
			deserializerFunc := fmt.Sprintf("deserialize_%s", strings.ToLower(param.Type))
			sb.WriteString(fmt.Sprintf("            %s = %s(test['input'][%d])\n", param.Name, deserializerFunc, i))
		} else {
			sb.WriteString(fmt.Sprintf("            %s = test['input'][%d]\n", param.Name, i))
		}
	}

	sb.WriteString("            \n")
	sb.WriteString("            # Call user function\n")
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("            result = %s(%s)\n", sig.FunctionName, strings.Join(paramNames, ", ")))

	// Serialize result if custom
	if s.isCustomType(sig.ReturnType) {
		serializerFunc := fmt.Sprintf("serialize_%s", strings.ToLower(sig.ReturnType))
		sb.WriteString(fmt.Sprintf("            result = %s(result)\n", serializerFunc))
	}

	sb.WriteString("            results.append({'output': result, 'error': None})\n")
	sb.WriteString("        except Exception as e:\n")
	sb.WriteString("            results.append({'output': None, 'error': str(e)})\n")
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
func (s *CodeGenService) generateJavaScriptHarness(sig ProblemSignature, userCode string, customTypes []string) (string, error) {
	var sb strings.Builder

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
	sb.WriteString("// Test harness\n")
	sb.WriteString("const testCases = JSON.parse(require('fs').readFileSync(0, 'utf-8'));\n")
	sb.WriteString("const results = [];\n\n")
	sb.WriteString("for (const test of testCases) {\n")
	sb.WriteString("    try {\n")

	// Parameter deserialization
	for i, param := range sig.Parameters {
		if param.IsCustom {
			deserializerFunc := fmt.Sprintf("deserialize%s", param.Type) // e.g. deserializeTreeNode
			sb.WriteString(fmt.Sprintf("        const %s = %s(test.input[%d]);\n", param.Name, deserializerFunc, i))
		} else {
			sb.WriteString(fmt.Sprintf("        const %s = test.input[%d];\n", param.Name, i))
		}
	}

	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("        let result = %s(%s);\n", sig.FunctionName, strings.Join(paramNames, ", ")))

	// Serialize result if custom
	if s.isCustomType(sig.ReturnType) {
		serializerFunc := fmt.Sprintf("serialize%s", sig.ReturnType) // e.g. serializeTreeNode
		sb.WriteString(fmt.Sprintf("        result = %s(result);\n", serializerFunc))
	}

	sb.WriteString("        results.push({ output: result, error: null });\n")
	sb.WriteString("    } catch (e) {\n")
	sb.WriteString("        results.push({ output: null, error: e.message });\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n\n")
	sb.WriteString("console.log(JSON.stringify(results));\n")
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
func (s *CodeGenService) generateJavaHarness(sig ProblemSignature, userCode string) (string, error) {
	var sb strings.Builder
	sb.WriteString("import java.util.*;\n")
	sb.WriteString("import com.fasterxml.jackson.databind.ObjectMapper;\n\n")
	sb.WriteString(userCode)
	sb.WriteString("\n\n")
	sb.WriteString("public class Main {\n")
	sb.WriteString("    public static void main(String[] args) throws Exception {\n")
	sb.WriteString("        Scanner sc = new Scanner(System.in);\n")
	sb.WriteString("        StringBuilder inputSb = new StringBuilder();\n")
	sb.WriteString("        while (sc.hasNextLine()) inputSb.append(sc.nextLine());\n")
	sb.WriteString("        \n")
	sb.WriteString("        ObjectMapper mapper = new ObjectMapper();\n")
	sb.WriteString("        List<Map<String, Object>> testCases = mapper.readValue(inputSb.toString(), List.class);\n")
	sb.WriteString("        List<Map<String, Object>> results = new ArrayList<>();\n")
	sb.WriteString("        Solution sol = new Solution();\n")
	sb.WriteString("        \n")
	sb.WriteString("        for (Map<String, Object> test : testCases) {\n")
	sb.WriteString("            try {\n")
	sb.WriteString("                List<Object> inputList = (List<Object>) test.get(\"input\");\n")
	for i, param := range sig.Parameters {
		sb.WriteString(fmt.Sprintf("                %s %s = (%s) inputList.get(%d);\n", s.mapTypeToJava(param.Type), param.Name, s.mapTypeToJava(param.Type), i))
	}
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("                Object result = sol.%s(%s);\n", sig.FunctionName, strings.Join(paramNames, ", ")))
	sb.WriteString("                Map<String, Object> res = new HashMap<>();\n")
	sb.WriteString("                res.put(\"output\", result);\n")
	sb.WriteString("                res.put(\"error\", null);\n")
	sb.WriteString("                results.add(res);\n")
	sb.WriteString("            } catch (Exception e) {\n")
	sb.WriteString("                Map<String, Object> res = new HashMap<>();\n")
	sb.WriteString("                res.put(\"output\", null);\n")
	sb.WriteString("                res.put(\"error\", e.getMessage());\n")
	sb.WriteString("                results.add(res);\n")
	sb.WriteString("            }\n")
	sb.WriteString("        }\n")
	sb.WriteString("        System.out.println(mapper.writeValueAsString(results));\n")
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
func (s *CodeGenService) generateCppHarness(sig ProblemSignature, userCode string) (string, error) {
	var sb strings.Builder
	sb.WriteString("#include <iostream>\n#include <vector>\n#include <string>\n#include <nlohmann/json.hpp>\n\nusing namespace std;\nusing json = nlohmann::json;\n\n")
	sb.WriteString(userCode)
	sb.WriteString("\n\n")
	sb.WriteString("int main() {\n")
	sb.WriteString("    json test_cases;\n")
	sb.WriteString("    cin >> test_cases;\n")
	sb.WriteString("    json results = json::array();\n")
	sb.WriteString("    Solution sol;\n\n")
	sb.WriteString("    for (auto& test : test_cases) {\n")
	sb.WriteString("        try {\n")
	for i, param := range sig.Parameters {
		sb.WriteString(fmt.Sprintf("            auto %s = test[\"input\"][%d].get<%s>();\n", param.Name, i, s.mapTypeToCpp(param.Type)))
	}
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("            auto result = sol.%s(%s);\n", sig.FunctionName, strings.Join(paramNames, ", ")))
	sb.WriteString("            results.push_back({{\"output\", result}, {\"error\", nullptr}});\n")
	sb.WriteString("        } catch (const exception& e) {\n")
	sb.WriteString("            results.push_back({{\"output\", nullptr}, {\"error\", e.what()}});\n")
	sb.WriteString("        }\n")
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
func (s *CodeGenService) generateGoHarness(sig ProblemSignature, userCode string) (string, error) {
	var sb strings.Builder
	sb.WriteString("package main\n\nimport (\n	\"encoding/json\"\n	\"fmt\"\n	\"io/ioutil\"\n	\"os\"\n)\n\n")
	sb.WriteString(userCode)
	sb.WriteString("\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    data, _ := ioutil.ReadAll(os.Stdin)\n")
	sb.WriteString("    var testCases []map[string]interface{}\n")
	sb.WriteString("    json.Unmarshal(data, &testCases)\n")
	sb.WriteString("    results := []map[string]interface{}{}\n\n")
	sb.WriteString("    for _, test := range testCases {\n")
	sb.WriteString("        input := test[\"input\"].([]interface{})\n")
	for i, param := range sig.Parameters {
		sb.WriteString(fmt.Sprintf("        %s := input[%d].(%s)\n", param.Name, i, s.mapTypeToGo(param.Type)))
	}
	paramNames := []string{}
	for _, param := range sig.Parameters {
		paramNames = append(paramNames, param.Name)
	}
	sb.WriteString(fmt.Sprintf("        result := %s(%s)\n", sig.FunctionName, strings.Join(paramNames, ", ")))
	sb.WriteString("        results = append(results, map[string]interface{}{\"output\": result, \"error\": nil})\n")
	sb.WriteString("    }\n")
	sb.WriteString("    out, _ := json.Marshal(results)\n")
	sb.WriteString("    fmt.Println(string(out))\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}
