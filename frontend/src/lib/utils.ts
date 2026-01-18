import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const filterEditorLanguage = (lang: string) => {

  switch (lang) {
    case "c++":
      lang = "cpp"
      break;
  }

  return lang || 'plaintext'
}