export function randomString(length = 6) {
  const chars = "abcdefghijklmnopqrstuvwxyz";
  return Array.from({ length }, () =>
    chars.charAt(Math.floor(Math.random() * chars.length)),
  ).join("");
}

export function randomEmail(domain = "example.com") {
  return `${randomString(5)}.${randomString(5)}@${domain}`;
}

export function randomPassword() {
  const randomNum = Math.floor(1000 + Math.random() * 9000);
  return `${randomString(6)}${randomNum}`;
}
