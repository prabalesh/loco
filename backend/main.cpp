#include <iostream>
#include <vector>
#include <string>
#include <chrono>
#include <sys/resource.h>
#include <signal.h>
#include <setjmp.h>
#include <algorithm>
#include <sstream>
#include <unistd.h>

using namespace std;

jmp_buf jump_buffer;
void timeout_handler(int sig) { longjmp(jump_buffer, 1); }

// Manual JSON serialization helpers
string escapeJSON(string s) {
    string res = "";
    for (char c : s) {
        if (c == '"') res += "\\\"";
        else if (c == '\\') res += "\\\\";
        else res += c;
    }
    return res;
}

string toJson(int v) { return to_string(v); }
string toJson(long v) { return to_string(v); }
string toJson(long long v) { return to_string(v); }
string toJson(bool v) { return v ? "true" : "false"; }
string toJson(string v) { return "\"" + escapeJSON(v) + "\""; }

template<typename T>
string toJson(const vector<T>& v) {
    string res = "[";
    for (size_t i = 0; i < v.size(); ++i) {
        res += toJson(v[i]);
        if (i < v.size() - 1) res += ",";
    }
    res += "]";
    return res;
}

// Minimal JSON Parser for Driver
struct JsonValue {
    string raw;
    vector<JsonValue> array;
    bool is_array = false;
};

JsonValue parseJson(istream& is) {
    JsonValue v; char c; while (is >> ws && is.get(c)) {
        if (c == '[') {
            v.is_array = true;
            while (is >> ws && is.peek() != ']') {
                v.array.push_back(parseJson(is));
                if (is >> ws && is.peek() == ',') is.get();
            }
            is.get(); return v;
        } else if (c == '{') {
            v.is_array = false; // Object as flat list of key-values in array
            while (is >> ws && is.peek() != '}') {
                v.array.push_back(parseJson(is)); // key
                if (is >> ws && is.peek() == ':') is.get();
                v.array.push_back(parseJson(is)); // value
                if (is >> ws && is.peek() == ',') is.get();
            }
            is.get(); return v;
        } else if (c == '"') {
            string s; char prev = 0;
            while (is.get(c)) {
                if (c == '"' && prev != '\\') break;
                s += c; prev = c;
            }
            v.raw = s; return v;
        } else {
            string s; s += c;
            while (is.peek() != EOF && !isspace(is.peek()) && is.peek() != ',' && is.peek() != ']' && is.peek() != '}') {
                is.get(c); s += c;
            }
            v.raw = s; return v;
        }
    }
    return v;
}

int asInt(JsonValue v) { return stoi(v.raw); }
bool asBool(JsonValue v) { return v.raw == "true"; }
string asString(JsonValue v) { return v.raw; }
vector<int> asIntArray(JsonValue v) {
    vector<int> res; for (auto& item : v.array) res.push_back(asInt(item)); return res;
}
vector<string> asStringArray(JsonValue v) {
    vector<string> res; for (auto& item : v.array) res.push_back(asString(item)); return res;
}

#include <iostream>
#include <vector>
#include <string>

using namespace std;

class Solution {
public:
    vector<int> twoSum(vector<int> nums, int target) {
        unordered_map<int, int> indexByValue;
        for (int i = 0; i < nums.size(); i++) {
            int complement = target - nums[i];

            if (indexByValue.count(complement)) {
                return { indexByValue[complement], i };
            }

            indexByValue[nums[i]] = i;
        }

        // Problem guarantees exactly one solution
        return {};
    }
};


struct TestResult {
    string status;
    long time_ms;
    long memory_kb;
    string output;
    string error;
    string input_description;
};

int main() {
    JsonValue root = parseJson(cin);
    if (!root.is_array) return 1;

    vector<TestResult> results;
    Solution sol;
    signal(SIGALRM, timeout_handler);

    for (auto& tc : root.array) {
        string status = "passed";
        string output_val = "";
        string error_msg = "";
        string input_desc = "";

        JsonValue inputObj, expectedVal;
        for(size_t k=0; k+1 < tc.array.size(); k+=2) {
            if(tc.array[k].raw == "input") inputObj = tc.array[k+1];
            if(tc.array[k].raw == "expected") expectedVal = tc.array[k+1];
        }

        vector<int> nums = asIntArray(inputObj.array[0]);
        input_desc += (input_desc.empty() ? "" : ", ") + toJson(nums);
        int target = asInt(inputObj.array[1]);
        input_desc += (input_desc.empty() ? "" : ", ") + toJson(target);
        vector<int> expected = asIntArray(expectedVal);
        auto start_time = chrono::high_resolution_clock::now();
        struct rusage usage_start, usage_end;
        getrusage(RUSAGE_SELF, &usage_start);

        alarm(5);
        if (setjmp(jump_buffer) == 0) {
            try {
                auto res = sol.twoSum(nums, target);
                output_val = toJson(res);
                if (res != expected) status = "failed";
            } catch (const exception& e) {
                status = "runtime_error";
                error_msg = e.what();
            } catch (...) {
                status = "runtime_error";
                error_msg = "Unknown error";
            }
            alarm(0);
        } else {
            status = "timeout";
        }

        auto end_time = chrono::high_resolution_clock::now();
        getrusage(RUSAGE_SELF, &usage_end);
        auto time_ms = chrono::duration_cast<chrono::milliseconds>(end_time - start_time).count();
        auto memory_kb = usage_end.ru_maxrss;

        results.push_back({status, (long)time_ms, (long)memory_kb, output_val, error_msg, "[" + input_desc + "]"});
    }

    // Standardized Verdict Aggregation
    string final_verdict = "ACCEPTED";
    long max_runtime = 0;
    long max_memory = 0;
    
    for (const auto& res : results) {
        if (res.status != "passed" && final_verdict == "ACCEPTED") {
            if (res.status == "timeout") final_verdict = "TLE";
            else if (res.status == "runtime_error") final_verdict = "RUNTIME_ERROR";
            else if (res.status == "failed") final_verdict = "WRONG_ANSWER";
            else final_verdict = res.status;
        }
        if (res.time_ms > max_runtime) max_runtime = res.time_ms;
        if (res.memory_kb > max_memory) max_memory = res.memory_kb;
    }

    // Output JSON manually
    cout << "{";
    cout << "\"verdict\":\"" + final_verdict + "\",";
    cout << "\"runtime\":" + to_string(max_runtime) + ",";
    cout << "\"memory\":" + to_string(max_memory) + ",";
    cout << "\"test_results\":[";
    for (size_t i = 0; i < results.size(); ++i) {
        cout << "{";
        cout << "\"passed\":" << (results[i].status == "passed" ? "true" : "false") << ",";
        cout << "\"input\":\"" << escapeJSON(results[i].input_description) << "\",";
        cout << "\"actual\":\"" << escapeJSON(results[i].output) << "\",";
        cout << "\"error\":\"" << escapeJSON(results[i].error) << "\"";
        cout << "}";
        if (i < results.size() - 1) cout << ",";
    }
    cout << "]}" << endl;
    return 0;
}
