import { useState, useEffect, useCallback } from "react";

type EditorTab = "function" | "main" | "solution";

type EditorContent = {
  function: string;
  main: string;
  solution: string;
};

type ProblemLanguage = {
  language_id: number;
  language_name: string;
  function_code: string;
  main_code: string;
  solution_code: string;
  is_validated: boolean;
};

export function useEditorState(
  selectedLangId: number | null,
  problemLanguages: ProblemLanguage[]
) {
  const [currentTab, setCurrentTab] = useState<EditorTab>("function");
  const [editorContent, setEditorContent] = useState<EditorContent>({
    function: "",
    main: "",
    solution: "",
  });

  useEffect(() => {
    if (!selectedLangId) return;

    const lang = problemLanguages.find((l) => l.language_id === selectedLangId);
    if (lang) {
      setEditorContent({
        function: lang.function_code,
        main: lang.main_code,
        solution: lang.solution_code,
      });
    }
  }, [selectedLangId, problemLanguages]);

  const updateEditorContent = useCallback((tab: EditorTab, code: string) => {
    setEditorContent((prev) => ({
      ...prev,
      [tab]: code,
    }));
  }, []);

  const resetEditorContent = useCallback((content: EditorContent) => {
    setEditorContent(content);
  }, []);

  return {
    currentTab,
    setCurrentTab,
    editorContent,
    updateEditorContent,
    resetEditorContent,
  };
}
