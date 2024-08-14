import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// 生成一个长度n的随机token
export function generateToken(n:number):string {
  const buffer = crypto.getRandomValues(new Uint8Array(n))
  return Array.from(buffer).map(byte => byte.toString(16).padStart(2, '0')).join('');
}
// 生成uuid
export function generateId():string {
  return crypto.randomUUID().toString();
}
