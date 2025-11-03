import { AuthContext } from "../types/auth_context";

/**
 * Builds authentication headers for GraphQL requests
 * @param accessToken - JWT token (for gateway requests)
 * @param authContext - Auth context (for direct router/service requests)
 * @param useDirectHeaders - Whether to include direct headers (x-user-sub, x-org-id, x-org-name)
 */
export function buildAuthHeaders(
  accessToken?: string,
  authContext?: AuthContext,
  useDirectHeaders: boolean = true
): Record<string, string> {
  const headers: Record<string, string> = {};

  // Add JWT token for gateway requests
  if (accessToken) {
    headers["Authorization"] = `Bearer ${accessToken}`;
  }

  // Add direct headers for router/service requests (when JWT isn't being processed by gateway)
  if (useDirectHeaders && authContext) {
    headers["x-org-id"] = authContext.organization.id;
    headers["x-user-sub"] = authContext.admin.id;
    headers["x-org-name"] = authContext.organization.name;
  }

  return headers;
}