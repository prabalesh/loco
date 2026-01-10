import { useMemo, useCallback } from "react";
import type { Language } from "../../../types";

type ProblemLanguage = {
  language_id: number;
  language_name: string;
  function_code: string;
  main_code: string;
  solution_code: string;
  is_validated: boolean;
};

export function useLanguageManager(
  availableLanguages: Language[],
  problemLanguages: ProblemLanguage[]
) {
  const unusedLanguages = useMemo(() => {
    return availableLanguages.filter(
      (al) => !problemLanguages.some((pl) => pl.language_id === al.id)
    );
  }, [availableLanguages, problemLanguages]);

  const getLanguageById = useCallback(
    (langId: number) => {
      return (
        problemLanguages.find((l) => l.language_id === langId) ||
        availableLanguages.find((l) => l.id === langId)
      );
    },
    [problemLanguages, availableLanguages]
  );

  const getEditorLanguage = useCallback((langName: string) => {
    const map: Record<string, string> = {
      c: "c",
      "c++": "cpp",
      cpp: "cpp",
      python: "python",
      java: "java",
    };
    return map[langName.toLowerCase()] || "plaintext";
  }, []);

  return {
    unusedLanguages,
    getLanguageById,
    getEditorLanguage,
  };
}
