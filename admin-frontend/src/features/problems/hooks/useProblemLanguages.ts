// hooks/useProblemLanguages.ts
import { useState, useEffect, useCallback } from "react";
import { adminProblemLanguagesApi } from "../../../lib/api/admin";
import type { CreateProblemLanguageRequest, UpdateProblemLanguageRequest } from "../../../types";
import type { ProblemLanguage as ProblemLanguageType } from "../../../types/problemLanguage";

export function useProblemLanguages(problemId: string | undefined) {
  const [problemLanguages, setProblemLanguages] = useState<ProblemLanguageType[]>([]);
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
        const response = await adminProblemLanguagesApi.getAll(String(problemId));
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
          String(problemId),
          langId,
          updateData
        );

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
          message: "Language updated successfully"
        };
      } else {
        // CREATE new language
        const createData: CreateProblemLanguageRequest = {
          language_id: langId,
          function_code: data.function_code,
          main_code: data.main_code,
          solution_code: data.solution_code,
        };

        const response = await adminProblemLanguagesApi.create(
          String(problemId),
          createData
        );

        // Add to local state
        setProblemLanguages((prev) => [...prev, response.data.data]);

        return {
          success: true,
          message: "Language added successfully"
        };
      }
    } catch (err: any) {
      console.error(err);
      const errorMessage = err.response?.data?.error?.message || err.message || "Failed to save language";
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
      await adminProblemLanguagesApi.delete(
        String(problemId),
        langId
      );

      // Remove from local state
      setProblemLanguages((prev) => prev.filter((l) => l.language_id !== langId));
      return {
        success: true,
        message: "Language removed successfully"
      };
    } catch (err: any) {
      console.error(err);
      const errorMessage = err.response?.data?.error?.message || err.message || "Failed to remove language";
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
        String(problemId),
        langId
      );

      const validationData = response.data.data;

      // Update validation status in local state
      setProblemLanguages((prev) =>
        prev.map((l) =>
          l.language_id === langId
            ? { ...l, is_validated: validationData.is_validated }
            : l
        )
      );

      return {
        success: true,
        message: "Validation completed",
        data: validationData
      };
    } catch (err: any) {
      console.error(err);
      const errorMessage = err.response?.data?.error?.message || err.message || "Failed to validate language";
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
