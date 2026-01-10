import { FormControl, InputLabel, Select, MenuItem, Box, Typography } from "@mui/material";
import type { Language } from "../../../types";

type Props = {
  unusedLanguages: Language[];
  onSelectLanguage: (langId: number) => void;
  hasAddedLanguages: boolean;
};

export function LanguageSelector({ unusedLanguages, onSelectLanguage, hasAddedLanguages }: Props) {
  return (
    <Box sx={{ display: "flex", gap: 2, alignItems: "center", mb: 2 }}>
      <FormControl sx={{ minWidth: 250 }}>
        <InputLabel>Add a language</InputLabel>
        <Select
          value=""
          onChange={(e) => {
            const langId = Number(e.target.value);
            if (langId) onSelectLanguage(langId);
          }}
          label="Add a language"
          disabled={unusedLanguages.length === 0}
        >
          {unusedLanguages.map((lang) => (
            <MenuItem key={lang.id} value={lang.id}>
              {lang.name} ({lang.version})
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {unusedLanguages.length === 0 && hasAddedLanguages && (
        <Typography variant="body2" color="text.secondary">
          All languages added
        </Typography>
      )}
    </Box>
  );
}
