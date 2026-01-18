export const filterEditorLanguage = (lang: string) => {

  switch (lang) {
    case "c++":
      lang = "cpp"
      break;
  }

  return lang || 'plaintext'
}