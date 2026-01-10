// hooks/useProblemLanguages.ts
import { useState, useEffect, useCallback } from "react";
import { adminProblemLanguagesApi } from "../../../api/adminApi";
import type { ProblemLanguage, CreateProblemLanguageRequest, UpdateProblemLanguageRequest } from "../../../types";

export function useProblemLanguages(problemId: string | undefined) {
  const [problemLanguages, setProblemLanguages] = useState<ProblemLanguage[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  // Fetch problem languages on mount
  useEffect(() => {
    if (!problemId) {
      setLoading(false);
      return;
    }

    const fetchProblemLanguages = async () => {
      try {
        const response = await adminProblemLanguagesApi.getAll(Number(problemId));
        setProblemLanguages(response.data.data || []);
      } catch (err) {
        console.error("Failed to load problem languages:", err);
        setError("Failed to load problem languages");
        setProblemLanguages([]);
      } finally {
        setLoading(false);
      }
    };

    fetchProblemLanguages();
  }, [problemId]);

  // Save language (create or update)
  const saveLanguage = useCallback(async (
    langId: number,
    data: { function_code: string; main_code: string; solution_code: string }
  ) => {
    if (!problemId) {
      return { success: false, message: "Problem ID is missing" };
    }

    setSaving(true);
    setError("");

    try {
      const existingLang = problemLanguages.find((l) => l.language_id === langId);

      if (existingLang) {
        // UPDATE existing language
        const updateData: UpdateProblemLanguageRequest = {
          function_code: data.function_code,
          main_code: data.main_code,
          solution_code: data.solution_code,
        };

        const response = await adminProblemLanguagesApi.update(
          Number(problemId),
          langId,
          updateData
        );

        if (response.data.success) {
          // Update local state
          setProblemLanguages((prev) =>
            prev.map((l) =>
              l.language_id === langId
                ? { ...l, ...response.data.data }
                : l
            )
          );

          return { 
            success: true, 
            message: response.data.message || "Language updated successfully" 
          };
        } else {
          throw new Error(response.data.message || "Failed to update");
        }
      } else {
        // CREATE new language
        const createData: CreateProblemLanguageRequest = {
          language_id: langId,
          function_code: data.function_code,
          main_code: data.main_code,
          solution_code: data.solution_code,
        };

        const response = await adminProblemLanguagesApi.create(
          Number(problemId),
          createData
        );

        if (response.data.success) {
          // Add to local state
          setProblemLanguages((prev) => [...prev, response.data.data]);

          return { 
            success: true, 
            message: response.data.message || "Language added successfully" 
          };
        } else {
          throw new Error(response.data.message || "Failed to create");
        }
      }
    } catch (err: any) {
      console.error(err);
      const errorMessage = err.response?.data?.message || err.message || "Failed to save language";
      setError(errorMessage);
      return { success: false, message: errorMessage };
    } finally {
      setSaving(false);
    }
  }, [problemId, problemLanguages]);

  // Delete language
  const deleteLanguage = useCallback(async (langId: number) => {
    if (!problemId) {
      return { success: false, message: "Problem ID is missing" };
    }

    try {
      const response = await adminProblemLanguagesApi.delete(
        Number(problemId),
        langId
      );

      if (response.data.success) {
        // Remove from local state
        setProblemLanguages((prev) => prev.filter((l) => l.language_id !== langId));
        return { 
          success: true, 
          message: response.data.message || "Language removed successfully" 
        };
      } else {
        throw new Error(response.data.message || "Failed to delete");
      }
    } catch (err: any) {
      console.error(err);
      const errorMessage = err.response?.data?.message || err.message || "Failed to remove language";
      return { success: false, message: errorMessage };
    }
  }, [problemId]);

  // Validate language
  const validateLanguage = useCallback(async (langId: number) => {
    if (!problemId) {
      return { success: false, message: "Problem ID is missing" };
    }

    try {
      const response = await adminProblemLanguagesApi.validate(
        Number(problemId),
        langId
      );

      if (response.data.success) {
        // Update validation status in local state
        setProblemLanguages((prev) =>
          prev.map((l) =>
            l.language_id === langId
              ? { ...l, is_validated: response.data.data.is_validated }
              : l
          )
        );

        return { 
          success: true, 
          message: "Validation completed",
          data: response.data.data
        };
      } else {
        throw new Error(response.data.message || "Validation failed");
      }
    } catch (err: any) {
      console.error(err);
      const errorMessage = err.response?.data?.message || err.message || "Failed to validate language";
      return { success: false, message: errorMessage };
    }
  }, [problemId]);

  return {
    problemLanguages,
    loading,
    saving,
    error,
    saveLanguage,
    deleteLanguage,
    validateLanguage,
  };
}
