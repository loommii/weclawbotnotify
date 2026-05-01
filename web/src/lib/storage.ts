export function getItem<T = string>(key: string): T | null {
  try {
    const raw = localStorage.getItem(key)
    if (raw === null) return null
    return JSON.parse(raw) as T
  } catch {
    return null
  }
}

export function setItem(key: string, value: unknown): void {
  try {
    localStorage.setItem(key, JSON.stringify(value))
  } catch {
    // storage full or disabled
  }
}

export function removeItem(key: string): void {
  try {
    localStorage.removeItem(key)
  } catch {
    // storage disabled
  }
}
