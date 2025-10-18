import http from "k6/http";
import { check } from "k6";
import { randomString, randomEmail, randomPassword } from "./random";
import { AuthContext } from "src/types/auth_context";

export function createOrganization(baseUrl: string): AuthContext {
  const orgName = `org_${randomString(6)}`;
  const username = randomString(8);
  const password = randomPassword();
  const email = randomEmail("mycompany.com");

  const payload = JSON.stringify({
    name: orgName,
    schema_name: orgName,
    admin: {
      username: username,
      email: email,
      password: password,
    },
  });

  const response = http.post(`${baseUrl}/api/organizations/create/`, payload, {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
  });

  check(response, {
    "organization created": (r) =>
      r.status === 201 && !!r.json("organization.id"),
  });

  const orgId = response.json("organization.id")?.toString();
  const schemaName = response
    .json("organization.schema_name")
    ?.toLocaleString();
  const adminId = response.json("admin_user.id")?.toString();

  if (!orgId || !adminId || !schemaName) {
    throw new Error("Failed to create organization");
  }

  return {
    organization: {
      id: orgId,
      name: orgName,
      schemaName: schemaName,
    },
    admin: {
      id: adminId,
      username: username,
      password: password,
      email: email,
    },
  };
}
