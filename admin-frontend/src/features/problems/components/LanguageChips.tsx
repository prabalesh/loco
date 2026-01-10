import { Box, Typography, Chip } from "@mui/material";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";

type ProblemLanguage = {
  language_id: number;
  language_name: string;
  is_validated: boolean;
};

type Props = {
  languages: ProblemLanguage[];
  selectedLangId: number | null;
  onSelectLanguage: (langId: number) => void;
  onDeleteLanguage: (langId: number) => void;
};

export function LanguageChips({
  languages,
  selectedLangId,
  onSelectLanguage,
  onDeleteLanguage,
}: Props) {
  if (languages.length === 0) return null;

  return (
    <Box>
      <Typography variant="subtitle2" gutterBottom>
        Added Languages:
      </Typography>
      <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
        {languages.map((lang) => (
          <Chip
            key={lang.language_id}
            label={lang.language_name}
            onClick={() => onSelectLanguage(lang.language_id)}
            onDelete={() => onDeleteLanguage(lang.language_id)}
            color={selectedLangId === lang.language_id ? "primary" : "default"}
            icon={lang.is_validated ? <CheckCircleOutlineIcon /> : undefined}
          />
        ))}
      </Box>
    </Box>
  );
}
