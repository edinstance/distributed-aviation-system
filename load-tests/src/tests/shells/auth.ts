import http from "k6/http";
import {
  randomString,
  randomEmail,
  randomPassword,
} from "../../helpers/random";
import { createOrganization } from "../../helpers/organization";
import { check } from "k6";

export function runAuthenticationScenario(baseUrl: string) {
  const authContext = createOrganization(baseUrl);

  let payload = JSON.stringify({
    username: authContext.admin.username,
    password: authContext.admin.password,
  });

  const loginRes = http.post(`${baseUrl}/api/auth/login/`, payload, {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      "X-Org-Id": authContext.organization.id,
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
        "X-Org-Id": authContext.organization.id,
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
      "X-Org-Id": authContext.organization.id,
    },
  });

  check(refreshTokenRes, {
    "tokens refreshed": (r) => r.status === 200 && !!r.json("access"),
  });

  return accessToken;
}
