export function parseBool(s: string): boolean {
  switch (s?.toLowerCase()?.trim()) {
    case "true":
    case "yes":
    case "1":
      return true;
    case "false":
    case "no":
    case "0":
    case null:
    case undefined:
      return false;
    default:
      return JSON.parse(s);
  }
}
