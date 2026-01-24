import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

export const calculateAcceptanceRate = (totalAccepted: number, totalSubmissions: number) => {
    return totalSubmissions === 0 ? 0 : (totalAccepted / totalSubmissions) * 100;
}