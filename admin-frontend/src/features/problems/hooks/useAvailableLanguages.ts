import { useState, useEffect } from "react";
import { adminLanguagesApi } from "../../../lib/api/admin";
import type { Language } from "../../../types";

export function useAvailableLanguages() {
  const [availableLanguages, setAvailableLanguages] = useState<Language[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchLanguages = async () => {
      try {
        const response = await adminLanguagesApi.getAllActive();
        setAvailableLanguages(response.data.data);
      } catch (err) {
        setError("Failed to load languages");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    console.log("rendering")
    fetchLanguages();
  }, []);

  return { availableLanguages, loading, error };
}
