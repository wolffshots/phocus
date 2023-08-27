// $lib/stringUtils.ts

export function camelCaseToWords(camelCase: string): string {
    return camelCase
    .replace(/([a-z])([A-Z])/g, '$1 $2') // Insert space between lowercase and uppercase
    .replace(/([A-Z])([A-Z][a-z])/g, '$1 $2'); // Insert space between consecutive uppercase
}