import http from "k6/http";
import {
  randomString,
  randomEmail,
  randomPassword,
} from "../../helpers/random";
import { check } from "k6";

export function runAuthenticationScenario(baseUrl: string) {
  const orgName = `org_${randomString(6)}`;
  const username = randomString(8);
  const password = randomPassword();

  let payload = JSON.stringify({
    name: orgName,
    schema_name: orgName,
    admin: {
      username: username,
      email: randomEmail("mycompany.com"),
      password: password,
    },
  });

  const createRes = http.post(`${baseUrl}/api/organizations/create/`, payload, {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
  });

  check(createRes, {
    "organization created": (r) =>
      r.status === 201 && !!r.json("organization.id"),
  });

  const orgId = createRes.json("organization.id");

  payload = JSON.stringify({ username: username, password: password });

  const loginRes = http.post(`${baseUrl}/api/auth/login/`, payload, {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      "X-Org-Id": orgId ? String(orgId) : "",
    },
  });

  check(loginRes, {
    "login successful": (r) => r.status === 200 && !!r.json("access"),
  });

  const accessToken = loginRes.json("access");
  const refreshToken = loginRes.json("refresh");

  payload = JSON.stringify({
    username: randomString(8),
    password: randomPassword(),
    email: randomEmail("mycompany.com"),
  });

  const registerUserRes = http.post(`${baseUrl}/api/users/create/`, payload, {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      Authorization: `Bearer ${accessToken}`,
    },
  });

  check(registerUserRes, {
    "user created": (r) => r.status === 201 && !!r.json("access"),
  });

  payload = JSON.stringify({
    token: accessToken,
  });

  const verifyTokenRes = http.post(
    `${baseUrl}/api/auth/verify-token/`,
    payload,
    {
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        "X-Org-Id": orgId ? String(orgId) : "",
      },
    },
  );

  check(verifyTokenRes, {
    "token is valid": (r) => r.status === 200 && !!r.json("valid") == true,
  });

  payload = JSON.stringify({
    refresh: refreshToken,
  });

  const refreshTokenRes = http.post(`${baseUrl}/api/auth/refresh/`, payload, {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      "X-Org-Id": orgId ? String(orgId) : "",
    },
  });

  check(refreshTokenRes, {
    "tokens refreshed": (r) => r.status === 200 && !!r.json("access"),
  });

  return accessToken;
}
